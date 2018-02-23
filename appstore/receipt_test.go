package appstore

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testReceipt1 = &Receipt{}
var testReceipt2 = &Receipt{}
var testReceipt3 = &Receipt{}

func init() {
	initializeReceipts()
}

func initializeReceipts() {
	list := []struct {
		receipt *Receipt
		data    string
	}{
		{testReceipt1, testReceiptString},
		{testReceipt2, testReceiptString2},
		{testReceipt3, testReceiptString3},
	}

	for _, r := range list {
		var result IAPResponseIOS7
		json.Unmarshal([]byte(r.data), &result)
		*r.receipt = *result.ToReceipt()
	}
}

func TestGetStatus(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt  *Receipt
		expected int
	}{
		{testReceipt1, 0},
		{testReceipt2, 21006},
		{testReceipt3, 21002},
	}

	for _, tt := range tests {
		assert.Equal(tt.expected, tt.receipt.GetStatus())
	}
}

func TestGetEnvironment(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt  *Receipt
		expected string
	}{
		{testReceipt1, "Sandbox"},
		{testReceipt2, "Production"},
		{testReceipt3, "Production"},
	}

	for _, tt := range tests {
		assert.Equal(tt.expected, tt.receipt.GetEnvironment())
	}
}

func TestLatestReceiptString(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt  *Receipt
		expected string
	}{
		{testReceipt1, "<latest receipt string>"},
		{testReceipt2, "<latest receipt string>"},
		{testReceipt3, "<latest receipt string>"},
	}

	for _, tt := range tests {
		assert.Equal(tt.expected, tt.receipt.LatestReceiptString())
	}
}

func TestIsValidReceipt(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt  *Receipt
		expected bool
	}{
		{testReceipt1, true},
		{testReceipt2, false},
		{testReceipt3, false},
	}

	for _, tt := range tests {
		assert.Equal(tt.expected, tt.receipt.IsValidReceipt())
	}
}

func TestIsAutoRenewable(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt  *Receipt
		expected bool
	}{
		{testReceipt1, true},
		{testReceipt2, true},
		{testReceipt3, false},
	}

	for _, tt := range tests {
		assert.Equal(tt.expected, tt.receipt.IsAutoRenewable())
	}
}

func TestHasError(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt *Receipt
		has     bool
	}{
		{testReceipt1, false},
		{testReceipt2, false},
		{testReceipt3, true},
	}

	for _, tt := range tests {
		if tt.has {
			assert.Error(tt.receipt.HasError())
		} else {
			assert.NoError(tt.receipt.HasError())
		}
	}
}

func TestHasExpired(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt  *Receipt
		expected bool
	}{
		{testReceipt1, false},
		{testReceipt2, true},
		{testReceipt3, false},
	}

	for _, tt := range tests {
		assert.Equal(tt.expected, tt.receipt.HasExpired())
	}
}

func TestGetTransactionIDs(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt     *Receipt
		expectedLen int
		expected1st int
	}{
		{testReceipt1, 242, 1000000068359169},
		{testReceipt2, 1, 1000000068359170},
		{testReceipt3, 1, 1000000181765148},
	}

	for _, tt := range tests {
		ids := tt.receipt.GetTransactionIDs()
		assert.EqualValues(tt.expected1st, ids[0])
		assert.Len(ids, tt.expectedLen)
	}
}

func TestGetTransactionIDsByProduct(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt     *Receipt
		product     string
		expectedLen int
		expected1st int
	}{
		{testReceipt1, "com.example.app.subscription_1.v2", 44, 1000000183885918},
		{testReceipt1, "com.example.app.subscription_3.v2", 90, 1000000146311029},
		{testReceipt1, "com.example.app.subscription_6.v2", 36, 1000000150322509},
		{testReceipt1, "com.example.app.subscription_12.v2", 18, 1000000181778917},
		{testReceipt1, "invalid_id", 0, 0},
		{testReceipt2, "com.example.app.subscription_1", 2, 1000000068358624},
		{testReceipt3, "com.example.app.consumable_10", 1, 1000000181765148},
	}

	for _, tt := range tests {
		ids := tt.receipt.GetTransactionIDsByProduct(tt.product)
		assert.Len(ids, tt.expectedLen)
		if tt.expectedLen > 0 {
			assert.EqualValues(tt.expected1st, ids[0])
		}
	}
}

func TestGetByTransactionID(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt         *Receipt
		txID            int64
		has             bool
		expectedProduct string
	}{
		{testReceipt1, 1000000117928271, true, "com.example.app.subscription_1.v2"},
		{testReceipt1, 1000000116297757, true, "com.example.app.subscription_3.v2"},
		{testReceipt1, 1000000112476475, true, "com.example.app.subscription_6.v2"},
		{testReceipt1, 1000000112589395, true, "com.example.app.subscription_12.v2"},
		{testReceipt1, 999, false, ""},
		{testReceipt2, 1000000068359170, true, "com.example.app.subscription_1"},
		{testReceipt3, 1000000181765148, true, "com.example.app.consumable_10"},
	}

	for _, tt := range tests {
		inApp := tt.receipt.GetByTransactionID(tt.txID)
		if !tt.has {
			assert.Nil(inApp)
			continue
		}
		assert.Equal(tt.txID, inApp.TransactionID)
		assert.Equal(tt.expectedProduct, inApp.ProductID)
	}
}

