package user

import (
	"data-api/internal/middleware"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

func (h UserHandler) SetupRoutes(api *gin.RouterGroup) {
	users := api.Group("/users") // Create a new route group for user-related endpoints.
	{
		users.GET("/", h.ListUsers) // GET /api/users - Retrieve all users.

		// GET /api/users/:id - Retrieve user data by ID.
		users.GET("/:id", h.GetUser)

		// POST /api/users - Create a new user.
		users.POST("/", middleware.JSONSchemaValidator("user"), h.CreateUser)
	}
}

func (h UserHandler) ListUsers(c *gin.Context) {
	// Retrieve all user data from Redis.
	data, err := h.Rdb.Keys(h.Ctx, "user:*").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the list of user keys.
	c.JSON(http.StatusOK, data)
}

func (h UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id") // Extract the user ID from the URL parameter.

	// Retrieve user data from Redis.
	data, err := h.Rdb.Get(h.Ctx, "user:"+id).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		// Respond with the user data.
		c.Data(http.StatusOK, "application/json", []byte(data))
	}
}

func (h UserHandler) CreateUser(c *gin.Context) {
	// Generate a unique ID for the user.
	uuidObj, err := uuid.NewV7()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate UUID"})
		return
	}

	input, err := h.GetInputFromContext(c)
	if err != nil {
		return
	}

	// Create a empty UserRegistered event.
	var event = UserCreated{
		ID:        uuidObj.String(),
		Name:      input["name"].(string),
		Email:     input["email"].(string),
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	data, err := sonic.Marshal(event)
	if err != nil {
		h.Logger.Errorw("Failed to marshal event",
			"error", err,
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal event"})
		return
	}

	// Publish the event to the NATS subject.
	_, err = h.Stream.Publish(h.Subject, data, nats.AckWait(5*time.Second))
	if err != nil {
		h.Logger.Errorw("Failed to publish event",
			"error", err,
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish event", "details": err.Error()})
		return
	}

	h.Logger.Debugw("Published event",
		"id", uuidObj.String(),
	)

	// Respond with the created event.
	c.JSON(http.StatusAccepted, event)
}
