[![website-categorization-go license](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![website-categorization-go made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://pkg.go.dev/github.com/whois-api-llc/website-categorization-go)
[![website-categorization-go test](https://github.com/whois-api-llc/website-categorization-go/workflows/Test/badge.svg)](https://github.com/whois-api-llc/website-categorization-go/actions/)

# Overview

The client library for
[Website Categorization API](https://website-categorization.whoisxmlapi.com/)
in Go language.

The minimum go version is 1.17.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/website-categorization-go
```

# Examples

Full API documentation available [here](https://website-categorization.whoisxmlapi.com/api/documentation/v3/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := websitecategorization.NewBasicClient(apiKey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := websitecategorization.NewClient(apiKey, websitecategorization.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

Website Categorization API lets you get all supported categories for websites.

```go

// Make request to get a list of categories by a domain name as a model instance.
wCategorizationResp, _, err := client.Get(ctx, "whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

for _, obj := range wCategorizationResp.Categories {
	log.Printf("ID: %d, Name: %s, Confidence: %f ", obj.ID, obj.Name, obj.Tier1.Confidence)
}

// Make request to get raw data in XML.
resp, err := client.GetRaw(context.Background(), "whoisxmlapi.com",
    websitecategorization.OptionOutputFormat("XML"))
if err != nil {
    log.Fatal(err)
}

log.Println(string(resp.Body))

```