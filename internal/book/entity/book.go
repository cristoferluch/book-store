package entity

import "book-store/internal/http_error"

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
	ID     int    `json:"id"`
}

type BookQuery struct {
	Author string `form:"author"`
	ISBN   string `form:"isbn"`
	Title  string `form:"title"`
	Page   int64  `form:"page"`
	Limit  int64  `form:"limit"`
	Id     int64  `form:"id"`
}

func IsValid(book Book) error {

	if book.Author == "" {
		return http_error.NewBadRequestError(http_error.ErrBookAuthorIsRequired)
	}

	if book.Title == "" {
		return http_error.NewBadRequestError(http_error.ErrBookTitleIsRequired)
	}

	if book.ISBN == "" {
		return http_error.NewBadRequestError(http_error.ErrBookIsbnIsRequired)
	}

	return nil
}
