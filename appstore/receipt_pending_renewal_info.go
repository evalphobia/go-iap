package appstore

// ReceiptPendingRenewalInfo is struct for pending_renewal_info field
type ReceiptPendingRenewalInfo struct {
	ExpirationIntent   int64  `json:"expiration_intent"`
	AutoRenewProductID string `json:"auto_renew_product_id"`
	RetryFlag          bool   `json:"is_in_billing_retry_period"`
	AutoRenewStatus    bool   `json:"auto_renew_status"`
	PriceConsentStatus bool   `json:"price_consent_status"`
	ProductID          string `json:"product_id"`
}

type ReceiptPendingRenewalInfos []*ReceiptPendingRenewalInfo

func (r ReceiptPendingRenewalInfos) IsAutoRenewable(productID string) bool {
	for _, v := range r {
		if v.ProductID == productID {
			return v.AutoRenewStatus
		}
	}
	return false
}
