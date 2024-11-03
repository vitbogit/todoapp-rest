package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/vitbog/titov-rest/internal/filters"
	"gitlab.com/vitbog/titov-rest/internal/token"
)

type CommentModelService struct {
	Id        int64     `json:"id" db:"id"`
	IdUser    int64     `json:"id_user" db:"id_user"`
	IdTask    int64     `json:"id_task" db:"id_task"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (s *Server) HandlerCommentsCreate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var comment CommentModelService
	err = json.Unmarshal(body, &comment)
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

	tok := r.Header.Get("Authorization")
	userLoginInterface, err := token.Field("login", tok, s.JWTSecretKey)
	userLogin, ok := userLoginInterface.(string)
	if !ok {
		log.Printf("Error: %s", "bad login in token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.Db.CommentCreate(taskId, userLogin, comment.Content)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandlerComments(w http.ResponseWriter, r *http.Request) {
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

	comments, err := s.Db.Comments(taskId, filtering)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(comments)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(comments) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandlerCommentsUpdate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) HandlerCommentsDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
