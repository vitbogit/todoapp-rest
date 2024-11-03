package db

import (
	"database/sql"
	"fmt"
	"time"

	"gitlab.com/vitbog/titov-rest/internal/filters"
)

type CommentModel struct {
	Id        int64     `json:"id" db:"id"`
	IdUser    int64     `json:"id_user" db:"id_user"`
	IdTask    int64     `json:"id_task" db:"id_task"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CommentDb struct {
	Id        int64          `json:"id" db:"id"`
	IdUser    int64          `json:"id_user" db:"id_user"`
	IdTask    int64          `json:"id_task" db:"id_task"`
	Content   sql.NullString `json:"content" db:"content"`
	CreatedAt sql.NullTime   `json:"created_at" db:"created_at"`
}

func (db *Db) CommentConvertFromDb(t CommentDb) (CommentModel, error) {

	return CommentModel{
		Id:        t.Id,
		IdUser:    t.IdUser,
		IdTask:    t.IdTask,
		Content:   t.Content.String,
		CreatedAt: t.CreatedAt.Time,
	}, nil
}

func (db *Db) CommentsConvertFromDb(comments []CommentDb) ([]CommentModel, error) {
	convertedComments := make([]CommentModel, 0, len(comments))
	for _, t := range comments {
		convertedComment, err := db.CommentConvertFromDb(t)
		if err != nil {
			return nil, err
		}

		convertedComments = append(convertedComments, convertedComment)
	}

	return convertedComments, nil
}

var (
	glCommentsAllowedColumns = []string{"id", "id_user", "id_task", "content", "created_at"}
)

func (db *Db) Comments(taskId int64, filt filters.Filtering) ([]CommentModel, error) {
	schema := "tasks"
	query := fmt.Sprintf("SELECT c.id, c.id_user, c.id_task, c.content, c.created_at from %s.comments_list($1) c", schema)
	narg := 1

	filterQuery, filterArgs, err := filt.Filter(query, 1, glCommentsAllowedColumns, "id", "id_user", "id_task", "content", "created_at")
	if err != nil {
		return nil, fmt.Errorf("error filtering: %v", err)
	}

	args := make([]interface{}, narg, len(filterArgs)+narg)
	args[0] = taskId
	args = append(args, filterArgs...)

	reply := []CommentDb{}
	err = db.Pg.Select(&reply, filterQuery, args...)
	if err != nil {
		return nil, err
	}

	converted, err := db.CommentsConvertFromDb(reply)
	if err != nil {
		return nil, err
	}

	return converted, nil
}

func (db *Db) CommentCreate(taskId int64, userLogin, content string) (int64, error) {
	schema := "tasks"
	query := fmt.Sprintf("SELECT %s.comment_create($1, $2, $3)", schema)
	var commentId int64

	err := db.Pg.Get(&commentId, query, taskId, userLogin, content)
	if err != nil {
		return 0, err
	}

	return commentId, nil
}
