package codegen

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/andreyvit/diff"
)

func TestConvertToGo(t *testing.T) {

	type args struct {
		schema  string
		queries []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"xxx",
			args{
				schema: `
type Project {
	name: String
	tagline: String
}

type Query {
	project(name: String): Project
}
					`,
				queries: []string{`
query testQuery {
		project(name: "GraphQL") {
		tagline
	}
}`},
			},
			`package graphql

import (
	"context"
	graphqlt "github.com/flexzuu/graphqlt"
)

type Client struct {
	*graphqlt.Client
}

// NewClient makes a new Client capable of making GraphQL requests.
func NewClient(endpoint string, opts ...graphqlt.ClientOption) *Client {
	c := &Client{graphqlt.NewClient(endpoint, opts...)}
	return c
}

type TestQueryResponse struct {
	Project struct {
		Tagline string
	}
}

func (c *Client) TestQuery(ctx context.Context, c string) (TestQueryResponse, error) {
	req := graphqlt.NewRequest("\nquery testQuery {\n\t\tproject(name: \"GraphQL\") {\n\t\ttagline\n\t}\n}")
	var respData TestQueryResponse
	err := c.Run(ctx, req, &respData)
	return respData, err
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmt.Sprintf("%#v", ConvertToGo(tt.args.schema, tt.args.queries...)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToGo() diff %s", diff.LineDiff(tt.want, got))
			}
		})
	}
}
