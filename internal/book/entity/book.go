package entity

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
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Sort   string `form:"sort"`
}
