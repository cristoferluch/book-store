package entity

import "context"

type BookService interface {
	GetBookById(ctx context.Context, id int64) (*Book, error)
}
