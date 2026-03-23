package service

import (
	"errors"
	"force-learning/internal/model"
	"force-learning/internal/repository"
	"time"

	"github.com/google/uuid"
)

type SyncService struct {
	recordRepo *repository.LearningRecordRepository
	userRepo   *repository.UserRepository
}

func NewSyncService(recordRepo *repository.LearningRecordRepository, userRepo *repository.UserRepository) *SyncService {
	return &SyncService{
		recordRepo: recordRepo,
		userRepo:   userRepo,
	}
}

type SyncResult struct {
	SyncedRecords []SyncedRecordResult `json:"synced_records"`
	ServerTime    string               `json:"server_time"`
	HasMore       bool                 `json:"has_more"`
}

type SyncedRecordResult struct {
	ClientID  string `json:"client_id"`
	ServerID  string `json:"server_id"`
	Synced    bool   `json:"synced"`
	LearnedAt string `json:"learned_at"`
}

func (s *SyncService) SyncLearningRecords(userID uuid.UUID, lastSyncTime string, records []model.SyncRecordDTO) (*SyncResult, error) {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	var syncedRecords []SyncedRecordResult
	now := time.Now()
	serverTime := now.Format(time.RFC3339)

	for _, record := range records {
		var existingRecords []model.LearningRecord
		s.recordRepo.FindByClientID(userID, record.ClientID, &existingRecords)

		if len(existingRecords) > 0 {
			syncedRecords = append(syncedRecords, SyncedRecordResult{
				ClientID:  record.ClientID,
				ServerID:  existingRecords[0].ID.String(),
				Synced:    true,
				LearnedAt: existingRecords[0].LearnedAt.Format(time.RFC3339),
			})
			continue
		}

		fileID, err := uuid.Parse(record.FileID)
		if err != nil {
			syncedRecords = append(syncedRecords, SyncedRecordResult{
				ClientID: record.ClientID,
				ServerID: "",
				Synced:   false,
			})
			continue
		}

		var learnedAt time.Time
		if record.LearnedAt != "" {
			learnedAt, _ = time.Parse(time.RFC3339, record.LearnedAt)
		} else {
			learnedAt = now
		}

		newRecord := &model.LearningRecord{
			UserID:          userID,
			FileID:          fileID,
			DurationSeconds: record.DurationSeconds,
			LearnedAt:       learnedAt,
		}

		if err := s.recordRepo.Create(newRecord); err != nil {
			syncedRecords = append(syncedRecords, SyncedRecordResult{
				ClientID: record.ClientID,
				ServerID: "",
				Synced:   false,
			})
			continue
		}

		syncedRecords = append(syncedRecords, SyncedRecordResult{
			ClientID:  record.ClientID,
			ServerID:  newRecord.ID.String(),
			Synced:    true,
			LearnedAt: newRecord.LearnedAt.Format(time.RFC3339),
		})
	}

	return &SyncResult{
		SyncedRecords: syncedRecords,
		ServerTime:    serverTime,
		HasMore:       false,
	}, nil
}

func (s *SyncService) GetUnsyncedRecords(userID uuid.UUID, since string) ([]model.LearningRecord, error) {
	var sinceTime time.Time
	var err error

	if since != "" {
		sinceTime, err = time.Parse(time.RFC3339, since)
		if err != nil {
			sinceTime = time.Time{}
		}
	} else {
		sinceTime = time.Time{}
	}

	return s.recordRepo.FindSince(userID, sinceTime)
}
