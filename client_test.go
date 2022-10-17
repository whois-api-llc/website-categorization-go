package websitecategorization

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathWCategorizationResponseOK         = "/WCategorization/ok"
	pathWCategorizationResponseError      = "/WCategorization/error"
	pathWCategorizationResponse500        = "/WCategorization/500"
	pathWCategorizationResponsePartial1   = "/WCategorization/partial"
	pathWCategorizationResponsePartial2   = "/WCategorization/partial2"
	pathWCategorizationResponseUnparsable = "/WCategorization/unparsable"

	pathWCategorizationCategoriesResponseOK         = "/WCategorization/ok/categories"
	pathWCategorizationCategoriesResponseError      = "/WCategorization/error/categories"
	pathWCategorizationCategoriesResponse500        = "/WCategorization/500/categories"
	pathWCategorizationCategoriesResponsePartial1   = "/WCategorization/partial/categories"
	pathWCategorizationCategoriesResponsePartial2   = "/WCategorization/partial2/categories"
	pathWCategorizationCategoriesResponseUnparsable = "/WCategorization/unparsable/categories"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

// dummyServer is the sample of the Website Categorization API server for testing.
func dummyServer(resp, respUnparsable string, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathWCategorizationResponseOK, pathWCategorizationCategoriesResponseOK:
		case pathWCategorizationResponseError, pathWCategorizationCategoriesResponseError:
			w.WriteHeader(499)
			response = respErr
		case pathWCategorizationResponse500, pathWCategorizationCategoriesResponse500:
			w.WriteHeader(500)
			response = respUnparsable
		case pathWCategorizationResponsePartial1, pathWCategorizationCategoriesResponsePartial1:
			response = response[:len(response)-10]
		case pathWCategorizationResponsePartial2, pathWCategorizationCategoriesResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		case pathWCategorizationResponseUnparsable, pathWCategorizationCategoriesResponseUnparsable:
			response = respUnparsable
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

// newAPI returns new Website Categorization API client for testing.
func newAPI(apiServer *httptest.Server, link string) *Client {
	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}

	apiURL.Path = link

	params := ClientParams{
		HTTPClient:             apiServer.Client(),
		WCategorizationBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

// TestWCategorizationGet tests the Get function.
func TestWCategorizationGet(t *testing.T) {
	checkResultRec := func(res *WCategorizationResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"categories":[{"tier1":{"confidence":0.9499015947838411,"id":"IAB-596",
"name":"Technology & Computing"},"tier2":{"confidence":0.8420597541031617,"id":"IAB-618",
"name":"Information and Network Security"}}],"domainName":"whoisxmlapi.com","websiteResponded":true}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory string
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathWCategorizationResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathWCategorizationResponse500,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathWCategorizationResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathWCategorizationResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathWCategorizationResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "API error: [499] Test error message.",
		},
		{
			name: "unparsable response",
			path: pathWCategorizationResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument1",
			path: pathWCategorizationResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "domainName" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.Get(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("WCategorization.Get() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("WCategorization.Get() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("WCategorization.Get() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestWCategorizationGetRaw tests the GetRaw function.
func TestWCategorizationGetRaw(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"categories":[{"tier1":{"confidence":0.9499015947838411,"id":"IAB-596",
"name":"Technology & Computing"},"tier2":{"confidence":0.8420597541031617,"id":"IAB-618",
"name":"Information and Network Security"}}],"domainName":"whoisxmlapi.com","websiteResponded":true}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory string
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathWCategorizationResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathWCategorizationResponse500,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathWCategorizationResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathWCategorizationResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathWCategorizationResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathWCategorizationResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "API failed with status code: 499",
		},
		{
			name: "invalid argument1",
			path: pathWCategorizationResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: `invalid argument: "domainName" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.GetRaw(tt.args.ctx, tt.args.options.mandatory)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("WCategorization.GetRaw() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("WCategorization.GetRaw() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}

// TestGetAllCategories tests the GetAllCategories function.
func TestGetAllCategories(t *testing.T) {
	checkResultRec := func(res []CategoryItem) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `[{"id":"IAB-1","name":"Automotive","parent":null},{"id":"CUS-1","name":"Trucks","parent":"IAB-1"},
{"id":"CUS-2","name":"Cars","parent":"IAB-1"},{"id":"IAB-25","name":"Car Culture","parent":"IAB-1"}]`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		option Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathWCategorizationResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					OptionOutputFormat("JSON"),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathWCategorizationResponse500,
			args: args{
				ctx: ctx,
				options: options{
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathWCategorizationResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathWCategorizationResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathWCategorizationResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "could not process request",
			path: pathWCategorizationResponseError,
			args: args{
				ctx: ctx,
				options: options{
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: json: cannot unmarshal object into Go value of type []websitecategorization.CategoryItem",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.GetAllCategories(tt.args.ctx, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("WCategorization.GetGetAlCategories() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("WCategorization.GetAlCategories() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("WCategorization.GetAllCategories() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestGetAllCategoriesRaw tests the GetAllCategoriesRaw function.
func TestGetAllCategoriesRaw(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `[{"id":"IAB-1","name":"Automotive","parent":null},{"id":"CUS-1","name":"Trucks","parent":"IAB-1"},
{"id":"CUS-2","name":"Cars","parent":"IAB-1"},{"id":"IAB-25","name":"Car Culture","parent":"IAB-1"}]`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"messages":"Test error message."}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		option Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathWCategorizationResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathWCategorizationResponse500,
			args: args{
				ctx: ctx,
				options: options{
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathWCategorizationResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathWCategorizationResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathWCategorizationResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathWCategorizationResponseError,
			args: args{
				ctx: ctx,
				options: options{
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "API failed with status code: 499",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.GetAllCategoriesRaw(tt.args.ctx, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("WCategorization.GetAllCategoriesRaw() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("WCategorization.GetAllCategoriesRaw() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}
