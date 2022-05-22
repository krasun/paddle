package paddle

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Users is an API to work with the Paddle subscription users.
type Users api

// ListUsersOptions represents options for getting subscription users.
type ListUsersOptions struct {
	SubscriptionID uint64
	PlanID         uint64
	State          string
	Page           int
	ResultsPerPage int
}

// encodeURLValues encodes users list options as URL parameters.
func (options *ListUsersOptions) encodeURLValues() (url.Values, error) {
	values := make(url.Values)
	if options.SubscriptionID != 0 {
		values.Set("subscription_id", strconv.FormatUint(options.SubscriptionID, 10))
	}
	if options.PlanID != 0 {
		values.Set("plan_id", strconv.FormatUint(options.PlanID, 10))
	}

	switch options.State {
	case "active":
	case "past_due":
	case "trialing":
	case "paused":
	case "deleted":
		values.Set("state", options.State)
	case "":
		break
	default:
		return nil, fmt.Errorf("\"status\" must be empty or one of \"active\", \"past_due\", \"trialing\", \"paused\", \"deleted\"")
	}

	if options.Page > 0 {
		values.Set("page", strconv.Itoa(options.Page))
	}
	if options.ResultsPerPage > 0 {
		values.Set("results_per_page", strconv.Itoa(options.ResultsPerPage))
	}

	return url.Values{}, nil
}

// User represents a Paddle user.
type User struct {
	SubscriptionID     int                 `json:"subscription_id,omitempty"`
	PlanID             int                 `json:"plan_id,omitempty"`
	UserID             int                 `json:"user_id,omitempty"`
	UserEmail          string              `json:"user_email,omitempty"`
	MarketingConsent   bool                `json:"marketing_consent,omitempty"`
	UpdateURL          string              `json:"update_url,omitempty"`
	CancelURL          string              `json:"cancel_url,omitempty"`
	State              string              `json:"state,omitempty"`
	SignupDate         string              `json:"signup_date,omitempty"`
	LastPayment        *UserPayment        `json:"last_payment,omitempty"`
	NextPayment        *UserPayment        `json:"next_payment,omitempty"`
	PaymentInformation *PaymentInformation `json:"payment_information,omitempty"`
	PausedAt           string              `json:"paused_at,omitempty"`
	PausedFrom         string              `json:"paused_from,omitempty"`
}

// UserPayment represents a user payment.
type UserPayment struct {
	Amount   float64 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`
	Date     string  `json:"date,omitempty"`
}

// PaymentInformation represents a user payment information.
type PaymentInformation struct {
	PaymentMethod  string `json:"payment_method,omitempty"`
	CardType       string `json:"card_type,omitempty"`
	LastFourDigits string `json:"last_four_digits,omitempty"`
	ExpiryDate     string `json:"expiry_date,omitempty"`
}

// List returns subscription users.
//
// Paddle docs: https://developer.paddle.com/api-reference/b3A6MzA3NDQ3MzA-list-users
func (users *Users) List(ctx context.Context, options *ListUsersOptions) ([]*User, *http.Response, error) {
	path := "2.0/subscription/users"

	if options == nil {
		options = new(ListUsersOptions)
	}
	request, err := newRequest(ctx, http.MethodPost, users.baseURL, path, users.authentication, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new request: %w", err)
	}

	response := new(response[[]*User])
	httpResponse, err := doRequest(users.httpClient, request, response)
	if err != nil {
		return nil, httpResponse, err
	}

	return *response.Response, httpResponse, nil
}

// UpdateUsersOptions represents options for update user subscription.
type UpdateUserOptions struct {
	SubscriptionID  uint64
	PlanID          uint64
	Prorate         bool
	BillImmediately bool
}

// encodeURLValues encodes options as URL parameters.
func (options *UpdateUserOptions) encodeURLValues() (url.Values, error) {
	values := make(url.Values)
	if options.SubscriptionID != 0 {
		values.Set("subscription_id", strconv.FormatUint(options.SubscriptionID, 10))
	} else {
		return nil, errors.New("\"subscription_id\" is required")
	}

	if options.PlanID != 0 {
		values.Set("plan_id", strconv.FormatUint(options.PlanID, 10))
	}
	values.Set("prorate", strconv.FormatBool(options.Prorate))
	values.Set("bill_immediately", strconv.FormatBool(options.BillImmediately))

	return values, nil
}

// UpdateUserResponse represents a response for the update user subscription request.
type UpdateUserResponse struct {
	SubscriptionID uint64       `json:"subscription_id,omitempty"`
	PlanID         uint64       `json:"plan_id,omitempty"`
	UserID         uint64       `json:"user_id,omitempty"`
	NextPayment    *UserPayment `json:"next_payment,omitempty"`
}

// Update updates user subscription.
//
// Paddle docs: https://developer.paddle.com/api-reference/b3A6MzA3NDQ3MzQ-update-user
func (users *Users) Update(ctx context.Context, options *UpdateUserOptions) (*UpdateUserResponse, *http.Response, error) {
	path := "2.0/subscription/users/update"

	if options == nil {
		options = new(UpdateUserOptions)
	}
	request, err := newRequest(ctx, http.MethodPost, users.baseURL, path, users.authentication, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new request: %w", err)
	}

	response := new(response[*UpdateUserResponse])
	httpResponse, err := doRequest(users.httpClient, request, response)
	if err != nil {
		return nil, httpResponse, err
	}

	return *response.Response, httpResponse, nil
}

// CancelUserOptions represents options for cancel user subscription.
type CancelUserOptions struct {
	SubscriptionID uint64
}

// encodeURLValues encodes options as URL parameters.
func (options *CancelUserOptions) encodeURLValues() (url.Values, error) {
	values := make(url.Values)
	if options.SubscriptionID != 0 {
		values.Set("subscription_id", strconv.FormatUint(options.SubscriptionID, 10))
	} else {
		return nil, errors.New("\"subscription_id\" is required")
	}

	return values, nil
}

// Cancel cancel the user subscription.
//
// Paddle docs: https://developer.paddle.com/api-reference/b3A6MzA3NDQ3MzU-cancel-user
func (users *Users) Cancel(ctx context.Context, options *CancelUserOptions) (*http.Response, error) {
	path := "2.0/subscription/users/cancel"

	if options == nil {
		options = new(CancelUserOptions)
	}
	request, err := newRequest(ctx, http.MethodPost, users.baseURL, path, users.authentication, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	httpResponse, err := doRequest(users.httpClient, request, new(response[interface{}]))
	if err != nil {
		return httpResponse, err
	}

	return httpResponse, nil
}
