package entity

import "context"

type BookRepository interface {
	GetBookById(ctx context.Context, id int64) (*Book, error)
}
