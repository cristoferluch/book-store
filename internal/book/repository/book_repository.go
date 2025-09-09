package repository

import (
	"book-store/internal/book/entity"
	"book-store/pkg/utils/http_errors"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type BookRepository struct {
	db *pgxpool.Pool
}

func NewBookRepository(db *pgxpool.Pool) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

func (r *BookRepository) GetBookById(ctx context.Context, id int64) (*entity.Book, error) {

	query := `
		SELECT 
    		b.author,
			b.id,
			b.title,
			b.isbn
    	FROM book b
    	WHERE id=$1
	`

	var book entity.Book
	err := r.db.QueryRow(ctx, query, id).Scan(
		&book.Author,
		&book.ID,
		&book.Title,
		&book.ISBN,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, http_errors.NewNotFoundError(http_errors.ErrBookNotFound)
		}
		slog.Error(err.Error())
		return nil, http_errors.NewUnexpectedError(http_errors.Unexpected)
	}

	return &book, nil
}
