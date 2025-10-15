package mail_api

import (
	"anemone_notes/internal/api/middlewares"
	"anemone_notes/internal/model/mail_model"
	"anemone_notes/internal/services/mail_services"
	"anemone_notes/internal/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type MailHandler struct {
	Service    *mail_services.MailService
	JWTManager *utils.JWTManager
}

func NewMailHandler(service *mail_services.MailService, jwtManager *utils.JWTManager) *MailHandler {
	return &MailHandler{
		Service:    service,
		JWTManager: jwtManager,
	}
}

func (h *MailHandler) RegisterRoutes(r *mux.Router) {
	api := r.PathPrefix("/api/v1/mail").Subrouter()

	api.HandleFunc("/addresses", h.generateAddress).Methods("POST")

	inbox := api.PathPrefix("/inbox").Subrouter()
	inbox.Use(middlewares.AuthMailMiddleware(h.JWTManager))
	inbox.HandleFunc("", h.getInbox).Methods("GET")
}

// generateAddress — POST /api/v1/mail/addresses
func (h *MailHandler) generateAddress(w http.ResponseWriter, r *http.Request) {
	response, err := h.Service.GenerateAddressAndToken()
	if err != nil {
		log.Printf("ERROR: could not generate address: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

// getInbox — GET /api/v1/mail/inbox
func (h *MailHandler) getInbox(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.AddressContextKey).(*utils.UserClaims)
	if !ok {
		http.Error(w, "Could not retrieve user from context", http.StatusInternalServerError)
		return
	}

	emails, err := h.Service.GetInboxForAddress(claims.UserID)
	if err != nil {
		log.Printf("ERROR: could not get inbox: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if emails == nil {
		emails = []mail_model.Email{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(emails)
}
