package handler

import (
	"book-store/internal/health/entity"
	"book-store/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) RegisterRouters(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", h.checkHealth)
	})
}

func (h *HealthHandler) checkHealth(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	status := "healthy"
	dbStatus := "connected"

	if err := h.db.Ping(ctx); err != nil {
		status = "unhealthy"
		dbStatus = "disconnected"
	}

	response := entity.HealthStatus{
		Status:    status,
		Timestamp: time.Now().UTC(),
		Database:  dbStatus,
	}

	if status == "unhealthy" {
		utils.SendJSON(w, response, http.StatusServiceUnavailable)
	} else {
		utils.SendJSON(w, response, http.StatusOK)
	}
}
