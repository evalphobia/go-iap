package appstore

import (
	"time"
)

// Receipt is struct for iap receipt data
type Receipt struct {
	responseVersion int
	rawReceipt      string
	Status          int
	Environment     string

	ReceiptType                string
	AdamID                     int64
	AppItemID                  int64
	BundleID                   string
	ApplicationVersion         string
	DownloadID                 int64
	OriginalApplicationVersion string
	RequestDate                time.Time
	OriginalPurchaseDate       time.Time
	InApps                     ReceiptInApps

	LatestReceiptInfo ReceiptInApps
	LatestReceipt     string

	PendingRenewalInfo ReceiptPendingRenewalInfos
	IsRetryable        bool
}

func (r *Receipt) String() string {
	return r.rawReceipt
}

// ResponseVersion returns receipt style;
// iOS6 style returns `6`
// iOS7 style returns `7`
func (r *Receipt) ResponseVersion() int {
	if r.responseVersion == 0 {
		return verIOS7
	}
	return r.responseVersion
}

// GetStatus returns status code of the receipt
// see: https://developer.apple.com/library/ios/releasenotes/General/ValidateAppStoreReceipt/Chapters/ValidateRemotely.html
func (r *Receipt) GetStatus() int {
	return r.Status
}

func (r *Receipt) GetEnvironment() string {
	return r.Environment
}

// LatestReceiptString returns raw receipt of `latest_receipt`
func (r *Receipt) LatestReceiptString() string {
	return r.LatestReceipt
}

// IsValidReceipt checks this receipt is valid receipt or not
// if this receipt is auto-renewable and iOS6 style, expired one returns false
func (r *Receipt) IsValidReceipt() bool {
	return r.Status == 0
}

// IsAutoRenewable checks this receipt is auto-renewable subscription or not
func (r *Receipt) IsAutoRenewable() bool {
	return r.InApps.IsAutoRenewable()
}

func (r *Receipt) HasError() error {
	return HandleError(r.Status)
}

// HasExpired checks this receipt is expired or not (only for iOS6 style)
func (r *Receipt) HasExpired() bool {
	return r.Status == 21006
}

// ShouldSendToTestEnvironment checks this receipt status is 21007
// this receipt is from the test environment, but it was sent to the production environment for verification
func (r *Receipt) ShouldSendToTestEnvironment() bool {
	return r.Status == 21007
}

// ShouldSendToProductionEnvironment checks this receipt status is 21008
// this receipt is from the production environment, but it was sent to the test environment for verification
func (r *Receipt) ShouldSendToProductionEnvironment() bool {
	return r.Status == 21008
}

// GetTransactionIDs returns all of transaction_id from `in_app`
func (r *Receipt) GetTransactionIDs() []int64 {
	return r.InApps.TransactionIDs()
}

// GetTransactionIDsByProduct returns all of transaction_id from `in_app` filtered by `product_id`
func (r *Receipt) GetTransactionIDsByProduct(product string) []int64 {
	checked := make(map[int64]bool)
	var matched []int64
	latest := r.LatestReceiptInfo.LastExpiresByProductIDForLatest(product)
	if latest != nil {
		txID := latest.TransactionID
		checked[txID] = true
		matched = append(matched, txID)
	}

	for _, v := range r.InApps {
		txID := v.TransactionID
		if checked[txID] {
			continue
		}
		checked[txID] = true
		if v.ProductID != product {
			continue
		}
		matched = append(matched, txID)
	}
	return matched
}

// GetTransactionIDsWithoutExpired returns all of transaction_id except expired
func (r *Receipt) GetTransactionIDsWithoutExpired() []int64 {
	checked := make(map[int64]bool)
	now := time.Now()
	var matched []int64
	matchedIDs := func(rc ReceiptInApps) {
		for _, v := range rc {
			if checked[v.TransactionID] {
				continue
			}
			checked[v.TransactionID] = true

			switch {
			case !v.ExpiresDate.IsZero() && v.ExpiresDate.Before(now):
				continue
			}
			matched = append(matched, v.TransactionID)
		}
	}
	matchedIDs(r.InApps)
	matchedIDs(r.LatestReceiptInfo)
	return matched
}

// GetTransactionIDsByProductWithoutExpired returns all of transaction_id filtered by `product_id` except expired
func (r *Receipt) GetTransactionIDsByProductWithoutExpired(product string) []int64 {
	checked := make(map[int64]bool)
	now := time.Now()
	var matched []int64
	matchedIDs := func(rc ReceiptInApps) {
		for _, v := range rc {
			if checked[v.TransactionID] {
				continue
			}
			checked[v.TransactionID] = true

			switch {
			case v.ProductID != product:
				continue
			case !v.ExpiresDate.IsZero() && v.ExpiresDate.Before(now):
				continue
			}
			matched = append(matched, v.TransactionID)
		}
	}
	matchedIDs(r.InApps)
	matchedIDs(r.LatestReceiptInfo)
	return matched
}

// GetByTransactionID returns receipt data by `transaction_id`
func (r *Receipt) GetByTransactionID(id int64) *ReceiptInApp {
	return r.InApps.ByTransactionID(id)
}

// GetLastExpiresByProductID returns latest expires receipt data by `product_id`
func (r *Receipt) GetLastExpiresByProductID(productID string) *ReceiptInApp {
	inAppLatest := r.LatestReceiptInfo.LastExpiresByProductIDForLatest(productID)
	inApp := r.InApps.LastExpiresByProductID(productID)
	switch {
	case inApp == nil:
		return inAppLatest
	case inAppLatest == nil:
		return inApp
	case inApp.ExpiresDate.After(inAppLatest.ExpiresDate):
		return inApp
	}
	return inAppLatest
}

// GetLastExpiresByTransactionIDs returns latest expires receipt data from `transaction_id` list
func (r *Receipt) GetLastExpiresByTransactionIDs(ids []int64) *ReceiptInApp {
	inAppLatest := r.LatestReceiptInfo.LastExpiresByTransactionIDsForLatest(ids)
	inApp := r.InApps.LastExpiresByTransactionIDs(ids)
	switch {
	case inApp == nil:
		return inAppLatest
	case inAppLatest == nil:
		return inApp
	case inApp.ExpiresDate.After(inAppLatest.ExpiresDate):
		return inApp
	}
	return inAppLatest
}
