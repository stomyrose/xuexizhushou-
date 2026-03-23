package handler

import (
	"force-learning/internal/model"
	"force-learning/internal/pkg/response"
	"force-learning/internal/service"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionService: subscriptionService}
}

func (h *SubscriptionHandler) GetPlans(c *gin.Context) {
	plans, err := h.subscriptionService.GetPlans()
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, plans)
}

func (h *SubscriptionHandler) Purchase(c *gin.Context) {
	var req model.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "user not found")
		return
	}

	subscription, err := h.subscriptionService.PurchaseByIDString(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, subscription)
}

func (h *SubscriptionHandler) GetCurrent(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "user not found")
		return
	}

	subscription, err := h.subscriptionService.GetCurrentByIDString(userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, subscription)
}
