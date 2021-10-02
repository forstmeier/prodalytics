//+build !test

package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/forstmeier/todalytics/pkg/evt"
	"github.com/forstmeier/todalytics/pkg/tbl"
)

func main() {
	evtClient := evt.New(os.Getenv("AUTHORIZATION_TOKEN"))

	sheetsClient, err := sheets.NewService(
		context.Background(),
		option.WithCredentialsFile("creds.json"),
	)
	if err != nil {
		log.Fatalf("error creating sheets client: %s", err.Error())
	}

	tblClient := tbl.New(sheetsClient.Spreadsheets.Values)

	lambda.Start(handler(evtClient, tblClient, os.Getenv("CLIENT_SECRET")))
}
