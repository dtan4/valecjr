package sts

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/dtan4/valecjr/aws/mock"
	"github.com/golang/mock/gomock"
)

func TestAssumeRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockSTSAPI(ctrl)
	api.EXPECT().AssumeRole(&sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(900),
		RoleArn:         aws.String("arn:aws:iam::123456789012:role/S3Access"),
		RoleSessionName: aws.String("testsession"),
	}).Return(&sts.AssumeRoleOutput{
		Credentials: &sts.Credentials{
			AccessKeyId:     aws.String("AKIAxxxxxxxxxxxxxxxx"),
			SecretAccessKey: aws.String("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"),
			SessionToken:    aws.String("yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"),
		},
	}, nil)

	client := &Client{
		api: api,
	}
	role := "arn:aws:iam::123456789012:role/S3Access"
	sessionName := "testsession"
	expected := &Credentials{
		AccessKeyID:     "AKIAxxxxxxxxxxxxxxxx",
		SecretAccessKey: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		SessionToken:    "yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy",
	}

	got, err := client.AssumeRole(role, sessionName)
	if err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Credentials does not match. expected: %+v, got: %+v", expected, got)
	}
}
