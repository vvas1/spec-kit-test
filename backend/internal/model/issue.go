package model

import "time"

const (
	StatusToDo       = "To Do"
	StatusInProgress = "In Progress"
	StatusReview     = "Review"
	StatusDone       = "Done"
)

var ValidStatuses = []string{StatusToDo, StatusInProgress, StatusReview, StatusDone}

// Issue represents an issue document.
type Issue struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Status      string    `json:"status" bson:"status"`
	AssigneeID  string    `json:"assigneeId" bson:"assigneeId,omitempty"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

// CreateIssueInput is the payload for creating an issue.
type CreateIssueInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	AssigneeID  string `json:"assigneeId"`
}
