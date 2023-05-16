package paddle

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Charges is an API to work with the Paddle subscription one-off charges.
type Charges api

// ChargeOptions represents options for charing.
type ChargeOptions struct {
	Amount     string
	ChargeName string
}

// encodeURLValues encodes options as URL parameters.
func (options *ChargeOptions) encodeURLValues() (url.Values, error) {
	values := make(url.Values)

	values.Set("amount", options.Amount)
	values.Set("charge_name", options.ChargeName)

	return values, nil
}

// ChargeResponse represents a response for the subscription charge.
type ChargeResponse struct {
	InvoiceID      uint64 `json:"invoice_id,omitempty"`
	SubscriptionID uint64 `json:"subscription_id,omitempty"`
	Amount         string `json:"amount,omitempty"`
	Currency       string `json:"currency,omitempty"`
	PaymentDate    string `json:"payment_date,omitempty"`
	ReceiptURL     string `json:"receipt_url,omitempty"`
	Status         string `json:"status,omitempty"`
}

// Charge charges.
//
// Paddle docs: https://developer.paddle.com/api-reference/23cf86225523f-create-one-off-charge
func (modifiers *Charges) Charge(ctx context.Context, subscriptionID uint64, options *ChargeOptions) (*ChargeResponse, *http.Response, error) {
	if subscriptionID == 0 {
		return nil, nil, errors.New("\"subscription_id\" can't be zero")
	}

	path := fmt.Sprintf("2.0/subscription/%d/charge", subscriptionID)

	if options == nil {
		options = new(ChargeOptions)
	}
	request, err := newRequest(ctx, http.MethodPost, modifiers.baseURL, path, modifiers.authentication, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new request: %w", err)
	}

	response := new(response[*ChargeResponse])
	httpResponse, err := doRequest(modifiers.httpClient, request, response)
	if err != nil {
		return nil, httpResponse, err
	}

	return *response.Response, httpResponse, nil
}
