package main

import "time"

type IssueStatus bool

const (
	ISSUE_OPENED IssueStatus = true
	ISSUE_CLOSED             = false
)

type Issue struct {
	Id       uint        `json:"id"`
	Title    string      `json:"title"`
	Comments []Comment   `json:"comments"`
	Status   IssueStatus `json:"status"`
}

func (i *Issue) Close() {
	i.Status = ISSUE_CLOSED
}

func (i *Issue) Comment(comment string) {
	c := NewComment(
		time.Now(),
		comment,
	)

	i.Comments = append(i.Comments, *c)
}

func NewIssue(id uint, title string, comments []Comment) *Issue {
	return &Issue{
		Id:       id,
		Title:    title,
		Comments: comments,
		Status:   ISSUE_OPENED,
	}
}
