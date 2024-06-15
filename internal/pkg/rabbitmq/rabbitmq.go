package rabbitmq

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strconv"
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
}

func NewRabbitMQConn(cfg *RabbitMQConfig, ctx context.Context) (*amqp.Connection, error) {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second       // Maximum time to retry
	maxRetries, err := strconv.Atoi(cfg.Retry) // Number of retries (including the initial attempt)
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

	log.Info("Connected to RabbitMQ")

	go func() {
		select {
		case <-ctx.Done():
			err := conn.Close()
			if err != nil {
				log.Error("Failed to close RabbitMQ connection")
			}
			log.Info("RabbitMQ connection is closed")
		}
	}()

	return conn, err
}
