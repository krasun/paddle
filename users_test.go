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
	"strings"
	"testing"
)

func TestUsersUpdateOnAPIError(t *testing.T) {
	expectedResponse := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(usersUpdateErrorJSON))),
		Header:     make(http.Header),
	}
	httpClient := newHTTPClient(func(req *http.Request) (*http.Response, error) {
		return expectedResponse, nil
	})

	u, _ := url.Parse(sandboxBaseURL)
	users := Users{httpClient: httpClient, baseURL: u, authentication: &Authentication{42, "123abc"}}

	_, _, err := users.Update(context.Background(), &UpdateUserOptions{42, 42, true, true})
	equals(t, err, &APIError{102, "Bad api key"})
}

func TestUsersUpdateOnValidationError(t *testing.T) {
	httpClient := newHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, nil
	})
	u, _ := url.Parse(sandboxBaseURL)
	users := Users{httpClient: httpClient, baseURL: u, authentication: &Authentication{42, "123abc"}}

	_, _, err := users.Update(context.Background(), &UpdateUserOptions{})
	errorred(t, err, "\"subscription_id\" is required")
}

func TestUsersUpdateOnSuccess(t *testing.T) {
	expectedResponse := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(usersUpdateJSON))),
		Header:     make(http.Header),
	}
	httpClient := newHTTPClient(func(req *http.Request) (*http.Response, error) {
		return expectedResponse, nil
	})

	u, _ := url.Parse(sandboxBaseURL)
	users := Users{httpClient: httpClient, baseURL: u, authentication: &Authentication{42, "123abc"}}

	result, actualResponse, err := users.Update(context.Background(), &UpdateUserOptions{42, 42, true, true})

	ok(t, err)
	equals(t, expectedResponse, actualResponse)
	equals(t, &UpdateUserResponse{12345, 525123, 425123, &UserPayment{144.06, "GBP", "2018-02-15"}}, result)
}

func TestUsersCancelOnAPIError(t *testing.T) {
	expectedResponse := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(usersCancelErrorJSON))),
		Header:     make(http.Header),
	}
	httpClient := newHTTPClient(func(req *http.Request) (*http.Response, error) {
		return expectedResponse, nil
	})

	u, _ := url.Parse(sandboxBaseURL)
	users := Users{httpClient: httpClient, baseURL: u, authentication: &Authentication{42, "123abc"}}

	_, err := users.Cancel(context.Background(), &CancelUserOptions{42})

	equals(t, err, &APIError{102, "Bad api key"})
}

func TestUsersCancelOnValidationError(t *testing.T) {
	httpClient := newHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, nil
	})
	u, _ := url.Parse(sandboxBaseURL)
	users := Users{httpClient: httpClient, baseURL: u, authentication: &Authentication{42, "123abc"}}

	_, err := users.Cancel(context.Background(), &CancelUserOptions{})
	errorred(t, err, "\"subscription_id\" is required")
}

func TestUsersCancelOnSuccess(t *testing.T) {
	expectedResponse := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(usersCancelJSON))),
		Header:     make(http.Header),
	}
	httpClient := newHTTPClient(func(req *http.Request) (*http.Response, error) {
		return expectedResponse, nil
	})

	u, _ := url.Parse(sandboxBaseURL)
	users := Users{httpClient: httpClient, baseURL: u, authentication: &Authentication{42, "123abc"}}

	_, err := users.Cancel(context.Background(), &CancelUserOptions{42})
	ok(t, err)
}

func TestUsersListOnAPIError(t *testing.T) {
	expectedResponse := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(usersListErrorJSON))),
		Header:     make(http.Header),
	}
	httpClient := newHTTPClient(func(req *http.Request) (*http.Response, error) {
		return expectedResponse, nil
	})

	u, _ := url.Parse(sandboxBaseURL)
	users := Users{httpClient: httpClient, baseURL: u, authentication: &Authentication{42, "123abc"}}

	result, _, err := users.List(context.Background(), &ListUsersOptions{})

	equals(t, err, &APIError{102, "Bad api key"})
	var actual []*User
	equals(t, result, actual)
}

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

const usersListErrorJSON = `{
    "success": false,
    "error": {
        "code": 102,
        "message": "Bad api key"
    }
}    
`

const usersUpdateErrorJSON = `{
    "success": false,
    "error": {
        "code": 102,
        "message": "Bad api key"
    }
}    
`
const usersUpdateJSON = `{
	"success": true,
	"response": {
		"subscription_id": 12345,
		"user_id": 425123,
		"plan_id": 525123,
		"next_payment": {
		"amount": 144.06,
		"currency": "GBP",
			"date": "2018-02-15"
		}
	}
}`

const usersCancelJSON = `{
    "success": true
}`

const usersCancelErrorJSON = `{
    "success": false,
    "error": {
        "code": 102,
        "message": "Bad api key"
    }
}`

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

// errorred fails the test if an err is nil or message is not found in the message string.
func errorred(tb testing.TB, err error, message string) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: expected error, but got nil\033[39m\n\n", filepath.Base(file), line)
		tb.FailNow()
		return
	}

	if !strings.Contains(err.Error(), message) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: \"%s\" not found in \"%s\"\033[39m\n\n", filepath.Base(file), line, message, err.Error())
		tb.FailNow()
		return
	}
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
