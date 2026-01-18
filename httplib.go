package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Header struct {Key, Value string}

type Request struct {
	// Are capital so that they may be public
	Method string
	Path string
	Headers []Header 
	Body string
}

type Response struct {
	StatusCode int
	Headers []Header
	Body string
}

func newRequest(method, path, host, body string) (*Request, error) {
	switch {
	case method == "":
		return nil, errors.New("missing required argument: method")	
	case path == "":
		return nil, errors.New("missing required argument: path")
	case !strings.HasPrefix(path, "/"):
		return nil, errors.New("path must start with: /")
	case host == "":
		return nil, errors.New("missing required argument: host")
	default:
		headers := []Header {
			{"Host", host},
		}
		if body != "" {
			headers = append(headers, Header{"Content-Length", fmt.Sprintf("%d", len(body))})
		}
		return &Request{Method: method, Path: path, Headers: headers, Body: body}, nil
	}
}

func newResponse(status int, body string) (*Response, error) {
	switch {
	case status < 100 || status > 599:
		return nil, errors.New("Invalid Status Code")
	default:
		if body == "" {
			body = http.StatusText(status)
		}
		headers := []Header {
			{"Content-Length", fmt.Sprintf("%d", len(body))},
		}
		return &Response{
			StatusCode: status,
			Headers: headers,
			Body: body,
		}, nil
	}
}
