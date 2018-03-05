package appstore

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIAPResponseIOS6(t *testing.T) {
	assert := assert.New(t)

	rc := "receipt"
	r := NewIAPResponseIOS6(rc)
	assert.Equal(rc, r.rawReceipt)
}

func TestIAPResponseIOS6(t *testing.T) {
	assert := assert.New(t)

	rc := "receipt"
	r := NewIAPResponseIOS6(rc)

	err := json.Unmarshal([]byte(testReceiptStringIOS6_1), r)
	assert.Nil(err)
	assert.Equal(0, r.Status)
	assert.Equal("dummy_latest_receipt", r.LatestReceipt)
	assert.Equal("900000001", r.Receipt.AppItemID)
	assert.Equal("com.example.app", r.Receipt.BundleID)
	assert.Equal("0.1", r.Receipt.ApplicationVersion)
	assert.Equal("", r.Receipt.OriginalApplicationVersion)
	assert.Equal("90000000000001", r.Receipt.OriginalTransactionID)
	assert.Equal("com.example.product.item", r.Receipt.ProductID)
	assert.Equal("1", r.Receipt.Quantity)
	assert.Equal("90000000000001", r.Receipt.TransactionID)
	assert.Equal("900000000", r.Receipt.VersionExternalIdentifier)
	assert.Equal("70000000000001", r.Receipt.WebOrderLineItemID)
	assert.Equal("2015-06-09 23:31:25 Etc/GMT", r.Receipt.ExpiresDate)
	assert.Equal("1433892685000", r.Receipt.ExpiresDateMS)
	assert.Equal("2015-06-09 16:31:25 America/Los_Angeles", r.Receipt.ExpiresDatePST)

	assert.Equal("", r.Receipt.RequestDate.RequestDate)
	assert.Equal("2015-05-09 23:31:25 Etc/GMT", r.Receipt.PurchaseDate.PurchaseDate)
	assert.Equal("2015-05-09 23:31:28 Etc/GMT", r.Receipt.OriginalPurchaseDate.OriginalPurchaseDate)
}

func TestToIOS7(t *testing.T) {
	assert := assert.New(t)

	baseResponse := IAPResponseIOS7{
		Status:          0,
		responseVersion: 6,
		rawReceipt:      "receipt",
		LatestReceipt:   "dummy_latest_receipt",
		Environment:     "",
		Receipt: ReceiptIOS7{
			ReceiptType:                "",
			AdamID:                     0,
			AppItemID:                  900000001,
			BundleID:                   "com.example.app",
			ApplicationVersion:         "0.1",
			DownloadID:                 0,
			OriginalApplicationVersion: "",
			RequestDate: RequestDate{
				RequestDate: "",
			},
			OriginalPurchaseDate: OriginalPurchaseDate{
				OriginalPurchaseDate:    "2015-05-09 23:31:28 Etc/GMT",
				OriginalPurchaseDateMS:  "1431214288317",
				OriginalPurchaseDatePST: "2015-05-09 16:31:28 America/Los_Angeles",
			},
		},
	}

	baseInApp := InApp{
		Quantity:                  "1",
		ProductID:                 "com.example.product.item",
		TransactionID:             "90000000000001",
		OriginalTransactionID:     "90000000000001",
		VersionExternalIdentifier: "900000000",
		WebOrderLineItemID:        "70000000000001",
		PurchaseDate: PurchaseDate{
			PurchaseDate:    "2015-05-09 23:31:25 Etc/GMT",
			PurchaseDateMS:  "1431214285000",
			PurchaseDatePST: "2015-05-09 16:31:25 America/Los_Angeles",
		},
		OriginalPurchaseDate: OriginalPurchaseDate{
			OriginalPurchaseDate:    "2015-05-09 23:31:28 Etc/GMT",
			OriginalPurchaseDateMS:  "1431214288317",
			OriginalPurchaseDatePST: "2015-05-09 16:31:28 America/Los_Angeles",
		},
		ExpiresDate: ExpiresDate{
			ExpiresDate:    "2015-06-09 23:31:25 Etc/GMT",
			ExpiresDateMS:  "1433892685000",
			ExpiresDatePST: "2015-06-09 16:31:25 America/Los_Angeles",
		},
		CancellationDate: CancellationDate{},
	}

	baseLatestReceiptInfo := InApp{
		Quantity:                  "1",
		ProductID:                 "com.example.product.item",
		TransactionID:             "90000000000002",
		OriginalTransactionID:     "10000010000000",
		VersionExternalIdentifier: "",
		WebOrderLineItemID:        "70000000000002",
		PurchaseDate: PurchaseDate{
			PurchaseDate:    "2015-11-10 00:31:25 Etc/GMT",
			PurchaseDateMS:  "1447115485000",
			PurchaseDatePST: "2015-11-09 16:31:25 America/Los_Angeles",
		},
		OriginalPurchaseDate: OriginalPurchaseDate{
			OriginalPurchaseDate:    "2015-05-09 23:31:28 Etc/GMT",
			OriginalPurchaseDateMS:  "1431214288000",
			OriginalPurchaseDatePST: "2015-05-09 16:31:28 America/Los_Angeles",
		},
		ExpiresDate: ExpiresDate{
			ExpiresDate:    "2015-12-10 00:31:25 Etc/GMT",
			ExpiresDateMS:  "1449707485000",
			ExpiresDatePST: "2015-12-09 16:31:25 America/Los_Angeles",
		},
		CancellationDate: CancellationDate{},
	}

	expectedResponse := baseResponse
	expectedResponse.Receipt.InApp = []InApp{baseInApp}
	expectedResponse.LatestReceiptInfo = []InApp{baseLatestReceiptInfo}

	tests := []struct {
		targetReceipt              string
		expectedResponse           *IAPResponseIOS7
		expectedPendingRenewalInfo *PendingRenewalInfo
	}{
		{
			targetReceipt:              testReceiptStringIOS6_1,
			expectedResponse:           &expectedResponse,
			expectedPendingRenewalInfo: nil,
		},
		{
			targetReceipt:    testReceiptStringIOS6_2,
			expectedResponse: &expectedResponse,
			expectedPendingRenewalInfo: &PendingRenewalInfo{
				AutoRenewStatus:    "1",
				AutoRenewProductID: "com.example.product.item",
			},
		},
		{
			targetReceipt:    testReceiptStringIOS6_3,
			expectedResponse: &expectedResponse,
			expectedPendingRenewalInfo: &PendingRenewalInfo{
				AutoRenewStatus:    "0",
				AutoRenewProductID: "com.example.product.item",
			},
		},
	}

	rc := "receipt"
	for _, tt := range tests {
		target := fmt.Sprintf("%+v", tt)
		if tt.expectedPendingRenewalInfo != nil {
			tt.expectedResponse.PendingRenewalInfo = []PendingRenewalInfo{*tt.expectedPendingRenewalInfo}
		}

		r6 := NewIAPResponseIOS6(rc)
		err := json.Unmarshal([]byte(tt.targetReceipt), r6)
		assert.Nil(err, target)

		r7 := r6.ToIOS7()
		assert.Equal(tt.expectedResponse, r7, target)
	}
}

func TestToIOS7ToReceipt(t *testing.T) {
	assert := assert.New(t)

	rc := "receipt"
	r6 := NewIAPResponseIOS6(rc)
	err := json.Unmarshal([]byte(testReceiptStringIOS6_1), r6)
	assert.Nil(err)

	r7 := r6.ToIOS7()
	r := r7.ToReceipt()
	assert.Equal(6, r.responseVersion)
}
