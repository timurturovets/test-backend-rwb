package consumer

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

func EnsureStream(ctx context.Context, js jetstream.JetStream, stream, subject string) error {
	_, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     stream,
		Subjects: []string{subject},
		MaxAge:   10 * 60 * 1e9, // 10 mins in nanosecs
	})
	return err
}
