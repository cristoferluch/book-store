package handler

import (
	"book-store/internal/health/entity"
	"book-store/pkg/config"
	"book-store/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"time"
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

	dbStatus := "up"

	if err := h.db.Ping(ctx); err != nil {
		dbStatus = "down"
	}

	response := entity.HealthResponse{
		Version:  config.Cfg.App.Version,
		Database: dbStatus,
		Uptime:   time.Since(config.Cfg.StartTime).Truncate(time.Second).String(),
	}

	utils.SendJSON(w, response, http.StatusOK)
}
