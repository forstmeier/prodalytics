//+build !test

package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/forstmeier/todalytics/pkg/evt"
	"github.com/forstmeier/todalytics/pkg/tbl"
)

func main() {
	evtClient := evt.New(os.Getenv("API_TOKEN"))

	newSession := session.New()

	tblClient := tbl.New(
		newSession,
		os.Getenv("TABLE_NAME"),
	)

	lambda.Start(handler(evtClient, tblClient, os.Getenv("CLIENT_SECRET")))
}
