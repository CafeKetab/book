package models

type Book struct {
	Id          uint64     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Publisher   Publisher  `json:"publisher"`
	Language    Language   `json:"language"`
	Authors     []Author   `json:"authors"`
	Categories  []Category `json:"categories"`
	CreatedAt   string     `json:"created_at,omitempty"`
}
