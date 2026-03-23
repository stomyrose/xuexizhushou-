package repository

import (
	"force-learning/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LearningRecordRepository struct {
	db *gorm.DB
}

func NewLearningRecordRepository(db *gorm.DB) *LearningRecordRepository {
	return &LearningRecordRepository{db: db}
}

func (r *LearningRecordRepository) Create(record *model.LearningRecord) error {
	return r.db.Create(record).Error
}

func (r *LearningRecordRepository) FindByUserID(userID uuid.UUID) ([]model.LearningRecord, error) {
	var records []model.LearningRecord
	err := r.db.Where("user_id = ?", userID).Order("learned_at DESC").Find(&records).Error
	return records, err
}

func (r *LearningRecordRepository) FindByUserIDAndDate(userID uuid.UUID, date string) ([]model.LearningRecord, error) {
	var records []model.LearningRecord
	err := r.db.Where("user_id = ? AND DATE(learned_at) = ?", userID, date).
		Order("learned_at DESC").Find(&records).Error
	return records, err
}

func (r *LearningRecordRepository) GetTotalDurationByUserID(userID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.Model(&model.LearningRecord{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(duration_seconds), 0)").
		Scan(&total).Error
	return total, err
}

func (r *LearningRecordRepository) BatchCreate(records []model.LearningRecord) error {
	return r.db.Create(&records).Error
}

func (r *LearningRecordRepository) FindByClientID(userID uuid.UUID, clientID string, result *[]model.LearningRecord) error {
	return r.db.Where("user_id = ? AND client_id = ?", userID, clientID).First(result).Error
}

func (r *LearningRecordRepository) FindSince(userID uuid.UUID, since time.Time) ([]model.LearningRecord, error) {
	var records []model.LearningRecord
	err := r.db.Where("user_id = ? AND learned_at > ?", userID, since).
		Order("learned_at ASC").Find(&records).Error
	return records, err
}
