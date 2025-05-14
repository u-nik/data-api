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

const StreamName = "data-api"

// Initialize sets up the NATS JetStream connection and ensures the stream exists with the correct subjects.
func Initialize(url string, handlerMap map[string]handlers.HandlerInterface) {
	// NATS Initialization
	nc, err := nats.Connect(url, nats.Timeout(5*time.Second)) // Connect to the NATS server.
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	zap.L().Sugar().Infoln("Connected to NATS server:", url)

	subjects := handlers.GetAllSubjects(handlerMap) // Get all subjects from the handler map.
	Context = SetupStream(nc, StreamName, subjects) // Set up the stream with the subjects.
}

// SetupStream creates or updates a JetStream stream with the given name and subjects.
// If the stream already exists, only the subjects are updated.
func SetupStream(nc *nats.Conn, name string, subjects []string) *nats.JetStreamContext {
	zap.L().Sugar().Infof("Setting up stream '%s' with subjects: %s", name, subjects)

	stream, err := nc.JetStream() // Initialize JetStream context.
	if err != nil {
		log.Fatalf("Failed to initialize JetStream: %v", err)
	}

	// Check if the stream already exists
	info, err := stream.StreamInfo(StreamName)
	if err == nil && info != nil {
		// Stream exists, update subjects if needed
		info.Config.Subjects = subjects
		_, err = stream.UpdateStream(&info.Config)
		if err != nil {
			log.Fatalf("Failed to update stream subjects: %v", err)
		}
		zap.L().Sugar().Infoln("Stream updated:", info.Config.Name)
		return &stream
	}

	// If stream does not exist, create it
	info, err = stream.AddStream(&nats.StreamConfig{
		Name:     StreamName, // Name of the stream.
		Subjects: subjects,   // Subjects associated with the stream.
	})
	if err != nil {
		log.Fatalf("Failed to add stream: %v", err)
	}

	zap.L().Sugar().Infoln("Stream created:", info.Config.Name)

	return &stream
}

// RegisterSubscribers subscribes to all handler subjects and stores received events in Redis as the read model.
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
					zap.L().Sugar().Errorf("Error unmarshaling event: %v", err)
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
					zap.L().Sugar().Errorf("Error marshaling event: %v", err)
					return
				}

				if err := rdb.Set(ctx, name+":"+id, val, 0).Err(); err != nil {
					zap.L().Sugar().Errorf("Error storing in Redis: %v", err)
					return
				}

				zap.L().Sugar().Debugf("Stored %s data in Redis: %s", name, id)
			}, nats.Durable("read-model-"+name))

			if err != nil {
				zap.L().Sugar().Errorf("Error subscribing to subject %s: %v", subject, err)
				return
			}
			zap.L().Sugar().Infof("Subscribed to subject %s with ID %s", subject, sub.Subject)
		}(handler, name)
	}
}
