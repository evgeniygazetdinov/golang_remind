package models

import "time"

// Task represents a SQL training task
type Task struct {
	ID          string    `json:"id" example:"task_1234567890"`
	Description string    `json:"description" example:"Найдите всех сотрудников с зарплатой выше средней"`
	TableName   string    `json:"table_name" example:"practice_table_1234"`
	Difficulty  string    `json:"difficulty" example:"medium"`
	CreatedAt   time.Time `json:"created_at"`
}
