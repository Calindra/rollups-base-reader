// Package transaction provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package transaction

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// Cartesi712 defines model for Cartesi712.
type Cartesi712 struct {
	Address   *string `json:"address,omitempty"`
	Signature *string `json:"signature,omitempty"`
	TypedData *struct {
		Account *string `json:"account,omitempty"`
		Domain  struct {
			ChainId           *int    `json:"chainId,omitempty"`
			Name              *string `json:"name,omitempty"`
			VerifyingContract *string `json:"verifyingContract,omitempty"`
			Version           *string `json:"version,omitempty"`
		} `json:"domain"`
		Message struct {
			App         string                                   `json:"app"`
			Data        string                                   `json:"data"`
			MaxGasPrice Cartesi712_TypedData_Message_MaxGasPrice `json:"max_gas_price"`
			Nonce       Cartesi712_TypedData_Message_Nonce       `json:"nonce"`
		} `json:"message"`
		PrimaryType string `json:"primaryType"`
		Types       struct {
			CartesiMessage *[]struct {
				Name *string `json:"name,omitempty"`
				Type *string `json:"type,omitempty"`
			} `json:"CartesiMessage,omitempty"`
			EIP712Domain *[]struct {
				Name *string `json:"name,omitempty"`
				Type *string `json:"type,omitempty"`
			} `json:"EIP712Domain,omitempty"`
		} `json:"types"`
	} `json:"typedData,omitempty"`
}

// Cartesi712TypedDataMessageMaxGasPrice0 defines model for .
type Cartesi712TypedDataMessageMaxGasPrice0 = string

// Cartesi712TypedDataMessageMaxGasPrice1 defines model for .
type Cartesi712TypedDataMessageMaxGasPrice1 = int

// Cartesi712_TypedData_Message_MaxGasPrice defines model for Cartesi712.TypedData.Message.MaxGasPrice.
type Cartesi712_TypedData_Message_MaxGasPrice struct {
	union json.RawMessage
}

// Cartesi712TypedDataMessageNonce0 defines model for .
type Cartesi712TypedDataMessageNonce0 = string

// Cartesi712TypedDataMessageNonce1 defines model for .
type Cartesi712TypedDataMessageNonce1 = uint64

// Cartesi712_TypedData_Message_Nonce defines model for Cartesi712.TypedData.Message.Nonce.
type Cartesi712_TypedData_Message_Nonce struct {
	union json.RawMessage
}

// GetNonce defines model for GetNonce.
type GetNonce struct {
	// AppContract App contract address
	AppContract string `json:"app_contract"`

	// MsgSender Message sender address
	MsgSender string `json:"msg_sender"`
}

// NonceResponse defines model for NonceResponse.
type NonceResponse struct {
	// Nonce Nonce number
	Nonce *int `json:"nonce,omitempty"`
}

// TransactionError defines model for TransactionError.
type TransactionError struct {
	// Message Detailed error message
	Message *string `json:"message,omitempty"`
}

// TransactionResponse defines model for TransactionResponse.
type TransactionResponse struct {
	// Id tx number
	Id *string `json:"id,omitempty"`
}

// GetNonceJSONRequestBody defines body for GetNonce for application/json ContentType.
type GetNonceJSONRequestBody = GetNonce

// SendCartesiTransactionJSONRequestBody defines body for SendCartesiTransaction for application/json ContentType.
type SendCartesiTransactionJSONRequestBody = Cartesi712

