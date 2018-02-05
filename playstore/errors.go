package playstore

import (
	"google.golang.org/api/googleapi"
)

const (
	errPurchaseTokenNoLongerValid = "purchaseTokenNoLongerValid"
)

// IsErrorCode410 checks if the error is 410 or not.
// Sometimes Google API returns 410 error "purchaseTokenNoLongerValid" when the token is invalid now.
// This often occurs when the user delete their Google account then their subscription are also deleted.
func IsErrorCode410(err error) bool {
	if e, ok := err.(*googleapi.Error); ok {
		return e.Code == 410
	}
	return false
}

// HasErrorPurchaseTokenNoLongerValid checks if the error contains "purchaseTokenNoLongerValid".
func HasErrorPurchaseTokenNoLongerValid(err error) bool {
	e, ok := err.(*googleapi.Error)
	if !ok {
		return false
	}

	for _, item := range e.Errors {
		if item.Reason == errPurchaseTokenNoLongerValid {
			return true
		}
	}
	return false
}
