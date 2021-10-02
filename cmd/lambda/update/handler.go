package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

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

func handler(evtClient evt.Eventer, tblClient tbl.Tabler, clientSecret string) func(ctx context.Context, event events.APIGatewayProxyRequest) error {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) error {
		util.Log("REQUEST", request)

		userAgentHeader, ok := request.Headers["User-Agent"]
		if !ok || userAgentHeader != "Todoist-Webhooks" {
			util.Log("INCORRECT_USER_AGENT_HEADER_VALUE", fmt.Sprintf("header check: %t, header value: %s", !ok, userAgentHeader))
			return errIncorrectUserAgentHeaderValue
		}

		hmacHeader, ok := request.Headers["X-Todoist-Hmac-SHA256"]
		if !ok {
			util.Log("NO_SECURITY_HEADER", fmt.Sprintf("header check: %t", !ok))
			return errNoSecurityHeader
		}

		receivedMAC, err := base64.StdEncoding.DecodeString(hmacHeader)
		if err != nil {
			util.Log("BASE64_DECODE_ERROR", err.Error())
			return errDecodeSecurityHeadr
		}

		h := hmac.New(sha256.New, []byte(clientSecret))
		h.Write([]byte(request.Body))
		expectedMAC := h.Sum(nil)

		if !hmac.Equal(receivedMAC, expectedMAC) {
			util.Log("INCORRECT_SECURITY_HEADER", "received and expected hmac values do not match")
			return errIncorrectSecurityHeaderValue
		}

		row, err := evtClient.Convert(ctx, []byte(request.Body))
		if err != nil {
			util.Log("CONVERT_ERROR", err.Error())
			return errConvertEvent
		}

		if err := tblClient.AppendRow(ctx, *row); err != nil {
			util.Log("APPEND_ROW_ERROR", err.Error())
			return errAppendRowData
		}

		return nil
	}
}
