package models

type Book struct {
	Id         uint64     `json:"id"`
	Name       string     `json:"name"`
	Title      string     `json:"title"`
	Authors    []Author   `json:"authors"`
	Publisher  Publisher  `json:"publisher"`
	Language   Language   `json:"language"`
	Categories []Category `json:"categories"`
	CreatedAt  string     `json:"created_at,omitempty"`
}
