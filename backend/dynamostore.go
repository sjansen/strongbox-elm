package main

import (
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sjansen/dynamostore"
)

var (
	dynamostoreCreateTable = os.Getenv("DYNAMOSTORE_CREATE_TABLE")
	dynamostoreEndpoint    = os.Getenv("DYNAMOSTORE_ENDPOINT")
)

func newDynamoStore(endpoint string) (scs.Store, error) {
	config := aws.NewConfig()
	if endpoint != "" {
		creds := credentials.NewStaticCredentials("id", "secret", "token")
		config = config.
			WithCredentials(creds).
			WithRegion("us-west-2").
			WithEndpoint(endpoint)
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	store := dynamostore.New(
		dynamodb.New(sess, config),
	)

	if dynamostoreCreateTable != "" {
		err = store.CreateTable()
		if err != nil {
			return nil, err
		}
	}

	return store, nil
}
