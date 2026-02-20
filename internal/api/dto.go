package api

type dbVideo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
	Type string `json:"type"` // "source", "lesson", "performance"
}