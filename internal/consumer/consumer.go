package consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Handler func(event SearchEvent)

type Consumer struct {
	nc      *nats.Conn
	js      jetstream.JetStream
	stream  string
	subject string
	handler Handler
	logger  *slog.Logger
}

func New(
	nc *nats.Conn,
	stream, subject string,
	handler Handler,
	logger *slog.Logger,
) (*Consumer, error) {
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		nc:      nc,
		js:      js,
		stream:  stream,
		subject: subject,
		handler: handler,
		logger:  logger,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	stream, err := c.js.Stream(ctx, c.stream)
	if err != nil {
		return err
	}

	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          "search_top_consumer",
		FilterSubject: c.subject,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxDeliver:    3,
	})
	if err != nil {
		return err
	}

	_, err = cons.Consume(func(msg jetstream.Msg) {
		var event SearchEvent
		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			c.logger.Warn("failed to unmarshal event",
				"error", err,
				"data", string(msg.Data()),
			)
			msg.Nak()
			return
		}

		if event.Query == "" || event.Timestamp == 0 {
			c.logger.Warn("invalid event: missing required fields")
			msg.Nak()
			return
		}

		age := time.Now().Unix() - event.Timestamp
		if age > 600 {
			c.logger.Warn("dropping stale event", "age_secs", age)
			msg.Term()
			return
		}

		c.handler(event)
		msg.Ack()
	})

	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
