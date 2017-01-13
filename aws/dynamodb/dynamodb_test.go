package dynamodb

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dtan4/valec/secret"
	"github.com/dtan4/valecjr/aws/mock"
	"github.com/golang/mock/gomock"
)

func TestListSecrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockDynamoDBAPI(ctrl)

	api.EXPECT().Query(&dynamodb.QueryInput{
		TableName: aws.String("valec"),
		KeyConditions: map[string]*dynamodb.Condition{
			"namespace": &dynamodb.Condition{
				ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
				AttributeValueList: []*dynamodb.AttributeValue{
					&dynamodb.AttributeValue{
						S: aws.String("test"),
					},
				},
			},
		},
	}).Return(&dynamodb.QueryOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("BAZ"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("1"),
				},
			},
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("FOO"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("bar"),
				},
			},
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("BAR"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("fuga"),
				},
			},
		},
	}, nil)
	client := &Client{
		api: api,
	}

	expected := []*secret.Secret{
		&secret.Secret{
			Key:   "BAZ",
			Value: "1",
		},
		&secret.Secret{
			Key:   "FOO",
			Value: "bar",
		},
		&secret.Secret{
			Key:   "BAR",
			Value: "fuga",
		},
	}

	table := "valec"
	namespace := "test"
	actual, err := client.ListSecrets(table, namespace)
	if err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Secrets does not match. expected: %v, actual: %v", expected, actual)
	}
}

func TestListNamespaces(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockDynamoDBAPI(ctrl)

	api.EXPECT().Scan(&dynamodb.ScanInput{
		TableName: aws.String("valec"),
	}).Return(&dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("BAZ"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("1"),
				},
			},
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test2"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("FOO"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("bar"),
				},
			},
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("BAR"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("fuga"),
				},
			},
			map[string]*dynamodb.AttributeValue{
				"namespace": &dynamodb.AttributeValue{
					S: aws.String("test3"),
				},
				"key": &dynamodb.AttributeValue{
					S: aws.String("FOO"),
				},
				"value": &dynamodb.AttributeValue{
					S: aws.String("fuga"),
				},
			},
		},
	}, nil)
	client := &Client{
		api: api,
	}

	expected := []string{
		"test",
		"test2",
		"test3",
	}

	table := "valec"
	actual, err := client.ListNamespaces(table)
	if err != nil {
		t.Errorf("Error should not be raised. error: %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Namespaces does not match. expected: %q, actual: %q", expected, actual)
	}
}
