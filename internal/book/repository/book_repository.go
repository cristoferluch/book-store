package repository

import (
	"book-store/internal/book/entity"
	"book-store/internal/http_error"
	"book-store/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
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

func (r *BookRepository) GetBooks(ctx context.Context, bookQuery entity.BookQuery) (*utils.PaginatedResponse[entity.Book], error) {

	var args []interface{}
	var conditions []string

	query := `
		SELECT
			b.id,
			b.title,
			b.author,
			b.isbn
		FROM book b
	`

	if bookQuery.Id != 0 {
		conditions = append(conditions, fmt.Sprintf("b.id = $%d", len(args)+1))
		args = append(args, bookQuery.Id)
	}

	if bookQuery.Author != "" {
		conditions = append(conditions, fmt.Sprintf("b.author ILIKE $%d", len(args)+1))
		args = append(args, "%"+bookQuery.Author+"%")
	}

	if bookQuery.Title != "" {
		conditions = append(conditions, fmt.Sprintf("b.title ILIKE $%d", len(args)+1))
		args = append(args, "%"+bookQuery.Title+"%")
	}

	if bookQuery.ISBN != "" {
		conditions = append(conditions, fmt.Sprintf("b.isbn ILIKE $%d", len(args)+1))
		args = append(args, "%"+bookQuery.ISBN+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	queryCount := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS BASE", query)

	var count int64
	if err := r.db.QueryRow(ctx, queryCount, args...).Scan(&count); err != nil {
		return nil, errors.Join(http_error.NewUnexpectedError(http_error.ErrUnexpected), fmt.Errorf("failed to count books"), err)
	}

	query += fmt.Sprintf(" ORDER BY b.id OFFSET $%d LIMIT $%d", len(args)+1, len(args)+2)

	offset := bookQuery.Limit * (bookQuery.Page - 1)
	args = append(args, offset, bookQuery.Limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Join(http_error.NewUnexpectedError(http_error.ErrUnexpected), fmt.Errorf("failed to get books"), err)
	}
	defer rows.Close()

	books := make([]entity.Book, 0)
	for rows.Next() {
		var book entity.Book

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
		)

		if err != nil {
			return nil, errors.Join(http_error.NewUnexpectedError(http_error.ErrUnexpected), fmt.Errorf("failed to scan rows"), err)
		}

		books = append(books, book)
	}

	return &utils.PaginatedResponse[entity.Book]{
		Items:      books,
		TotalItems: count,
		Page:       bookQuery.Page,
		PageSize:   bookQuery.Limit,
	}, nil
}

func (r *BookRepository) DeleteBook(ctx context.Context, id int64) error {

	query := "DELETE FROM book WHERE id = $1"

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return http_error.NewNotFoundError(http_error.ErrBookNotFound)
		}
		return errors.Join(http_error.NewUnexpectedError(http_error.ErrUnexpected), fmt.Errorf("failed to delete book"), err)
	}

	return nil
}

func (r *BookRepository) UpdateBook(ctx context.Context, book entity.Book, id int64) (*entity.Book, error) {

	query := `
		UPDATE book SET title=$1, author=$2, isbn=$3 WHERE id=$4 RETURNING id
	`

	err := r.db.QueryRow(ctx, query, book.Title, book.Author, book.ISBN, id).Scan(&book.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, http_error.NewNotFoundError(http_error.ErrBookNotFound)
		}
		return nil, errors.Join(http_error.NewUnexpectedError(http_error.ErrUnexpected), fmt.Errorf("failed to update book"), err)
	}

	return &book, nil
}
