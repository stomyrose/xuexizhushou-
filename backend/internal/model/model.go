package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email         *string    `gorm:"uniqueIndex;size:255" json:"email"`
	Phone         *string    `gorm:"uniqueIndex;size:20" json:"phone"`
	PasswordHash  string     `gorm:"size:255;not null" json:"-"`
	CreatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	RemainingDays int        `gorm:"default:3" json:"remaining_days"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type SubscriptionPlan struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	DurationDays int       `gorm:"not null" json:"duration_days"`
	Price        float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (sp *SubscriptionPlan) BeforeCreate(tx *gorm.DB) error {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}

type Subscription struct {
	ID        uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	PlanID    uuid.UUID        `gorm:"type:uuid;not null" json:"plan_id"`
	StartDate time.Time        `gorm:"not null" json:"start_date"`
	EndDate   time.Time        `gorm:"not null" json:"end_date"`
	Status    string           `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	User      User             `gorm:"foreignKey:UserID" json:"-"`
	Plan      SubscriptionPlan `gorm:"foreignKey:PlanID" json:"-"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type KnowledgeFile struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Filename   string    `gorm:"size:255;not null" json:"filename"`
	FilePath   string    `gorm:"size:512;not null" json:"file_path"`
	FileType   string    `gorm:"size:10;not null" json:"file_type"`
	Category   *string   `gorm:"size:100" json:"category"`
	IsVisible  bool      `gorm:"default:true" json:"is_visible"`
	UploadedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"uploaded_at"`
}

func (kf *KnowledgeFile) BeforeCreate(tx *gorm.DB) error {
	if kf.ID == uuid.Nil {
		kf.ID = uuid.New()
	}
	return nil
}

type LearningRecord struct {
	ID              uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	FileID          uuid.UUID     `gorm:"type:uuid;not null" json:"file_id"`
	LearnedAt       time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"learned_at"`
	DurationSeconds int           `json:"duration_seconds"`
	User            User          `gorm:"foreignKey:UserID" json:"-"`
	File            KnowledgeFile `gorm:"foreignKey:FileID" json:"-"`
}

func (lr *LearningRecord) BeforeCreate(tx *gorm.DB) error {
	if lr.ID == uuid.Nil {
		lr.ID = uuid.New()
	}
	return nil
}
