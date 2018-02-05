package playstore

import (
	"errors"
	"fmt"
	"testing"

	"google.golang.org/api/googleapi"
)

func TestIsErrorCode410(t *testing.T) {
	tests := []struct {
		expected bool
		err      error
	}{
		{false, errors.New("error")},
		{false, nil},
		{false, &googleapi.Error{}},
		{false, &googleapi.Error{
			Code: 409,
		}},
		{true, &googleapi.Error{
			Code: 410,
		}},
		{false, &googleapi.Error{
			Code: 411,
		}},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("expected=%t, err=%+v", tt.expected, tt.err)
		actual := IsErrorCode410(tt.err)
		if actual != tt.expected {
			t.Errorf("got %#v\nwant %#v, data=%s", actual, tt.expected, target)
		}
	}
}

func TestHasErrorPurchaseTokenNoLongerValid(t *testing.T) {
	tests := []struct {
		expected bool
		err      error
	}{
		{false, errors.New("error")},
		{false, nil},
		{false, &googleapi.Error{}},
		{false, &googleapi.Error{
			Code: 410,
		}},
		{false, &googleapi.Error{
			Errors: []googleapi.ErrorItem{},
		}},
		{false, &googleapi.Error{
			Errors: []googleapi.ErrorItem{{Reason: "error"}},
		}},
		{true, &googleapi.Error{
			Errors: []googleapi.ErrorItem{{Reason: errPurchaseTokenNoLongerValid}},
		}},
	}

	for _, tt := range tests {
		target := fmt.Sprintf("expected=%t, err=%+v", tt.expected, tt.err)
		actual := HasErrorPurchaseTokenNoLongerValid(tt.err)
		if actual != tt.expected {
			t.Errorf("got %#v\nwant %#v, data=%s", actual, tt.expected, target)
		}
	}
}
