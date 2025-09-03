package entity

import "errors"

var ErrBookNotFound = errors.New("book not found")

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
	ID     int    `json:"id"`
}
