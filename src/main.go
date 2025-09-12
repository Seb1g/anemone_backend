package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func healthCheckHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "OK")
}

type TaskInput struct {
	Title string `json:"title"`
}

func createNewTask(db *sql.DB, title string) (int, error) {
	var id int
	err := db.QueryRow(`
	INSERT INTO tasks(title) VALUES($1) RETURNING id
	`, title).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("error while creating: %v", err)
	}

	return id, nil
}

func createNewTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var taskInput TaskInput
		err := json.NewDecoder(req.Body).Decode(&taskInput)

		if err != nil {
			http.Error(res, "Unused json format", http.StatusBadRequest)
			return
		}

		taskID, err := createNewTask(db, taskInput.Title)

		if err != nil {
			http.Error(res, "Close", http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		fmt.Fprintf(res, `{"id": %d, "message": "Created"}`, taskID)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	connStr := "user=seb1glory password=881tfRqr4D0Z host=7ty2ryz3.ru port= 5432 dbname=trello sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("Ошибка при открытии соединения:", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		is_completed BOOLEAN NOT NULL DEFAULT FALSE
		);
		`)

	if err != nil {
		fmt.Println("Ошибка при создании таблицы:", err)
		return
	}

	fmt.Println("Таблица 'tasks' успешно создана.")

	router.HandleFunc("/notes/create_new_notes", createNewTaskHandler(db)).Methods("POST")

	fmt.Println("Успешное подключение к базе данных!")
	fmt.Println("Server work on 8080 port")
	http.ListenAndServe(":8080", nil)
}
