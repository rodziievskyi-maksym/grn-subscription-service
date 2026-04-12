package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/domain"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/github"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/repository"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/usecase"
)

type SubscriptionHandler struct {
	useCase   usecase.SubscriptionUseCaseContract
	validator *validator.Validate
}

func NewSubscriptionHandler(
	useCase usecase.SubscriptionUseCaseContract,
	validator *validator.Validate,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		useCase:   useCase,
		validator: validator,
	}
}

type subscribeRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Repository string `json:"repository" binding:"required"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// Subscribe godoc
// @Summary      Subscribe to repository updates
// @Description  Creates a new subscription or reactivates an existing one to track GitHub releases
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request   body      subscribeRequest  true  "Subscription details"
// @Param        X-API-KEY header    string            true  "API Key"
// @Success      201       {object}  domain.Subscription
// @Failure      400       {object}  map[string]string "Invalid JSON or email format"
// @Failure      401       {object}  map[string]string "Unauthorized"
// @Failure      404       {object}  map[string]string "Repository not found on GitHub"
// @Failure      500       {object}  map[string]string "Internal server error"
// @Router       /api/v1/subscribe [post]
func (h *SubscriptionHandler) Subscribe(ctx *gin.Context) {
	var req subscribeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	subscription, err := h.useCase.Subscribe(ctx.Request.Context(), req.Email, req.Repository)
	if err != nil {
		if errors.Is(err, github.ErrRepositoryNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))

			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	ctx.JSON(http.StatusCreated, subscription)
}

// Unsubscribe godoc
// @Summary      Unsubscribe from repository updates
// @Description  Deactivates an active subscription (Soft delete)
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request   body      subscribeRequest  true  "Subscription details"
// @Param        X-API-KEY header    string            true  "API Key"
// @Success      204       "No Content - Successfully unsubscribed"
// @Failure      400       {object}  map[string]string "Invalid JSON"
// @Failure      401       {object}  map[string]string "Unauthorized"
// @Failure      404       {object}  map[string]string "Subscription not found"
// @Failure      500       {object}  map[string]string "Internal server error"
// @Router       /api/v1/subscribe [delete]
func (h *SubscriptionHandler) Unsubscribe(ctx *gin.Context) {
	var req subscribeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	if err := h.useCase.Unsubscribe(ctx.Request.Context(), req.Email, req.Repository); err != nil {
		if errors.Is(err, repository.ErrSubscriptionNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))

			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetSubscriptions godoc
// @Summary      Get user subscriptions
// @Description  Returns a list of active repository subscriptions for a given email
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        email  query     string  true  "User Email"
// @Param        X-API-KEY header    string  true  "API Key"
// @Success      200  {array}   domain.Subscription
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/subscriptions [get]
func (h *SubscriptionHandler) GetSubscriptions(c *gin.Context) {
	email := c.Query("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("email parameter is required")))

		return
	}

	subs, err := h.useCase.GetSubscriptionsByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("failed to retrieve subscriptions")))

		return
	}

	if subs == nil {
		subs = []domain.Subscription{}
	}

	c.JSON(http.StatusOK, subs)
}
