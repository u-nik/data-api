package user

import (
	_ "data-api/api"
	"data-api/internal/auth"
	"data-api/internal/schema"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

func (h UserHandler) SetupRoutes(api *gin.RouterGroup) {
	users := api.Group("/admin/users")
	users.Use(auth.Auth())
	users.Use(auth.RequireScope("admin.users.read")) // Create a new route group for user-related endpoints.
	{
		// Retrieve all users.
		users.GET("/", h.ListUsers)

		// Retrieve user data by ID.
		users.GET("/:id", h.GetUser)

		// Create a new user.
		users.POST("/", auth.RequireScope("admin.users.create"), schema.JSONSchemaValidator("user"), h.CreateUser)
	}
}

// @Summary      List all users
// @Description  List all users from Redis database
// @Tags         admin, users
// @Security     BearerAuth
// @Router       /admin/users [get]
// @Success      200  {string} string  "ok"
// @Failure      401  {string} string  "unauthorized"
// @Failure      403  {string} string  "forbidden"
// @Failure      500  {string} string  "internal server error"
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

// @Summary      Get user by ID
// @Description  Get user data by ID from Redis database
// @Tags         admin, users
// @Security     BearerAuth
// @Param        id  path      string  true  "User ID"
// @Success      200  {string} string  "ok"
// @Failure      404  {string} string  "not found"
// @Failure      401  {string} string  "unauthorized"
// @Failure      403  {string} string  "forbidden"
// @Failure      500  {string} string  "internal server error"
// @Router       /admin/users/{id} [get]
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

// @Summary      Create a new user
// @Description  Create a new user and publish an event to NATS
// @Tags         admin, users
// @Security     BearerAuth
// @Param        user  body      UserInput  true  "User data"
// @Success      202  {object} UserCreatedEvent
// @Failure      400  {string} string  "bad request"
// @Failure      401  {string} string  "unauthorized"
// @Failure      403  {string} string  "forbidden"
// @Failure      500  {string} string  "internal server error"
// @Router       /admin/users [post]
func (h UserHandler) CreateUser(c *gin.Context) {
	// Generate a unique ID for the user.
	uuidObj, err := uuid.NewV7()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate UUID"})
		return
	}

	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a empty UserRegistered event.
	var event = UserCreatedEvent{
		ID:        uuidObj.String(),
		UserInput: input,
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
