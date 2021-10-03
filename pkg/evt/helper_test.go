package evt

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func Test_getExtraValues(t *testing.T) {
	tests := []struct {
		description             string
		mockGetDataErrorURLPath string
		mockGetDataError        error
		extraValues             *extraValues
		error                   string
	}{
		{
			description:             "error get project data",
			mockGetDataErrorURLPath: "projects",
			mockGetDataError:        errors.New("mock get project error"),
			extraValues:             nil,
			error:                   "mock get project error",
		},
		{
			description:             "error get notes data",
			mockGetDataErrorURLPath: "comments",
			mockGetDataError:        errors.New("mock get notes error"),
			extraValues:             nil,
			error:                   "mock get notes error",
		},
		{
			description:             "error get section data",
			mockGetDataErrorURLPath: "sections",
			mockGetDataError:        errors.New("mock get section error"),
			extraValues:             nil,
			error:                   "mock get section error",
		},
		{
			description:             "error get labels data",
			mockGetDataErrorURLPath: "labels",
			mockGetDataError:        errors.New("mock get labels error"),
			extraValues:             nil,
			error:                   "mock get labels error",
		},
		{
			description:             "successful get extra data invocation",
			mockGetDataErrorURLPath: "non_path",
			mockGetDataError:        nil,
			extraValues: &extraValues{
				projectName: "project_name",
				notes:       []string{"note content"},
				sectionName: "section_name",
				labelNames:  []string{"label"},
			},
			error: "",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			h := &help{
				authorizationHeader: "authorizationHeader",
				httpClient:          http.Client{},
				getData: func(ctx context.Context, httpClient http.Client, url, authorizationHeader string, data interface{}) error {
					if strings.Contains(url, test.mockGetDataErrorURLPath) {
						return test.mockGetDataError
					}

					switch dt := data.(type) {
					case *project:
						dt.Name = "project_name"
					case *[]note:
						*dt = []note{
							{
								Content: "note content",
							},
						}
					case *section:
						dt.Name = "section_name"
					case *label:
						dt.Name = "label"
					default:
						t.Fatalf("incorrect type received in mock function [%v]", dt)
					}

					return nil
				},
			}

			projectID := 1
			itemID := 2
			sectionID := 3
			labelIDs := []int{4}
			extraValues, err := h.getExtraValues(context.Background(), projectID, itemID, &sectionID, labelIDs)

			if err != nil {
				if err.Error() != test.error {
					t.Errorf("incorrect error, received: %s, expected: %s", err.Error(), test.error)
				}
			} else {
				if !reflect.DeepEqual(extraValues, test.extraValues) {
					t.Errorf("incorrect extra values, \nreceived: %+v,\nexpected: %+v", extraValues, test.extraValues)
				}
			}
		})
	}
}
