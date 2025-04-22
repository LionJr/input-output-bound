package handlers

import (
	"github.com/LionJr/input-output-bound/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	logger  *zap.Logger
	service *services.TaskManagerService
}

func NewTaskManagerHandler(logger *zap.Logger, service *services.TaskManagerService) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) CreateTask(c *gin.Context) {
	task, err := h.service.CreateTask()
	if err != nil {
		h.logger.Error("failed to create task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
	}

	c.JSON(http.StatusAccepted, task)
}

func (h *Handler) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, ok := h.service.GetTask(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}
