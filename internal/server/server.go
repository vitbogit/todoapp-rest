package server

import (
	"net/http"
	"time"

	"gitlab.com/vitbog/titov-rest/internal/db"
)

type Server struct {
	HTTP *http.Server
	Db   *db.Db
	Config
}

type Config struct {
	JWTSecretKey  string
	JWTAccessTime time.Duration
}

func ConfigConvert(cfgFile ConfigFile) (Config, error) {
	JWTAccessTime, err := time.ParseDuration(cfgFile.JWTAccessTime)
	if err != nil {
		return Config{}, err
	}

	return Config{
		JWTSecretKey:  cfgFile.JWTSecretKey,
		JWTAccessTime: JWTAccessTime,
	}, nil
}

type ConfigFile struct {
	JWTSecretKey  string `json:"secret"`
	JWTAccessTime string `json:"access_time"`
}

func (s *Server) SetupDb(pgConnectionString string) error {
	db, err := db.New(pgConnectionString)
	if err != nil {
		panic(err)
	}

	s.Db = db

	return nil
}

func (s *Server) SetupHTTP(serverAddress string, r http.Handler) error {
	s.HTTP = &http.Server{Addr: serverAddress, Handler: r}

	return nil
}

func (s *Server) Run() error {
	err := s.HTTP.ListenAndServe()
	return err
}
