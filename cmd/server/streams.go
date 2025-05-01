package main

import (
	"context"
	"data-api/internal/handlers"
	"encoding/json"
	"log"

	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
)

func SetupStream(url string) nats.JetStreamContext {
	// NATS Initialization
	nc, err := nats.Connect(url) // Connect to the NATS server.
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	stream, err := nc.JetStream() // Initialize JetStream context.
	if err != nil {
		log.Fatalf("Failed to initialize JetStream: %v", err)
	}

	// Ensure the NATS stream exists.
	_, err = stream.AddStream(&nats.StreamConfig{
		Name:     "DATA-API",    // Name of the stream.
		Subjects: []string{"*"}, // Subjects associated with the stream.
	})
	if err != nil {
		log.Fatalf("Failed to add stream: %v", err)
	}

	return stream
}

func RegisterSubscribers(
	stream nats.JetStreamContext,
	ctx context.Context,
	rdb *redis.Client,
	handlerMap map[string]handlers.HandlerInterface,
) {
	for name, handler := range handlerMap {
		go func(h handlers.HandlerInterface, name string) {
			subject := h.GetSubject()
			log.Println("Register message subscriber for package:", name)
			stream.Subscribe(subject, func(msg *nats.Msg) {
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
		}(handler, name)
	}
}
