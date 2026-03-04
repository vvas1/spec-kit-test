package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"issue-tracker/backend/internal/model"
)

// UserStore handles user (assignee list) persistence.
type UserStore struct {
	c *mongo.Collection
}

// NewUserStore returns a new UserStore for the given database.
func NewUserStore(db *mongo.Database) *UserStore {
	return &UserStore{c: db.Collection("users")}
}

// List returns all users.
func (s *UserStore) List(ctx context.Context) ([]model.User, error) {
	cur, err := s.c.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var users []model.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	if users == nil {
		users = []model.User{}
	}
	return users, nil
}
