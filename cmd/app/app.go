package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gitlab.com/vitbog/titov-rest/internal/server"
)

func main() {
	fmt.Println("hello world")

	cfgFile := server.ConfigFile{}
	raw, err := os.ReadFile("./cfg/app.json")
	if err != nil {
		log.Fatalf("Error reading config file: %s\n", err)
	}
	err = json.Unmarshal(raw, &cfgFile)
	if err != nil {
		log.Fatalf("Error unmarshalling config: %s\n", err)
	}
	cfg, err := server.ConfigConvert(cfgFile)
	if err != nil {
		log.Fatalf("Error converting config: %s\n", err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pgConnectionString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"))

	s := &server.Server{Config: cfg}
	err = s.SetupDb(pgConnectionString)
	if err != nil {
		log.Fatalf("Error connecting db: %s\n", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
		w.WriteHeader(http.StatusOK)
	}))

	r.Handle("/auth", http.HandlerFunc(s.Auth)).Methods(http.MethodPost)
	r.Handle("/register", http.HandlerFunc(s.Register)).Methods(http.MethodPost)

	r.Handle("/tasks/create", s.Middleware(http.HandlerFunc(s.HandlerTasksCreate))).Methods(http.MethodPost)
	r.Handle("/tasks/list", s.Middleware(http.HandlerFunc(s.HandlerTasks))).Methods(http.MethodPost)
	r.Handle("/tasks/update/{id_task}", s.Middleware(http.HandlerFunc(s.HandlerTasksUpdate))).Methods(http.MethodPut)
	r.Handle("/tasks/delete", s.Middleware(http.HandlerFunc(s.HandlerTasksDelete))).Methods(http.MethodDelete)
	r.Handle("/tasks/csv", s.Middleware(http.HandlerFunc(s.HandlerTasksCSV))).Methods(http.MethodPost)

	r.Handle("/comments/create/{id_task}", s.Middleware(http.HandlerFunc(s.HandlerCommentsCreate))).Methods(http.MethodPost)
	r.Handle("/comments/list/{id_task}", s.Middleware(http.HandlerFunc(s.HandlerComments))).Methods(http.MethodPost)
	r.Handle("/comments/update/{id_comment}", s.Middleware(http.HandlerFunc(s.HandlerCommentsUpdate))).Methods(http.MethodPut)
	r.Handle("/comments/delete/{id_comment}", s.Middleware(http.HandlerFunc(s.HandlerCommentsDelete))).Methods(http.MethodDelete)

	s.SetupHTTP("0.0.0.0:8080", r)

	fmt.Println("Starting server...")
	s.Run()

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}
