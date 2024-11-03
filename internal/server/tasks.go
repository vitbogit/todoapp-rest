package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/gorilla/mux"
	"gitlab.com/vitbog/titov-rest/internal/db"
	"gitlab.com/vitbog/titov-rest/internal/filters"
)

type TaskModelService struct {
	Id          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (s *Server) TaskModelServiceConvertFromModel(t db.TaskModel) (TaskModelService, error) {
	return TaskModelService{
		Id:          t.Id,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		CreatedAt:   t.CreatedAt,
		DueDate:     t.DueDate,
		UpdatedAt:   t.UpdatedAt,
	}, nil
}

func (s *Server) TasksModelServiceConvertFromModel(tasks []db.TaskModel) ([]TaskModelService, error) {
	convertedTasks := make([]TaskModelService, 0, len(tasks))
	for _, t := range tasks {
		convertedTask, err := s.TaskModelServiceConvertFromModel(t)
		if err != nil {
			return nil, err
		}

		convertedTasks = append(convertedTasks, convertedTask)
	}

	return convertedTasks, nil
}

func (s *Server) HandlerTasksCreate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var task TaskModelService
	err = json.Unmarshal(body, &task)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.Db.TasksCreate(task.Title, task.Description, task.Status, task.DueDate)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandlerTasks(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var filtering filters.Filtering
	err = json.Unmarshal(body, &filtering)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tasks, err := s.Db.Tasks(filtering)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(tasks)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(tasks) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandlerTasksCSV(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var filtering filters.Filtering
	err = json.Unmarshal(body, &filtering)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tasks, err := s.Db.Tasks(filtering)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(tasks) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Add("Content-Disposition", `attachment; filename="test.csv"`)
	gocsv.Marshal(tasks, w)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandlerTasksUpdate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var task TaskModelService
	err = json.Unmarshal(body, &task)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	taskIdStr, ok := mux.Vars(r)["id_task"]
	if !ok {
		log.Printf("Error: %s", "no id specified")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	taskId, err := strconv.ParseInt(taskIdStr, 10, 64)
	if err != nil {
		log.Printf("Error: %s", "bad id specified")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.Db.TasksUpdate(taskId, task.Title, task.Description, task.Status, task.DueDate)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandlerTasksDelete(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var Ids struct {
		Ids []int64 `json:"ids"`
	}
	err = json.Unmarshal(body, &Ids)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.Db.TasksDelete(Ids.Ids)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
