package paddle

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Modifiers is an API to work with the Paddle subscription modifiers.
type Modifiers api

// CreateModifierOptions represents options for creating a modifier.
type CreateModifierOptions struct {
	SubscriptionID      uint64
	ModifierRecurring   bool
	ModifierAmount      string
	ModifierDescription string
}

// encodeURLValues encodes options as URL parameters.
func (options *CreateModifierOptions) encodeURLValues() (url.Values, error) {
	values := make(url.Values)
	if options.SubscriptionID == 0 {
		return nil, errors.New("\"subscription_id\" is required")
	}
	if options.ModifierAmount == "" {
		return nil, errors.New("\"modifier_amount\" is required")
	}

	values.Set("subscription_id", strconv.FormatUint(options.SubscriptionID, 10))
	values.Set("modifier_amount", options.ModifierAmount)
	values.Set("modifier_recurring", strconv.FormatBool(options.ModifierRecurring))
	values.Set("modifier_description", options.ModifierDescription)

	return values, nil
}

// CreateModifierResponse represents a response for the create subscription modifier.
type CreateModifierResponse struct {
	SubscriptionID uint64 `json:"subscription_id,omitempty"`
	ModifierID     uint64 `json:"modifier_id,omitempty"`
}

// Create creates new modifier.
//
// Paddle docs: https://developer.paddle.com/api-reference/dc2b0c06f0481-create-modifier
func (modifiers *Modifiers) Create(ctx context.Context, options *CreateModifierOptions) (*CreateModifierResponse, *http.Response, error) {
	path := "2.0/subscription/modifiers/create"

	if options == nil {
		options = new(CreateModifierOptions)
	}
	request, err := newRequest(ctx, http.MethodPost, modifiers.baseURL, path, modifiers.authentication, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new request: %w", err)
	}

	response := new(response[*CreateModifierResponse])
	httpResponse, err := doRequest(modifiers.httpClient, request, response)
	if err != nil {
		return nil, httpResponse, err
	}

	return *response.Response, httpResponse, nil
}
