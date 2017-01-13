package dynamodb

import (
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/dtan4/valec/secret"
	"github.com/pkg/errors"
)

// Client represents a wrapper of DynamoDB API client
type Client struct {
	api dynamodbiface.DynamoDBAPI
}

// NewClient creates new Client object
func NewClient(api dynamodbiface.DynamoDBAPI) *Client {
	return &Client{
		api: api,
	}
}

// ListSecrets returns all secrets in the given table and namespace
func (c *Client) ListSecrets(table, namespace string) ([]*secret.Secret, error) {
	keyConditions := map[string]*dynamodb.Condition{
		"namespace": &dynamodb.Condition{
			ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
			AttributeValueList: []*dynamodb.AttributeValue{
				&dynamodb.AttributeValue{
					S: aws.String(namespace),
				},
			},
		},
	}
	params := &dynamodb.QueryInput{
		TableName:     aws.String(table),
		KeyConditions: keyConditions,
	}

	resp, err := c.api.Query(params)
	if err != nil {
		return []*secret.Secret{}, errors.Wrapf(err, "Failed to list up secrets. namespace=%s", namespace)
	}

	secrets := []*secret.Secret{}

	for _, item := range resp.Items {
		secret := &secret.Secret{
			Key:   *item["key"].S,
			Value: *item["value"].S,
		}

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

// ListNamespaces returns all namespaces
func (c *Client) ListNamespaces(table string) ([]string, error) {
	resp, err := c.api.Scan(&dynamodb.ScanInput{
		TableName: aws.String(table),
	})
	if err != nil {
		return []string{}, errors.Wrapf(err, "Failed to retrieve items from DynamoDB table. table=%s", table)
	}

	nsmap := map[string]bool{}

	for _, item := range resp.Items {
		nsmap[*item["namespace"].S] = true
	}

	namespaces := []string{}

	for k := range nsmap {
		namespaces = append(namespaces, k)
	}

	sort.Strings(namespaces)

	return namespaces, nil
}
