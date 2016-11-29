go-iap
======

[![Build Status](https://travis-ci.org/evalphobia/go-iap.svg?branch=master)](https://travis-ci.org/evalphobia/go-iap)
[![codecov.io](https://codecov.io/github/evalphobia/go-iap/coverage.svg?branch=master)](https://codecov.io/github/evalphobia/go-iap?branch=master)

go-iap verifies the purchase receipt via AppStore or GooglePlayStore.

Current API Documents:

* AppStore: [![GoDoc](https://godoc.org/github.com/evalphobia/go-iap/appstore?status.svg)](https://godoc.org/github.com/evalphobia/go-iap/appstore)
* GooglePlay: [![GoDoc](https://godoc.org/github.com/evalphobia/go-iap/playstore?status.svg)](https://godoc.org/github.com/evalphobia/go-iap/playstore)

# Differences from original

This repository is forked from [dogenzaka/go-iap](https://github.com/dogenzaka/go-iap)

- supports iOS6 Style receipt
- some api for iap receipts

# Dependencies
```bash
go get github.com/parnurzeal/gorequest
go get golang.org/x/net/context
go get golang.org/x/oauth2
go get google.golang.org/api/androidpublisher/v2
```

# Installation
```bash
go get github.com/evalphobia/go-iap/appstore
go get github.com/evalphobia/go-iap/playstore
```


# Quick Start

### In App Purchase (via App Store)

```go
import(
	"log"

	"github.com/evalphobia/go-iap/appstore"
)

func buy() {
	client := appstore.NewWithConfig(appstore.Config{
		TimeOut:      30 * time.Second,
		IsProduction: true,
		Retry:        true,  // retry when HTTP error status
		Debug:        false, // show HTTP request
	})

	// call apple store api to check receipt
	resp, err := client.Verify(appstore.IAPRequest{
		ReceiptData: `<your receipt data encoded by base64>`,
		Password:    `<your app's shared secret (a hexadecimal string).>`,
	})
	switch {
	case err != nil:
		log.Errof("error occured on api call: %s", err.Error())
		return
	case !resp.IsValidReceipt():
		log.Errof("invalid receipt status: %d", resp.Status)
		return
	case resp.BundleID != "<my app bundle id>":
		log.Errof("invalid bundle id: %s", resp.BundleID)
	}

	// check new receipt or not
	productID := `<prodct id>`
	transactionIDs := resp.GetTransactionIDsByProduct(productID)
	transactionIDs = filterNeverUsedTransactionIDs(transactionIDs) // check if already used one or not by your own logic
	if len(transactionIDs) == 0 {
		log.Errof("all of trnasaction id already used")
		return
	}

	var inApp *appstore.ReceiptInApp
	switch {
	case resp.IsAutoRenewable():
		// for auto-renew: last expires
		inApp = resp.GetLastExpiresByProductID(productID)
	default:
		// for consume: first element
		inApp = resp.GetByTransactionID(transactionIDs[0])
	}

	// unknown error
	if inApp == nil {
		log.Errof("cannot find valid in_app data by unknown error")
		return
	}

	// save iap data of valid receipt...
}
```

### In App Billing (via GooglePlay)

```go
import(
	"log"

	"golang.org/x/oauth2"
	"github.com/evalphobia/go-iap/playstore"
)

// check receipt
func buy() {
	// You need to prepare a public key for your Android app's in app billing
	// at https://console.developers.google.com.
	client := playstore.NewWithParams(`developer's private key data`, `developer's email`)

	resp, err := client.Verify(`<package name>`, `<subscriptionID>`, `<purchaseToken>`)
	switch {
	case err != nil:
		log.Errof("error occured on api call: %s", err.Error())
		return
	case !resp.IsValidReceipt():
		log.Errof("purchase state is invalid")
		return
	case resp.IsExpired():
		log.Errof("subscription date is expired")
		return
	}

	// save iab data of valid receipt...
	switch {
	case resp.IsValidProduct():
		// for consumed item
	case resp.IsValidSubscription():
		// for subscription item
		// autoRenew := resp.AutoRenewing
		// expires := resp.ExpiryTimeMillis
	}
}

// cancel subscription
func cancel() {
	client := playstore.NewWithParams(`developer's private key data`, `developer's email`)

	// check currenct status
	resp, err := client.Verify(`<package name>`, `<subscriptionID>`, `<purchaseToken>`)
	switch {
	case err != nil:
		log.Errof("error occured on verify api call: %s", err.Error())
		return
	case !resp.IsActive():
		// already cancelled
		return
	}

	// execute cancel operation
	err = client.CancelSubscription(`<package name>`, `<subscriptionID>`, `<purchaseToken>`)
	if err != nil {
		log.Errof("error occured on cancel api call: %s", err.Error())
		return
	}

	// confirm cancel status
	resp, err = client.Verify(`<package name>`, `<subscriptionID>`, `<purchaseToken>`)
	switch {
	case err != nil:
		log.Errof("error occured on verify api call: %s", err.Error())
		return
	case !resp.IsValidSubscription():
		log.Errof("purchase state is invalid on cancel")
		return
	case resp.IsActive():
		log.Errof("subscription is still active")
		return
	}

	// successfully cancelled
}
```

# Support

### In App Purchase
This validator supports the receipt type for both iOS6 style or newer style.

### In App Billing
This validator uses [Version 3 API](http://developer.android.com/google/play/billing/api.html).

# License
go-iap is licensed under the MIT.
