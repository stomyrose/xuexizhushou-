package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"force-learning/internal/model"
	"force-learning/internal/repository"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	subscriptionRepo *repository.SubscriptionRepository
	userRepo         *repository.UserRepository
	alipayConfig     *AlipayConfig
	wxpayConfig      *WxpayConfig
}

type AlipayConfig struct {
	AppID           string
	PrivateKey      string
	AlipayPublicKey string
	NotifyURL       string
	Sandbox         bool
}

type WxpayConfig struct {
	AppID     string
	MchID     string
	APIKey    string
	CertPath  string
	KeyPath   string
	NotifyURL string
	Sandbox   bool
}

type PaymentResult struct {
	Success      bool
	OrderID      string
	PaymentURL   string
	PaymentData  map[string]interface{}
	ErrorMessage string
}

type PayOrder struct {
	OrderID       string     `json:"order_id"`
	UserID        uuid.UUID  `json:"user_id"`
	PlanID        uuid.UUID  `json:"plan_id"`
	Amount        float64    `json:"amount"`
	PaymentMethod string     `json:"payment_method"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
}

func NewPaymentService(
	subscriptionRepo *repository.SubscriptionRepository,
	userRepo *repository.UserRepository,
	alipayConfig *AlipayConfig,
	wxpayConfig *WxpayConfig,
) *PaymentService {
	return &PaymentService{
		subscriptionRepo: subscriptionRepo,
		userRepo:         userRepo,
		alipayConfig:     alipayConfig,
		wxpayConfig:      wxpayConfig,
	}
}

func (s *PaymentService) CreatePayment(userIDStr, planIDStr, paymentMethod string) (*PaymentResult, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		return nil, errors.New("invalid plan ID")
	}

	return s.CreatePaymentOrder(userID, planID, paymentMethod)
}

func (s *PaymentService) CreatePaymentOrder(userID, planID uuid.UUID, paymentMethod string) (*PaymentResult, error) {
	plan, err := s.subscriptionRepo.FindPlanByID(planID)
	if err != nil {
		return nil, errors.New("plan not found")
	}

	if !plan.IsActive {
		return nil, errors.New("plan is not available")
	}

	order := &PayOrder{
		OrderID:       generateOrderID(),
		UserID:        userID,
		PlanID:        planID,
		Amount:        plan.Price,
		PaymentMethod: paymentMethod,
		Status:        "pending",
		CreatedAt:     time.Now(),
	}

	var paymentURL string
	var paymentData map[string]interface{}

	switch paymentMethod {
	case "alipay":
		paymentURL, paymentData, err = s.createAlipayOrder(order)
	case "wxpay":
		paymentURL, paymentData, err = s.createWxpayOrder(order)
	default:
		return nil, errors.New("unsupported payment method")
	}

	if err != nil {
		return &PaymentResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	return &PaymentResult{
		Success:     true,
		OrderID:     order.OrderID,
		PaymentURL:  paymentURL,
		PaymentData: paymentData,
	}, nil
}

func (s *PaymentService) createAlipayOrder(order *PayOrder) (string, map[string]interface{}, error) {
	if s.alipayConfig == nil {
		return "", nil, errors.New("alipay not configured")
	}

	bizContent := map[string]interface{}{
		"out_trade_no":    order.OrderID,
		"product_code":    "FAST_INSTANT_TRADE_PAY",
		"total_amount":    fmt.Sprintf("%.2f", order.Amount),
		"subject":         fmt.Sprintf("Force Learning Subscription - %d days", s.getPlanDuration(order.PlanID)),
		"timeout_express": "30m",
	}

	bizContentJSON, _ := json.Marshal(bizContent)

	sign := s.generateAlipaySign(bizContentJSON)

	data := map[string]interface{}{
		"app_id":      s.alipayConfig.AppID,
		"method":      "alipay.trade.app.pay",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"notify_url":  s.alipayConfig.NotifyURL,
		"biz_content": string(bizContentJSON),
		"sign":        sign,
	}

	var payURL strings.Builder
	payURL.WriteString("https://openapi.alipay.com/gateway.do?")
	for key, value := range data {
		payURL.WriteString(fmt.Sprintf("%s=%s&", key, fmt.Sprintf("%v", value)))
	}

	return payURL.String(), data, nil
}

func (s *PaymentService) createWxpayOrder(order *PayOrder) (string, map[string]interface{}, error) {
	if s.wxpayConfig == nil {
		return "", nil, errors.New("wechat pay not configured")
	}

	nonceStr := generateNonceStr(32)

	payRequest := map[string]interface{}{
		"appid":            s.wxpayConfig.AppID,
		"mch_id":           s.wxpayConfig.MchID,
		"nonce_str":        nonceStr,
		"body":             fmt.Sprintf("Force Learning Subscription - %d days", s.getPlanDuration(order.PlanID)),
		"out_trade_no":     order.OrderID,
		"total_fee":        int(order.Amount * 100),
		"spbill_create_ip": "10.0.0.1",
		"notify_url":       s.wxpayConfig.NotifyURL,
		"trade_type":       "APP",
	}

	sign := s.generateWxpaySign(payRequest)
	payRequest["sign"] = sign

	return "", payRequest, nil
}

func (s *PaymentService) generateAlipaySign(bizContent []byte) string {
	if s.alipayConfig == nil {
		return ""
	}

	sortedParams := []string{
		fmt.Sprintf("app_id=%s", s.alipayConfig.AppID),
		fmt.Sprintf("biz_content=%s", string(bizContent)),
		fmt.Sprintf("charset=utf-8"),
		fmt.Sprintf("method=alipay.trade.app.pay"),
		fmt.Sprintf("sign_type=RSA2"),
		fmt.Sprintf("timestamp=%s", time.Now().Format("2006-01-02 15:04:05")),
		fmt.Sprintf("version=1.0"),
	}
	sort.Strings(sortedParams)

	signString := strings.Join(sortedParams, "&")

	hash := md5.Sum([]byte(signString + s.alipayConfig.PrivateKey))
	return hex.EncodeToString(hash[:])
}

func (s *PaymentService) generateWxpaySign(params map[string]interface{}) string {
	if s.wxpayConfig == nil {
		return ""
	}

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signParts []string
	for _, k := range keys {
		signParts = append(signParts, fmt.Sprintf("%s=%v", k, params[k]))
	}
	signString := strings.Join(signParts, "&") + "&key=" + s.wxpayConfig.APIKey

	hash := md5.Sum([]byte(signString))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func (s *PaymentService) getPlanDuration(planID uuid.UUID) int {
	plan, err := s.subscriptionRepo.FindPlanByID(planID)
	if err != nil {
		return 0
	}
	return plan.DurationDays
}

func (s *PaymentService) VerifyAlipayCallback(params map[string]string) (bool, string) {
	if params["trade_status"] != "TRADE_SUCCESS" {
		return false, params["trade_status"]
	}

	sign := params["sign"]
	delete(params, "sign")
	delete(params, "sign_type")

	expectedSign := s.generateAlipaySignFromParams(params)
	return sign == expectedSign, params["out_trade_no"]
}

func (s *PaymentService) generateAlipaySignFromParams(params map[string]string) string {
	if s.alipayConfig == nil {
		return ""
	}

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signParts []string
	for _, k := range keys {
		signParts = append(signParts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	signString := strings.Join(signParts, "&")

	hash := md5.Sum([]byte(signString + s.alipayConfig.PrivateKey))
	return hex.EncodeToString(hash[:])
}

func (s *PaymentService) VerifyWxpayCallback(params map[string]string) (bool, string) {
	sign := params["sign"]
	delete(params, "sign")

	expectedSign := s.generateWxpaySignFromParams(params)
	if sign != expectedSign {
		return false, ""
	}

	if params["trade_state"] != "SUCCESS" {
		return false, params["trade_state"]
	}

	return true, params["out_trade_no"]
}

func (s *PaymentService) generateWxpaySignFromParams(params map[string]string) string {
	if s.wxpayConfig == nil {
		return ""
	}

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signParts []string
	for _, k := range keys {
		signParts = append(signParts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	signString := strings.Join(signParts, "&") + "&key=" + s.wxpayConfig.APIKey

	hash := md5.Sum([]byte(signString))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func (s *PaymentService) HandlePaymentCallback(orderID string) error {
	order, err := s.subscriptionRepo.FindOrderByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status == "paid" {
		return nil
	}

	order.Status = "paid"
	now := time.Now()
	order.PaidAt = &now

	if err := s.subscriptionRepo.UpdateOrder(order); err != nil {
		return err
	}

	plan, err := s.subscriptionRepo.FindPlanByID(order.PlanID)
	if err != nil {
		return err
	}

	existingSubscription, _ := s.subscriptionRepo.FindActiveByUserID(order.UserID)
	if existingSubscription != nil {
		existingSubscription.Status = "expired"
		s.subscriptionRepo.Update(existingSubscription)
	}

	subscription := &model.Subscription{
		UserID:    order.UserID,
		PlanID:    order.PlanID,
		StartDate: now,
		EndDate:   now.AddDate(0, 0, plan.DurationDays),
		Status:    "active",
	}

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return err
	}

	s.userRepo.UpdateRemainingDays(order.UserID, plan.DurationDays)

	return nil
}

func generateOrderID() string {
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), generateNonceStr(8))
}

func generateNonceStr(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
	}
	return string(result)
}
