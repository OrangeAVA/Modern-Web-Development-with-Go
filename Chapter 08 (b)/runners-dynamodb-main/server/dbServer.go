package server

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/spf13/viper"
)

func InitDatabase(config *viper.Viper) *dynamodb.DynamoDB {
	connectionString := config.GetString("database.connection_string")
	awsRegion := config.GetString("database.aws_region")
	awsAccessKeyID := config.GetString("database.aws_access_key_id")
	awsSecretAccessKey := config.GetString("database.aws_secret_access_key")

	if connectionString == "" {
		log.Fatalf("Database connectin string is missing")
	}

	session := session.Must(session.NewSession(
		&aws.Config{
			Region:      aws.String(awsRegion),
			Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
			Endpoint:    aws.String(connectionString),
		},
	))

	return dynamodb.New(session)
}
