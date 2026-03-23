package service

import (
	"errors"
	"force-learning/internal/model"
	"force-learning/internal/repository"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	subscriptionRepo *repository.SubscriptionRepository
	userRepo         *repository.UserRepository
}

func NewSubscriptionService(subscriptionRepo *repository.SubscriptionRepository, userRepo *repository.UserRepository) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		userRepo:         userRepo,
	}
}

func (s *SubscriptionService) GetPlans() ([]model.SubscriptionPlan, error) {
	plans, err := s.subscriptionRepo.FindAllPlans()
	if err != nil {
		return nil, errors.New("failed to get plans")
	}
	return plans, nil
}

func (s *SubscriptionService) Purchase(userID uuid.UUID, req *model.PurchaseRequest) (*model.Subscription, error) {
	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return nil, errors.New("invalid plan ID")
	}

	plan, err := s.subscriptionRepo.FindPlanByID(planID)
	if err != nil {
		return nil, errors.New("plan not found")
	}

	if !plan.IsActive {
		return nil, errors.New("plan is not available")
	}

	existingSubscription, _ := s.subscriptionRepo.FindActiveByUserID(userID)
	if existingSubscription != nil {
		existingSubscription.Status = "expired"
		s.subscriptionRepo.Update(existingSubscription)
	}

	now := time.Now()
	subscription := &model.Subscription{
		UserID:    userID,
		PlanID:    planID,
		StartDate: now,
		EndDate:   now.AddDate(0, 0, plan.DurationDays),
		Status:    "active",
	}

	if err := s.subscriptionRepo.Create(subscription); err != nil {
		return nil, errors.New("failed to create subscription")
	}

	s.userRepo.UpdateRemainingDays(userID, plan.DurationDays)

	return subscription, nil
}

func (s *SubscriptionService) GetCurrent(userID uuid.UUID) (*model.Subscription, error) {
	subscription, err := s.subscriptionRepo.FindActiveByUserID(userID)
	if err != nil {
		return nil, errors.New("no active subscription found")
	}
	return subscription, nil
}

func (s *SubscriptionService) PurchaseByIDString(userIDStr string, req *model.PurchaseRequest) (*model.Subscription, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	return s.Purchase(userID, req)
}

func (s *SubscriptionService) GetCurrentByIDString(userIDStr string) (*model.Subscription, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	return s.GetCurrent(userID)
}

func (s *SubscriptionService) CreatePlan(name string, durationDays int, price float64) (*model.SubscriptionPlan, error) {
	plan := &model.SubscriptionPlan{
		Name:         name,
		DurationDays: durationDays,
		Price:        price,
		IsActive:     true,
	}

	if err := s.subscriptionRepo.CreatePlan(plan); err != nil {
		return nil, errors.New("failed to create plan")
	}

	return plan, nil
}
