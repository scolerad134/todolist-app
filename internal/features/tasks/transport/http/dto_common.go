package tasks_transport_http

import (
	"time"

	"github.com/scolerad134/todolist-app/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id"             example:"15"`
	Version      int        `json:"version"        example:"3"`
	Title        string     `json:"title"          example:"Тренировка"`
	Description  *string    `json:"description"    example:"Начало тренировки в 19:30"`
	Completed    bool       `json:"completed"      example:"false"`
	CreatedAt    time.Time  `json:"created_at"     example:"2026-02-26T10:30:00Z"`
	CompletedAt  *time.Time `json:"completed_at"   example:"null"`
	AuthorUserID int        `json:"author_user_id" example:"5"`
}

func TaskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func TaskDTOsFromDomains(tasks []domain.Task) []TaskDTOResponse {
	tasksDTO := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		tasksDTO[i] = TaskDTOFromDomain(task)
	}

	return tasksDTO
}
