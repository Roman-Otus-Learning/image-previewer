package client

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

var _ Client = (*HTTPClient)(nil)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

type Client interface {
	GetWithHeaders(ctx context.Context, url string, headers http.Header) (*http.Response, error)
}

type HTTPClient struct {
	client *http.Client
}

func CreateHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *HTTPClient) WithRoundTripFunc(f RoundTripperFunc) *HTTPClient {
	c.client.Transport = f

	return c
}

func (c *HTTPClient) GetWithHeaders(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	rq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "HTTPClient create request")
	}

	for key, hh := range headers {
		for _, h := range hh {
			rq.Header.Add(key, h)
		}
	}

	rsp, err := c.client.Do(rq)
	if err != nil {
		return nil, errors.Wrap(err, "HTTPClient do request")
	}

	return rsp, nil
}
