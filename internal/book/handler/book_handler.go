package handler

import (
	"book-store/internal/book/entity"
	"book-store/pkg/utils"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BookHandler struct {
	bookService entity.BookService
}

func NewBookHandler(bookService entity.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) RegisterRouters(r *chi.Mux) {

	r.Route("/api/book", func(r chi.Router) {
		r.Get("/{id}", h.getBookById)
		r.Get("/", h.getBooks)
		r.Post("/", h.saveBook)
		r.Put("/{id}", h.updateBook)
		r.Delete("/{id}", h.deleteBook)
	})
}

func (h *BookHandler) getBookById(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	p := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		return
	}

	response, err := h.bookService.GetBookById(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Erro in getBookById", slog.String("error", err.Error()))
		utils.SendError(w, err)
		return
	}

	utils.SendJSON(w, response, http.StatusOK)
}

func (h *BookHandler) getBooks(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *BookHandler) saveBook(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
}
