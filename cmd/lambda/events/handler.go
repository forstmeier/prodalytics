package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"github.com/forstmeier/todalytics/pkg/evt"
	"github.com/forstmeier/todalytics/pkg/tbl"
	"github.com/forstmeier/todalytics/util"
)

var (
	errIncorrectUserAgentHeaderValue = errors.New("incorrect user agent header value")
	errNoSecurityHeader              = errors.New("no security header received")
	errDecodeSecurityHeadr           = errors.New("error decode security header value")
	errIncorrectSecurityHeaderValue  = errors.New("incorrect security header value")
	errConvertEvent                  = errors.New("error converting event")
	errAppendRowData                 = errors.New("error appending row data")
)

func handler(evtClient evt.Eventer, tblClient tbl.Tabler, clientSecret string) func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		util.Log("REQUEST", request)

		userAgentHeader, ok := request.Headers["User-Agent"]
		if !ok || userAgentHeader != "Todoist-Webhooks" {
			util.Log("INCORRECT_USER_AGENT_HEADER_VALUE", fmt.Sprintf("header check: %t, header value: %s", !ok, userAgentHeader))
			return events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				Body:            fmt.Sprintf(`{"error": "%s"}`, errIncorrectUserAgentHeaderValue.Error()),
				IsBase64Encoded: false,
			}, errIncorrectUserAgentHeaderValue
		}

		hmacHeader, ok := request.Headers["x-todoist-hmac-sha256"]
		if !ok {
			util.Log("NO_SECURITY_HEADER", fmt.Sprintf("header check: %t", !ok))
			return events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				Body:            fmt.Sprintf(`{"error": "%s"}`, errNoSecurityHeader.Error()),
				IsBase64Encoded: false,
			}, errNoSecurityHeader
		}

		receivedMAC, err := base64.StdEncoding.DecodeString(hmacHeader)
		if err != nil {
			util.Log("BASE64_DECODE_ERROR", err.Error())
			return events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				Body:            fmt.Sprintf(`{"error": "%s"}`, errDecodeSecurityHeadr.Error()),
				IsBase64Encoded: false,
			}, errDecodeSecurityHeadr
		}

		h := hmac.New(sha256.New, []byte(clientSecret))
		h.Write([]byte(request.Body))
		expectedMAC := h.Sum(nil)

		if !hmac.Equal(receivedMAC, expectedMAC) {
			util.Log("INCORRECT_SECURITY_HEADER", "received and expected hmac values do not match")
			return events.APIGatewayProxyResponse{
				StatusCode:      http.StatusBadRequest,
				Body:            fmt.Sprintf(`{"error": "%s"}`, errIncorrectSecurityHeaderValue.Error()),
				IsBase64Encoded: false,
			}, errIncorrectSecurityHeaderValue
		}

		row, err := evtClient.Convert(ctx, []byte(request.Body))
		if err != nil {
			util.Log("CONVERT_ERROR", err.Error())
			return events.APIGatewayProxyResponse{
				StatusCode:      http.StatusInternalServerError,
				Body:            fmt.Sprintf(`{"error": "%s"}`, errConvertEvent.Error()),
				IsBase64Encoded: false,
			}, errConvertEvent
		}

		if err := tblClient.AppendRow(ctx, *row); err != nil {
			util.Log("APPEND_ROW_ERROR", err.Error())
			return events.APIGatewayProxyResponse{
				StatusCode:      http.StatusInternalServerError,
				Body:            fmt.Sprintf(`{"error": "%s"}`, errAppendRowData.Error()),
				IsBase64Encoded: false,
			}, errAppendRowData
		}

		return events.APIGatewayProxyResponse{
			StatusCode:      http.StatusOK,
			Body:            `{"message": "success"}`,
			IsBase64Encoded: false,
		}, nil
	}
}
