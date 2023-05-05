package models

type Category struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// func (c *Category) ToBeScanned(required bool) []any {
// 	result := []any{&c.Id, &c.Name, &c.Title}
// 	if !required {
// 		result = append(result, &c.Description)
// 	}
// 	return result
// }
