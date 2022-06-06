package ydb

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

func New() *dynamo.DB {
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess,
		&aws.Config{
			Region:   aws.String(os.Getenv("AWS_DEFAULT_REGION")),
			Endpoint: aws.String(os.Getenv("YDB_ENDPOINT")),
		},
	)
	return db
}
