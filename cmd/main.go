package main

import (
	"anemone_notes/internal/api/notes_api"
	"anemone_notes/internal/repository/notes_repository"
	"anemone_notes/internal/services/notes_services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db, err_db := notes_repository.Database_conn()

	if err_db != nil {
		log.Fatalf("Error db: %v", err_db)
	}
	defer db.Close()

	pageRepo := notes_repository.NewPageRepo(db)
	pageSvc := notes_services.NewPageService(pageRepo)
	pageHandler := notes_api.NewPageHandler(pageSvc)

	r := mux.NewRouter()
	pageHandler.RegisterRoutes(r)

	log.Println("server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
