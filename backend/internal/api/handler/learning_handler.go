package handler

import (
	"force-learning/internal/model"
	"force-learning/internal/pkg/response"
	"force-learning/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LearningHandler struct {
	learningService *service.LearningService
}

func NewLearningHandler(learningService *service.LearningService) *LearningHandler {
	return &LearningHandler{learningService: learningService}
}

func (h *LearningHandler) CreateRecord(c *gin.Context) {
	var req struct {
		FileID          string `json:"file_id" binding:"required"`
		DurationSeconds int    `json:"duration_seconds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "user not found")
		return
	}

	fileID, err := uuid.Parse(req.FileID)
	if err != nil {
		response.BadRequest(c, "invalid file ID")
		return
	}

	record, err := h.learningService.CreateRecord(uuid.MustParse(userID), fileID, req.DurationSeconds)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, record)
}

func (h *LearningHandler) GetRecords(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "user not found")
		return
	}

	date := c.Query("date")
	if date != "" {
		records, err := h.learningService.GetUserRecordsByDate(uuid.MustParse(userID), date)
		if err != nil {
			response.InternalServerError(c, err.Error())
			return
		}
		response.Success(c, records)
		return
	}

	records, err := h.learningService.GetUserRecords(uuid.MustParse(userID))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, records)
}

func (h *LearningHandler) GetStatistics(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "user not found")
		return
	}

	totalDuration, err := h.learningService.GetTotalDuration(uuid.MustParse(userID))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"total_duration_seconds":   totalDuration,
		"total_duration_formatted": formatDuration(totalDuration),
	})
}

func formatDuration(seconds int64) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60
	return time.Date(0, 0, 0, int(hours), int(minutes), int(secs), 0, time.UTC).Format("15:04:05")
}

func (h *LearningHandler) BatchCreate(c *gin.Context) {
	var req struct {
		Records []struct {
			FileID          string `json:"file_id" binding:"required"`
			DurationSeconds int    `json:"duration_seconds" binding:"required"`
			LearnedAt       string `json:"learned_at"`
		} `json:"records" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "user not found")
		return
	}

	var records []model.LearningRecord
	for _, r := range req.Records {
		fileID, err := uuid.Parse(r.FileID)
		if err != nil {
			continue
		}
		record := model.LearningRecord{
			UserID:          uuid.MustParse(userID),
			FileID:          fileID,
			DurationSeconds: r.DurationSeconds,
		}
		if r.LearnedAt != "" {
			if t, err := time.Parse(time.RFC3339, r.LearnedAt); err == nil {
				record.LearnedAt = t
			}
		}
		records = append(records, record)
	}

	if err := h.learningService.BatchCreateRecords(records); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, gin.H{"created": len(records)})
}
