package repository

import (
	"book-store/internal/book/entity"
	"book-store/internal/http_error"
	"context"
	"errors"
	"fmt"
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
			return nil, http_error.NewNotFoundError(http_error.ErrBookNotFound)
		}
		return nil, errors.Join(http_error.NewUnexpectedError(http_error.ErrUnexpected), fmt.Errorf("failed to get book"), err)
	}

	return &book, nil
}

func (r *BookRepository) SaveBook(ctx context.Context, book entity.Book) (*entity.Book, error) {

	query := `
		INSERT INTO book (title, author, isbn) values ($1, $2, $3) returning id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		book.Title,
		book.Author,
		book.ISBN,
	).Scan(
		&book.ID,
	)

	if err != nil {
		return nil, errors.Join(http_error.NewUnexpectedError(http_error.ErrUnexpected), fmt.Errorf("failed to save book"), err)
	}

	return &book, nil
}
