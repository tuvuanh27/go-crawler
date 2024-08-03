package rabbitmq

import (
	"context"
	"flag"
	"fmt"
	"github.com/iancoleman/strcase"
	jsoniter "github.com/json-iterator/go"
	"github.com/streadway/amqp"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/pkg/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"runtime"
	"strings"
	"time"
)

//go:generate mockery --name IConsumer
type IConsumer[T any] interface {
	ConsumeMessage(msg interface{}, dependencies T) error
	worker(deliveries <-chan amqp.Delivery, dependencies T, consumerHandlerName, queueName string)
}

var consumerTag = flag.String("consumer-tag", "simple-consumer", "consumer tag")

type Consumer[T any] struct {
	cfg          *RabbitMQConfig
	mq           IRabbitMQ
	log          logger.ILogger
	handler      func(queue string, msg amqp.Delivery, dependencies T) error
	jaegerTracer trace.Tracer
	ctx          context.Context
	concurrency  int
}

func (c Consumer[T]) ConsumeMessage(msg interface{}, dependencies T) error {

	strName := strings.Split(runtime.FuncForPC(reflect.ValueOf(c.handler).Pointer()).Name(), ".")
	var consumerHandlerName = strName[len(strName)-1]

	conn := c.mq.GetConn()

	ch, err := conn.Channel()
	if err != nil {
		c.log.Error("Error in opening channel to consume message")
		return err
	}

	typeName := reflect.TypeOf(msg).Name()
	snakeTypeName := strcase.ToSnake(typeName)

	err = ch.ExchangeDeclare(
		snakeTypeName, // name
		c.cfg.Kind,    // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		c.log.Error("Error in declaring exchange to consume message")
		return err
	}

	q, err := ch.QueueDeclare(
		fmt.Sprintf("%s_%s", snakeTypeName, "queue"), // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		c.log.Error("Error in declaring queue to consume message")
		return err
	}

	err = ch.QueueBind(
		q.Name,        // queue name
		snakeTypeName, // routing key
		snakeTypeName, // exchange
		false,
		nil)
	if err != nil {
		c.log.Error("Error in binding queue to consume message")
		return err
	}

	deliveries, err := ch.Consume(
		q.Name,       // queue
		*consumerTag, // consumer
		false,        // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)

	if err != nil {
		c.log.Error("Error in consuming message")
		return err
	}

	c.log.Infof("Waiting for messages in queue :%s. To exit press CTRL+C", q.Name)

	// Create a worker pool
	for i := 0; i < c.concurrency; i++ {
		go c.worker(deliveries, dependencies, consumerHandlerName, q.Name)
	}

	// Block forever
	select {
	case <-c.ctx.Done():
		defer func(ch *amqp.Channel) {
			err := ch.Close()
			if err != nil {
				c.log.Errorf("failed to close channel for queue: %s", q.Name)
			}
		}(ch)
		c.log.Infof("channel closed for queue: %s", q.Name)
	}

	return nil
}

func (c Consumer[T]) worker(deliveries <-chan amqp.Delivery, dependencies T, consumerHandlerName string, queueName string) {

	for {
		select {
		case <-c.ctx.Done():
			c.log.Infof("context done for queue: %s", queueName)
			return

		case delivery, ok := <-deliveries:
			{
				if !ok {
					c.log.Errorf("NOT OK deliveries channel closed for queue: %s", queueName)
					return
				}

				// Extract headers
				c.ctx = otel.ExtractAMQPHeaders(c.ctx, delivery.Headers)

				err := c.handler(queueName, delivery, dependencies)
				if err != nil {
					c.log.Error(err.Error())
				}

				_, span := c.jaegerTracer.Start(c.ctx, consumerHandlerName)

				h, err := jsoniter.Marshal(delivery.Headers)

				if err != nil {
					c.log.Errorf("Error in marshalling headers in consumer: %v", string(h))
				}

				span.SetAttributes(attribute.Key("message-id").String(delivery.MessageId))
				span.SetAttributes(attribute.Key("correlation-id").String(delivery.CorrelationId))
				span.SetAttributes(attribute.Key("queue").String(queueName))
				span.SetAttributes(attribute.Key("exchange").String(delivery.Exchange))
				span.SetAttributes(attribute.Key("routing-key").String(delivery.RoutingKey))
				span.SetAttributes(attribute.Key("ack").Bool(true))
				span.SetAttributes(attribute.Key("timestamp").String(delivery.Timestamp.String()))
				span.SetAttributes(attribute.Key("body").String(string(delivery.Body)))
				span.SetAttributes(attribute.Key("headers").String(string(h)))

				// Cannot use defer inside a for loop
				time.Sleep(1 * time.Millisecond)
				span.End()

				err = delivery.Ack(false)
				c.log.Infof("Ack for delivery: %v", string(delivery.Body))
				if err != nil {
					c.log.Errorf("We didn't get a ack for delivery: %v", string(delivery.Body))
				}
			}
		}
	}
}

func NewConsumer[T any](ctx context.Context, cfg *RabbitMQConfig, mq IRabbitMQ, log logger.ILogger, jaegerTracer trace.Tracer, handler func(queue string, msg amqp.Delivery, dependencies T) error) IConsumer[T] {
	return &Consumer[T]{ctx: ctx, cfg: cfg, mq: mq, log: log, jaegerTracer: jaegerTracer, handler: handler, concurrency: cfg.Concurrency}
}
