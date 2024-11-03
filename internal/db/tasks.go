package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"gitlab.com/vitbog/titov-rest/internal/filters"
)

type TaskStatus string

const (
	TaskStatusFrozen     = "frozen"
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in-progress"
	TaskStatusCompleted  = "completed"
)

func TaskStatusIsValid(candidate string) bool {
	return candidate == TaskStatusFrozen || candidate == TaskStatusPending || candidate == TaskStatusInProgress || candidate == TaskStatusCompleted
}

type TaskModel struct {
	Id          int64      `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Status      TaskStatus `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	DueDate     time.Time  `json:"due_date" db:"due_date"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type TaskDb struct {
	Id          int64          `json:"id" db:"id"`
	Title       sql.NullString `json:"title" db:"title"`
	Description sql.NullString `json:"description" db:"description"`
	Status      sql.NullString `json:"status" db:"status"`
	CreatedAt   sql.NullTime   `json:"created_at" db:"created_at"`
	DueDate     sql.NullTime   `json:"due_date" db:"due_date"`
	UpdatedAt   sql.NullTime   `json:"updated_at" db:"updated_at"`
}

func (db *Db) TaskConvertFromDb(t TaskDb) (TaskModel, error) {
	taskStatus := t.Status.String
	if !TaskStatusIsValid(taskStatus) {
		return TaskModel{}, errors.New("task status is not valid")
	}

	return TaskModel{
		Id:          t.Id,
		Title:       t.Title.String,
		Description: t.Description.String,
		Status:      TaskStatus(t.Status.String),
		CreatedAt:   t.CreatedAt.Time,
		DueDate:     t.DueDate.Time,
		UpdatedAt:   t.UpdatedAt.Time,
	}, nil
}

func (db *Db) TasksConvertFromDb(tasks []TaskDb) ([]TaskModel, error) {
	convertedTasks := make([]TaskModel, 0, len(tasks))
	for _, t := range tasks {
		convertedTask, err := db.TaskConvertFromDb(t)
		if err != nil {
			return nil, err
		}

		convertedTasks = append(convertedTasks, convertedTask)
	}

	return convertedTasks, nil
}

var (
	glTasksAllowedColumns = []string{"id", "title", "description", "status", "created_at", "updated_at", "due_date", "updated_at"}
)

func (db *Db) TasksCreate(taskTitle, taskDescription, taskStatus string, DueDate time.Time) (int64, error) {
	schema := "tasks"
	query := fmt.Sprintf("SELECT %s.tasks_create($1, $2, $3, $4)", schema)
	var taskId int64

	err := db.Pg.Get(&taskId, query, taskTitle, taskDescription, taskStatus, DueDate)
	if err != nil {
		return 0, err
	}

	return taskId, nil
}

func (db *Db) Tasks(filt filters.Filtering) ([]TaskModel, error) {
	schema := "tasks"
	query := fmt.Sprintf("SELECT t.id, t.title, t.description, t.status, t.created_at, t.due_date, t.updated_at from %s.tasks_list() t", schema)

	filterQuery, args, err := filt.Filter(query, 0, glTasksAllowedColumns, "id", "title", "description", "status", "created_at", "updated_at", "due_date", "updated_at")
	if err != nil {
		return nil, fmt.Errorf("error filtering: %v", err)
	}

	reply := []TaskDb{}
	err = db.Pg.Select(&reply, filterQuery, args...)
	if err != nil {
		return nil, err
	}

	converted, err := db.TasksConvertFromDb(reply)
	if err != nil {
		return nil, err
	}

	return converted, nil
}

func (db *Db) TasksUpdate(taskId int64, taskTitle, taskDescription, taskStatus string, DueDate time.Time) error {
	schema := "tasks"
	query := fmt.Sprintf("CALL %s.tasks_update($1, $2, $3, $4, $5)", schema)

	_, err := db.Pg.Exec(query, taskId, taskTitle, taskDescription, taskStatus, DueDate)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) TasksDelete(ids []int64) error {
	schema := "tasks"
	query := fmt.Sprintf("CALL %s.tasks_delete($1)", schema)

	_, err := db.Pg.Exec(query, pq.Array(ids))
	if err != nil {
		return err
	}

	return nil
}
