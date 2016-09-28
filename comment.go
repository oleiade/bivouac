package main

import "time"

// Comment represents a comment added to an issue
type Comment struct {
	CreatedAt time.Time `json:"created_at"`
	Body      string    `json:"body"`
}

// NewComment creates a new comment instance
func NewComment(createdAt time.Time, body string) *Comment {
	return &Comment{
		CreatedAt: createdAt,
		Body:      body,
	}
}
