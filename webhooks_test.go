package paddle

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestWebhooksErrorsOnInvalidHeader(t *testing.T) {
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyEncoded)
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
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyEncoded)
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
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyEncoded)
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
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyEncoded)
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

const publicKeyEncoded = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUEySEJEWjgycHZqY1dzVzRYQ2RLRApUeGYxcUp3ZjJ0MFhUOHcyUlVLVnd4QXVzWEJrM0huZWFIZkRPT1ZNWEUyODRDYmNZOWQvajREVlVQU0p3c2ZkCjZ1dyt6OERYb3lFWWRBVEU1eXBTVlVtNXByV0ZNMzJ4K3dVVWh1REw1MnBQbGpjKzcrYTdXL3o1OUc3V1pPK3MKaTlnQTFVbXBDRWhySWlWbk85OThBem9NUS9WemQ0Sm05ajhlN0dWSnUwR1lMMXF3eDVGeHV5SGEySnZ5L1RlYwpMejBYbVNzbzZLM3pRclYzVkNvYzJUd1N0RFFDMldLK01EQ3B3SmcwQi9FcCtIMktub043NFpDcEpkaGVFWGxoCkRJTkFyZy8yRERNNUUrQnNyS2czZEZyU2pjbnFsVTA4akRnSnVmMzdEQld4ZFNMa09nL2pTVlNCdHhqMlBtTE4KWThWME9rMy85czVybVgwdW9LaG9md2VXdER2T2JNMWE5d21saHlRWlNoc2tvWWJKbDQzM081YTMxY1ozYStKWgpnSU5TajdMSHMyMnNvQjRXQ0cvY25JQVcxbUhraU5tUnY1ZUxKeXIyZS8vSDdnSEhGRmh6ZkM5MnVabVQ4RWRuCkdhUjRWTDBMMjhnQW9pTktqUXc0RGdQZFJxRk1QNXkzR1loVm1rdk14a2VXaWQwekVvcFFFZ240akpiMkNLMUUKWkdtb3RYQUpGVXFndGM4NDJhdGZvK2pscjE5MGljUEJpcEM0Ykg4bUhpcU1yTzRwMGRocVZKS3kyQzJsMkkxOAo1cE0va0t0SCtiWitYUnR3RTlTWk5UUjJvU29hcEFlSEhSMy9kMlZub2JoOC9sbTBpRVJUM3N6K1k2NTh4THE5CjdoU0Z3Vk1uQ3pZb0wrV2ZxZFpNQUFFQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ=="
const subscriptionPaymentSucceededPostBody = "alert_id=1651572&alert_name=subscription_payment_succeeded&balance_currency=USD&balance_earnings=23.25&balance_fee=1.75&balance_gross=25&balance_tax=0&checkout_id=675737-chre19a993fdd63-1c4273c652&country=IL&coupon=&currency=USD&customer_name=chromium&earnings=23.25&email=qa%40screenshotone.com&event_time=2022-05-03+13%3A36%3A29&fee=1.75&initial_payment=1&instalments=1&marketing_consent=0&next_bill_date=2022-06-03&next_payment_amount=25&order_id=317366-1864142&passthrough=REhmFjVR3YcJI8D15lDi9whlBEwdcGWoyFMNvo4cSbjqMWWC%2F5TLOVA%3D&payment_method=card&payment_tax=0&plan_name=Essentials&quantity=1&receipt_url=http%3A%2F%2Fsandbox-my.paddle.com%2Freceipt%2F317366-1864142%2F675737-chre19a993fdd63-1c4273c652&sale_gross=25&status=active&subscription_id=250148&subscription_payment_id=1864142&subscription_plan_id=26279&unit_price=25.00&user_id=176032&p_signature=mPMiot4fcxKYeomJL8wpDiC6hYbssfs2nQaYfPfowd%2F3lui48feTzqBuXrY1Pp1Kwj6BG1TOkFNth%2FwErDQ0TnbiAVmOy0eWhaQvYTVGOoS0qA8cadvq4uJCHaDNbmvnidGiVEvxlaimXYcyZzHrsi123lmfLEddwvvTX%2BbsvXixI5l6vDneFuWqhNTyE1tlRCJHm9C2zDubNlvABVTCGe5%2FYQ4CS5dR5Iq8jYTZkEB6%2F0idppPUxRWnH6d1Aj3lP%2BQoGf%2FpzMYkzAFgGbK3Wa5go2bponf6T1dPKWq4hNKiKFeAQjexu9VHKcZguUFA1WgsEY4Gh5ahtASaSGoOIz2JhupzfCMMJF2i7WePf0MFun%2BCpmDBSo%2B1hDcOecj0z06P6nziE%2F3ylrw0TxjPslgT9Shuc%2FFeFeSd42fp5fWGrhUjXAyomJG5ZDBMVvwhTqxOa1ZK0y5rH2q%2BAn9wZnH3NKFblGRoVTcnfMk4WAZ%2BmPaN19ZQOz%2F24f29mMxk1p2JezahtNPptraYMt%2BXbHnp6OjX4eka3Auysr1WSL90T%2FSYLPcbzDYlhxNNAukHUW62dZZVMt0oulb1VAcl8O96EmXJoRRnOMUr%2BGRjJa1wl8onjuTMcos0b%2Bms8sZSOd7mSa2nB6KXmEO%2BS4Iy093CosjSngU1qJvD8neXHhE%3D"
