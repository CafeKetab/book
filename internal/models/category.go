package models

type Category struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
