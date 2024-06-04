package utils

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"time"
)

const RetryCount = 3

type retryableTransport struct {
	transport http.RoundTripper
}

// backoff calculates the backoff duration based on the number of retries.
func backoff(retries int) time.Duration {
	return time.Duration(math.Pow(2, float64(retries))) * time.Second
}

// shouldRetry determines if the request should be retried.
func shouldRetry(err error, resp *http.Response) bool {
	if err != nil {
		return true
	}

	// Retry on 5xx status codes
	if resp.StatusCode == http.StatusBadGateway ||
		resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode == http.StatusGatewayTimeout {
		return true
	}
	return false
}

// drain body to use the same connection
func drainBody(resp *http.Response) {
	if resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

/*
By default, the Golang HTTP client will close the request body after a request is sent.
This can cause issues when retrying requests since the body may have already been closed.
To prevent this from happening, we can create a custom RoundTripper
that wraps the default Transport and prevents the request body from being closed.
*/
func (t *retryableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request body
	var bodyBytes []byte
	if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	resp, err := t.transport.RoundTrip(req)

	retries := 0
	for shouldRetry(err, resp) && retries < RetryCount {
		// Wait for the specified backoff period
		time.Sleep(backoff(retries))
		drainBody(resp)

		// Clone the request body again
		if req.Body != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Retry the request
		resp, err = t.transport.RoundTrip(req)
		retries++
	}
	return resp, err
}

func NewRetryableClient() *http.Client {
	transport := &retryableTransport{
		transport: &http.Transport{},
	}

	return &http.Client{
		Transport: transport,
	}
}
