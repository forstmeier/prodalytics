package evt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/forstmeier/todalytics/pkg/tbl"
)

var _ Eventer = &Client{}

// Client implements the evt.Eventer methods using HTTP requests.
type Client struct {
	helper helper
}

// New generates a Client pointer instance.
func New(authorizationToken string) *Client {
	return &Client{
		helper: &help{
			authorizationHeader: fmt.Sprintf("Bearer %s", authorizationToken),
			httpClient:          http.Client{},
			getData:             getData,
		},
	}
}

// Convert implements the evt.Eventer.Convert interface method.
func (c *Client) Convert(ctx context.Context, data []byte) (*tbl.Row, error) {
	event := Event{}
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, &ConvertError{
			err:      err,
			function: "json unmarshal",
		}
	}

	row := tbl.Row{
		ID:          uuid.NewString(),
		ItemID:      event.EventData.ID,
		Event:       event.EventName,
		UserID:      event.EventData.UserID,
		UserEmail:   event.Initiator.Email,
		ProjectID:   event.EventData.ProjectID,
		Content:     event.EventData.Content,
		Description: event.EventData.Description,
		Priority:    event.EventData.Priority,
		ParentID:    event.EventData.ParentID,
		SectionID:   event.EventData.SectionID,
		LabelIDs:    event.EventData.Labels,
	}

	extraValues, err := c.helper.getExtraValues(
		ctx,
		event.EventData.ProjectID,
		event.EventData.ID,
		event.EventData.SectionID,
		event.EventData.Labels,
	)
	if err != nil {
		return nil, &ConvertError{
			err:      err,
			function: "get extra values",
		}
	}

	row.ProjectName = extraValues.projectName
	row.Notes = extraValues.notes
	row.SectionName = &extraValues.sectionName
	row.LabelNames = extraValues.labelNames

	dateAdded, dateCompleted, err := parseTimes(
		ctx,
		event.EventData.DateAdded,
		event.EventData.DateCompleted,
	)
	if err != nil {
		return nil, &ConvertError{
			err:      err,
			function: "parse times",
		}
	}

	row.DateAdded = *dateAdded
	if !dateCompleted.IsZero() {
		row.DateCompleted = dateCompleted
	}

	row.Checked = parseChecked(ctx, event.EventData.Checked)

	return &row, nil
}