func TestGetLastExpiresByProductID(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt         *Receipt
		product         string
		has             bool
		expectedTxID    int64
		expiectedExpire string
	}{
		{testReceipt1, "com.example.app.subscription_1.v2", true, 1000000183885918, "2015-12-07T23:49:14Z"},
		{testReceipt1, "com.example.app.subscription_3.v2", true, 1000000146311029, "2015-03-09T05:02:26Z"},
		{testReceipt1, "com.example.app.subscription_6.v2", true, 1000000150322509, "2015-04-03T09:30:20Z"},
		{testReceipt1, "com.example.app.subscription_12.v2", true, 1000000181778917, "2015-11-25T07:42:15Z"},
		{testReceipt1, "invalid_id", false, 0, ""},
		{testReceipt2, "com.example.app.subscription_1", true, 1000000068359170, "2013-03-18T06:22:47Z"},
		{testReceipt3, "com.example.app.consumable_10", true, 1000000181765148, "0001-01-01T00:00:00Z"},
	}

	for _, tt := range tests {
		inApp := tt.receipt.GetLastExpiresByProductID(tt.product)
		if !tt.has {
			assert.Nil(inApp)
			continue
		}
		assert.Equal(tt.product, inApp.ProductID)
		assert.Equal(tt.expectedTxID, inApp.TransactionID)
		expire, _ := time.Parse(time.RFC3339, tt.expiectedExpire)
		assert.Equal(expire.Unix(), inApp.ExpiresDate.Unix(), inApp.ExpiresDate.String())
	}
}

func TestLastExpiresByProductID(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("com.example.app", testReceipt1.BundleID)

	inApp := testReceipt1.InApps.LastExpiresByProductID("com.example.app.subscription_1.v2")
	expectedExpire, _ := time.Parse(time.RFC3339, "2015-12-07T23:24:14Z")
	assert.Equal("com.example.app.subscription_1.v2", inApp.ProductID)
	assert.EqualValues(1000000183882899, inApp.TransactionID)
	assert.Equal(expectedExpire.Unix(), inApp.ExpiresDate.Unix())

	inApp = testReceipt1.LatestReceiptInfo.LastExpiresByProductID("com.example.app.subscription_1.v2")
	expectedExpire, _ = time.Parse(time.RFC3339, "2015-12-07T23:49:14Z")
	assert.Equal("com.example.app.subscription_1.v2", inApp.ProductID)
	assert.EqualValues(1000000183885918, inApp.TransactionID)
	assert.Equal(expectedExpire.Unix(), inApp.ExpiresDate.Unix())
}

func TestLastExpiresByProductIDForLatest(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("com.example.app", testReceipt1.BundleID)

	inApp := testReceipt1.LatestReceiptInfo.LastExpiresByProductIDForLatest("com.example.app.subscription_1.v2")
	expectedExpire, _ := time.Parse(time.RFC3339, "2015-12-07T23:49:14Z")
	assert.Equal("com.example.app.subscription_1.v2", inApp.ProductID)
	assert.EqualValues(1000000183885918, inApp.TransactionID)
	assert.Equal(expectedExpire.Unix(), inApp.ExpiresDate.Unix())
}

func TestGetTransactionIDsWithoutExpired(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt     *Receipt
		expectedLen int
		expected1st int
	}{
		{testReceipt1, 2, 1000000183886962},
		{testReceipt2, 0, 1000000068359170},
		{testReceipt3, 1, 1000000181765148},
	}

	for _, tt := range tests {
		ids := tt.receipt.GetTransactionIDsWithoutExpired()
		assert.Len(ids, tt.expectedLen)
		if tt.expectedLen > 0 {
			assert.EqualValues(tt.expected1st, ids[0])
		}
	}
}

func TestGetTransactionIDsByProductWithoutExpired(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt     *Receipt
		product     string
		expectedLen int
		expected1st int
	}{
		{testReceipt1, "com.example.app.subscription_1.v2", 0, 1000000117928271},
		{testReceipt1, "com.example.app.subscription_3.v2", 0, 1000000116297757},
		{testReceipt1, "com.example.app.subscription_6.v2", 0, 1000000112476475},
		{testReceipt1, "com.example.app.subscription_12.v2", 0, 1000000112589395},
		{testReceipt1, "com.example.app.subscription_long", 1, 1000000183886963},
		{testReceipt1, "invalid_id", 0, 0},
		{testReceipt2, "com.example.app.subscription_1", 0, 1000000068359170},
		{testReceipt3, "com.example.app.consumable_10", 1, 1000000181765148},
	}

	for _, tt := range tests {
		ids := tt.receipt.GetTransactionIDsByProductWithoutExpired(tt.product)
		assert.Len(ids, tt.expectedLen)
		if tt.expectedLen > 0 {
			assert.EqualValues(tt.expected1st, ids[0])
		}
	}
}

func TestReceiptPendingRenewalInfoIsAutoRenewable(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		receipt   *Receipt
		productID string
		expected  bool
	}{
		{testReceipt1, "com.example.app.subscription_1", true},
		{testReceipt1, "com.example.app.subscription_2", false},
	}

	for _, tt := range tests {
		assert.Equal(
			tt.expected,
			tt.receipt.PendingRenewalInfo.IsAutoRenewable(tt.productID),
		)
	}
}
