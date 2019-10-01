package playstore

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/androidpublisher/v3"
)

const (
	scope = "https://www.googleapis.com/auth/androidpublisher"

	defaultTimeout = time.Second * 5
)

var timeout = defaultTimeout

// SetTimeout sets dial timeout duration
func SetTimeout(t time.Duration) {
	timeout = t
}

// The IABClient type is an interface to verify purchase token
type IABClient interface {
	VerifySubscription(string, string, string) (*androidpublisher.SubscriptionPurchase, error)
	VerifyProduct(string, string, string) (*androidpublisher.ProductPurchase, error)
}

// The Client type implements VerifySubscription method
type Client struct {
	httpClient *http.Client
}

// New returns http client which includes the credentials to access androidpublisher API.
// You should create a service account for your project at
// https://console.developers.google.com and download a JSON key file to set this argument.
func New(jsonKey []byte) (Client, error) {
	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, &http.Client{
		Timeout: timeout,
	})

	conf, err := google.JWTConfigFromJSON(jsonKey, scope)

	return Client{conf.Client(ctx)}, err
}

func NewWithParams(key, email string) Client {
	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, &http.Client{
		Timeout: timeout,
	})
	conf := &jwt.Config{
		Email:      email,
		PrivateKey: []byte(key),
		Scopes:     []string{scope},
		TokenURL:   google.JWTTokenURL,
	}
	return Client{conf.Client(ctx)}
}

// Verify retrieves product and subscription status from GooglePlay API
func (c *Client) Verify(packageName, productID, token string) (*IABResponse, error) {
	resp, err := c.VerifyProduct(packageName, productID, token)
	if err == nil {
		return resp, nil
	}
	resp.SubscriptionPurchase, err = c.verifySubscription(packageName, productID, token)
	return resp, err
}

// VerifySubscription retrieves product status from GooglePlay API
func (c *Client) VerifySubscription(packageName, productID, token string) (*IABResponse, error) {
	result, err := c.verifySubscription(packageName, productID, token)
	return &IABResponse{SubscriptionPurchase: result}, err
}

func (c *Client) verifySubscription(packageName, subscriptionID, token string) (*androidpublisher.SubscriptionPurchase, error) {
	service, err := androidpublisher.New(c.httpClient)
	if err != nil {
		return nil, err
	}

	ps := androidpublisher.NewPurchasesSubscriptionsService(service)
	result, err := ps.Get(packageName, subscriptionID, token).Do()

	return result, err
}

// VerifyProduct retrieves product status from GooglePlay API
func (c *Client) VerifyProduct(packageName, productID, token string) (*IABResponse, error) {
	result, err := c.verifyProduct(packageName, productID, token)
	return &IABResponse{ProductPurchase: result}, err
}

func (c *Client) verifyProduct(packageName, productID, token string) (*androidpublisher.ProductPurchase, error) {
	service, err := androidpublisher.New(c.httpClient)
	if err != nil {
		return nil, err
	}

	ps := androidpublisher.NewPurchasesProductsService(service)
	result, err := ps.Get(packageName, productID, token).Do()

	return result, err
}

// CancelSubscription cancels recurring payment of given subscription
func (c *Client) CancelSubscription(packageName, subscriptionID, token string) error {
	service, err := androidpublisher.New(c.httpClient)
	if err != nil {
		return err
	}

	ps := androidpublisher.NewPurchasesSubscriptionsService(service)
	return ps.Cancel(packageName, subscriptionID, token).Do()
}
