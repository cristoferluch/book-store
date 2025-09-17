package service

import (
	"book-store/internal/book/entity"
	"book-store/internal/http_error"
	"book-store/pkg/utils"
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

	if id == 0 {
		return nil, http_error.NewNotFoundError(http_error.ErrBookNotFound)
	}

	return s.bookRepository.GetBookById(ctx, id)
}

func (s *BookService) SaveBook(ctx context.Context, book entity.Book) (*entity.Book, error) {

	if err := entity.IsValid(book); err != nil {
		return nil, err
	}

	return s.bookRepository.SaveBook(ctx, book)
}

func (s *BookService) GetBooks(ctx context.Context, bookQuery entity.BookQuery) (*utils.PaginatedResponse[entity.Book], error) {

	if bookQuery.Page <= 0 {
		bookQuery.Page = 1
	}

	if bookQuery.Limit <= 0 {
		bookQuery.Limit = 10
	}

	return s.bookRepository.GetBooks(ctx, bookQuery)
}

func (s *BookService) DeleteBook(ctx context.Context, id int64) error {

	if id == 0 {
		return http_error.NewNotFoundError(http_error.ErrBookNotFound)
	}

	return s.bookRepository.DeleteBook(ctx, id)
}

func (s *BookService) UpdateBook(ctx context.Context, book entity.Book, id int64) (*entity.Book, error) {

	if id == 0 {
		return nil, http_error.NewNotFoundError(http_error.ErrBookNotFound)
	}

	if err := entity.IsValid(book); err != nil {
		return nil, err
	}

	return s.bookRepository.UpdateBook(ctx, book, id)
}
