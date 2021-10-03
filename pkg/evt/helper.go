package evt

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var _ helper = &help{}

type helper interface {
	getExtraValues(ctx context.Context, projectID, itemID int, sectionID *int, labelIDs []int) (*extraValues, error)
}

type help struct {
	authorizationHeader string
	httpClient          http.Client
	getData             func(ctx context.Context, httpClient http.Client, url, authorizationHeader string, data interface{}) error
}

type project struct {
	Name string `json:"name"`
}

type note struct {
	Content string `json:"content"`
}

type section struct {
	Name string `json:"name"`
}

type label struct {
	Name string `json:"name"`
}

type extraValues struct {
	projectName string
	notes       []string
	sectionName string
	labelNames  []string
}

func (h *help) getExtraValues(ctx context.Context, projectID, itemID int, sectionID *int, labelIDs []int) (*extraValues, error) {
	output := &extraValues{}

	projectData := project{}
	if err := h.getData(ctx, h.httpClient, fmt.Sprintf("https://api.todoist.com/rest/v1/projects/%d", projectID), h.authorizationHeader, &projectData); err != nil {
		return nil, err
	}
	output.projectName = projectData.Name

	notesData := []note{}
	if err := h.getData(ctx, h.httpClient, fmt.Sprintf("https://api.todoist.com/rest/v1/comments?task_id=%d", itemID), h.authorizationHeader, &notesData); err != nil {
		return nil, err
	}

	notesContent := make([]string, len(notesData))
	for i, noteData := range notesData {
		notesContent[i] = noteData.Content
	}
	output.notes = notesContent

	if sectionID != nil {
		sectionData := section{}
		if err := h.getData(ctx, h.httpClient, fmt.Sprintf("https://api.todoist.com/rest/v1/sections/%d", sectionID), h.authorizationHeader, &sectionData); err != nil {
			return nil, err
		}
		output.sectionName = sectionData.Name
	}

	labelsData := label{}
	labelNames := make([]string, len(labelIDs))
	for i, labelID := range labelIDs {
		if err := h.getData(ctx, h.httpClient, fmt.Sprintf("https://api.todoist.com/rest/v1/labels/%d", labelID), h.authorizationHeader, &labelsData); err != nil {
			return nil, err
		}

		labelNames[i] = labelsData.Name
	}
	output.labelNames = labelNames

	return output, nil
}

func getData(ctx context.Context, httpClient http.Client, url, authorizationHeader string, data interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", authorizationHeader)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return err
	}

	return nil
}

func parseTimes(ctx context.Context, dateAdded string, dateCompleted *string) (*time.Time, *time.Time, error) {
	dateAddedTime, err := time.Parse(time.RFC3339, dateAdded)
	if err != nil {
		return nil, nil, err
	}

	var dateCompletedTime time.Time
	if dateCompleted != nil {
		dateCompletedTime, err = time.Parse(time.RFC3339, *dateCompleted)
		if err != nil {
			return nil, nil, err
		}
	}

	return &dateAddedTime, &dateCompletedTime, nil
}

func parseChecked(ctx context.Context, checked int) bool {
	checkedBool := false
	if checked == 1 {
		checkedBool = true
	}

	return checkedBool
}
