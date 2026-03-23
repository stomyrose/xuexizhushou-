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

type AlipayCallback struct {
	TradeStatus string `json:"trade_status"`
	OutTradeNo  string `json:"out_trade_no"`
	TradeNo     string `json:"trade_no"`
	TotalAmount string `json:"total_amount"`
	Sign        string `json:"sign"`
}

type WxpayCallback struct {
	TradeState    string `json:"trade_state"`
	OutTradeNo    string `json:"out_trade_no"`
	TransactionID string `json:"transaction_id"`
	TotalFee      string `json:"total_fee"`
	Sign          string `json:"sign"`
}

type SyncRequest struct {
	LastSyncTime string          `json:"last_sync_time"`
	Records      []SyncRecordDTO `json:"records"`
}

type SyncRecordDTO struct {
	ClientID        string `json:"client_id"`
	FileID          string `json:"file_id"`
	DurationSeconds int    `json:"duration_seconds"`
	LearnedAt       string `json:"learned_at"`
}

type SyncResponse struct {
	SyncedRecords []SyncedRecordDTO `json:"synced_records"`
	ServerTime    string            `json:"server_time"`
	HasMore       bool              `json:"has_more"`
}

type SyncedRecordDTO struct {
	ClientID  string `json:"client_id"`
	ServerID  string `json:"server_id"`
	Synced    bool   `json:"synced"`
	LearnedAt string `json:"learned_at"`
}
