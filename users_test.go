package paddle

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

const usersListJSON = `{
    "success": true,
    "response": [
        {
            "subscription_id": 232564,
            "plan_id": 26100,
            "user_id": 176032,
            "user_email": "qa@screenshotone.com",
            "update_url": "https://sandbox-subscription-management.paddle.com/subscription/232564/hash/ff7900e4181d56c1a733f8d11c1663267676f392050cff4c82171855a617f43a/update",
            "cancel_url": "https://sandbox-subscription-management.paddle.com/subscription/232564/hash/ff7900e4181d56c1a733f8d11c1663267676f392050cff4c82171855a617f43a/cancel",
            "state": "active",
            "signup_date": "2022-03-30 16:06:05",
            "last_payment": {
                "amount": 7,
                "currency": "USD",
                "date": "2022-04-30"
            },
            "next_payment": {
                "amount": 7,
                "currency": "USD",
                "date": "2022-05-30"
            },
            "payment_information": {
                "payment_method": "card",
                "card_type": "visa",
                "last_four_digits": "4242",
                "expiry_date": "12/2023"
            }
        },
        {
            "subscription_id": 232566,
            "plan_id": 26100,
            "user_id": 176032,
            "user_email": "qa@screenshotone.com",
            "update_url": "https://sandbox-subscription-management.paddle.com/subscription/232566/hash/44709b85411754493e4b881c0916eb5969616d3d412e9e04c8262aee0fee6bee/update",
            "cancel_url": "https://sandbox-subscription-management.paddle.com/subscription/232566/hash/44709b85411754493e4b881c0916eb5969616d3d412e9e04c8262aee0fee6bee/cancel",
            "state": "active",
            "signup_date": "2022-03-30 16:11:10",
            "last_payment": {
                "amount": 7,
                "currency": "USD",
                "date": "2022-04-30"
            },
            "next_payment": {
                "amount": 7,
                "currency": "USD",
                "date": "2022-05-30"
            },
            "payment_information": {
                "payment_method": "card",
                "card_type": "visa",
                "last_four_digits": "4242",
                "expiry_date": "12/2022"
            }
        }
    ]
}
`

func TestUsersListOnSuccess(t *testing.T) {
	expectedResponse := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(usersListJSON))),
		Header:     make(http.Header),
	}
	httpClient := newHTTPClient(func(req *http.Request) (*http.Response, error) {
		return expectedResponse, nil
	})

	u, _ := url.Parse(sandboxBaseURL)
	users := Users{httpClient: httpClient, baseURL: u, authentication: &Authentication{42, "123abc"}}

	result, actualResponse, err := users.List(context.Background(), &ListUsersOptions{})
	ok(t, err)

	equals(t, expectedResponse, actualResponse)

	equals(t, 2, len(result))
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

type roundTripper func(req *http.Request) (*http.Response, error)

func (rt roundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return rt(request)
}

func newHTTPClient(roundTripper roundTripper) *http.Client {
	return &http.Client{
		Transport: roundTripper,
	}
}
