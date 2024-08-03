package rabbitmq

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/streadway/amqp"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"strconv"
	"sync"
	"time"
)

type RabbitMQConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	ExchangeName string
	Kind         string
	Retry        string
	Concurrency  int
}

type rabbitMQ struct {
	cfg            *RabbitMQConfig
	ctx            context.Context
	Conn           *amqp.Connection
	log            logger.ILogger
	lockConnection sync.Mutex
}

type IRabbitMQ interface {
	NewRabbitMQConn(ctx context.Context, log logger.ILogger) error
	GetConn() *amqp.Connection
}

func NewRabbitMQ(cfg *RabbitMQConfig, log logger.ILogger, ctx context.Context) IRabbitMQ {
	return &rabbitMQ{
		cfg: cfg,
		log: log,
		ctx: ctx,
	}
}

func (r *rabbitMQ) GetConn() *amqp.Connection {
	return r.Conn
}

func (r *rabbitMQ) NewRabbitMQConn(ctx context.Context, log logger.ILogger) error {
	r.lockConnection.Lock()
	defer r.lockConnection.Unlock()

	if r.Conn != nil && !r.Conn.IsClosed() {
		return nil
	}

	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		r.cfg.User,
		r.cfg.Password,
		r.cfg.Host,
		r.cfg.Port,
	)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second         // Maximum time to retry
	maxRetries, err := strconv.Atoi(r.cfg.Retry) // Number of retries (including the initial attempt)
	if err != nil {
		maxRetries = 5
	}

	var conn *amqp.Connection

	err = backoff.Retry(func() error {

		conn, err = amqp.Dial(connAddr)
		if err != nil {
			log.Errorf("Failed to connect to RabbitMQ: %v. Connection information: %s", err, connAddr)
			return err
		}

		return nil
	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1)))

	log.Debug("Connected to RabbitMQ")
	r.Conn = conn

	go func() {
		select {
		case <-ctx.Done():
			err = conn.Close()
			if err != nil {
				log.Error("Failed to close RabbitMQ connection")
			}
			log.Debug("RabbitMQ connection is closed")
		}
	}()

	return err
}
