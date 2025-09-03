package service

import (
	"book-store/internal/book/entity"
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
