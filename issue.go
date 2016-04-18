package main

import "time"

type IssueStatus bool

const (
	ISSUE_OPENED IssueStatus = true
	ISSUE_CLOSED             = false
)

type Issue struct {
	Id          uint        `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	ClosedAt    time.Time   `json:"closed_at"`
	Comments    []Comment   `json:"comments"`
	Status      IssueStatus `json:"status"`
}

func (i *Issue) Close() {
	i.Status = ISSUE_CLOSED
	i.ClosedAt = time.Now()
}

func (i *Issue) Comment(comment string) {
	c := NewComment(
		time.Now(),
		comment,
	)

	i.Comments = append(i.Comments, *c)
}

func NewIssue(id uint, title, description string) *Issue {
	return &Issue{
		Id:          id,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		Status:      ISSUE_OPENED,
	}
}
