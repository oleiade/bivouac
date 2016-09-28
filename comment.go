package main

// Comment represents a comment added to an issue
type Comment struct {
	CreatedAt int64  `json:"created_at"`
	Body      string `json:"body"`
}

// NewComment creates a new comment instance
func NewComment(createdAt int64, body string) *Comment {
	return &Comment{
		CreatedAt: createdAt,
		Body:      body,
	}
}
