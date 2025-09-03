package handler

import (
	"book-store/internal/book/entity"
	"book-store/pkg/utils"
	"errors"
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
		utils.SendJSON(w, utils.APIResponse{Error: "invalid id"}, http.StatusBadRequest)
		return
	}

	response, err := h.bookService.GetBookById(ctx, id)
	if err != nil {
		slog.Error("Erro in getBookById", slog.String("error", err.Error()))
		if errors.Is(err, entity.ErrBookNotFound) {
			utils.SendJSON(w, utils.APIResponse{Error: entity.ErrBookNotFound.Error()}, http.StatusNotFound)
			return
		}
		utils.SendJSON(w, utils.APIResponse{Error: "something went wrong"}, http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, utils.APIResponse{Data: response}, http.StatusOK)
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
