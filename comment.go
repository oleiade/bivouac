package main

type Comment struct {
	CreatedAt string
	Body      string
}

func NewComment(createdAt, body string) *Comment {
	return &Comment{
		CreatedAt: createdAt,
		Body:      body,
	}
}
