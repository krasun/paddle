package paddle

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	productionBaseURL = "https://vendors.paddle.com/api/"
	sandboxBaseURL    = "https://sandbox-vendors.paddle.com/api/"

	vendorID       = "vendor_id"
	vendorAuthCode = "vendor_auth_code"
)

// api is type shared by different Paddle API sections implementations, like
// Webhooks, Payments, etc.
// It is defined for the usability of the exposed Paddle client.
type api struct {
	authentication *Authentication
	baseURL        *url.URL
	httpClient     *http.Client
}

// Client is a Paddle client.
type Client struct {
	// Users represents an API for working with subscription users.
	Users *Users
	// Modifiers represents an API for working with subscription modifiers.
	Modifiers *Modifiers
	// Charges represents an API for working with subscription charges.
	Charges *Charges
}

// Authentication represents credentials for working with the Paddle API.
type Authentication struct {
	// VendorID identifies the seller account, can be found in Developer Tools > Authentication.
	VendorID int
	// VendorAuthCode is a private API key for authenticating API requests, should never be used in client side code or shared publicly.
	// And can be found in Developer Tools > Authentication.
	VendorAuthCode string
}

// NewProductionClient creates a new Paddle production client.
func NewProductionClient(authentication Authentication) (*Client, error) {
	baseURL, err := url.Parse(productionBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL %s: %w", productionBaseURL, err)
	}

	return newClient(&authentication, baseURL), nil
}

// NewSandboxClient creates a new Paddle sandbox client.
func NewSandboxClient(authentication Authentication) (*Client, error) {
	baseURL, err := url.Parse(sandboxBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL %s: %w", sandboxBaseURL, err)
	}

	return newClient(&authentication, baseURL), nil
}

// newClient instantiates a new Paddle client with the specified base URL.
func newClient(authentication *Authentication, baseURL *url.URL) *Client {
	httpClient := http.DefaultClient

	return &Client{
		&Users{httpClient: httpClient, baseURL: baseURL, authentication: authentication},
		&Modifiers{httpClient: httpClient, baseURL: baseURL, authentication: authentication},
		&Charges{httpClient: httpClient, baseURL: baseURL, authentication: authentication},
	}
}

// prepareURL copies base URL with a new path parameter.
func prepareURL(baseURL *url.URL, path string) (*url.URL, error) {
	u, err := baseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	return u, nil
}

// newRequest prepare a new HTTP request instance.
func newRequest(ctx context.Context, method string, baseURL *url.URL, path string, authentication *Authentication, options urlValuesEncoder) (*http.Request, error) {
	u, err := prepareURL(baseURL, path)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare a request URL: %w", err)
	}

	var values url.Values
	if v, err := options.encodeURLValues(); err != nil {
		return nil, fmt.Errorf("failed to encode URL values: %w", err)
	} else if v != nil {
		values = v
	} else if v == nil {
		values = make(url.Values)
	}

	if authentication != nil {
		values.Set(vendorID, strconv.Itoa(authentication.VendorID))
		values.Set(vendorAuthCode, authentication.VendorAuthCode)
	}

	body := strings.NewReader(values.Encode())
	request, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate new request: %w", err)
	}

	if body.Size() > 0 {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return request, nil
}

// urlValuesEncoder encodes URL values.
type urlValuesEncoder interface {
	encodeURLValues() (url.Values, error)
}

// APIError represents a Paddle API error.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error formats the error as a string.
func (e *APIError) Error() string {
	return fmt.Sprintf("Paddle API error: code=%d, message=%s", e.Code, e.Message)
}

// response represents a deserialized response from the Paddle API.
type response[T any] struct {
	Success  bool      `json:"bool"`
	Error    *APIError `json:"error"`
	Response *T        `json:"response"`
}

// doRequest executes an HTTP request, decodes response and returns both decoded and HTTP responses.
func doRequest[T any](httpClient *http.Client, request *http.Request, paddleResponse *response[T]) (*http.Response, error) {
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the request: %w", err)
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(data, paddleResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if paddleResponse.Error != nil {
		return response, paddleResponse.Error
	}

	return response, nil
}
