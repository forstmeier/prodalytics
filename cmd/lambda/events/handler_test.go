package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/forstmeier/todalytics/pkg/tbl"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

type mockEVTClient struct {
	mockEVTClientOutput *tbl.Row
	mockEVTClientError  error
}

func (m *mockEVTClient) Convert(ctx context.Context, data []byte) (*tbl.Row, error) {
	return m.mockEVTClientOutput, m.mockEVTClientError
}

type mockTBLClient struct {
	row                tbl.Row
	mockTBLClientError error
}

func (m *mockTBLClient) AppendRow(ctx context.Context, row tbl.Row) error {
	m.row = row

	return m.mockTBLClientError
}

func Test_handler(t *testing.T) {
	clientSecret := "clientSecret"

	body := `{
		"event_name": "item:added",
		"event_data": {
			"checked": 0,
			"content": "A new task",
			"description": "",
			"date_added": "2021-02-10T10:33:38Z",
			"date_completed": null,
			"id": 2995104339,
			"labels": [],
			"parent_id": null,
			"priority": 1,
			"project_id": 2203306141,
			"section_id": null
		},
		"initiator": {
		  "email": "alice@example.com"
		}
	}`

	key := []byte(clientSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(body))
	securityHeaderValue := h.Sum(nil)

	tests := []struct {
		description         string
		mockEVTClientOutput *tbl.Row
		mockEVTClientError  error
		mockTBLClientError  error
		request             events.APIGatewayProxyRequest
		row                 tbl.Row
		error               error
	}{
		{
			description:         "no user agent header provided",
			mockEVTClientOutput: nil,
			mockEVTClientError:  nil,
			mockTBLClientError:  nil,
			request: events.APIGatewayProxyRequest{
				Headers: map[string]string{},
			},
			row:   tbl.Row{},
			error: errIncorrectUserAgentHeaderValue,
		},
		{
			description:         "no security header provided",
			mockEVTClientOutput: nil,
			mockEVTClientError:  nil,
			mockTBLClientError:  nil,
			request: events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"User-Agent": "Todoist-Webhooks",
				},
			},
			row:   tbl.Row{},
			error: errNoSecurityHeader,
		},
		{
			description:         "incorrect security header value",
			mockEVTClientOutput: nil,
			mockEVTClientError:  nil,
			mockTBLClientError:  nil,
			request: events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"User-Agent":            "Todoist-Webhooks",
					"x-todoist-hmac-sha256": base64.StdEncoding.EncodeToString([]byte("incorrect security header")),
				},
			},
			row:   tbl.Row{},
			error: errIncorrectSecurityHeaderValue,
		},
		{
			description:         "error convert event to row",
			mockEVTClientOutput: nil,
			mockEVTClientError:  errors.New("mock convert error"),
			mockTBLClientError:  nil,
			request: events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"User-Agent":            "Todoist-Webhooks",
					"x-todoist-hmac-sha256": base64.StdEncoding.EncodeToString(securityHeaderValue),
				},
				Body: body,
			},
			row:   tbl.Row{},
			error: errConvertEvent,
		},
		{
			description:         "error append row data",
			mockEVTClientOutput: &tbl.Row{},
			mockEVTClientError:  nil,
			mockTBLClientError:  errors.New("mock append row error"),
			request: events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"User-Agent":            "Todoist-Webhooks",
					"x-todoist-hmac-sha256": base64.StdEncoding.EncodeToString(securityHeaderValue),
				},
				Body: body,
			},
			row:   tbl.Row{},
			error: errAppendRowData,
		},
		{
			description:         "successful handler invocation",
			mockEVTClientOutput: &tbl.Row{},
			mockEVTClientError:  nil,
			mockTBLClientError:  nil,
			request: events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"User-Agent":            "Todoist-Webhooks",
					"x-todoist-hmac-sha256": base64.StdEncoding.EncodeToString(securityHeaderValue),
				},
				Body: body,
			},
			row:   tbl.Row{},
			error: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			evtClient := &mockEVTClient{
				mockEVTClientOutput: test.mockEVTClientOutput,
				mockEVTClientError:  test.mockEVTClientError,
			}

			tblClient := &mockTBLClient{
				mockTBLClientError: test.mockTBLClientError,
			}

			handlerFunc := handler(evtClient, tblClient, clientSecret)

			err := handlerFunc(context.Background(), test.request)

			if err != test.error {
				t.Errorf("incorrect error, received: %v, expected: %v", err, test.error)
			}

			if !reflect.DeepEqual(tblClient.row, test.row) {
				t.Errorf("incorrect row, received: %+v, expected: %+v", tblClient.row, test.row)
			}
		})
	}
}
