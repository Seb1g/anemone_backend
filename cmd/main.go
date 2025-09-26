package main

import (
	"anemone_notes/internal/api/notes_api"
	"anemone_notes/internal/api/auth_api"
	"anemone_notes/internal/repository/notes_repository"
	"anemone_notes/internal/repository/auth_repository"
	"anemone_notes/internal/services/notes_services"
	"anemone_notes/internal/services/auth_services"
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

	userRepo := auth_repository.NewUserRepo(db)
	refreshRepo := auth_repository.NewRefreshRepo(db)
	authSvc := auth_services.NewAuthService(userRepo, refreshRepo)
	authHandler := auth_api.NewAuthHandler(authSvc)

	pageRepo := notes_repository.NewPageRepo(db)
	pageSvc := notes_services.NewPageService(pageRepo)
	pageHandler := notes_api.NewPageHandler(pageSvc, authSvc)

	folderRepo := notes_repository.NewFolderRepo(db)
	folderSvc := notes_services.NewFolderService(folderRepo)
	folderHandler := notes_api.NewFolderHandler(folderSvc, authSvc)

	r := mux.NewRouter()
	
	authHandler.RegisterRoutes(r)
	pageHandler.PagesRoutes(r)
	folderHandler.FolderRoutes(r)

	log.Println("server started at :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
