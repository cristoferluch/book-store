package entity

import (
	"book-store/pkg/utils"
	"context"
)

type BookService interface {
	GetBookById(ctx context.Context, id int64) (*Book, error)
	SaveBook(ctx context.Context, book Book) (*Book, error)
	GetBooks(ctx context.Context, bookQuery BookQuery) (*utils.PaginatedResponse[Book], error)
	DeleteBook(ctx context.Context, id int64) error
	UpdateBook(ctx context.Context, book Book, id int64) (*Book, error)
}
