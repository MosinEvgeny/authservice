package http

import (
	"encoding/json"
	"net/http"

	"github.com/MosinEvgeny/authservice/internal/domain"
	"github.com/gorilla/mux"
)

type Handler struct {
	authService domain.Auth
}

func NewHandler(authService domain.Auth) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) GenerateTokens(w http.ResponseWriter, r *http.Request) {
	//  ... (получение userID из запроса)
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	clientIP := r.RemoteAddr

	tokenPair, err := h.authService.GenerateTokens(r.Context(), userID, clientIP)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)

}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// ... (получение refresh токена из запроса)
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	clientIP := r.RemoteAddr
	tokenPair, err := h.authService.RefreshToken(r.Context(), request.RefreshToken, clientIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(tokenPair)
}

func (h *Handler) InitRoutes(router *mux.Router) {
	router.HandleFunc("/generate_tokens", h.GenerateTokens).Methods(http.MethodGet)
	router.HandleFunc("/refresh_token", h.RefreshToken).Methods(http.MethodPost)
}
