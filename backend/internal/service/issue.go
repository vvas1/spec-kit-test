package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"issue-tracker/backend/internal/model"
	"issue-tracker/backend/internal/store/mongodb"
)

const (
	MaxTitleLength       = 200
	MaxDescriptionLength = 10000
)

// IssueService handles issue business logic.
type IssueService struct {
	store *mongodb.IssueStore
}

// NewIssueService returns a new IssueService.
func NewIssueService(store *mongodb.IssueStore) *IssueService {
	return &IssueService{store: store}
}

// ValidateCreateInput validates CreateIssueInput. Returns a descriptive error.
func ValidateCreateInput(in *model.CreateIssueInput) error {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return errors.New("title is required")
	}
	if len(title) > MaxTitleLength {
		return errors.New("title must be at most 200 characters")
	}
	if len(in.Description) > MaxDescriptionLength {
		return errors.New("description must be at most 10000 characters")
	}
	status := in.Status
	if status == "" {
		status = model.StatusToDo
	}
	valid := false
	for _, s := range model.ValidStatuses {
		if status == s {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("invalid status")
	}
	return nil
}

// Create creates an issue with validation and default status.
func (s *IssueService) Create(ctx context.Context, in *model.CreateIssueInput) (*model.Issue, error) {
	if err := ValidateCreateInput(in); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	status := strings.TrimSpace(in.Status)
	if status == "" {
		status = model.StatusToDo
	}
	issue := &model.Issue{
		ID:          uuid.New().String(),
		Title:       strings.TrimSpace(in.Title),
		Description: strings.TrimSpace(in.Description),
		Status:      status,
		AssigneeID:  strings.TrimSpace(in.AssigneeID),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.Create(ctx, issue); err != nil {
		return nil, err
	}
	return issue, nil
}
