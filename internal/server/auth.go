package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gitlab.com/vitbog/titov-rest/internal/token"
)

type UserInfo struct {
	Id               int64  `json:"id"`
	Login            string `json:"login"`
	FName            string `json:"fname"`
	LName            string `json:"lname"`
	Role             string `json:"role"`
	DateRegistration string `json:"date_registration"`
}

type UserCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Server) Auth(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userCredentials UserCredentials
	err = json.Unmarshal(body, &userCredentials)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	roleName, err := s.Db.Auth(userCredentials.Login, userCredentials.Password)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessToken, err := token.Generate(userCredentials.Login, roleName, s.JWTSecretKey, s.JWTAccessTime)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: accessToken,
	}
	resultString, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resultString)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userCredentials UserCredentials
	err = json.Unmarshal(body, &userCredentials)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userInfo UserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.Db.Register(userCredentials.Login, userCredentials.Password, DefaultUserRoleName, userInfo.FName, userInfo.LName)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := r.Header.Get("Authorization")
		_, err := token.Verify(t, s.JWTSecretKey)
		if err != nil {
			log.Printf("Error: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
