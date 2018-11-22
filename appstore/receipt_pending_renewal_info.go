package appstore

// ReceiptPendingRenewalInfo is struct for pending_renewal_info field.
type ReceiptPendingRenewalInfo struct {
	ExpirationIntent   int64  `json:"expiration_intent"`
	AutoRenewProductID string `json:"auto_renew_product_id"`
	RetryFlag          bool   `json:"is_in_billing_retry_period"`
	AutoRenewStatus    bool   `json:"auto_renew_status"`
	PriceConsentStatus bool   `json:"price_consent_status"`
	ProductID          string `json:"product_id"`
}

// IsDifferentAutoRenewProductID checks that AutoRenewProductID is changed from ProductID.
func (r ReceiptPendingRenewalInfo) IsDifferentAutoRenewProductID() bool {
	return r.ProductID != r.AutoRenewProductID
}

type ReceiptPendingRenewalInfos []*ReceiptPendingRenewalInfo

// GetRenewalInfo returns ReceiptPendingRenewalInfo of given productID.
func (r ReceiptPendingRenewalInfos) GetRenewalInfo(productID string) *ReceiptPendingRenewalInfo {
	for _, v := range r {
		if v.AutoRenewProductID == productID {
			return v
		}
	}
	return nil
}

// IsAutoRenewStatusOn confirms `auto_renew_status` is enabled for given product id.
func (r ReceiptPendingRenewalInfos) IsAutoRenewStatusOn(productID string) bool {
	for _, v := range r {
		if v.AutoRenewProductID == productID {
			return v.AutoRenewStatus
		}
	}
	return false
}

// IsAutoRenewStatusOff confirms `auto_renew_status` is disabled for given product id.
func (r ReceiptPendingRenewalInfos) IsAutoRenewStatusOff(productID string) bool {
	for _, v := range r {
		if v.AutoRenewProductID == productID {
			return !v.AutoRenewStatus
		}
	}
	return false
}