// AsCartesi712TypedDataMessageMaxGasPrice0 returns the union data inside the Cartesi712_TypedData_Message_MaxGasPrice as a Cartesi712TypedDataMessageMaxGasPrice0
func (t Cartesi712_TypedData_Message_MaxGasPrice) AsCartesi712TypedDataMessageMaxGasPrice0() (Cartesi712TypedDataMessageMaxGasPrice0, error) {
	var body Cartesi712TypedDataMessageMaxGasPrice0
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromCartesi712TypedDataMessageMaxGasPrice0 overwrites any union data inside the Cartesi712_TypedData_Message_MaxGasPrice as the provided Cartesi712TypedDataMessageMaxGasPrice0
func (t *Cartesi712_TypedData_Message_MaxGasPrice) FromCartesi712TypedDataMessageMaxGasPrice0(v Cartesi712TypedDataMessageMaxGasPrice0) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeCartesi712TypedDataMessageMaxGasPrice0 performs a merge with any union data inside the Cartesi712_TypedData_Message_MaxGasPrice, using the provided Cartesi712TypedDataMessageMaxGasPrice0
func (t *Cartesi712_TypedData_Message_MaxGasPrice) MergeCartesi712TypedDataMessageMaxGasPrice0(v Cartesi712TypedDataMessageMaxGasPrice0) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsCartesi712TypedDataMessageMaxGasPrice1 returns the union data inside the Cartesi712_TypedData_Message_MaxGasPrice as a Cartesi712TypedDataMessageMaxGasPrice1
func (t Cartesi712_TypedData_Message_MaxGasPrice) AsCartesi712TypedDataMessageMaxGasPrice1() (Cartesi712TypedDataMessageMaxGasPrice1, error) {
	var body Cartesi712TypedDataMessageMaxGasPrice1
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromCartesi712TypedDataMessageMaxGasPrice1 overwrites any union data inside the Cartesi712_TypedData_Message_MaxGasPrice as the provided Cartesi712TypedDataMessageMaxGasPrice1
func (t *Cartesi712_TypedData_Message_MaxGasPrice) FromCartesi712TypedDataMessageMaxGasPrice1(v Cartesi712TypedDataMessageMaxGasPrice1) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeCartesi712TypedDataMessageMaxGasPrice1 performs a merge with any union data inside the Cartesi712_TypedData_Message_MaxGasPrice, using the provided Cartesi712TypedDataMessageMaxGasPrice1
func (t *Cartesi712_TypedData_Message_MaxGasPrice) MergeCartesi712TypedDataMessageMaxGasPrice1(v Cartesi712TypedDataMessageMaxGasPrice1) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t Cartesi712_TypedData_Message_MaxGasPrice) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *Cartesi712_TypedData_Message_MaxGasPrice) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsCartesi712TypedDataMessageNonce0 returns the union data inside the Cartesi712_TypedData_Message_Nonce as a Cartesi712TypedDataMessageNonce0
func (t Cartesi712_TypedData_Message_Nonce) AsCartesi712TypedDataMessageNonce0() (Cartesi712TypedDataMessageNonce0, error) {
	var body Cartesi712TypedDataMessageNonce0
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromCartesi712TypedDataMessageNonce0 overwrites any union data inside the Cartesi712_TypedData_Message_Nonce as the provided Cartesi712TypedDataMessageNonce0
func (t *Cartesi712_TypedData_Message_Nonce) FromCartesi712TypedDataMessageNonce0(v Cartesi712TypedDataMessageNonce0) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeCartesi712TypedDataMessageNonce0 performs a merge with any union data inside the Cartesi712_TypedData_Message_Nonce, using the provided Cartesi712TypedDataMessageNonce0
func (t *Cartesi712_TypedData_Message_Nonce) MergeCartesi712TypedDataMessageNonce0(v Cartesi712TypedDataMessageNonce0) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

// AsCartesi712TypedDataMessageNonce1 returns the union data inside the Cartesi712_TypedData_Message_Nonce as a Cartesi712TypedDataMessageNonce1
func (t Cartesi712_TypedData_Message_Nonce) AsCartesi712TypedDataMessageNonce1() (Cartesi712TypedDataMessageNonce1, error) {
	var body Cartesi712TypedDataMessageNonce1
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromCartesi712TypedDataMessageNonce1 overwrites any union data inside the Cartesi712_TypedData_Message_Nonce as the provided Cartesi712TypedDataMessageNonce1
func (t *Cartesi712_TypedData_Message_Nonce) FromCartesi712TypedDataMessageNonce1(v Cartesi712TypedDataMessageNonce1) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeCartesi712TypedDataMessageNonce1 performs a merge with any union data inside the Cartesi712_TypedData_Message_Nonce, using the provided Cartesi712TypedDataMessageNonce1
func (t *Cartesi712_TypedData_Message_Nonce) MergeCartesi712TypedDataMessageNonce1(v Cartesi712TypedDataMessageNonce1) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t Cartesi712_TypedData_Message_Nonce) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *Cartesi712_TypedData_Message_Nonce) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetNonceWithBody request with any body
	GetNonceWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	GetNonce(ctx context.Context, body GetNonceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SendCartesiTransactionWithBody request with any body
	SendCartesiTransactionWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SendCartesiTransaction(ctx context.Context, body SendCartesiTransactionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetNonceWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetNonceRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetNonce(ctx context.Context, body GetNonceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetNonceRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SendCartesiTransactionWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSendCartesiTransactionRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SendCartesiTransaction(ctx context.Context, body SendCartesiTransactionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSendCartesiTransactionRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetNonceRequest calls the generic GetNonce builder with application/json body
func NewGetNonceRequest(server string, body GetNonceJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewGetNonceRequestWithBody(server, "application/json", bodyReader)
}

// NewGetNonceRequestWithBody generates requests for GetNonce with any type of body
func NewGetNonceRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/nonce")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewSendCartesiTransactionRequest calls the generic SendCartesiTransaction builder with application/json body
func NewSendCartesiTransactionRequest(server string, body SendCartesiTransactionJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewSendCartesiTransactionRequestWithBody(server, "application/json", bodyReader)
}

// NewSendCartesiTransactionRequestWithBody generates requests for SendCartesiTransaction with any type of body
func NewSendCartesiTransactionRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/submit")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
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

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
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

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetNonceWithBodyWithResponse request with any body
	GetNonceWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetNonceResponse, error)

	GetNonceWithResponse(ctx context.Context, body GetNonceJSONRequestBody, reqEditors ...RequestEditorFn) (*GetNonceResponse, error)

	// SendCartesiTransactionWithBodyWithResponse request with any body
	SendCartesiTransactionWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SendCartesiTransactionResponse, error)

	SendCartesiTransactionWithResponse(ctx context.Context, body SendCartesiTransactionJSONRequestBody, reqEditors ...RequestEditorFn) (*SendCartesiTransactionResponse, error)
}

type GetNonceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *NonceResponse
}

