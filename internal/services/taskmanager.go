package services

import (
	"fmt"
	"github.com/LionJr/input-output-bound/internal/models"
	"github.com/LionJr/input-output-bound/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type TaskManagerService struct {
	logger *zap.Logger
	repo   *repositories.TaskManagerRepository
}

func NewTaskManagerService(logger *zap.Logger, repo *repositories.TaskManagerRepository) *TaskManagerService {
	return &TaskManagerService{
		logger: logger,
		repo:   repo,
	}
}

func (tm *TaskManagerService) CreateTask() (*models.Task, error) {
	id := uuid.NewString()
	task := &models.Task{
		Id:     id,
		Status: models.Pending,
	}

	if err := tm.repo.AddTask(id, task); err != nil {
		return nil, err
	}

	go tm.executeTask(id)

	return task, nil
}

func (tm *TaskManagerService) executeTask(id string) {
	task, err := tm.repo.GetTask(id)
	if err != nil {
		tm.logger.Warn("task execution failed", zap.String("id", id), zap.Error(err))
		return
	}

	task.Status = models.Running
	if err = tm.repo.AddTask(id, task); err != nil {
		tm.logger.Warn("task execution failed", zap.String("id", id), zap.Error(err))
		return
	}

	time.Sleep(3 * time.Minute)

	task.Status = models.Done
	task.Result = fmt.Sprintf("Task %s completed successfully!!!", id)

	if err = tm.repo.AddTask(id, task); err != nil {
		tm.logger.Warn("task execution failed", zap.String("id", id), zap.Error(err))
		return
	}
}

func (tm *TaskManagerService) GetTask(id string) (*models.Task, bool) {
	task, err := tm.repo.GetTask(id)
	if err != nil {
		return nil, false
	}

	return task, true
}
