package main

import "time"

// IssueStatus is a boolean representing a statuses state: open or closed
type IssueStatus bool

const (
	// IssueOpen represents an issue open state
	IssueOpen IssueStatus = true

	// IssueClosed represents an issue closed state
	IssueClosed = false
)

// Issue represents an issue
type Issue struct {
	ID          uint        `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	CreatedAt   int64       `json:"created_at"`
	ClosedAt    int64       `json:"closed_at"`
	Comments    []Comment   `json:"comments"`
	Status      IssueStatus `json:"status"`
}

// Close sets issue status to close and fulfills the ClosedAt attribute
func (i *Issue) Close() {
	i.Status = IssueClosed
	i.ClosedAt = time.Now().Unix()
}

// Comment adds a comment to the issue
func (i *Issue) Comment(comment string) {
	c := NewComment(
		time.Now(),
		comment,
	)

	i.Comments = append(i.Comments, *c)
}

// NewIssue creates a new issue object
func NewIssue(id uint, title, description string) *Issue {
	return &Issue{
		ID:          id,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now().Unix(),
		Status:      IssueOpen,
	}
}
