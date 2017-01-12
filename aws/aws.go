package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	stsapi "github.com/aws/aws-sdk-go/service/sts"
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
		return errors.Wrap(err, "Failed to create new AWS session.")
	}

	STS = sts.NewClient(stsapi.New(sess))

	return nil
}
