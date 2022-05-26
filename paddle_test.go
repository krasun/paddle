package paddle

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestUsersSandboxURL(t *testing.T) {
	client, err := NewSandboxClient(Authentication{VendorID: 42, VendorAuthCode: "abc"})
	ok(t, err)

	equals(t, "https://sandbox-vendors.paddle.com/api/", client.Users.baseURL.String())
}

func TestUsersProductionURL(t *testing.T) {
	client, err := NewProductionClient(Authentication{VendorID: 42, VendorAuthCode: "abc"})
	ok(t, err)

	equals(t, "https://vendors.paddle.com/api/", client.Users.baseURL.String())
}

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
