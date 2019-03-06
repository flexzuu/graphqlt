# graphqlt - simple graphql client

This is a striped down version of [maschinebox/graphql](https://github.com/machinebox/graphql)

## Installation
Make sure you have a working Go environment. To install graphql, simply run:

```
$ go get github.com/flexzuu/graphqlt
```

## Usage

```go
import "context"

// create a client (safe to share across requests)
client := graphqlt.NewClient("https://yourserver.example/graphql")

// make a request
req := graphqlt.NewRequest(`
    query ($key: String!) {
        items (id:$key) {
            field1
            field2
            field3
        }
    }
`)

// set any variables
req.Var("key", "value")

// set header fields
req.Header.Set("Cache-Control", "no-cache")

// define a Context for the request
ctx := context.Background()

// run it and capture the response
var respData struct {
    Data struct {
        items []entity.Item
    }
}
if err := client.Run(ctx, req, &respData); err != nil {
    log.Fatal(err)
}
```
