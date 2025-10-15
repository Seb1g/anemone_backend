package main

import (
	"anemone_notes/internal/api/auth_api"
	"anemone_notes/internal/api/mail_api"
	"anemone_notes/internal/api/notes_api"
	"anemone_notes/internal/config"
	"anemone_notes/internal/database"
	"anemone_notes/internal/repository/auth_repository"
	"anemone_notes/internal/repository/mail_repository"
	"anemone_notes/internal/repository/notes_repository"
	"anemone_notes/internal/services/auth_services"
	"anemone_notes/internal/services/mail_services"
	"anemone_notes/internal/services/notes_services"
	"anemone_notes/internal/smtp_server"
	"anemone_notes/internal/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

func main() {
	cfg := config.Load()

	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("FATAL: database connection failed: %v", err)
	}
	defer db.Close()
	log.Println("INFO: Database connection successful")

	jwtManager := utils.NewJWTManager(cfg.JWTSecret)

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

	mailRepo := mail_repository.New(db)
	mailService := mail_services.New(mailRepo, jwtManager, cfg.DomainName)
	mailHandler := mail_api.NewMailHandler(mailService, jwtManager)

	r := mux.NewRouter()

	authHandler.RegisterRoutes(r)
	pageHandler.PagesRoutes(r)
	folderHandler.FolderRoutes(r)
	mailHandler.RegisterRoutes(r)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		smtpServer := smtp_server.NewServer(cfg, mailRepo)
		smtpServer.Start()
	}()

	go func() {
		defer wg.Done()
		log.Printf("INFO: Starting HTTP server on port %s", cfg.HTTPPort)
		if err := http.ListenAndServe(":"+cfg.HTTPPort, r); err != nil {
			log.Fatalf("FATAL: failed to start HTTP server: %v", err)
		}
	}()

	log.Println("INFO: All services are running")

	wg.Wait()
}
