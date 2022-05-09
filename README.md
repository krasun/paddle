# paddle

[![Build](https://github.com/krasun/paddle/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/krasun/paddle/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/krasun/paddle/branch/main/graph/badge.svg?token=rh8BDdHc2v)](https://codecov.io/gh/krasun/paddle)
[![Go Report Card](https://goreportcard.com/badge/github.com/krasun/paddle)](https://goreportcard.com/report/github.com/krasun/paddle)
[![GoDoc](https://godoc.org/https://godoc.org/github.com/krasun/paddle?status.svg)](https://godoc.org/github.com/krasun/paddle)

A [Paddle API](https://www.paddle.com/) client for Go. 

## Installation

```shell
go get github.com/krasun/paddle
```

## Usage

Import the library: 
```go
import "github.com/krasun/paddle"
```

Calling API methods: 
```go
var paddleClient *paddle.Client
switch environment {
case "sandbox":
    paddleClient, err = paddle.NewSandboxClient(paddle.Authentication{paddleVendorID, paddleVendorAuthCode})
    if err != nil {
        log.Fatalf("failed to instantiate paddle %s client: %s", environment, err)
        return
    }
case "production":
    paddleClient, err = paddle.NewProductionClient(paddle.Authentication{paddleVendorID, paddleVendorAuthCode})
    if err != nil {
        log.Fatalf("failed to instantiate paddle %s client: %s", environment, err)
        return
    }
default:
    log.Fatalf("unsupported Paddle environment: %s", environment)
    return
}

options := &paddle.UpdateUserOptions{
    // ... 
    Prorate:         true,
    BillImmediately: true,
    // ...
}
response, _, err := paddleClient.Users.Update(ctx, options)
if err != nil {
    log.Error(err)
    return
}
```

Handling webhooks:
```go
webhooks, err := paddle.NewWebhooks(paddlePublicKey)
if err != nil {
    log.Fatalf("failed to instantiate Paddle webhooks client: %s", err)
    return
}

func handlePaddleWebhooks(webhooks *paddle.Webhooks, payments *payments) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		alert, err := webhooks.ParseRequest(r)
		if err != nil {
			log.Error(err)
			http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
			return
		}

		switch alert := alert.(type) {
        // ... 
		case *paddle.SubscriptionCreatedAlert:
			err := payments.processSubscriptionCreated(ctx, alert)
			if err != nil {
				log.Error(err)
				http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
				return
			}
		case *paddle.SubscriptionUpdatedAlert:
            // ... 
			return
		case *paddle.SubscriptionCancelledAlert:
			// ... 
			return
        // ...
		}
	}
}
```

## Tests 

To run tests, just execute: 
```
$ go test . 
```

## License 

`paddle` is released under [the MIT license](LICENSE).
