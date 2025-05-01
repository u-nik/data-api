package user

import (
	"data-api/internal/schema"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func (h UserHandler) SetupRoutes(api *gin.RouterGroup) {
	users := api.Group("/users") // Create a new route group for user-related endpoints.
	{
		users.GET("/:id", h.GetUser)  // GET /api/users/:id - Retrieve user data by ID.
		users.POST("/", h.CreateUser) // POST /api/users - Create a new user.
	}
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
	var input struct {
		Email string `json:"email"` // Input structure for user email.
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a unique ID for the user.
	uuidObj, err := uuid.NewV7()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate UUID"})
		return
	}
	id := uuidObj.String()

	// Create a UserRegistered event.
	event := UserRegistered{ID: id, Email: input.Email}
	data, _ := sonic.Marshal(event)

	h.Logger.Debugw("Validating user",
		"data", data,
	)

	err = schema.GetManager().ValidateJSON("user", data) // Validate the input against the JSON schema.
	if err != nil {
		h.Logger.Errorw("Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.Logger.Debugw("Publishing event",
		"data", string(data),
	)

	// Publish the event to the NATS subject.
	h.Stream.Publish(h.Subject, data)

	// Respond with the created event.
	c.JSON(http.StatusAccepted, event)
}
