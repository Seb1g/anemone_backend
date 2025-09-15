package notes_api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"anemone_notes/internal/api"
	"anemone_notes/internal/services/notes_services"
	"github.com/gorilla/mux"
)

type PageHandler struct {
	Service *notes_services.PageService
}

func NewPageHandler(s *notes_services.PageService) *PageHandler {
	return &PageHandler{Service: s}
}

func (h *PageHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/pages", h.createPage).Methods("POST");
	r.HandleFunc("/pages/{id}", h.getPage).Methods("GET");
	r.HandleFunc("/health", api.HealthCheckHandler).Methods("GET");
}

func (h *PageHandler) createPage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  int    `json:"user_id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := h.Service.CreatePage(r.Context(), req.UserID, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *PageHandler) getPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	p, err := h.Service.GetPage(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
