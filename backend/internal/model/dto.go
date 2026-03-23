package model

type RegisterRequest struct {
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Password string  `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Password string  `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserStatusResponse struct {
	ID              string `json:"id"`
	Email           string `json:"email,omitempty"`
	Phone           string `json:"phone,omitempty"`
	RemainingDays   int    `json:"remaining_days"`
	IsActive        bool   `json:"is_active"`
	HasSubscription bool   `json:"has_subscription"`
}

type PurchaseRequest struct {
	PlanID string `json:"plan_id" binding:"required"`
}
