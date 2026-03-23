package service

import (
	"errors"
	"force-learning/internal/model"
	"force-learning/internal/repository"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type KnowledgeService struct {
	knowledgeRepo *repository.KnowledgeRepository
	uploadPath    string
}

func NewKnowledgeService(knowledgeRepo *repository.KnowledgeRepository, uploadPath string) *KnowledgeService {
	return &KnowledgeService{
		knowledgeRepo: knowledgeRepo,
		uploadPath:    uploadPath,
	}
}

func (s *KnowledgeService) ListFiles(category string) ([]model.KnowledgeFile, error) {
	if category != "" {
		return s.knowledgeRepo.FindByCategory(category)
	}
	return s.knowledgeRepo.FindAll()
}

func (s *KnowledgeService) GetRandom() (*model.KnowledgeFile, error) {
	file, err := s.knowledgeRepo.FindRandom()
	if err != nil {
		return nil, errors.New("no content available")
	}
	return file, nil
}

func (s *KnowledgeService) Download(id string) (*model.KnowledgeFile, error) {
	fileID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid file ID")
	}

	file, err := s.knowledgeRepo.FindByID(fileID)
	if err != nil {
		return nil, errors.New("file not found")
	}

	return file, nil
}

func (s *KnowledgeService) Upload(file *model.KnowledgeFile, reader io.Reader) error {
	if err := os.MkdirAll(s.uploadPath, 0755); err != nil {
		return errors.New("failed to create upload directory")
	}

	fullPath := filepath.Join(s.uploadPath, file.ID.String()+filepath.Ext(file.Filename))
	dst, err := os.Create(fullPath)
	if err != nil {
		return errors.New("failed to create file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		return errors.New("failed to save file")
	}

	file.FilePath = fullPath
	if err := s.knowledgeRepo.Create(file); err != nil {
		os.Remove(fullPath)
		return errors.New("failed to save file record")
	}

	return nil
}

func (s *KnowledgeService) Delete(id string) error {
	fileID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid file ID")
	}

	file, err := s.knowledgeRepo.FindByID(fileID)
	if err != nil {
		return errors.New("file not found")
	}

	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
	}

	return s.knowledgeRepo.Delete(fileID)
}
