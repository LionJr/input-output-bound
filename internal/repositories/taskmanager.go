package repositories

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"

	"github.com/LionJr/input-output-bound/internal/models"
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
	data, err := jsoniter.Marshal(task)
	if err != nil {
		return err
	}
	return tr.db.Set(context.TODO(), id, data, 0).Err()
}

func (tr *TaskManagerRepository) GetTask(id string) (*models.Task, error) {
	var task models.Task
	data, err := tr.db.Get(context.TODO(), id).Result()
	if err != nil {
		return nil, err
	}
	return &task, jsoniter.UnmarshalFromString(data, &task)
}
