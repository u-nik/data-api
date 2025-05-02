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
)

var Context *nats.JetStreamContext

func Initialize(url string) {
	// NATS Initialization
	nc, err := nats.Connect(url, nats.Timeout(5*time.Second)) // Connect to the NATS server.
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	log.Println("Connected to NATS server:", url)

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
	log.Println("Stream info:", info)
}

func RegisterSubscribers(
	ctx context.Context,
	rdb *redis.Client,
	handlerMap map[string]handlers.HandlerInterface,
) {
	for name, handler := range handlerMap {
		go func(h handlers.HandlerInterface, name string) {
			subject := h.GetSubject()
			log.Println("Register message subscriber for package:", name)
			sub, err := (*Context).Subscribe(subject, func(msg *nats.Msg) {
				var evt map[string]interface{}
				if err := json.Unmarshal(msg.Data, &evt); err != nil {
					log.Printf("Error unmarshaling event: %v", err)
					return
				}

				id, ok := evt["id"].(string)
				if !ok {
					log.Println("Event has no ID field or ID is not a string")
					return
				}

				// Store the event data in Redis as the read model.
				val, err := sonic.Marshal(evt)
				if err != nil {
					log.Printf("Error marshaling event: %v", err)
					return
				}

				if err := rdb.Set(ctx, name+":"+id, val, 0).Err(); err != nil {
					log.Printf("Error storing in Redis: %v", err)
					return
				}

				log.Printf("Stored %s data in Redis: %s", name, id)
			}, nats.Durable("read-model-"+name))

			if err != nil {
				log.Printf("Error subscribing to subject %s: %v", subject, err)
				return
			}
			log.Printf("Subscribed to subject %s with ID %s", subject, sub.Subject)
		}(handler, name)
	}
}
