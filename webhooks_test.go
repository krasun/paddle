package paddle

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestWebhooksErrorsOnInvalidHeader(t *testing.T) {
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyEncodedForPaymentSucceeded)
	if err != nil {
		t.Fatalf("failed to parse public key: %s", err)
		return
	}

	webhooks, err := NewWebhooks(publicKey)
	if err != nil {
		t.Fatalf("failed to instantiate webhooks: %s", err)
		return
	}

	r, err := http.NewRequest("POST", "https://example.com/hooks", strings.NewReader(subscriptionPaymentSucceededPostBody))
	if err != nil {
		t.Fatalf("failed to instantiate new request: %s", err)
		return
	}

	_, err = webhooks.ParseRequest(r)
	errorred(t, err, "webhook request has unsupported")
}

func TestWebhooksErrorsOnBrokenSignatureEncoding(t *testing.T) {
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyEncodedForPaymentSucceeded)
	if err != nil {
		t.Fatalf("failed to parse public key: %s", err)
		return
	}

	webhooks, err := NewWebhooks(publicKey)
	if err != nil {
		t.Fatalf("failed to instantiate webhooks: %s", err)
		return
	}

	query, err := url.ParseQuery(subscriptionPaymentSucceededPostBody)
	if err != nil {
		t.Fatalf("failed to parse query string: %s", err)
		return
	}
	query.Set("p_signature", "invalid signature")

	r, err := http.NewRequest("POST", "https://example.com/hooks", strings.NewReader(query.Encode()))
	if err != nil {
		t.Fatalf("failed to instantiate new request: %s", err)
		return
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = webhooks.ParseRequest(r)
	errorred(t, err, "failed to decode the signature")
}

func TestWebhooksErrorsOnBrokenSignatureValue(t *testing.T) {
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyEncodedForPaymentSucceeded)
	if err != nil {
		t.Fatalf("failed to parse public key: %s", err)
		return
	}

	webhooks, err := NewWebhooks(publicKey)
	if err != nil {
		t.Fatalf("failed to instantiate webhooks: %s", err)
		return
	}

	query, err := url.ParseQuery(subscriptionPaymentSucceededPostBody)
	if err != nil {
		t.Fatalf("failed to parse query string: %s", err)
		return
	}
	query.Set("p_signature", base64.StdEncoding.EncodeToString([]byte("invalid signature")))

	r, err := http.NewRequest("POST", "https://example.com/hooks", strings.NewReader(query.Encode()))
	if err != nil {
		t.Fatalf("failed to instantiate new request: %s", err)
		return
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = webhooks.ParseRequest(r)
	errorred(t, err, "failed to verify the signature")
}

func TestSubscriptionPaymentSucceededIsParsed(t *testing.T) {
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyEncodedForPaymentSucceeded)
	if err != nil {
		t.Fatalf("failed to parse public key: %s", err)
		return
	}

	webhooks, err := NewWebhooks(publicKey)
	if err != nil {
		t.Fatalf("failed to instantiate webhooks: %s", err)
		return
	}

	r, err := http.NewRequest("POST", "https://example.com/hooks", strings.NewReader(subscriptionPaymentSucceededPostBody))
	if err != nil {
		t.Fatalf("failed to instantiate new request: %s", err)
		return
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	alert, err := webhooks.ParseRequest(r)
	if err != nil {
		t.Fatalf("failed to parse request: %s", err)
		return
	}

	_, ok := alert.(*SubscriptionPaymentSucceededAlert)
	if !ok {
		t.Fatalf("alert is not of type *SubscriptionPaymentSucceededAlert")
		return
	}
}

func TestSubscriptionPaymentRefundedIsParsed(t *testing.T) {
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyEncodedForSubscriptionPaymentRefunded)
	if err != nil {
		t.Fatalf("failed to parse public key: %s", err)
		return
	}

	webhooks, err := NewWebhooks(publicKey)
	if err != nil {
		t.Fatalf("failed to instantiate webhooks: %s", err)
		return
	}

	r, err := http.NewRequest("POST", "https://example.com/hooks", strings.NewReader(subscriptionPaymentRefundedPostBody))
	if err != nil {
		t.Fatalf("failed to instantiate new request: %s", err)
		return
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	alert, err := webhooks.ParseRequest(r)
	if err != nil {
		t.Fatalf("failed to parse request: %s", err)
		return
	}

	_, ok := alert.(*SubscriptionPaymentRefundedAlert)
	if !ok {
		t.Fatalf("alert is not of type *SubscriptionPaymentRefundedAlert")
		return
	}
}

