package server

import (
	bookHdl "book-store/internal/book/handler"
	bookRepo "book-store/internal/book/repository"
	bookSvc "book-store/internal/book/service"
	healthHdl "book-store/internal/health/handler"
	"book-store/pkg/middlewares"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer(db *pgxpool.Pool) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.LoggingMiddleware)

	healthHandler := healthHdl.NewHealthHandler(db)
	bookRepository := bookRepo.NewBookRepository(db)
	bookService := bookSvc.NewBookService(bookRepository)
	bookHandler := bookHdl.NewBookHandler(bookService)

	bookHandler.RegisterRouters(r)
	healthHandler.RegisterRouters(r)

	return r
}
