package paddle

import "time"

// SubscriptionPaymentSucceededAlert is fired when a subscription payment is received successfully.
// Docs: https://developer.paddle.com/webhook-reference/subscription-alerts/subscription-payment-succeeded
type SubscriptionPaymentSucceededAlert struct {
	AlertName             string    `schema:"alert_name"`
	AlertID               uint64    `schema:"alert_id"`
	BalanceCurrency       *string   `schema:"balance_currency"`
	BalanceEarnings       *string   `schema:"balance_earnings"`
	BalanceFee            *string   `schema:"balance_fee"`
	BalanceGross          *string   `schema:"balance_gross"`
	BalanceTax            *string   `schema:"balance_tax"`
	CheckoutID            *string   `schema:"checkout_id"`
	Country               *string   `schema:"country"`
	Coupon                *string   `schema:"coupon"`
	Currency              *string   `schema:"currency"`
	CustomerName          *string   `schema:"customer_name"`
	Earnings              *string   `schema:"earnings"`
	Email                 *string   `schema:"email"`
	EventTime             time.Time `schema:"event_time"`
	Fee                   *string   `schema:"fee"`
	InitialPayment        bool      `schema:"initial_payment"`
	Instalments           int       `schema:"instalments"`
	MarketingConsent      bool      `schema:"marketing_consent"`
	NextBillDate          *string   `schema:"next_bill_date"`
	NextPaymentAmount     *string   `schema:"next_payment_amount"`
	OrderID               *string   `schema:"order_id"`
	Passthrough           *string   `schema:"passthrough"`
	PaymentMethod         *string   `schema:"payment_method"`
	PaymentTax            *string   `schema:"payment_tax"`
	PlanName              *string   `schema:"plan_name"`
	Quantity              *string   `schema:"quantity"`
	ReceiptURL            *string   `schema:"receipt_url"`
	SaleGross             *string   `schema:"sale_gross"`
	Status                *string   `schema:"status"`
	SubscriptionID        uint64    `schema:"subscription_id"`
	SubscriptionPaymentID uint64    `schema:"subscription_payment_id"`
	SubscriptionPlanID    uint64    `schema:"subscription_plan_id"`
	UnitPrice             *string   `schema:"unit_price"`
	UserID                uint64    `schema:"user_id"`
}

// SubscriptionCreatedAlert is fired a new subscription is created, and a customer has successfully subscribed.
// Docs: https://developer.paddle.com/webhook-reference/subscription-alerts/subscription-created
type SubscriptionCreatedAlert struct {
	AlertName          string    `schema:"alert_name"`
	AlertID            uint64    `schema:"alert_id"`
	CancelURL          string    `schema:"cancel_url"`
	CheckoutID         *string   `schema:"checkout_id"`
	Currency           *string   `schema:"currency"`
	Email              *string   `schema:"email"`
	EventTime          time.Time `schema:"event_time"`
	MarketingConsent   bool      `schema:"marketing_consent"`
	NextBillDate       time.Time `schema:"next_bill_date"`
	Passthrough        *string   `schema:"passthrough"`
	Quantity           *string   `schema:"quantity"`
	Source             *string   `schema:"source"`
	Status             *string   `schema:"status"`
	SubscriptionID     uint64    `schema:"subscription_id"`
	SubscriptionPlanID uint64    `schema:"subscription_plan_id"`
	UnitPrice          *string   `schema:"unit_price"`
	UserID             uint64    `schema:"user_id"`
	UpdateURL          *string   `schema:"update_url"`
}

// SubscriptionUpdatedAlert is fired when the plan, price, quantity, status of an existing subscription changes, or if the payment date is rescheduled manually.
// Docs: https://developer.paddle.com/webhook-reference/subscription-alerts/subscription-updated
type SubscriptionUpdatedAlert struct {
	AlertName             string    `schema:"alert_name"`
	AlertID               uint64    `schema:"alert_id"`
	CancelURL             string    `schema:"cancel_url"`
	CheckoutID            *string   `schema:"checkout_id"`
	Email                 *string   `schema:"email"`
	EventTime             time.Time `schema:"event_time"`
	MarketingConsent      bool      `schema:"marketing_consent"`
	NewPrice              *string   `schema:"new_price"`
	NewQuantity           *string   `schema:"new_quantity"`
	NewUnitPrice          *string   `schema:"new_unit_price"`
	NextBillDate          time.Time `schema:"next_bill_date"`
	OldPrice              *string   `schema:"old_price"`
	OldQuantity           *string   `schema:"old_quantity"`
	OldUnitPrice          *string   `schema:"old_unit_price"`
	Currency              *string   `schema:"currency"`
	Passthrough           *string   `schema:"passthrough"`
	Status                *string   `schema:"status"`
	SubscriptionID        uint64    `schema:"subscription_id"`
	SubscriptionPlanID    uint64    `schema:"subscription_plan_id"`
	UserID                uint64    `schema:"user_id"`
	UpdateURL             string    `schema:"update_url"`
	OldNextBillDate       time.Time `schema:"old_next_bill_date"`
	OldStatus             *string   `schema:"old_status"`
	OldSubscriptionPlanID *string   `schema:"old_subscription_plan_id"`
	PausedAt              *string   `schema:"paused_at"`
	PausedFrom            *string   `schema:"paused_from"`
	PausedReason          *string   `schema:"paused_reason"`
}

