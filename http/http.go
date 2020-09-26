package http

import (
	"time"
	"net"
	"net/http"
)

const (
	userAgent = "Mozilla/5.0 (compatible; cf-bypass); https://github.com/amourha/cf-bypass)"
)

type Client struct {
	clientObj 	*http.Client
	maxRetries 	uint
}

func NewHTTPClient(maxRetries uint) (*Client) {
	// Setup http client
	client := &http.Client{
		Timeout: time.Second * 15,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	return &Client {
		clientObj: client,
		maxRetries: maxRetries,
	}
}

func (httpClient *Client) DoRequest(url string) (*http.Response, error) {
	var lastError error

	for i := httpClient.maxRetries; i >= 0; i-- {
		req, err := http.NewRequest("GET", url, nil)
		if nil != err {
			lastError = err
			continue
		}

		req.Header.Add("User-Agent", userAgent)
		resp, err := httpClient.clientObj.Do(req)
		if nil == err {
			// If it got here we succeeded
			return resp, nil
		}
	}

	// If it got here we failed
	return nil, lastError
}