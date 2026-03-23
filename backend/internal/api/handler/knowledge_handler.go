package handler

import (
	"force-learning/internal/model"
	"force-learning/internal/pkg/response"
	"force-learning/internal/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type KnowledgeHandler struct {
	knowledgeService *service.KnowledgeService
}

func NewKnowledgeHandler(knowledgeService *service.KnowledgeService) *KnowledgeHandler {
	return &KnowledgeHandler{knowledgeService: knowledgeService}
}

func (h *KnowledgeHandler) ListFiles(c *gin.Context) {
	category := c.Query("category")

	files, err := h.knowledgeService.ListFiles(category)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, files)
}

func (h *KnowledgeHandler) GetRandom(c *gin.Context) {
	file, err := h.knowledgeService.GetRandom()
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, file)
}

func (h *KnowledgeHandler) Download(c *gin.Context) {
	id := c.Param("id")

	file, err := h.knowledgeService.Download(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		response.NotFound(c, "file not found on disk")
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+file.Filename)
	c.File(file.FilePath)
}

func (h *KnowledgeHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	category := c.PostForm("category")
	fileType := c.PostForm("file_type")

	knowledgeFile := &model.KnowledgeFile{
		Filename:  file.Filename,
		FileType:  fileType,
		Category:  &category,
		IsVisible: true,
	}

	openedFile, err := file.Open()
	if err != nil {
		response.InternalServerError(c, "failed to open file")
		return
	}

	if err := h.knowledgeService.Upload(knowledgeFile, openedFile); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, knowledgeFile)
}

func (h *KnowledgeHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.knowledgeService.Delete(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *KnowledgeHandler) CreateFileRecord(filename, fileType, category string) *model.KnowledgeFile {
	return &model.KnowledgeFile{
		ID:        uuid.New(),
		Filename:  filename,
		FileType:  fileType,
		Category:  &category,
		IsVisible: true,
	}
}