// Status returns HTTPResponse.Status
func (r GetNonceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetNonceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SendCartesiTransactionResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *TransactionResponse
	JSON400      *TransactionError
}

// Status returns HTTPResponse.Status
func (r SendCartesiTransactionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SendCartesiTransactionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetNonceWithBodyWithResponse request with arbitrary body returning *GetNonceResponse
func (c *ClientWithResponses) GetNonceWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*GetNonceResponse, error) {
	rsp, err := c.GetNonceWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetNonceResponse(rsp)
}

func (c *ClientWithResponses) GetNonceWithResponse(ctx context.Context, body GetNonceJSONRequestBody, reqEditors ...RequestEditorFn) (*GetNonceResponse, error) {
	rsp, err := c.GetNonce(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetNonceResponse(rsp)
}

// SendCartesiTransactionWithBodyWithResponse request with arbitrary body returning *SendCartesiTransactionResponse
func (c *ClientWithResponses) SendCartesiTransactionWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SendCartesiTransactionResponse, error) {
	rsp, err := c.SendCartesiTransactionWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSendCartesiTransactionResponse(rsp)
}

func (c *ClientWithResponses) SendCartesiTransactionWithResponse(ctx context.Context, body SendCartesiTransactionJSONRequestBody, reqEditors ...RequestEditorFn) (*SendCartesiTransactionResponse, error) {
	rsp, err := c.SendCartesiTransaction(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSendCartesiTransactionResponse(rsp)
}

// ParseGetNonceResponse parses an HTTP response from a GetNonceWithResponse call
func ParseGetNonceResponse(rsp *http.Response) (*GetNonceResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetNonceResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest NonceResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseSendCartesiTransactionResponse parses an HTTP response from a SendCartesiTransactionWithResponse call
func ParseSendCartesiTransactionResponse(rsp *http.Response) (*SendCartesiTransactionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SendCartesiTransactionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest TransactionResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest TransactionError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get Nonce
	// (POST /nonce)
	GetNonce(ctx echo.Context) error

	// (POST /submit)
	SendCartesiTransaction(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetNonce converts echo context to params.
func (w *ServerInterfaceWrapper) GetNonce(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetNonce(ctx)
	return err
}

// SendCartesiTransaction converts echo context to params.
func (w *ServerInterfaceWrapper) SendCartesiTransaction(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SendCartesiTransaction(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/nonce", wrapper.GetNonce)
	router.POST(baseURL+"/submit", wrapper.SendCartesiTransaction)

}
