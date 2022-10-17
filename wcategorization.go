package websitecategorization

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// WCategorizationService is an interface for Website Categorization API.
type WCategorizationService interface {
	// Get returns parsed Website Categorization API response.
	Get(ctx context.Context, domainName string, opts ...Option) (*WCategorizationResponse, *Response, error)

	// GetRaw returns raw Website Categorization API response as the Response struct with Body saved as a byte slice.
	GetRaw(ctx context.Context, domainName string, opts ...Option) (*Response, error)

	// GetAllCategories returns all possible categories.
	GetAllCategories(ctx context.Context, opts ...Option) (categories []CategoryItem, response *Response, err error)

	// GetAllCategoriesRaw returns all possible categories as a raw API response.
	GetAllCategoriesRaw(ctx context.Context, opts ...Option) (response *Response, err error)
}

// Response is the http.Response wrapper with Body saved as a byte slice.
type Response struct {
	*http.Response

	// Body is the byte slice representation of http.Response Body
	Body []byte
}

// wCategorizationServiceOp is the type implementing the WCategorization interface.
type wCategorizationServiceOp struct {
	client  *Client
	baseURL *url.URL
}

var _ WCategorizationService = &wCategorizationServiceOp{}

// newRequest creates the API request with default parameters and the specified apiKey.
func (service wCategorizationServiceOp) newRequest() (*http.Request, error) {
	req, err := service.client.NewRequest(http.MethodGet, service.baseURL, nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("apiKey", service.client.apiKey)

	req.URL.RawQuery = query.Encode()

	return req, nil
}

// apiResponse is used for parsing Website Categorization API response as a model instance.
type apiResponse struct {
	WCategorizationResponse
	ErrorMessage
}

// requestCategories returns intermediate API response for the /categories path.
func (service wCategorizationServiceOp) requestCategories(ctx context.Context, opts ...Option) (*Response, error) {
	categoriesURL := service.baseURL
	categoriesURL.Path += "/categories"

	req, err := service.client.NewRequest(http.MethodGet, categoriesURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := service.request(ctx, req, opts...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// requestBase returns intermediate API response for the base path.
func (service wCategorizationServiceOp) requestBase(ctx context.Context, domainName string, opts ...Option) (*Response, error) {
	if domainName == "" {
		return nil, &ArgError{"domainName", "can not be empty"}
	}

	req, err := service.newRequest()
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("domainName", domainName)
	req.URL.RawQuery = q.Encode()

	resp, err := service.request(ctx, req, opts...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// request returns intermediate API response for further actions.
func (service wCategorizationServiceOp) request(ctx context.Context, req *http.Request, opts ...Option) (*Response, error) {
	q := req.URL.Query()
	for _, opt := range opts {
		opt(q)
	}
	req.URL.RawQuery = q.Encode()

	var b bytes.Buffer

	resp, err := service.client.Do(ctx, req, &b)
	if err != nil {
		return &Response{
			Response: resp,
			Body:     b.Bytes(),
		}, err
	}

	return &Response{
		Response: resp,
		Body:     b.Bytes(),
	}, nil
}

// parse parses raw Website Categorization API response.
func parse(raw []byte) (*apiResponse, error) {
	var response apiResponse

	err := json.NewDecoder(bytes.NewReader(raw)).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	return &response, nil
}

// parseCategories parses raw Website Categorization API response with an array of categories.
func parseCategories(raw []byte) ([]CategoryItem, error) {
	var respCategories []CategoryItem

	err := json.NewDecoder(bytes.NewReader(raw)).Decode(&respCategories)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	return respCategories, nil
}

// Get returns parsed Website Categorization API response.
func (service wCategorizationServiceOp) Get(
	ctx context.Context,
	domainName string,
	opts ...Option,
) (wCategorizationResponse *WCategorizationResponse, resp *Response, err error) {
	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionOutputFormat("JSON"))

	resp, err = service.requestBase(ctx, domainName, optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	wCategorizationResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if wCategorizationResp.Message != "" || wCategorizationResp.Code != 0 {
		return nil, nil, &ErrorMessage{
			Code:    wCategorizationResp.Code,
			Message: wCategorizationResp.Message,
		}
	}

	return &wCategorizationResp.WCategorizationResponse, resp, nil
}

// GetRaw returns raw Website Categorization API response as the Response struct with Body saved as a byte slice.
func (service wCategorizationServiceOp) GetRaw(
	ctx context.Context,
	domainName string,
	opts ...Option,
) (resp *Response, err error) {
	resp, err = service.requestBase(ctx, domainName, opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// GetAllCategories returns all possible categories.
func (service wCategorizationServiceOp) GetAllCategories(ctx context.Context, opts ...Option) (
	categories []CategoryItem, resp *Response, err error) {
	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionOutputFormat("JSON"))

	resp, err = service.requestCategories(ctx, optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	respCategories, err := parseCategories(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	return respCategories, resp, nil
}

// GetAllCategoriesRaw returns all possible categories as a raw API response.
func (service wCategorizationServiceOp) GetAllCategoriesRaw(ctx context.Context, opts ...Option) (
	resp *Response, err error) {

	resp, err = service.requestCategories(ctx, opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// ArgError is the argument error.
type ArgError struct {
	Name    string
	Message string
}

// Error returns error message as a string.
func (a *ArgError) Error() string {
	return `invalid argument: "` + a.Name + `" ` + a.Message
}
