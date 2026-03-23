package handler

import (
	"force-learning/internal/model"
	"force-learning/internal/pkg/response"
	"force-learning/internal/service"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
	paymentService      *service.PaymentService
}

func NewSubscriptionHandler(subscriptionService *service.SubscriptionService, paymentService *service.PaymentService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		paymentService:      paymentService,
	}
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

func (h *SubscriptionHandler) CreatePayment(c *gin.Context) {
	var req struct {
		PlanID        string `json:"plan_id" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		response.Unauthorized(c, "user not found")
		return
	}

	result, err := h.paymentService.CreatePayment(userIDStr, req.PlanID, req.PaymentMethod)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if !result.Success {
		response.BadRequest(c, result.ErrorMessage)
		return
	}

	response.Success(c, result)
}

func (h *SubscriptionHandler) AlipayCallback(c *gin.Context) {
	var params model.AlipayCallback
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	paramsMap := map[string]string{
		"trade_status": params.TradeStatus,
		"out_trade_no": params.OutTradeNo,
		"trade_no":     params.TradeNo,
		"total_amount": params.TotalAmount,
	}

	valid, orderID := h.paymentService.VerifyAlipayCallback(paramsMap)
	if !valid {
		response.BadRequest(c, "invalid signature")
		return
	}

	if err := h.paymentService.HandlePaymentCallback(orderID); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	c.JSON(200, gin.H{"code": "success", "msg": "payment processed"})
}

func (h *SubscriptionHandler) WxpayCallback(c *gin.Context) {
	var params model.WxpayCallback
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	paramsMap := map[string]string{
		"trade_state":    params.TradeState,
		"out_trade_no":   params.OutTradeNo,
		"transaction_id": params.TransactionID,
		"total_fee":      params.TotalFee,
	}

	valid, orderID := h.paymentService.VerifyWxpayCallback(paramsMap)
	if !valid {
		response.BadRequest(c, "invalid signature")
		return
	}

	if err := h.paymentService.HandlePaymentCallback(orderID); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	c.JSON(200, gin.H{"code": "success", "msg": "payment processed"})
}
