package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type UserRegistered struct {
	ID    string `json:"id" jsonschema:"uuid"`
	Email string `json:"email" jsonschema:"email"`
}

var (
	ctx     = context.Background()
	rdb     *redis.Client
	js      nats.JetStreamContext
	subject = "user.events"
)

func main() {
	r := gin.Default()

	// Redis Init
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	// NATS Init
	nc, _ := nats.Connect("nats:4222")
	js, _ = nc.JetStream()

	// Ensure Stream exists
	js.AddStream(&nats.StreamConfig{
		Name:     "USER",
		Subjects: []string{subject},
	})

	// Commands (Write)
	r.POST("/users", func(c *gin.Context) {
		var input struct {
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		uuidObj, err := uuid.NewV7()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate UUID"})
			return
		}
		id := uuidObj.String()
		event := UserRegistered{ID: id, Email: input.Email}
		data, _ := json.Marshal(event)
		js.Publish(subject, data)

		c.JSON(http.StatusAccepted, event)
	})

	// Queries (Read)
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		data, err := rdb.Get(ctx, "user:"+id).Result()
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.Data(http.StatusOK, "application/json", []byte(data))
		}
	})

	// Event-Handler (in-process)
	go func() {
		js.Subscribe(subject, func(msg *nats.Msg) {
			var evt UserRegistered
			json.Unmarshal(msg.Data, &evt)

			// Apply event to read model (Redis)
			val, _ := json.Marshal(evt)
			rdb.Set(ctx, "user:"+evt.ID, val, 0)
		}, nats.Durable("read-model"))
	}()

	r.Run(":8080")
}
