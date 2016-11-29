package appstore

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http/httputil"
	"time"

	"gopkg.in/h2non/gentleman-retry.v1"
	"gopkg.in/h2non/gentleman.v1"
	"gopkg.in/h2non/gentleman.v1/context"
	"gopkg.in/h2non/gentleman.v1/plugin"
	"gopkg.in/h2non/gentleman.v1/plugins/body"
	"gopkg.in/h2non/gentleman.v1/plugins/timeout"
)

// post sends POST request with option.
func post(url string, opt option) (*response, error) {
	opt.URL = url
	return call(opt)
}

func call(opt option) (*response, error) {
	cli := gentleman.New()
	cli.URL(opt.URL)

	req := cli.Request()
	req.Method("POST")

	// Set timeout
	if opt.hasTimeout() {
		req.Use(timeout.Request(opt.Timeout))
	}
	// Set retry (3 times)
	if opt.Retry {
		req.Use(retry.New(retry.ConstantBackoff))
	}
	// Set POST parameter
	if opt.hasPayload() {
		req.Use(body.JSON(opt.Payload))
	}

	// show debug request
	if opt.Debug {
		req.Use(debugRequest())
	}

	resp, err := req.Send()
	if opt.Debug {
		showDebugResponse(resp, err)
	}
	if err != nil {
		return nil, err
	}
	return &response{resp}, nil
}

// option is wrapper struct of http option
type option struct {
	URL     string
	Timeout time.Duration
	Retry   bool
	Debug   bool

	// POST Parameter
	Payload interface{}
}

func (o option) hasTimeout() bool {
	return o.Timeout > 0
}

func (o option) hasPayload() bool {
	return o.Payload != nil
}

// response is wrapper struct of *gentleman.Response
type response struct {
	*gentleman.Response
}

func debugRequest() plugin.Plugin {
	p := plugin.New()
	p.SetHandler("before dial", func(ctx *context.Context, h context.Handler) {
		req := ctx.Request
		body, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
		dump, _ := httputil.DumpRequest(req, false)

		fmt.Printf("---> [HTTP Request] %s[Request Body]\n%s\n", string(dump), body)
		h.Next(ctx)
	})
	return p
}

func showDebugResponse(resp *gentleman.Response, err error) {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		fmt.Printf("<--- [HTTP Response ] timeout\n\n")
		return
	}
	if resp == nil {
		return
	}

	res := resp.RawResponse
	body, _ := ioutil.ReadAll(res.Body)
	res.Body = ioutil.NopCloser(bytes.NewReader(body))
	dump, _ := httputil.DumpResponse(res, false)

	fmt.Printf("<--- [HTTP Response] %s[Response Body]\n%s\n", string(dump), body)
}
