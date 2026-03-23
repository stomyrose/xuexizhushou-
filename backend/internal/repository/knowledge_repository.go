package repository

import (
	"force-learning/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KnowledgeRepository struct {
	db *gorm.DB
}

func NewKnowledgeRepository(db *gorm.DB) *KnowledgeRepository {
	return &KnowledgeRepository{db: db}
}

func (r *KnowledgeRepository) Create(file *model.KnowledgeFile) error {
	return r.db.Create(file).Error
}

func (r *KnowledgeRepository) FindAll() ([]model.KnowledgeFile, error) {
	var files []model.KnowledgeFile
	err := r.db.Where("is_visible = ?", true).Order("uploaded_at DESC").Find(&files).Error
	return files, err
}

func (r *KnowledgeRepository) FindByCategory(category string) ([]model.KnowledgeFile, error) {
	var files []model.KnowledgeFile
	err := r.db.Where("category = ? AND is_visible = ?", category, true).
		Order("uploaded_at DESC").Find(&files).Error
	return files, err
}

func (r *KnowledgeRepository) FindByID(id uuid.UUID) (*model.KnowledgeFile, error) {
	var file model.KnowledgeFile
	err := r.db.Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *KnowledgeRepository) FindRandom() (*model.KnowledgeFile, error) {
	var file model.KnowledgeFile
	err := r.db.Where("is_visible = ?", true).
		Order("RANDOM()").First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *KnowledgeRepository) Update(file *model.KnowledgeFile) error {
	return r.db.Save(file).Error
}

func (r *KnowledgeRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.KnowledgeFile{}).Error
}
