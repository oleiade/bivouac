package main

import (
	"fmt"
	"time"
)

type Issue struct {
	Id       uint      `json:"id"`
	Title    string    `json:"title"`
	Comments []Comment `json:"comments"`
}

func (i *Issue) Comment(comment string) {
	c := NewComment(
		time.Now().String(),
		comment,
	)

	i.Comments = append(i.Comments, *c)
	fmt.Println(len(i.Comments))
}

func NewIssue(id uint, title string, comments []Comment) *Issue {
	return &Issue{
		Id:       id,
		Title:    title,
		Comments: comments,
	}
}
