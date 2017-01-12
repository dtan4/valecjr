package sts

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/pkg/errors"
)

const (
	// DurationSeconds must be in 900 - 3600 seconds
	// https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html
	defaultDurationSeconds = 900
)

// Client represents a wrapper of STS API client
type Client struct {
	api stsiface.STSAPI
}

// Credentials represents temporary session credentials
type Credentials struct {
	// AccessKeyID represents AWS_ACCESS_KEY_ID for temporary session
	AccessKeyID string
	// SecretAccessKey represents AWS_SECRET_ACCESS_KEY for temporary session
	SecretAccessKey string
	// SessionToken represents AWS_SESSION_TOKEN for temporary session
	SessionToken string
}

// NewClient creates new Client object
func NewClient(api stsiface.STSAPI) *Client {
	return &Client{
		api: api,
	}
}

// AssumeRole obtains temporary security credentials of the given IAM role
func (c *Client) AssumeRole(roleARN, sessionName string) (*Credentials, error) {
	resp, err := c.api.AssumeRole(&sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(defaultDurationSeconds),
		RoleArn:         aws.String(roleARN),
		RoleSessionName: aws.String(sessionName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute AssumeRole API.")
	}

	return &Credentials{
		AccessKeyID:     aws.StringValue(resp.Credentials.AccessKeyId),
		SecretAccessKey: aws.StringValue(resp.Credentials.SecretAccessKey),
		SessionToken:    aws.StringValue(resp.Credentials.SessionToken),
	}, nil
}