// SubscriptionCancelledAlert is triggered whenever a user cancel a subscription.
// Docs: https://developer.paddle.com/webhook-reference/subscription-alerts/subscription-cancelled
type SubscriptionCancelledAlert struct {
	AlertName                 string    `schema:"alert_name"`
	AlertID                   uint64    `schema:"alert_id"`
	CancellationEffectiveDate time.Time `schema:"cancellation_effective_date"`
	CheckoutID                *string   `schema:"checkout_id"`
	Currency                  *string   `schema:"currency"`
	Email                     *string   `schema:"email"`
	EventTime                 time.Time `schema:"event_time"`
	MarketingConsent          bool      `schema:"marketing_consent"`
	Passthrough               *string   `schema:"passthrough"`
	Quantity                  *string   `schema:"quantity"`
	Status                    *string   `schema:"status"`
	SubscriptionID            *string   `schema:"subscription_id"`
	SubscriptionPlanID        *string   `schema:"subscription_plan_id"`
	UnitPrice                 *string   `schema:"unit_price"`
	UserID                    uint64    `schema:"user_id"`
}

// SubscriptionPaymentFailedAlert is fired when a payment for an existing subscription fails.
// Docs: https://developer.paddle.com/webhook-reference/subscription-alerts/subscription-payment-failed
type SubscriptionPaymentFailedAlert struct {
	AlertName             string    `schema:"alert_name"`
	AlertID               uint64    `schema:"alert_id"`
	Amount                *string   `schema:"amount"`
	CancelURL             string    `schema:"cancel_url"`
	CheckoutID            *string   `schema:"checkout_id"`
	Currency              *string   `schema:"currency"`
	Email                 *string   `schema:"email"`
	EventTime             time.Time `schema:"event_time"`
	MarketingConsent      bool      `schema:"marketing_consent"`
	NextRetryDate         time.Time `schema:"next_retry_date"`
	Passthrough           *string   `schema:"passthrough"`
	Quantity              *string   `schema:"quantity"`
	Status                *string   `schema:"status"`
	SubscriptionID        uint64    `schema:"subscription_id"`
	SubscriptionPlanID    uint64    `schema:"subscription_plan_id"`
	UnitPrice             *string   `schema:"unit_price"`
	UpdateURL             string    `schema:"update_url"`
	SubscriptionPaymentID uint64    `schema:"subscription_payment_id"`
	Instalments           int       `schema:"instalments"`
	OrderID               *string   `schema:"order_id"`
	UserID                uint64    `schema:"user_id"`
	AttemptNumber         *string   `schema:"attempt_number"`
}

// SubscriptionPaymentRefundedAlert is fired when a refund for an existing subscription is issued.
// Docs: https://developer.paddle.com/webhook-reference/4974afe939abc-subscription-payment-refunded
type SubscriptionPaymentRefundedAlert struct {
	AlertName               string     `schema:"alert_name"`
	AlertID                 uint64     `schema:"alert_id"`
	Amount                  *string    `schema:"amount"`
	BalanceCurrency         *string    `schema:"balance_currency"`
	BalanceEarningsDecrease *string    `schema:"balance_earnings_decrease"`
	BalanceFeeRefund        *string    `schema:"balance_fee_refund"`
	BalanceGrossRefund      *string    `schema:"balance_gross_refund"`
	BalanceTaxRefund        *string    `schema:"balance_tax_refund"`
	CheckoutID              *string    `schema:"checkout_id"`
	Currency                *string    `schema:"currency"`
	CustomData              *string    `schema:"custom_data"`
	EarningsDecrease        *string    `schema:"earnings_decrease"`
	Email                   *string    `schema:"email"`
	EventTime               time.Time  `schema:"event_time"`
	FeeRefund               *string    `schema:"fee_refund"`
	GrossRefund             *string    `schema:"gross_refund"`
	InitialPayment          bool       `schema:"initial_payment"`
	Instalments             int        `schema:"instalments"`
	MarketingConsent        bool       `schema:"marketing_consent"`
	OrderID                 *string    `schema:"order_id"`
	Passthrough             *string    `schema:"passthrough"`
	Quantity                *string    `schema:"quantity"`
	RefundReason            *string    `schema:"refund_reason"`
	RefundType              RefundType `schema:"refund_type"`
	Status                  *string    `schema:"status"`
	SubscriptionID          uint64     `schema:"subscription_id"`
	SubscriptionPaymentID   uint64     `schema:"subscription_payment_id"`
	SubscriptionPlanID      uint64     `schema:"subscription_plan_id"`
	TaxRefund               *string    `schema:"tax_refund"`
	UnitPrice               *string    `schema:"unit_price"`
	UserID                  uint64     `schema:"user_id"`
}
