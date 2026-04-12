package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rodziievskyi-maksym/go-genesis-case-task/internal/infrastructure/github"
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