func TestSubscriptionSubscriptionCreatedIsParsed(t *testing.T) {
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyEncodedForSubscriptionCreated)
	if err != nil {
		t.Fatalf("failed to parse public key: %s", err)
		return
	}

	webhooks, err := NewWebhooks(publicKey)
	if err != nil {
		t.Fatalf("failed to instantiate webhooks: %s", err)
		return
	}

	r, err := http.NewRequest("POST", "https://example.com/hooks", strings.NewReader(subscriptionCreatedPostBody))
	if err != nil {
		t.Fatalf("failed to instantiate new request: %s", err)
		return
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	alert, err := webhooks.ParseRequest(r)
	if err != nil {
		t.Fatalf("failed to parse request: %s", err)
		return
	}

	_, ok := alert.(*SubscriptionCreatedAlert)
	if !ok {
		t.Fatalf("alert is not of type *SubscriptionCreatedAlert")
		return
	}
}

const publicKeyEncodedForPaymentSucceeded = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUEySEJEWjgycHZqY1dzVzRYQ2RLRApUeGYxcUp3ZjJ0MFhUOHcyUlVLVnd4QXVzWEJrM0huZWFIZkRPT1ZNWEUyODRDYmNZOWQvajREVlVQU0p3c2ZkCjZ1dyt6OERYb3lFWWRBVEU1eXBTVlVtNXByV0ZNMzJ4K3dVVWh1REw1MnBQbGpjKzcrYTdXL3o1OUc3V1pPK3MKaTlnQTFVbXBDRWhySWlWbk85OThBem9NUS9WemQ0Sm05ajhlN0dWSnUwR1lMMXF3eDVGeHV5SGEySnZ5L1RlYwpMejBYbVNzbzZLM3pRclYzVkNvYzJUd1N0RFFDMldLK01EQ3B3SmcwQi9FcCtIMktub043NFpDcEpkaGVFWGxoCkRJTkFyZy8yRERNNUUrQnNyS2czZEZyU2pjbnFsVTA4akRnSnVmMzdEQld4ZFNMa09nL2pTVlNCdHhqMlBtTE4KWThWME9rMy85czVybVgwdW9LaG9md2VXdER2T2JNMWE5d21saHlRWlNoc2tvWWJKbDQzM081YTMxY1ozYStKWgpnSU5TajdMSHMyMnNvQjRXQ0cvY25JQVcxbUhraU5tUnY1ZUxKeXIyZS8vSDdnSEhGRmh6ZkM5MnVabVQ4RWRuCkdhUjRWTDBMMjhnQW9pTktqUXc0RGdQZFJxRk1QNXkzR1loVm1rdk14a2VXaWQwekVvcFFFZ240akpiMkNLMUUKWkdtb3RYQUpGVXFndGM4NDJhdGZvK2pscjE5MGljUEJpcEM0Ykg4bUhpcU1yTzRwMGRocVZKS3kyQzJsMkkxOAo1cE0va0t0SCtiWitYUnR3RTlTWk5UUjJvU29hcEFlSEhSMy9kMlZub2JoOC9sbTBpRVJUM3N6K1k2NTh4THE5CjdoU0Z3Vk1uQ3pZb0wrV2ZxZFpNQUFFQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ=="
const publicKeyEncodedForSubscriptionCreated = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUEzcyt6SzR0MjJGWm81SjdWb2QvbQpLKzlPV1BrdkpyeGZvazA0aktNdVk4N3BHU1hyeWxzVWdPZUFQV2NvMDduODROT0o0c0xLTm9FaGJHaVhVZVlxCnB5NUp6UmJsUU9JUnZtQS8yZFlrcFd2WUYxL051aTNPSWZ3Ui9tWGhad3FxcUo0Nk9FU3pxUTBJZ2xCRC92dVMKTzNZbzUxS1BGb0dCTXFGYkRoODVFc0VLaWtiQmpPWjk4M1ZWTkVSUEpuV3p4dDBteVZFZ2l1ZWRZVEFiM3RyQQp6RDFKaTdaeDVDRjA5SGhRK0J6eVg4SW9UdytrQW5Sc3RqYVpEK0hLYVc3aVAzdnNPeW9uOHk4b1dZVlYwTnZxCmJBMXIwNHFpTnBuQ0dSTzdXQ1BWOGhPWXUrRXVUbUlqZ0JFNWNqbk1QRWVSMlpFSGZhTXBIUWZudk1kZlVIZU0KYU5jWkpVUEJQRFRqRDNwVGpZMXpZbHllZjFiU2llNTNSK3NUTnE1ZjVmbjFURmorUko3TmloamNYQ0habnlyawpUTDM2aUdNTkkvWnNTbk80c0NJOW5nTStZeHVDTktlbUgrbk9CTWRqYWlXL0RkVm96U0hXWXhjeGhxMW0vck03Cm9NcW9ZbitlMWhNS0I0SU02bjltN1RqTnhKVm10MGtFV3BVSVlDbE9tQTJ6bWw1ZFdQVjZNYTlqRjZDcHFSR3YKcDVObEZZMWJjUkU5L3FxeVNnNWdSMEJFK2R1TWthaWdyMUJsOWVWNXpFZDNPYmZaNm9xanpkMnZyTTM1TWJjegp3bU5sdmptMjRRUSt5ZHRSMXdvQVgyLzRsOFBqK05IV0JpOGN0WHZhTDAxaDF4c28vQ0R0NG8rNCtOL3liNDU0CmdiZ2M0NktyUmF1YnpnZlRkMkphVFBzQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo="
const publicKeyEncodedForSubscriptionPaymentRefunded = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUEzcyt6SzR0MjJGWm81SjdWb2QvbQpLKzlPV1BrdkpyeGZvazA0aktNdVk4N3BHU1hyeWxzVWdPZUFQV2NvMDduODROT0o0c0xLTm9FaGJHaVhVZVlxCnB5NUp6UmJsUU9JUnZtQS8yZFlrcFd2WUYxL051aTNPSWZ3Ui9tWGhad3FxcUo0Nk9FU3pxUTBJZ2xCRC92dVMKTzNZbzUxS1BGb0dCTXFGYkRoODVFc0VLaWtiQmpPWjk4M1ZWTkVSUEpuV3p4dDBteVZFZ2l1ZWRZVEFiM3RyQQp6RDFKaTdaeDVDRjA5SGhRK0J6eVg4SW9UdytrQW5Sc3RqYVpEK0hLYVc3aVAzdnNPeW9uOHk4b1dZVlYwTnZxCmJBMXIwNHFpTnBuQ0dSTzdXQ1BWOGhPWXUrRXVUbUlqZ0JFNWNqbk1QRWVSMlpFSGZhTXBIUWZudk1kZlVIZU0KYU5jWkpVUEJQRFRqRDNwVGpZMXpZbHllZjFiU2llNTNSK3NUTnE1ZjVmbjFURmorUko3TmloamNYQ0habnlyawpUTDM2aUdNTkkvWnNTbk80c0NJOW5nTStZeHVDTktlbUgrbk9CTWRqYWlXL0RkVm96U0hXWXhjeGhxMW0vck03Cm9NcW9ZbitlMWhNS0I0SU02bjltN1RqTnhKVm10MGtFV3BVSVlDbE9tQTJ6bWw1ZFdQVjZNYTlqRjZDcHFSR3YKcDVObEZZMWJjUkU5L3FxeVNnNWdSMEJFK2R1TWthaWdyMUJsOWVWNXpFZDNPYmZaNm9xanpkMnZyTTM1TWJjegp3bU5sdmptMjRRUSt5ZHRSMXdvQVgyLzRsOFBqK05IV0JpOGN0WHZhTDAxaDF4c28vQ0R0NG8rNCtOL3liNDU0CmdiZ2M0NktyUmF1YnpnZlRkMkphVFBzQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ=="
const subscriptionPaymentSucceededPostBody = "alert_id=1651572&alert_name=subscription_payment_succeeded&balance_currency=USD&balance_earnings=23.25&balance_fee=1.75&balance_gross=25&balance_tax=0&checkout_id=675737-chre19a993fdd63-1c4273c652&country=IL&coupon=&currency=USD&customer_name=chromium&earnings=23.25&email=qa%40screenshotone.com&event_time=2022-05-03+13%3A36%3A29&fee=1.75&initial_payment=1&instalments=1&marketing_consent=0&next_bill_date=2022-06-03&next_payment_amount=25&order_id=317366-1864142&passthrough=REhmFjVR3YcJI8D15lDi9whlBEwdcGWoyFMNvo4cSbjqMWWC%2F5TLOVA%3D&payment_method=card&payment_tax=0&plan_name=Essentials&quantity=1&receipt_url=http%3A%2F%2Fsandbox-my.paddle.com%2Freceipt%2F317366-1864142%2F675737-chre19a993fdd63-1c4273c652&sale_gross=25&status=active&subscription_id=250148&subscription_payment_id=1864142&subscription_plan_id=26279&unit_price=25.00&user_id=176032&p_signature=mPMiot4fcxKYeomJL8wpDiC6hYbssfs2nQaYfPfowd%2F3lui48feTzqBuXrY1Pp1Kwj6BG1TOkFNth%2FwErDQ0TnbiAVmOy0eWhaQvYTVGOoS0qA8cadvq4uJCHaDNbmvnidGiVEvxlaimXYcyZzHrsi123lmfLEddwvvTX%2BbsvXixI5l6vDneFuWqhNTyE1tlRCJHm9C2zDubNlvABVTCGe5%2FYQ4CS5dR5Iq8jYTZkEB6%2F0idppPUxRWnH6d1Aj3lP%2BQoGf%2FpzMYkzAFgGbK3Wa5go2bponf6T1dPKWq4hNKiKFeAQjexu9VHKcZguUFA1WgsEY4Gh5ahtASaSGoOIz2JhupzfCMMJF2i7WePf0MFun%2BCpmDBSo%2B1hDcOecj0z06P6nziE%2F3ylrw0TxjPslgT9Shuc%2FFeFeSd42fp5fWGrhUjXAyomJG5ZDBMVvwhTqxOa1ZK0y5rH2q%2BAn9wZnH3NKFblGRoVTcnfMk4WAZ%2BmPaN19ZQOz%2F24f29mMxk1p2JezahtNPptraYMt%2BXbHnp6OjX4eka3Auysr1WSL90T%2FSYLPcbzDYlhxNNAukHUW62dZZVMt0oulb1VAcl8O96EmXJoRRnOMUr%2BGRjJa1wl8onjuTMcos0b%2Bms8sZSOd7mSa2nB6KXmEO%2BS4Iy093CosjSngU1qJvD8neXHhE%3D"
const subscriptionCreatedPostBody = "alert_id=1790992&alert_name=subscription_created&cancel_url=https%3A%2F%2Fsandbox-subscription-management.paddle.com%2Fsubscription%2F264546%2Fhash%2Fea6e498c94d13147ff9c89f44918b21c8cb23f23ecdc37c9f385a875795ce6be%2Fcancel&checkout_id=731012-chrecf4c85d4606-416a3fdb5a&currency=USD&email=qa%40screenshotone.com&event_time=2022-05-25+15%3A01%3A13&linked_subscriptions=&marketing_consent=0&next_bill_date=2022-06-25&passthrough=VxD6wDdRlQPANFV5i%2FTgqacBGK%2B0CHXg0ybpOLtrbyVwYzSfVwH1aIQ%3D&quantity=1&source=localhost%3A3000+%2F+localhost%3A3000&status=active&subscription_id=264546&subscription_plan_id=29418&unit_price=7.00&update_url=https%3A%2F%2Fsandbox-subscription-management.paddle.com%2Fsubscription%2F264546%2Fhash%2Fea6e498c94d13147ff9c89f44918b21c8cb23f23ecdc37c9f385a875795ce6be%2Fupdate&user_id=176032&p_signature=qlKsmHd8zb4zVC%2Bdjs%2B2HOo%2BZLDR1ViMLBCpl2wJjolEvSLCvGf8ncT6qT1w5MA1T96GsIa8fvXcgYmO17WZamw%2BcFCD3zCA%2BKm%2BABVIXoN8pdEQj30Zc36uhSpScyYIgdHE1dHYcsE%2BF4K51hjUccAAhbrcelvYG%2FG%2FtmSqfyAknNlt2pAt23J%2BMWlk6yDMZGny40oxHfBglE7OuViupFY4tJI2sNGtu73iK0U3JdXqsX1DHkFgD7hh16ouhyjvi%2BW5cCFqg7P5jrsRn2wu%2Bdy5Qi14%2Bh1bsNdz2GE4sVwVcn6CryKLAk8PwkfAOG5xMT99dwHkq6ppgmDmZGOSmGSdCr4YcDJdjPWji26hnxqyYSXz13UK53U9h3saVeSWe374kL0yzrfNjFDUSpWW4ClQLjIZ%2F3WIE47jW3SEIRVhM4fZPiB4s2ZBovsY%2BiTBFSpQCVNYn9veWwJ2dZS3WeAO1YCpv01ioT7xSRoHgRyoO1bQPk3r61qzBJigNHArT7Vq1quBmGhmh7O08PJTuCh0Q2lKOBjXuD1B3xQCqM3x%2Ft%2FHvdBo%2BmtEYHClvRnItMCliJ32GU%2F7%2FAg9G5Rv7vWPCKt8nYr4et1%2FY7Z%2F3gw1ZvvSx7diZg0pFP0TkQkiz9TuVgrOz6FqeWNWrZxIKxZub4Bj%2BJaLg8RF3T%2Bu5bw%3D"
const subscriptionPaymentRefundedPostBody = "alert_id=1877471131&alert_name=subscription_payment_refunded&amount=15.17&balance_currency=EUR&balance_earnings_decrease=0.99&balance_fee_refund=0.48&balance_gross_refund=0.51&balance_tax_refund=0.15&checkout_id=5-575a57596786393-49632268db&currency=EUR&custom_data=custom_data&earnings_decrease=0.61&email=ppollich%40example.com&event_time=2022-09-27%2010%3A32%3A58&fee_refund=0.28&gross_refund=0.39&initial_payment=false&instalments=8&marketing_consent=&order_id=2&passthrough=Example%20String&quantity=10&refund_reason=refund_reason&refund_type=full&status=trialing&subscription_id=1&subscription_payment_id=4&subscription_plan_id=4&tax_refund=0&unit_price=unit_price&user_id=3&p_signature=h%2FQPN4aNU2JRsBi1Vv5UVSgzFJLagFNkJcAvh5y70S6Rp%2BVBh%2B%2FuSGmeZhdKwGzYeiRIf5L5Fmv8ZmpaneZ3DBjtU3KbJTLp6FNsYr0Y7TPkbZ%2B0oZ0ONI3rAf3YPXmvB5Oxg0fF0t9QWaKRkXLHTEmTcDnuYwfU5iwJV6AaRDqTxlZzOShj84RS0ZB4NO14hKDPgjd7fOqm8ODv4%2BuUSR1JZcgQFteF0KIXC7hZMJ6YPpZh2hzi3pVkrZINWy0gFxEQv8Ex7%2BR4Jd86hsjV79YJIw2FajUVxgPcC05RhpBFdS4SVJKurVrudh0Ji9x53JW9nO4EnRcgJQnAqk%2Bca3NdZIjbU%2F9x0cFn165uEQUCBTahfheoBUuCTm4MDzGOmMqxqzaUhrCkL%2FVX5NgopKwnRrXrSIl7SP2ooOqZpRTOji%2B9qJfhcH9LOgAQ1FGP0FKeMEJurwDZ%2FNFv4%2F129PkKoMVm6CjNrgMwnlngIfsfq9DcfYRmk3lP2MzKdNUj4NVyIR6uddy7FTiNrr7fo%2Fpha5E%2FS34XhIZtfls0toUOk2ErODIPyl1mZrFAJlH7ghebxi5DKdNhHZowK4iRM%2BW%2FhC8YrKRcG7ZLtwZxWstGoVPq%2FW8SxSlz1Gl6NGlxDNip3z4jjHY9ZTggx4kljFgcIb9L2F%2BbDU1chMws2ZE%3D"
