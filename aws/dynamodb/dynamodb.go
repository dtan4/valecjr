package dynamodb

import (
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
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
