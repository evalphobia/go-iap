package appstore

const verIOS7 = 7

// The IAPResponse type has the response properties
type IAPResponseIOS7 struct {
	responseVersion    int                  `json:"-"`
	rawReceipt         string               `json:"-"`
	Status             int                  `json:"status"`
	Environment        string               `json:"environment"`
	Receipt            ReceiptIOS7          `json:"receipt"`
	LatestReceiptInfo  []InApp              `json:"latest_receipt_info"`
	LatestReceipt      string               `json:"latest_receipt"`
	PendingRenewalInfo []PendingRenewalInfo `json:"pending_renewal_info"`
	IsRetryable        bool                 `json:"is_retryable"`
}

func NewIAPResponseIOS7(rc string) *IAPResponseIOS7 {
	return &IAPResponseIOS7{rawReceipt: rc}
}

func (r *IAPResponseIOS7) ToReceipt() *Receipt {
	rr := r.Receipt
	receipt := &Receipt{
		responseVersion:            r.responseVersion,
		rawReceipt:                 r.rawReceipt,
		Status:                     r.Status,
		Environment:                r.Environment,
		ReceiptType:                rr.ReceiptType,
		AdamID:                     rr.AdamID,
		AppItemID:                  rr.AppItemID,
		BundleID:                   rr.BundleID,
		ApplicationVersion:         rr.ApplicationVersion,
		DownloadID:                 rr.DownloadID,
		OriginalApplicationVersion: rr.OriginalApplicationVersion,
		RequestDate:                ToTime(rr.RequestDate.RequestDateMS),
		OriginalPurchaseDate:       ToTime(rr.OriginalPurchaseDate.OriginalPurchaseDateMS),
		LatestReceipt:              r.LatestReceipt,
		IsRetryable:                r.IsRetryable,
	}
	receipt.InApps = ToReceiptInApps(rr.InApp)
	receipt.LatestReceiptInfo = ToReceiptInApps(r.LatestReceiptInfo)
	receipt.PendingRenewalInfo = ToReceiptPendingRenewalInfos(r.PendingRenewalInfo)
	return receipt
}

// The Receipt type has whole data of receipt
type ReceiptIOS7 struct {
	ReceiptType                string  `json:"receipt_type"`
	AdamID                     int64   `json:"adam_id"`
	AppItemID                  int64   `json:"app_item_id"`
	BundleID                   string  `json:"bundle_id"`
	ApplicationVersion         string  `json:"application_version"`
	DownloadID                 int64   `json:"download_id"`
	OriginalApplicationVersion string  `json:"original_application_version"`
	InApp                      []InApp `json:"in_app"`
	RequestDate
	OriginalPurchaseDate
}

// The InApp type has the receipt attributes
type InApp struct {
	Quantity                  string `json:"quantity"`
	ProductID                 string `json:"product_id"`
	TransactionID             string `json:"transaction_id"`
	OriginalTransactionID     string `json:"original_transaction_id"`
	IsTrialPeriod             string `json:"is_trial_period"`
	IsInIntroOfferPeriod      string `json:"is_in_intro_offer_period"`
	AppItemID                 string `json:"app_item_id"`
	VersionExternalIdentifier string `json:"version_external_identifier"`
	WebOrderLineItemID        string `json:"web_order_line_item_id"`
	PurchaseDate
	OriginalPurchaseDate
	ExpiresDate
	CancellationDate
}

// PendingRenewalInfo auto-renewable subscriptions
type PendingRenewalInfo struct {
	ExpirationIntent   string `json:"expiration_intent"`
	AutoRenewProductID string `json:"auto_renew_product_id"`
	RetryFlag          string `json:"is_in_billing_retry_period"`
	AutoRenewStatus    string `json:"auto_renew_status"`
	PriceConsentStatus string `json:"price_consent_status"`
	ProductID          string `json:"product_id"`
}
