package gen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type RequestEditorFn func(ctx context.Context, req *http.Request) error

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	Server string

	Client HttpRequestDoer

	RequestEditors []RequestEditorFn
}

type ClientOption func(*Client) error

func NewClient(server string, opts ...ClientOption) (*Client, error) {

	client := Client{
		Server: server,
	}

	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}

	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}

	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

type ClientInterface interface {
	LoginWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	Login(ctx context.Context, body LoginJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	RegisterWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	Register(ctx context.Context, body RegisterJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetCatalog(ctx context.Context, params *GetCatalogParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateMovieWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateMovie(ctx context.Context, body CreateMovieJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	DeleteMovie(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetMovieById(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateMovieWithBody(ctx context.Context, movieId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateMovie(ctx context.Context, movieId openapi_types.UUID, body UpdateMovieJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetCollection(ctx context.Context, params *GetCollectionParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	AddMovieToCollectionWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	AddMovieToCollection(ctx context.Context, body AddMovieToCollectionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	DeleteMovieFromCollection(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetMovieFromCollection(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateCollectionItemWithBody(ctx context.Context, movieId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateCollectionItem(ctx context.Context, movieId openapi_types.UUID, body UpdateCollectionItemJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetUsers(ctx context.Context, params *GetUsersParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	DeleteUser(ctx context.Context, userId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetUserById(ctx context.Context, userId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateUserBlockStatusWithBody(ctx context.Context, userId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateUserBlockStatus(ctx context.Context, userId openapi_types.UUID, body UpdateUserBlockStatusJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) LoginWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewLoginRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Login(ctx context.Context, body LoginJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewLoginRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) RegisterWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRegisterRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Register(ctx context.Context, body RegisterJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewRegisterRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetCatalog(ctx context.Context, params *GetCatalogParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetCatalogRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateMovieWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateMovieRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateMovie(ctx context.Context, body CreateMovieJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateMovieRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteMovie(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteMovieRequest(c.Server, movieId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetMovieById(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetMovieByIdRequest(c.Server, movieId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateMovieWithBody(ctx context.Context, movieId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateMovieRequestWithBody(c.Server, movieId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateMovie(ctx context.Context, movieId openapi_types.UUID, body UpdateMovieJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateMovieRequest(c.Server, movieId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetCollection(ctx context.Context, params *GetCollectionParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetCollectionRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AddMovieToCollectionWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddMovieToCollectionRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AddMovieToCollection(ctx context.Context, body AddMovieToCollectionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddMovieToCollectionRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteMovieFromCollection(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteMovieFromCollectionRequest(c.Server, movieId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetMovieFromCollection(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetMovieFromCollectionRequest(c.Server, movieId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateCollectionItemWithBody(ctx context.Context, movieId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateCollectionItemRequestWithBody(c.Server, movieId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateCollectionItem(ctx context.Context, movieId openapi_types.UUID, body UpdateCollectionItemJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateCollectionItemRequest(c.Server, movieId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetUsers(ctx context.Context, params *GetUsersParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetUsersRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteUser(ctx context.Context, userId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteUserRequest(c.Server, userId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetUserById(ctx context.Context, userId openapi_types.UUID, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetUserByIdRequest(c.Server, userId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateUserBlockStatusWithBody(ctx context.Context, userId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateUserBlockStatusRequestWithBody(c.Server, userId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateUserBlockStatus(ctx context.Context, userId openapi_types.UUID, body UpdateUserBlockStatusJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateUserBlockStatusRequest(c.Server, userId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func NewLoginRequest(server string, body LoginJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewLoginRequestWithBody(server, "application/json", bodyReader)
}

func NewLoginRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/auth/login")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func NewRegisterRequest(server string, body RegisterJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewRegisterRequestWithBody(server, "application/json", bodyReader)
}

func NewRegisterRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/auth/register")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func NewGetCatalogRequest(server string, params *GetCatalogParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/catalog")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {

		queryValues := queryURL.Query()

		var rawQueryFragments []string

		if params.Genre != nil {

			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "genre", *params.Genre, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationQuery, Type: "string", Format: ""}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}

		}

		if params.Search != nil {

			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "search", *params.Search, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationQuery, Type: "string", Format: ""}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}

		}

		if params.Limit != nil {

			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "limit", *params.Limit, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationQuery, Type: "integer", Format: ""}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}

		}

		if params.Offset != nil {

			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "offset", *params.Offset, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationQuery, Type: "integer", Format: ""}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}

		}

		if encoded := queryValues.Encode(); encoded != "" {
			rawQueryFragments = append(rawQueryFragments, encoded)
		}
		queryURL.RawQuery = strings.Join(rawQueryFragments, "&")
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewCreateMovieRequest(server string, body CreateMovieJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateMovieRequestWithBody(server, "application/json", bodyReader)
}

func NewCreateMovieRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/catalog")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func NewDeleteMovieRequest(server string, movieId openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "movie_id", movieId, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "string", Format: "uuid"})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/catalog/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewGetMovieByIdRequest(server string, movieId openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "movie_id", movieId, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "string", Format: "uuid"})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/catalog/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewUpdateMovieRequest(server string, movieId openapi_types.UUID, body UpdateMovieJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateMovieRequestWithBody(server, movieId, "application/json", bodyReader)
}

func NewUpdateMovieRequestWithBody(server string, movieId openapi_types.UUID, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "movie_id", movieId, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "string", Format: "uuid"})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/catalog/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func NewGetCollectionRequest(server string, params *GetCollectionParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/collection")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {

		queryValues := queryURL.Query()

		var rawQueryFragments []string

		if params.Status != nil {

			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "status", *params.Status, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationQuery, Type: "string", Format: ""}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}

		}

		if params.Limit != nil {

			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "limit", *params.Limit, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationQuery, Type: "integer", Format: ""}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}

		}

		if params.Offset != nil {

			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "offset", *params.Offset, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationQuery, Type: "integer", Format: ""}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}

		}

		if encoded := queryValues.Encode(); encoded != "" {
			rawQueryFragments = append(rawQueryFragments, encoded)
		}
		queryURL.RawQuery = strings.Join(rawQueryFragments, "&")
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewAddMovieToCollectionRequest(server string, body AddMovieToCollectionJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAddMovieToCollectionRequestWithBody(server, "application/json", bodyReader)
}

func NewAddMovieToCollectionRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/collection")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func NewDeleteMovieFromCollectionRequest(server string, movieId openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "movie_id", movieId, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "string", Format: "uuid"})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/collection/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewGetMovieFromCollectionRequest(server string, movieId openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "movie_id", movieId, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "string", Format: "uuid"})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/collection/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewUpdateCollectionItemRequest(server string, movieId openapi_types.UUID, body UpdateCollectionItemJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateCollectionItemRequestWithBody(server, movieId, "application/json", bodyReader)
}

func NewUpdateCollectionItemRequestWithBody(server string, movieId openapi_types.UUID, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "movie_id", movieId, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "string", Format: "uuid"})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/collection/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func NewGetUsersRequest(server string, params *GetUsersParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/users")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {

		queryValues := queryURL.Query()

		var rawQueryFragments []string

		if params.Limit != nil {

			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "limit", *params.Limit, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationQuery, Type: "integer", Format: ""}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}

		}

		if params.Offset != nil {

			if queryFrag, err := runtime.StyleParamWithOptions("form", true, "offset", *params.Offset, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationQuery, Type: "integer", Format: ""}); err != nil {
				return nil, err
			} else {
				for _, qp := range strings.Split(queryFrag, "&") {
					rawQueryFragments = append(rawQueryFragments, qp)
				}
			}

		}

		if encoded := queryValues.Encode(); encoded != "" {
			rawQueryFragments = append(rawQueryFragments, encoded)
		}
		queryURL.RawQuery = strings.Join(rawQueryFragments, "&")
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewDeleteUserRequest(server string, userId openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "user_id", userId, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "string", Format: "uuid"})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/users/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewGetUserByIdRequest(server string, userId openapi_types.UUID) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "user_id", userId, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "string", Format: "uuid"})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/users/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewUpdateUserBlockStatusRequest(server string, userId openapi_types.UUID, body UpdateUserBlockStatusJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateUserBlockStatusRequestWithBody(server, userId, "application/json", bodyReader)
}

func NewUpdateUserBlockStatusRequestWithBody(server string, userId openapi_types.UUID, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithOptions("simple", false, "user_id", userId, runtime.StyleParamOptions{ParamLocation: runtime.ParamLocationPath, Type: "string", Format: "uuid"})
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/users/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

type ClientWithResponses struct {
	ClientInterface
}

func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

type ClientWithResponsesInterface interface {
	LoginWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*LoginHTTPResponse, error)

	LoginWithResponse(ctx context.Context, body LoginJSONRequestBody, reqEditors ...RequestEditorFn) (*LoginHTTPResponse, error)

	RegisterWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RegisterHTTPResponse, error)

	RegisterWithResponse(ctx context.Context, body RegisterJSONRequestBody, reqEditors ...RequestEditorFn) (*RegisterHTTPResponse, error)

	GetCatalogWithResponse(ctx context.Context, params *GetCatalogParams, reqEditors ...RequestEditorFn) (*GetCatalogHTTPResponse, error)

	CreateMovieWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateMovieHTTPResponse, error)

	CreateMovieWithResponse(ctx context.Context, body CreateMovieJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateMovieHTTPResponse, error)

	DeleteMovieWithResponse(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*DeleteMovieHTTPResponse, error)

	GetMovieByIdWithResponse(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*GetMovieByIdHTTPResponse, error)

	UpdateMovieWithBodyWithResponse(ctx context.Context, movieId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateMovieHTTPResponse, error)

	UpdateMovieWithResponse(ctx context.Context, movieId openapi_types.UUID, body UpdateMovieJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateMovieHTTPResponse, error)

	GetCollectionWithResponse(ctx context.Context, params *GetCollectionParams, reqEditors ...RequestEditorFn) (*GetCollectionHTTPResponse, error)

	AddMovieToCollectionWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AddMovieToCollectionHTTPResponse, error)

	AddMovieToCollectionWithResponse(ctx context.Context, body AddMovieToCollectionJSONRequestBody, reqEditors ...RequestEditorFn) (*AddMovieToCollectionHTTPResponse, error)

	DeleteMovieFromCollectionWithResponse(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*DeleteMovieFromCollectionHTTPResponse, error)

	GetMovieFromCollectionWithResponse(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*GetMovieFromCollectionHTTPResponse, error)

	UpdateCollectionItemWithBodyWithResponse(ctx context.Context, movieId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateCollectionItemHTTPResponse, error)

	UpdateCollectionItemWithResponse(ctx context.Context, movieId openapi_types.UUID, body UpdateCollectionItemJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateCollectionItemHTTPResponse, error)

	GetUsersWithResponse(ctx context.Context, params *GetUsersParams, reqEditors ...RequestEditorFn) (*GetUsersHTTPResponse, error)

	DeleteUserWithResponse(ctx context.Context, userId openapi_types.UUID, reqEditors ...RequestEditorFn) (*DeleteUserHTTPResponse, error)

	GetUserByIdWithResponse(ctx context.Context, userId openapi_types.UUID, reqEditors ...RequestEditorFn) (*GetUserByIdHTTPResponse, error)

	UpdateUserBlockStatusWithBodyWithResponse(ctx context.Context, userId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateUserBlockStatusHTTPResponse, error)

	UpdateUserBlockStatusWithResponse(ctx context.Context, userId openapi_types.UUID, body UpdateUserBlockStatusJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateUserBlockStatusHTTPResponse, error)
}

type LoginHTTPResponse429Headers struct {
	RetryAfter *int
}

type LoginHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *LoginResponse

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *LoginHTTPResponse429Headers
}

func (r LoginHTTPResponse) GetJSON200() *LoginResponse {
	return r.JSON200
}

func (r LoginHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r LoginHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r LoginHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r LoginHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r LoginHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r LoginHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r LoginHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r LoginHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r LoginHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type RegisterHTTPResponse429Headers struct {
	RetryAfter *int
}

type RegisterHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON201 *User

	JSON400 *BadRequest

	JSON409 *Conflict

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *RegisterHTTPResponse429Headers
}

func (r RegisterHTTPResponse) GetJSON201() *User {
	return r.JSON201
}

func (r RegisterHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r RegisterHTTPResponse) GetJSON409() *Conflict {
	return r.JSON409
}

func (r RegisterHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r RegisterHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r RegisterHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r RegisterHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r RegisterHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r RegisterHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type GetCatalogHTTPResponse429Headers struct {
	RetryAfter *int
}

type GetCatalogHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *MoviePage

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *GetCatalogHTTPResponse429Headers
}

func (r GetCatalogHTTPResponse) GetJSON200() *MoviePage {
	return r.JSON200
}

func (r GetCatalogHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r GetCatalogHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r GetCatalogHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r GetCatalogHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r GetCatalogHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r GetCatalogHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r GetCatalogHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r GetCatalogHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r GetCatalogHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type CreateMovieHTTPResponse429Headers struct {
	RetryAfter *int
}

type CreateMovieHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON201 *Movie

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *CreateMovieHTTPResponse429Headers
}

func (r CreateMovieHTTPResponse) GetJSON201() *Movie {
	return r.JSON201
}

func (r CreateMovieHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r CreateMovieHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r CreateMovieHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r CreateMovieHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r CreateMovieHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r CreateMovieHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r CreateMovieHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r CreateMovieHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r CreateMovieHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type DeleteMovieHTTPResponse429Headers struct {
	RetryAfter *int
}

type DeleteMovieHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON404 *NotFound

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *DeleteMovieHTTPResponse429Headers
}

func (r DeleteMovieHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r DeleteMovieHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r DeleteMovieHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r DeleteMovieHTTPResponse) GetJSON404() *NotFound {
	return r.JSON404
}

func (r DeleteMovieHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r DeleteMovieHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r DeleteMovieHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r DeleteMovieHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r DeleteMovieHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r DeleteMovieHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type GetMovieByIdHTTPResponse429Headers struct {
	RetryAfter *int
}

type GetMovieByIdHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *Movie

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON404 *NotFound

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *GetMovieByIdHTTPResponse429Headers
}

func (r GetMovieByIdHTTPResponse) GetJSON200() *Movie {
	return r.JSON200
}

func (r GetMovieByIdHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r GetMovieByIdHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r GetMovieByIdHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r GetMovieByIdHTTPResponse) GetJSON404() *NotFound {
	return r.JSON404
}

func (r GetMovieByIdHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r GetMovieByIdHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r GetMovieByIdHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r GetMovieByIdHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r GetMovieByIdHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r GetMovieByIdHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type UpdateMovieHTTPResponse429Headers struct {
	RetryAfter *int
}

type UpdateMovieHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *Movie

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON404 *NotFound

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *UpdateMovieHTTPResponse429Headers
}

func (r UpdateMovieHTTPResponse) GetJSON200() *Movie {
	return r.JSON200
}

func (r UpdateMovieHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r UpdateMovieHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r UpdateMovieHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r UpdateMovieHTTPResponse) GetJSON404() *NotFound {
	return r.JSON404
}

func (r UpdateMovieHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r UpdateMovieHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r UpdateMovieHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r UpdateMovieHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r UpdateMovieHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r UpdateMovieHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type GetCollectionHTTPResponse429Headers struct {
	RetryAfter *int
}

type GetCollectionHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *CollectionPage

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *GetCollectionHTTPResponse429Headers
}

func (r GetCollectionHTTPResponse) GetJSON200() *CollectionPage {
	return r.JSON200
}

func (r GetCollectionHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r GetCollectionHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r GetCollectionHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r GetCollectionHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r GetCollectionHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r GetCollectionHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r GetCollectionHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r GetCollectionHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r GetCollectionHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type AddMovieToCollectionHTTPResponse429Headers struct {
	RetryAfter *int
}

type AddMovieToCollectionHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON201 *CollectionItem

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON404 *NotFound

	JSON409 *Conflict

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *AddMovieToCollectionHTTPResponse429Headers
}

func (r AddMovieToCollectionHTTPResponse) GetJSON201() *CollectionItem {
	return r.JSON201
}

func (r AddMovieToCollectionHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r AddMovieToCollectionHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r AddMovieToCollectionHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r AddMovieToCollectionHTTPResponse) GetJSON404() *NotFound {
	return r.JSON404
}

func (r AddMovieToCollectionHTTPResponse) GetJSON409() *Conflict {
	return r.JSON409
}

func (r AddMovieToCollectionHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r AddMovieToCollectionHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r AddMovieToCollectionHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r AddMovieToCollectionHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r AddMovieToCollectionHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r AddMovieToCollectionHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type DeleteMovieFromCollectionHTTPResponse429Headers struct {
	RetryAfter *int
}

type DeleteMovieFromCollectionHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON404 *NotFound

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *DeleteMovieFromCollectionHTTPResponse429Headers
}

func (r DeleteMovieFromCollectionHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r DeleteMovieFromCollectionHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r DeleteMovieFromCollectionHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r DeleteMovieFromCollectionHTTPResponse) GetJSON404() *NotFound {
	return r.JSON404
}

func (r DeleteMovieFromCollectionHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r DeleteMovieFromCollectionHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r DeleteMovieFromCollectionHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r DeleteMovieFromCollectionHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r DeleteMovieFromCollectionHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r DeleteMovieFromCollectionHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type GetMovieFromCollectionHTTPResponse429Headers struct {
	RetryAfter *int
}

type GetMovieFromCollectionHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *CollectionItem

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON404 *NotFound

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *GetMovieFromCollectionHTTPResponse429Headers
}

func (r GetMovieFromCollectionHTTPResponse) GetJSON200() *CollectionItem {
	return r.JSON200
}

func (r GetMovieFromCollectionHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r GetMovieFromCollectionHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r GetMovieFromCollectionHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r GetMovieFromCollectionHTTPResponse) GetJSON404() *NotFound {
	return r.JSON404
}

func (r GetMovieFromCollectionHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r GetMovieFromCollectionHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r GetMovieFromCollectionHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r GetMovieFromCollectionHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r GetMovieFromCollectionHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r GetMovieFromCollectionHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type UpdateCollectionItemHTTPResponse429Headers struct {
	RetryAfter *int
}

type UpdateCollectionItemHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *CollectionItem

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON404 *NotFound

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *UpdateCollectionItemHTTPResponse429Headers
}

func (r UpdateCollectionItemHTTPResponse) GetJSON200() *CollectionItem {
	return r.JSON200
}

func (r UpdateCollectionItemHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r UpdateCollectionItemHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r UpdateCollectionItemHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r UpdateCollectionItemHTTPResponse) GetJSON404() *NotFound {
	return r.JSON404
}

func (r UpdateCollectionItemHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r UpdateCollectionItemHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r UpdateCollectionItemHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r UpdateCollectionItemHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r UpdateCollectionItemHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r UpdateCollectionItemHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type GetUsersHTTPResponse429Headers struct {
	RetryAfter *int
}

type GetUsersHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *UserPage

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *GetUsersHTTPResponse429Headers
}

func (r GetUsersHTTPResponse) GetJSON200() *UserPage {
	return r.JSON200
}

func (r GetUsersHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r GetUsersHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r GetUsersHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r GetUsersHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r GetUsersHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r GetUsersHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r GetUsersHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r GetUsersHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r GetUsersHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type DeleteUserHTTPResponse429Headers struct {
	RetryAfter *int
}

type DeleteUserHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON404 *NotFound

	JSON409 *ErrorResponse

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *DeleteUserHTTPResponse429Headers
}

func (r DeleteUserHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r DeleteUserHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r DeleteUserHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r DeleteUserHTTPResponse) GetJSON404() *NotFound {
	return r.JSON404
}

func (r DeleteUserHTTPResponse) GetJSON409() *ErrorResponse {
	return r.JSON409
}

func (r DeleteUserHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r DeleteUserHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r DeleteUserHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r DeleteUserHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r DeleteUserHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r DeleteUserHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type GetUserByIdHTTPResponse429Headers struct {
	RetryAfter *int
}

type GetUserByIdHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *User

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON404 *NotFound

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *GetUserByIdHTTPResponse429Headers
}

func (r GetUserByIdHTTPResponse) GetJSON200() *User {
	return r.JSON200
}

func (r GetUserByIdHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r GetUserByIdHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r GetUserByIdHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r GetUserByIdHTTPResponse) GetJSON404() *NotFound {
	return r.JSON404
}

func (r GetUserByIdHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r GetUserByIdHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r GetUserByIdHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r GetUserByIdHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r GetUserByIdHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r GetUserByIdHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

type UpdateUserBlockStatusHTTPResponse429Headers struct {
	RetryAfter *int
}

type UpdateUserBlockStatusHTTPResponse struct {
	Body         []byte
	HTTPResponse *http.Response

	JSON200 *User

	JSON400 *BadRequest

	JSON401 *Unauthorized

	JSON403 *Forbidden

	JSON404 *NotFound

	JSON409 *ErrorResponse

	JSON429 *TooManyRequests

	JSON500 *InternalServerError

	Headers429 *UpdateUserBlockStatusHTTPResponse429Headers
}

func (r UpdateUserBlockStatusHTTPResponse) GetJSON200() *User {
	return r.JSON200
}

func (r UpdateUserBlockStatusHTTPResponse) GetJSON400() *BadRequest {
	return r.JSON400
}

func (r UpdateUserBlockStatusHTTPResponse) GetJSON401() *Unauthorized {
	return r.JSON401
}

func (r UpdateUserBlockStatusHTTPResponse) GetJSON403() *Forbidden {
	return r.JSON403
}

func (r UpdateUserBlockStatusHTTPResponse) GetJSON404() *NotFound {
	return r.JSON404
}

func (r UpdateUserBlockStatusHTTPResponse) GetJSON409() *ErrorResponse {
	return r.JSON409
}

func (r UpdateUserBlockStatusHTTPResponse) GetJSON429() *TooManyRequests {
	return r.JSON429
}

func (r UpdateUserBlockStatusHTTPResponse) GetJSON500() *InternalServerError {
	return r.JSON500
}

func (r UpdateUserBlockStatusHTTPResponse) GetBody() []byte {
	return r.Body
}

func (r UpdateUserBlockStatusHTTPResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

func (r UpdateUserBlockStatusHTTPResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

func (r UpdateUserBlockStatusHTTPResponse) ContentType() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Header.Get("Content-Type")
	}
	return ""
}

func (c *ClientWithResponses) LoginWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*LoginHTTPResponse, error) {
	rsp, err := c.LoginWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseLoginHTTPResponse(rsp)
}

func (c *ClientWithResponses) LoginWithResponse(ctx context.Context, body LoginJSONRequestBody, reqEditors ...RequestEditorFn) (*LoginHTTPResponse, error) {
	rsp, err := c.Login(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseLoginHTTPResponse(rsp)
}

func (c *ClientWithResponses) RegisterWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*RegisterHTTPResponse, error) {
	rsp, err := c.RegisterWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRegisterHTTPResponse(rsp)
}

func (c *ClientWithResponses) RegisterWithResponse(ctx context.Context, body RegisterJSONRequestBody, reqEditors ...RequestEditorFn) (*RegisterHTTPResponse, error) {
	rsp, err := c.Register(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseRegisterHTTPResponse(rsp)
}

func (c *ClientWithResponses) GetCatalogWithResponse(ctx context.Context, params *GetCatalogParams, reqEditors ...RequestEditorFn) (*GetCatalogHTTPResponse, error) {
	rsp, err := c.GetCatalog(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCatalogHTTPResponse(rsp)
}

func (c *ClientWithResponses) CreateMovieWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateMovieHTTPResponse, error) {
	rsp, err := c.CreateMovieWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateMovieHTTPResponse(rsp)
}

func (c *ClientWithResponses) CreateMovieWithResponse(ctx context.Context, body CreateMovieJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateMovieHTTPResponse, error) {
	rsp, err := c.CreateMovie(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateMovieHTTPResponse(rsp)
}

func (c *ClientWithResponses) DeleteMovieWithResponse(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*DeleteMovieHTTPResponse, error) {
	rsp, err := c.DeleteMovie(ctx, movieId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteMovieHTTPResponse(rsp)
}

func (c *ClientWithResponses) GetMovieByIdWithResponse(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*GetMovieByIdHTTPResponse, error) {
	rsp, err := c.GetMovieById(ctx, movieId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetMovieByIdHTTPResponse(rsp)
}

func (c *ClientWithResponses) UpdateMovieWithBodyWithResponse(ctx context.Context, movieId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateMovieHTTPResponse, error) {
	rsp, err := c.UpdateMovieWithBody(ctx, movieId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateMovieHTTPResponse(rsp)
}

func (c *ClientWithResponses) UpdateMovieWithResponse(ctx context.Context, movieId openapi_types.UUID, body UpdateMovieJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateMovieHTTPResponse, error) {
	rsp, err := c.UpdateMovie(ctx, movieId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateMovieHTTPResponse(rsp)
}

func (c *ClientWithResponses) GetCollectionWithResponse(ctx context.Context, params *GetCollectionParams, reqEditors ...RequestEditorFn) (*GetCollectionHTTPResponse, error) {
	rsp, err := c.GetCollection(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCollectionHTTPResponse(rsp)
}

func (c *ClientWithResponses) AddMovieToCollectionWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AddMovieToCollectionHTTPResponse, error) {
	rsp, err := c.AddMovieToCollectionWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddMovieToCollectionHTTPResponse(rsp)
}

func (c *ClientWithResponses) AddMovieToCollectionWithResponse(ctx context.Context, body AddMovieToCollectionJSONRequestBody, reqEditors ...RequestEditorFn) (*AddMovieToCollectionHTTPResponse, error) {
	rsp, err := c.AddMovieToCollection(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddMovieToCollectionHTTPResponse(rsp)
}

func (c *ClientWithResponses) DeleteMovieFromCollectionWithResponse(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*DeleteMovieFromCollectionHTTPResponse, error) {
	rsp, err := c.DeleteMovieFromCollection(ctx, movieId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteMovieFromCollectionHTTPResponse(rsp)
}

func (c *ClientWithResponses) GetMovieFromCollectionWithResponse(ctx context.Context, movieId openapi_types.UUID, reqEditors ...RequestEditorFn) (*GetMovieFromCollectionHTTPResponse, error) {
	rsp, err := c.GetMovieFromCollection(ctx, movieId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetMovieFromCollectionHTTPResponse(rsp)
}

func (c *ClientWithResponses) UpdateCollectionItemWithBodyWithResponse(ctx context.Context, movieId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateCollectionItemHTTPResponse, error) {
	rsp, err := c.UpdateCollectionItemWithBody(ctx, movieId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateCollectionItemHTTPResponse(rsp)
}

func (c *ClientWithResponses) UpdateCollectionItemWithResponse(ctx context.Context, movieId openapi_types.UUID, body UpdateCollectionItemJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateCollectionItemHTTPResponse, error) {
	rsp, err := c.UpdateCollectionItem(ctx, movieId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateCollectionItemHTTPResponse(rsp)
}

func (c *ClientWithResponses) GetUsersWithResponse(ctx context.Context, params *GetUsersParams, reqEditors ...RequestEditorFn) (*GetUsersHTTPResponse, error) {
	rsp, err := c.GetUsers(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetUsersHTTPResponse(rsp)
}

func (c *ClientWithResponses) DeleteUserWithResponse(ctx context.Context, userId openapi_types.UUID, reqEditors ...RequestEditorFn) (*DeleteUserHTTPResponse, error) {
	rsp, err := c.DeleteUser(ctx, userId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteUserHTTPResponse(rsp)
}

func (c *ClientWithResponses) GetUserByIdWithResponse(ctx context.Context, userId openapi_types.UUID, reqEditors ...RequestEditorFn) (*GetUserByIdHTTPResponse, error) {
	rsp, err := c.GetUserById(ctx, userId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetUserByIdHTTPResponse(rsp)
}

func (c *ClientWithResponses) UpdateUserBlockStatusWithBodyWithResponse(ctx context.Context, userId openapi_types.UUID, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateUserBlockStatusHTTPResponse, error) {
	rsp, err := c.UpdateUserBlockStatusWithBody(ctx, userId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateUserBlockStatusHTTPResponse(rsp)
}

func (c *ClientWithResponses) UpdateUserBlockStatusWithResponse(ctx context.Context, userId openapi_types.UUID, body UpdateUserBlockStatusJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateUserBlockStatusHTTPResponse, error) {
	rsp, err := c.UpdateUserBlockStatus(ctx, userId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateUserBlockStatusHTTPResponse(rsp)
}

func ParseLoginHTTPResponse(rsp *http.Response) (*LoginHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &LoginHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest LoginResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers LoginHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseRegisterHTTPResponse(rsp *http.Response) (*RegisterHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &RegisterHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest User
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest Conflict
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON409 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers RegisterHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseGetCatalogHTTPResponse(rsp *http.Response) (*GetCatalogHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetCatalogHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest MoviePage
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers GetCatalogHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseCreateMovieHTTPResponse(rsp *http.Response) (*CreateMovieHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateMovieHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest Movie
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers CreateMovieHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseDeleteMovieHTTPResponse(rsp *http.Response) (*DeleteMovieHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteMovieHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case rsp.StatusCode == 204:
		break

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers DeleteMovieHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseGetMovieByIdHTTPResponse(rsp *http.Response) (*GetMovieByIdHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetMovieByIdHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Movie
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers GetMovieByIdHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseUpdateMovieHTTPResponse(rsp *http.Response) (*UpdateMovieHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateMovieHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Movie
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers UpdateMovieHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseGetCollectionHTTPResponse(rsp *http.Response) (*GetCollectionHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetCollectionHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest CollectionPage
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers GetCollectionHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseAddMovieToCollectionHTTPResponse(rsp *http.Response) (*AddMovieToCollectionHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AddMovieToCollectionHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest CollectionItem
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest Conflict
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON409 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers AddMovieToCollectionHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseDeleteMovieFromCollectionHTTPResponse(rsp *http.Response) (*DeleteMovieFromCollectionHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteMovieFromCollectionHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case rsp.StatusCode == 204:
		break

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers DeleteMovieFromCollectionHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseGetMovieFromCollectionHTTPResponse(rsp *http.Response) (*GetMovieFromCollectionHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetMovieFromCollectionHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest CollectionItem
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers GetMovieFromCollectionHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseUpdateCollectionItemHTTPResponse(rsp *http.Response) (*UpdateCollectionItemHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateCollectionItemHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest CollectionItem
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers UpdateCollectionItemHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseGetUsersHTTPResponse(rsp *http.Response) (*GetUsersHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetUsersHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest UserPage
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers GetUsersHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseDeleteUserHTTPResponse(rsp *http.Response) (*DeleteUserHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteUserHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case rsp.StatusCode == 204:
		break

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON409 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers DeleteUserHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseGetUserByIdHTTPResponse(rsp *http.Response) (*GetUserByIdHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetUserByIdHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest User
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers GetUserByIdHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}

func ParseUpdateUserBlockStatusHTTPResponse(rsp *http.Response) (*UpdateUserBlockStatusHTTPResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateUserBlockStatusHTTPResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest User
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest BadRequest
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Unauthorized
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest Forbidden
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest NotFound
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 409:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON409 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 429:
		var dest TooManyRequests
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON429 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest InternalServerError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	switch {
	case rsp.StatusCode == 429:
		var headers UpdateUserBlockStatusHTTPResponse429Headers
		if values := rsp.Header.Values("Retry-After"); len(values) > 0 {
			var value int
			if err := runtime.BindStyledParameterWithOptions("simple", "Retry-After", values[0], &value, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false, Type: "integer", Format: ""}); err != nil {
				return nil, err
			}
			headers.RetryAfter = &value
		}
		response.Headers429 = &headers
	}

	return response, nil
}
