package invitations

import (
	_ "data-api/api"
	"data-api/internal/auth"
	"data-api/internal/const/scopes"
	"data-api/internal/const/subjects"
	"data-api/internal/db"
	"data-api/internal/entities/invitation"
	"data-api/internal/events"
	"data-api/internal/schema"
	"data-api/internal/stream"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
)

func (h InvitationsHandler) SetupRoutes(api *gin.RouterGroup) {
	// The "/invitations/accept" endpoint is used to accept an invitation.
	// It must not be protected by authentication.
	api.POST(
		"/invitations/accept",
		schema.JSONSchemaValidator("invitations-accept"),
		h.Accept,
	)

	g := api.Group("/invitations")
	g.Use(auth.Auth())
	{
		// Retrieve all users.
		g.GET(
			"/",
			auth.RequireScope(scopes.Invitations.List),
			h.List,
		)

		// Retrieve user data by ID.
		g.GET(
			"/:id",
			auth.RequireScope(scopes.Invitations.Read),
			h.Get,
		)

		// Invite a new user.
		g.POST(
			"/",
			auth.RequireScope(scopes.Invitations.Create),
			schema.JSONSchemaValidator("invitations-create"),
			h.Create,
		)
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
func (h InvitationsHandler) List(c *gin.Context) {
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
func (h InvitationsHandler) Get(c *gin.Context) {
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
func (h InvitationsHandler) Create(c *gin.Context) {
	var input events.InvitationCreateData
	if err := schema.ShouldBindValidInput(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind valid input", "details": err.Error()})
		return
	}

	var event = events.EventFactory(func(b events.BaseEvent) events.InvitationCreateEvent {
		return events.InvitationCreateEvent{
			BaseEvent: b,
			Data:      input,
		}
	})

	data, err := sonic.Marshal(event)
	if err != nil {
		h.Logger.Errorw("Failed to marshal event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal event"})
		return
	}

	// Publish the event to the NATS subject.
	_, err = (*stream.Context).Publish(subjects.Invitations.Create, data, nats.AckWait(5*time.Second))
	if err != nil {
		h.Logger.Errorw("Failed to publish event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish event", "details": err.Error()})
		return
	}

	h.Logger.Debugw("Published event", "id", event.ID, "subject", subjects.Invitations.Create)

	// Respond with the created event.
	c.JSON(http.StatusAccepted, event)
}

// @Summary      Accept an invitation
// @Description  Accept an invitation and publish an event to NATS
// @Tags         admin, users
// @Security     BearerAuth
// @Param        id  path      string  true  "User ID"
// @Param        user  body      UserInput  true  "User data"
// @Success      202  {object} UserCreatedEvent
// @Failure      400  {string} string  "bad request"
// @Failure      401  {string} string  "unauthorized"
// @Failure      403  {string} string  "forbidden"
// @Failure      500  {string} string  "internal server error"
// @Router       /admin/users/{id}/accept [post]
func (h InvitationsHandler) Accept(c *gin.Context) {
	var input events.InvitationAcceptData
	if err := schema.ShouldBindValidInput(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind valid input", "details": err.Error()})
		return
	}

	// Validate the token instead of bearer token.
	inviteRepo := invitation.NewInvitationRepository(db.DB)
	_, err := inviteRepo.FindByToken(h.Ctx, input.Token)
	if err != nil {
		h.Logger.Errorw("Failed to find invitation by token", "token", input.Token, "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
		return
	}

	event := events.EventFactory(func(b events.BaseEvent) events.InvitationAcceptEvent {
		return events.InvitationAcceptEvent{
			BaseEvent: b,
			Data:      input,
		}
	})

	data, err := sonic.Marshal(event)
	if err != nil {
		h.Logger.Errorw("Failed to marshal event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal event"})
		return
	}

	// Create NATS message with headers for idempotent publishing
	msg := nats.NewMsg(subjects.Invitations.Accept)
	msg.Data = data
	msg.Header = nats.Header{}
	msg.Header.Set(nats.MsgIdHdr, input.Token) // Use token as idempotency key

	// Publish the event to the NATS subject
	_, err = (*stream.Context).PublishMsg(msg, nats.AckWait(5*time.Second))
	if err != nil {
		h.Logger.Errorw("Failed to publish event", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish event", "details": err.Error()})
		return
	}

	h.Logger.Debugw("Published invitation accept event", "subject", subjects.Invitations.Accept)
	c.JSON(http.StatusAccepted, gin.H{"message": "Invitation accepted, processing..."})
}
