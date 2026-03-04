package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"issue-tracker/backend/internal/model"
)

// IssueStore handles issue persistence.
type IssueStore struct {
	c *mongo.Collection
}

// NewIssueStore returns a new IssueStore for the given database.
func NewIssueStore(db *mongo.Database) *IssueStore {
	return &IssueStore{c: db.Collection("issues")}
}

// Create inserts an issue and returns it with ID and timestamps set.
func (s *IssueStore) Create(ctx context.Context, issue *model.Issue) error {
	_, err := s.c.InsertOne(ctx, issue)
	return err
}

// GetByID returns an issue by ID, or nil if not found.
func (s *IssueStore) GetByID(ctx context.Context, id string) (*model.Issue, error) {
	var issue model.Issue
	err := s.c.FindOne(ctx, bson.M{"_id": id}).Decode(&issue)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &issue, nil
}
