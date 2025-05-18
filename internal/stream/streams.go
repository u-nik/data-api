package stream

import (
	"context"
	"data-api/internal/handlers"
	"log"
	"time"

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
			logger := zap.L().Sugar()
			if h.Subscribe(ctx, *Context, logger) {
				logger.Infof("Subscribers registered for handler '%s'", name)
				return
			}
		}(handler, name)
	}
}
