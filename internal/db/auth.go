package db

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (db *Db) Auth(login, password string) (string, error) {
	schema := "users"
	query := fmt.Sprintf("SELECT a.password, a.role from %s.auth($1) a", schema)
	reply := []struct {
		Password string `json:"password" db:"password"`
		RoleName string `json:"role" db:"role"`
	}{}

	err := db.Pg.Select(&reply, query, login)
	if err != nil {
		return "", err
	}
	if len(reply) == 0 {
		return "", fmt.Errorf("user not found")
	}

	hashedPassword := reply[0].Password
	roleName := reply[0].RoleName

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		return roleName, nil
	} else {
		return "", fmt.Errorf("wrong password")
	}
}

func (db *Db) Register(login, password, roleName, fname, lname string) (int64, error) {
	schema := "users"
	query := fmt.Sprintf("SELECT %s.register($1, $2, $3, $4, $5)", schema)
	var userId int64

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("failed to hash password")
	}

	err = db.Pg.Get(&userId, query, login, hashedPassword, roleName, fname, lname)
	if err != nil {
		return 0, err
	}

	return userId, nil

}
