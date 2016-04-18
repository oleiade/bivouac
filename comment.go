package main

import "time"

type Comment struct {
	CreatedAt time.Time
	Body      string
}

func NewComment(createdAt time.Time, body string) *Comment {
	return &Comment{
		CreatedAt: createdAt,
		Body:      body,
	}
}
