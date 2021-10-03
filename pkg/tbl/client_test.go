package tbl

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type mockDynamoDBClient struct {
	mockDynamoDBClientError error
}

func (m *mockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, m.mockDynamoDBClientError
}

func TestNew(t *testing.T) {
	client := New(session.New(), "tableName")
	if client == nil {
		t.Error("error creating tbl client")
	}
}

func TestAppendRow(t *testing.T) {
	tests := []struct {
		description             string
		mockDynamoDBClientError error
		error                   error
	}{
		{
			description:             "error append row",
			mockDynamoDBClientError: errors.New("mock put item error"),
			error:                   &AppendRowError{},
		},
		{
			description:             "successful append row invocation",
			mockDynamoDBClientError: nil,
			error:                   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			client := &Client{
				dynamoDBClient: &mockDynamoDBClient{
					mockDynamoDBClientError: test.mockDynamoDBClientError,
				},
			}

			parentID := 4
			sectionID := 5
			sectionName := "section name"
			dateAdded := time.Now()
			dateCompleted := dateAdded.Add(time.Hour * 1)

			err := client.AppendRow(context.Background(), Row{
				ID:            "id",
				ItemID:        1,
				Event:         "item:added",
				UserID:        2,
				UserEmail:     "test@email.com",
				ProjectID:     3,
				ProjectName:   "project name",
				Content:       "task content",
				Description:   "task description",
				Notes:         []string{"task note"},
				Priority:      1,
				ParentID:      &parentID,
				SectionID:     &sectionID,
				SectionName:   &sectionName,
				LabelIDs:      []int{6},
				LabelNames:    []string{"label name"},
				Checked:       true,
				DateAdded:     dateAdded,
				DateCompleted: &dateCompleted,
			})

			if err != nil {
				switch e := test.error.(type) {
				case *AppendRowError:
					if !errors.As(err, &e) {
						t.Errorf("incorrect error, received: %v, expected: %v", err, e)
					}
				default:
					t.Fatalf("unexpected error type: %v", err)
				}
			}
		})
	}
}
