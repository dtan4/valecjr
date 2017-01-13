package aws

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	dynamodbapi "github.com/aws/aws-sdk-go/service/dynamodb"
	stsapi "github.com/aws/aws-sdk-go/service/sts"
	"github.com/dtan4/valec/aws/dynamodb"
	"github.com/dtan4/valecjr/aws/sts"
	"github.com/pkg/errors"
)

var (
	// AccessKeyID represents AWS_ACCESS_KEY_ID
	AccessKeyID string
	// SecretAccessKey represents AWS_SECRET_ACCESS_KEY
	SecretAccessKey string
	// Region represents AWS_REGION
	Region string
	// IAMRoleARN represents IAM Role ARN to use
	IAMRoleARN string

	// DynamoDB represents DynamoDB API client
	DynamoDB *dynamodb.Client
	// STS represents STS API client
	STS *sts.Client
)

// Initialize initialized AWS API Clients
func Initialize() error {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(AccessKeyID, SecretAccessKey, ""),
		Region:      aws.String(Region),
	})
	if err != nil {
		return errors.Wrap(err, "Failed to create new AWS session with embedded credentials.")
	}

	STS = sts.NewClient(stsapi.New(sess))
	creds, err := STS.AssumeRole(IAMRoleARN, sessionName())
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve temporary credentials.")
	}

	tmpSess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(creds.AccessKeyID, creds.SecretAccessKey, creds.SessionToken),
		Region:      aws.String(Region),
	})
	if err != nil {
		return errors.Wrap(err, "Failed to create new AWS session with temporary credentials.")
	}

	DynamoDB = dynamodb.NewClient(dynamodbapi.New(tmpSess))

	return nil
}

func sessionName() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
