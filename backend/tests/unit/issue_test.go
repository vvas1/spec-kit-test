package unit

import (
	"testing"

	"issue-tracker/backend/internal/model"
	"issue-tracker/backend/internal/service"
)

func TestValidateCreateInput(t *testing.T) {
	tests := []struct {
		name    string
		in      model.CreateIssueInput
		wantErr bool
	}{
		{"empty title", model.CreateIssueInput{Title: ""}, true},
		{"whitespace title", model.CreateIssueInput{Title: "   "}, true},
		{"valid minimal", model.CreateIssueInput{Title: "a"}, false},
		{"title too long", model.CreateIssueInput{Title: string(make([]byte, 201))}, true},
		{"title 200 ok", model.CreateIssueInput{Title: string(make([]byte, 200))}, false},
		{"description too long", model.CreateIssueInput{Title: "x", Description: string(make([]byte, 10001))}, true},
		{"description 10k ok", model.CreateIssueInput{Title: "x", Description: string(make([]byte, 10000))}, false},
		{"invalid status", model.CreateIssueInput{Title: "x", Status: "Invalid"}, true},
		{"valid status To Do", model.CreateIssueInput{Title: "x", Status: "To Do"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateCreateInput(&tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCreateInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
