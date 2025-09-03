package repository

import (
	"book-store/internal/book/entity"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
			return nil, entity.ErrBookNotFound
		}
		return nil, err
	}

	return &book, nil
}
