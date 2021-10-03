package tbl

import (
	"context"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var _ Tabler = &Client{}

// Client implements the tbl.Tabler methods using
// AWS DynamoDB.
type Client struct {
	tableName      string
	dynamoDBClient dynamoDBClient
}

type dynamoDBClient interface {
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
}

// New generates a Client pointer instance.
func New(newSession *session.Session, tableName string) *Client {
	dynamoDBClient := dynamodb.New(newSession)

	return &Client{
		tableName:      tableName,
		dynamoDBClient: dynamoDBClient,
	}
}

// AppendRow implements the tbl.Tabler.AppendRow
// interface method.
func (c *Client) AppendRow(ctx context.Context, row Row) error {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(row.ID),
			},
			"item_id": {
				N: aws.String(strconv.Itoa(row.ItemID)),
			},
			"event": {
				S: aws.String(row.Event),
			},
			"user_id": {
				N: aws.String(strconv.Itoa(row.UserID)),
			},
			"user_email": {
				S: aws.String(row.UserEmail),
			},
			"project_id": {
				N: aws.String(strconv.Itoa(row.ProjectID)),
			},
			"project_name": {
				S: aws.String(row.ProjectName),
			},
			"content": {
				S: aws.String(row.Content),
			},
			"description": {
				S: aws.String(row.Description),
			},
			"notes": {
				S: aws.String(strings.Join(row.Notes, "\n\n")),
			},
			"priority": {
				N: aws.String(strconv.Itoa(row.Priority)),
			},
			"checked": {
				BOOL: aws.Bool(row.Checked),
			},
			"date_added": {
				S: aws.String(row.DateAdded.String()),
			},
		},
		TableName: aws.String(c.tableName),
	}

	if row.ParentID != nil {
		input.Item["parent_id"] = &dynamodb.AttributeValue{
			N: aws.String(strconv.Itoa(*row.ParentID)),
		}
	}

	if row.SectionID != nil {
		input.Item["section_id"] = &dynamodb.AttributeValue{
			N: aws.String(strconv.Itoa(*row.SectionID)),
		}
	}

	if row.SectionName != nil {
		input.Item["section_name"] = &dynamodb.AttributeValue{
			S: aws.String(*row.SectionName),
		}
	}

	if len(row.LabelIDs) > 0 {
		labelIDs := []*string{}
		for i := range row.LabelIDs {
			labelIDString := strconv.Itoa(row.LabelIDs[i])
			labelIDs = append(labelIDs, &labelIDString)
		}
		input.Item["label_ids"] = &dynamodb.AttributeValue{
			NS: labelIDs,
		}
	}

	if len(row.LabelNames) > 0 {
		labelNames := []*string{}
		for i := range row.LabelNames {
			labelNames = append(labelNames, &row.LabelNames[i])
		}
		input.Item["label_names"] = &dynamodb.AttributeValue{
			SS: labelNames,
		}
	}

	if row.DateCompleted != nil {
		input.Item["date_completed"] = &dynamodb.AttributeValue{
			S: aws.String(row.DateCompleted.String()),
		}
	}

	_, err := c.dynamoDBClient.PutItem(input)
	if err != nil {
		return &AppendRowError{err: err}
	}

	return nil
}
