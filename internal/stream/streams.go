package stream

import (
	"context"
	"data-api/internal/handlers"
	"encoding/json"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var Context *nats.JetStreamContext

func Initialize(url string) {
	// NATS Initialization
	nc, err := nats.Connect(url, nats.Timeout(5*time.Second)) // Connect to the NATS server.
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	zap.L().Sugar().Infoln("Connected to NATS server:", url)

	stream, err := nc.JetStream() // Initialize JetStream context.
	if err != nil {
		log.Fatalf("Failed to initialize JetStream: %v", err)
	}

	// Ensure the NATS stream exists.
	info, err := stream.AddStream(&nats.StreamConfig{
		Name:     "data-api",              // Name of the stream.
		Subjects: []string{"user.events"}, // Subjects associated with the stream.
	})
	if err != nil {
		log.Fatalf("Failed to add stream: %v", err)
	}

	Context = &stream // Store the JetStream context for later use.
	zap.L().Sugar().Infoln("Stream info:", info.Sources)
}

func RegisterSubscribers(
	ctx context.Context,
	rdb *redis.Client,
	handlerMap map[string]handlers.HandlerInterface,
) {
	for name, handler := range handlerMap {
		go func(h handlers.HandlerInterface, name string) {
			subject := h.GetSubject()
			zap.L().Sugar().Infoln("Register message subscriber for package:", name)
			sub, err := (*Context).Subscribe(subject, func(msg *nats.Msg) {
				var evt map[string]interface{}
				if err := json.Unmarshal(msg.Data, &evt); err != nil {
					zap.L().Sugar().Errorln("Error unmarshaling event: %v", err)
					return
				}

				id, ok := evt["id"].(string)
				if !ok {
					zap.L().Sugar().Infoln("Event has no ID field or ID is not a string")
					return
				}

				// Store the event data in Redis as the read model.
				val, err := sonic.Marshal(evt)
				if err != nil {
					zap.L().Sugar().Errorln("Error marshaling event: %v", err)
					return
				}

				if err := rdb.Set(ctx, name+":"+id, val, 0).Err(); err != nil {
					zap.L().Sugar().Errorln("Error storing in Redis: %v", err)
					return
				}

				zap.L().Sugar().Debugln("Stored %s data in Redis: %s", name, id)
			}, nats.Durable("read-model-"+name))

			if err != nil {
				zap.L().Sugar().Errorln("Error subscribing to subject %s: %v", subject, err)
				return
			}
			zap.L().Sugar().Infoln("Subscribed to subject %s with ID %s", subject, sub.Subject)
		}(handler, name)
	}
}
