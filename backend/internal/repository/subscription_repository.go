package repository

import (
	"force-learning/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(subscription *model.Subscription) error {
	return r.db.Create(subscription).Error
}

func (r *SubscriptionRepository) FindByUserID(userID uuid.UUID) ([]model.Subscription, error) {
	var subscriptions []model.Subscription
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&subscriptions).Error
	return subscriptions, err
}

func (r *SubscriptionRepository) FindActiveByUserID(userID uuid.UUID) (*model.Subscription, error) {
	var subscription model.Subscription
	err := r.db.Where("user_id = ? AND status = 'active' AND end_date > NOW()", userID).
		Preload("Plan").First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepository) Update(subscription *model.Subscription) error {
	return r.db.Save(subscription).Error
}

func (r *SubscriptionRepository) FindAllPlans() ([]model.SubscriptionPlan, error) {
	var plans []model.SubscriptionPlan
	err := r.db.Where("is_active = ?", true).Find(&plans).Error
	return plans, err
}

func (r *SubscriptionRepository) FindPlanByID(id uuid.UUID) (*model.SubscriptionPlan, error) {
	var plan model.SubscriptionPlan
	err := r.db.Where("id = ?", id).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *SubscriptionRepository) CreatePlan(plan *model.SubscriptionPlan) error {
	return r.db.Create(plan).Error
}

func (r *SubscriptionRepository) UpdatePlan(plan *model.SubscriptionPlan) error {
	return r.db.Save(plan).Error
}
