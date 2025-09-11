package service

import (
	"book-store/internal/book/entity"
	"book-store/internal/http_error"
	"context"
)

type BookService struct {
	bookRepository entity.BookRepository
}

func NewBookService(bookRepository entity.BookRepository) *BookService {
	return &BookService{
		bookRepository: bookRepository,
	}
}

func (s *BookService) GetBookById(ctx context.Context, id int64) (*entity.Book, error) {
	return s.bookRepository.GetBookById(ctx, id)
}

func (s *BookService) SaveBook(ctx context.Context, book entity.Book) (*entity.Book, error) {

	if book.Author == "" {
		return nil, http_error.NewBadRequestError(http_error.ErrBookAuthorIsRequired)
	}

	if book.Title == "" {
		return nil, http_error.NewBadRequestError(http_error.ErrBookTitleIsRequired)
	}

	if book.ISBN == "" {
		return nil, http_error.NewBadRequestError(http_error.ErrBookIsbnIsRequired)
	}

	return s.bookRepository.SaveBook(ctx, book)
}
