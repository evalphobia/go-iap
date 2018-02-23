package appstore

import (
	"strconv"
	"time"
)

func ToReceiptInApps(aps []InApp) ReceiptInApps {
	var rc ReceiptInApps
	for _, ap := range aps {
		rc = append(rc, ToReceiptInApp(ap))
	}
	return rc
}

func ToReceiptInApp(ap InApp) *ReceiptInApp {
	return &ReceiptInApp{
		Quantity:                  ToInt64(ap.Quantity),
		ProductID:                 ap.ProductID,
		TransactionID:             ToInt64(ap.TransactionID),
		OriginalTransactionID:     ToInt64(ap.OriginalTransactionID),
		IsTrialPeriod:             ToBool(ap.IsTrialPeriod),
		AppItemID:                 ToInt64(ap.AppItemID),
		VersionExternalIdentifier: ToInt64(ap.VersionExternalIdentifier),
		WebOrderLineItemID:        ToInt64(ap.WebOrderLineItemID),
		PurchaseDate:              ToTime(ap.PurchaseDate.PurchaseDateMS),
		OriginalPurchaseDate:      ToTime(ap.OriginalPurchaseDate.OriginalPurchaseDateMS),
		ExpiresDate:               ToTime(ap.ExpiresDate.ExpiresDateMS),
		CancellationDate:          ToTime(ap.CancellationDate.CancellationDateMS),
	}
}

func ToInt64(v string) int64 {
	intVal, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0
	}
	return intVal
}

func ToBool(v string) bool {
	boolVal, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return boolVal
}

func ToTime(msString string) time.Time {
	ms, err := strconv.Atoi(msString)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(int64(ms/1000), 0)
}

func ToReceiptPendingRenewalInfos(pris []PendingRenewalInfo) ReceiptPendingRenewalInfos {
	var rpris ReceiptPendingRenewalInfos
	for _, pri := range pris {
		rpris = append(rpris, ToReceiptPendingRenewalInfo(pri))
	}
	return rpris
}

func ToReceiptPendingRenewalInfo(pri PendingRenewalInfo) *ReceiptPendingRenewalInfo {
	return &ReceiptPendingRenewalInfo{
		ExpirationIntent:   ToInt64(pri.ExpirationIntent),
		AutoRenewProductID: pri.AutoRenewProductID,
		RetryFlag:          ToBool(pri.RetryFlag),
		AutoRenewStatus:    ToBool(pri.AutoRenewStatus),
		PriceConsentStatus: ToBool(pri.PriceConsentStatus),
		ProductID:          pri.ProductID,
	}
}
