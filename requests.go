package gotiktoklive

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type reqOptions struct {
	// Endpoint is the request path of tiktok api
	Endpoint string

	// IsPost set to true will send request with POST method.
	//
	// By default this option is false.
	IsPost bool

	// Compress post form data with gzip
	Gzip bool

	// Query is the parameters of the request
	//
	// This parameters are independents of the request method (POST|GET)
	Query map[string]string

	// List of headers to ignore
	IgnoreHeaders []string

	// Extra headers to add
	ExtraHeaders map[string]string

	// Use base tiktok URi instead of webcast api
	OmitAPI bool

	// Specifiy base URI
	URI string
}

func (t *TikTok) sendRequest(o *reqOptions) (body []byte, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Failed to send request to %s: %w", o.Endpoint, err)
		}
	}()
	method := "GET"
	if o.IsPost {
		method = "POST"
	}

	uri := tiktokAPIUrl
	if o.OmitAPI {
		uri = tiktokBaseUrl
	}
	if o.URI != "" {
		uri = o.URI
	}
	uri = uri + o.Endpoint

	u, err := url.Parse(uri)
	if err != nil {
		return
	}

	vs := url.Values{}
	for k, v := range o.Query {
		if v != "" {
			vs.Add(k, v)
		}
	}

	reqData := bytes.NewBuffer([]byte{})
	if o.IsPost {
		reqData.WriteString(vs.Encode())
	} else {
		u.RawQuery = vs.Encode()
	}

	ua := userAgent
	fullUrl := u.String()
	if !o.OmitAPI && o.URI == "" && o.Endpoint == urlRoomData {
		signed, err := t.signURL(fullUrl)
		if err != nil {
			return nil, err
		}
		ua = signed.UserAgent
		fullUrl = signed.SignedURL
	}

	var req *http.Request
	req, err = http.NewRequest(method, fullUrl, reqData)
	if err != nil {
		return
	}

	ignoreHeader := func(h string) bool {
		for _, k := range o.IgnoreHeaders {
			if k == h {
				return true
			}
		}
		return false
	}

	setHeaders := func(h map[string]string) {
		for k, v := range h {
			if v != "" && !ignoreHeader(k) {
				req.Header.Set(k, v)
			}
		}
	}

	headers := map[string]string{
		"Connection": "keep-alive",
		// "Connection":      "close",
		"Cache-Control":   "max-age=0",
		"User-Agent":      ua,
		"Accept":          "text/html,application/json,application/protobuf",
		"Referer":         referer,
		"Origin":          origin,
		"Accept-Language": "en-US,en;q=0.9",
		"Accept-Encoding": "gzip, deflate",
	}

	setHeaders(headers)
	setHeaders(o.ExtraHeaders)

	resp, err := t.c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var bb bytes.Buffer
	_, err = io.Copy(&bb, resp.Body)
	if err != nil {
		return
	}
	body = bb.Bytes()

	if resp.StatusCode == 429 {
		err = ErrRateLimitExceeded
		return
	} else if resp.StatusCode >= 400 {
		err = fmt.Errorf("received status code %d", resp.StatusCode)
		return
	}

	// Decode gzip encoded responses
	encoding := resp.Header.Get("Content-Encoding")
	if encoding != "" && encoding == "gzip" {
		buf := bytes.NewBuffer(body)
		zr, err := gzip.NewReader(buf)
		if err != nil {
			return nil, err
		}
		var bb bytes.Buffer
		_, err = io.Copy(&bb, zr)
		body = bb.Bytes()
		if err != nil {
			return body, err
		}
		if err := zr.Close(); err != nil {
			return body, err
		}
	}

	// Log complete response body
	if t.LogRequests {
		r := map[string]interface{}{
			"status":   resp.StatusCode,
			"endpoint": o.Endpoint,
			"body":     string(body),
		}

		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return nil, err
		}
		t.debugHandler(string(b))
	}
	return
}
