package service

import (
	"errors"
	"force-learning/internal/model"
	"force-learning/internal/repository"
	"time"

	"github.com/google/uuid"
)

type LearningService struct {
	recordRepo *repository.LearningRecordRepository
	userRepo   *repository.UserRepository
}

func NewLearningService(recordRepo *repository.LearningRecordRepository, userRepo *repository.UserRepository) *LearningService {
	return &LearningService{
		recordRepo: recordRepo,
		userRepo:   userRepo,
	}
}

func (s *LearningService) CreateRecord(userID, fileID uuid.UUID, durationSeconds int, clientID string) (*model.LearningRecord, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.RemainingDays <= 0 {
		return nil, errors.New("no remaining learning days")
	}

	record := &model.LearningRecord{
		UserID:          userID,
		FileID:          fileID,
		DurationSeconds: durationSeconds,
		LearnedAt:       time.Now(),
		ClientID:        clientID,
	}

	if err := s.recordRepo.Create(record); err != nil {
		return nil, errors.New("failed to create learning record")
	}

	return record, nil
}

func (s *LearningService) GetUserRecords(userID uuid.UUID) ([]model.LearningRecord, error) {
	return s.recordRepo.FindByUserID(userID)
}

func (s *LearningService) GetUserRecordsByDate(userID uuid.UUID, date string) ([]model.LearningRecord, error) {
	return s.recordRepo.FindByUserIDAndDate(userID, date)
}

func (s *LearningService) GetTotalDuration(userID uuid.UUID) (int64, error) {
	return s.recordRepo.GetTotalDurationByUserID(userID)
}

func (s *LearningService) BatchCreateRecords(records []model.LearningRecord) error {
	return s.recordRepo.BatchCreate(records)
}
