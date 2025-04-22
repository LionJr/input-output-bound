package repositories

import (
	"context"
	"github.com/LionJr/input-output-bound/internal/models"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
)

type TaskManagerRepository struct {
	db *redis.Client
}

func NewTaskManagerRepository(db *redis.Client) *TaskManagerRepository {
	return &TaskManagerRepository{
		db: db,
	}
}

func (tr *TaskManagerRepository) AddTask(id string, task *models.Task) error {
	ctx := context.TODO()

	data, err := jsoniter.Marshal(task)
	if err != nil {
		return err
	}

	return tr.db.Set(ctx, id, data, 0).Err()
}

func (tr *TaskManagerRepository) GetTask(id string) (*models.Task, error) {
	ctx := context.TODO()

	var task models.Task

	data, err := tr.db.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}

	return &task, jsoniter.UnmarshalFromString(data, &task)
}
