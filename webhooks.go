package paddle

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"time"

	"github.com/gorilla/schema"
)

// SubscriptionStatus represents subscription status: active, trialing, past due, paused or deleted.
type SubscriptionStatus string

const (
	subscriptionUnknown SubscriptionStatus = "unknown"
	// SubscriptionActive represents active subscription status.
	SubscriptionActive SubscriptionStatus = "active"
	// SubscriptionTrialing represents trialing subscription status.
	SubscriptionTrialing SubscriptionStatus = "trialing"
	// SubscriptionPastDue represents past due subscription status.
	SubscriptionPastDue SubscriptionStatus = "past_due"
	// SubscriptionPaused represents paused subscription status.
	SubscriptionPaused SubscriptionStatus = "paused"
	// SubscriptionDeleted represents deleted subscription status.
	SubscriptionDeleted SubscriptionStatus = "deleted"
)

// RefundType represents refund type: full, vat or partial.
type RefundType string

const (
	refundUnknown RefundType = "unknown"
	// RefundFull represents full refund.
	RefundFull RefundType = "full"
	// RefundVAT represents VAT refund.
	RefundVAT RefundType = "vat"
	// RefundPartial represents partial refund.
	RefundPartial RefundType = "partial"
)

// Webhooks validates and parses webhook alerts.
type Webhooks struct {
	signKey *rsa.PublicKey
	decoder *schema.Decoder
}

// NewWebhooks returns a new instance of the webhooks.
func NewWebhooks(publicKey []byte) (*Webhooks, error) {
	pemBlock, _ := pem.Decode(publicKey)
	if pemBlock == nil {
		return nil, errors.New("failed to locate public key PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	signKey, isRSAPublicKey := pub.(*rsa.PublicKey)
	if !isRSAPublicKey {
		return nil, fmt.Errorf("invalid key format, expected RSA public key")
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.ZeroEmpty(true)
	decoder.RegisterConverter(time.Time{}, convertTime)
	decoder.RegisterConverter(subscriptionUnknown, convertSubscriptionStatus)
	decoder.RegisterConverter(refundUnknown, convertRefundType)

	return &Webhooks{signKey: signKey, decoder: decoder}, nil
}

func convertTime(value string) reflect.Value {
	if v, err := time.Parse("2006-01-02", value); err == nil {
		return reflect.ValueOf(v)
	}

	if v, err := time.Parse("2006-01-02 15:04:05", value); err == nil {
		return reflect.ValueOf(v)
	}

	return reflect.Value{}
}

func convertSubscriptionStatus(value string) reflect.Value {
	switch value {
	case "active":
		return reflect.ValueOf(SubscriptionActive)
	case "trialing":
		return reflect.ValueOf(SubscriptionTrialing)
	case "past_due":
		return reflect.ValueOf(SubscriptionPastDue)
	case "paused":
		return reflect.ValueOf(SubscriptionPaused)
	case "deleted":
		return reflect.ValueOf(SubscriptionDeleted)
	}

	return reflect.Value{}
}

func convertRefundType(value string) reflect.Value {
	switch value {
	case "full":
		return reflect.ValueOf(RefundFull)
	case "vat":
		return reflect.ValueOf(RefundVAT)
	case "partial":
		return reflect.ValueOf(RefundPartial)
	}

	return reflect.Value{}
}

// ParseRequest validates the Paddle webhook request and returns typed alert in case of success,
// otherwise it returns an error.
func (webhooks *Webhooks) ParseRequest(r *http.Request) (interface{}, error) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/x-www-form-urlencoded" {
		return nil, fmt.Errorf("webhook request has unsupported \"Content-Type\": %s", contentType)
	}

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	signature, err := base64.StdEncoding.DecodeString(r.Form.Get("p_signature"))
	if err != nil {
		fmt.Println("signature")
		fmt.Println(r.Form.Get("p_signature"))
		fmt.Println("signature")
		

		return nil, fmt.Errorf("failed to decode the signature: %w", err)
	}

	var keys []string
	for key := range r.Form {
		if key != "p_signature" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	serialized := fmt.Sprintf("a:%d:{", len(keys))
	for _, k := range keys {
		serialized += fmt.Sprintf("s:%d:\"%s\";s:%d:\"%s\";", len(k), k, len(r.Form.Get(k)), r.Form.Get(k))
	}
	serialized += "}"

	checksum := sha1.Sum([]byte(serialized))

	err = rsa.VerifyPKCS1v15(webhooks.signKey, crypto.SHA1, checksum[:], signature)
	if err != nil {
		return nil, fmt.Errorf("failed to verify the signature: %w", err)
	}

	var alert interface{}
	alertName := r.Form.Get("alert_name")

	switch alertName {
	case "subscription_created":
		alert = &SubscriptionCreatedAlert{}
	case "subscription_updated":
		alert = &SubscriptionUpdatedAlert{}
	case "subscription_cancelled":
		alert = &SubscriptionCancelledAlert{}
	case "subscription_payment_succeeded":
		alert = &SubscriptionPaymentSucceededAlert{}
	case "subscription_payment_failed":
		alert = &SubscriptionPaymentFailedAlert{}
	case "subscription_payment_refunded":
		alert = &SubscriptionPaymentRefundedAlert{}
	case "payment_succeeded":
	case "payment_refunded":
	case "locker_processed":
	case "payment_dispute_created":
	case "payment_dispute_closed":
	case "high_risk_transaction_created":
	case "high_risk_transaction_updated":
	case "transfer_created":
	case "transfer_paid":
	case "new_audience_member":
	case "update_audience_member":
	case "invoice_paid":
	case "invoice_sent":
	case "invoice_overdue":
		return nil, fmt.Errorf("not implemented \"alert_name\": %v", alertName)
	default:
		return nil, fmt.Errorf("unknown \"alert_name\": %v", alertName)
	}

	err = webhooks.decoder.Decode(alert, r.Form)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the form values: %w", err)
	}

	return alert, nil
}
