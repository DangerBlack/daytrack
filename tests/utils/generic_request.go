package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	ACCESS_TOKEN_NAME = "_access"
	AUTHORIZATION     = "Authorization"
	BEARER            = "Bearer"
	TTL               = "ttl"
	CONCURRENCY       = "concurrency"
)

type RequestBody = map[string]string

type RequestOptions struct {
	method    string
	headers   map[string]string
	body      RequestBody
	bodyBytes []byte
	status    int
}
type RequestModifier = func(*RequestOptions, *http.Response) error

var sem = make(chan int, 5)

func FreeSemaphore() {
	<-sem
}

func WaitSemaphore() {
	sem <- 1
}

func DoRequest(url string, opts ...RequestModifier) error {
	var err error
	var reqJSONBody []byte
	opt := &RequestOptions{
		method:    http.MethodGet,
		headers:   make(map[string]string),
		body:      make(RequestBody),
		bodyBytes: nil,
		status:    http.StatusOK,
	}

	for _, modifier := range opts {
		if err = modifier(opt, nil); err != nil {
			return fmt.Errorf("error while applying modifier to request: %w", err)
		}
	}

	if reqJSONBody, err = json.Marshal(opt.body); err != nil {
		return fmt.Errorf("error while marshalling request: %w", err)
	}

	reqBody := bytes.NewBuffer(reqJSONBody)

	if opt.bodyBytes != nil {
		reqBody = bytes.NewBuffer(opt.bodyBytes)
	}

	var req *http.Request
	var timeBoundedClient = &http.Client{
		Timeout: time.Second * 100,
	}
	if req, err = http.NewRequest(opt.method, url, reqBody); err != nil {
		return fmt.Errorf("error while creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = int64(reqBody.Len())

	for key, value := range opt.headers {
		switch key {
		case TTL:
			var ttl int
			if ttl, err = strconv.Atoi(value); err != nil {
				ttl = 20
			}

			timeBoundedClient = &http.Client{
				Timeout: time.Second * time.Duration(ttl),
			}
		case CONCURRENCY:
			req.Close = true
			WaitSemaphore()
			defer FreeSemaphore()
		default:
			req.Header.Set(key, value)
		}
	}

	var res *http.Response
	if res, err = timeBoundedClient.Do(req); err != nil {
		return fmt.Errorf("error while performing the request: %w", err)
	}

	defer res.Body.Close()
	if opt.status != -1 && res.StatusCode != opt.status {
		return fmt.Errorf("error while performing the request status code expected %d but received %d instead", opt.status, res.StatusCode)
	}

	for _, modifier := range opts {
		if err = modifier(nil, res); err != nil {
			return fmt.Errorf("error while applying modifier after request: %w", err)
		}
	}

	return nil
}

func WithRequestMethod(method string) RequestModifier {
	return func(opt *RequestOptions, res *http.Response) error {
		if opt == nil {
			return nil
		}

		opt.method = method

		return nil
	}
}

func WithAccessToken(accessToken string) RequestModifier {
	return func(opt *RequestOptions, res *http.Response) error {
		if opt == nil {
			return nil
		}

		opt.headers[AUTHORIZATION] = BEARER + " " + accessToken

		return nil
	}
}

func WithTTL(ttl int) RequestModifier {
	return func(opt *RequestOptions, res *http.Response) error {
		if opt == nil {
			return nil
		}

		opt.headers[TTL] = strconv.Itoa(ttl)

		return nil
	}
}

func WithoutConcurrency() RequestModifier {
	return func(opt *RequestOptions, res *http.Response) error {
		if opt == nil {
			return nil
		}

		opt.headers[CONCURRENCY] = "false"

		return nil
	}
}

func WithRequestBody(body RequestBody) RequestModifier {
	return func(opt *RequestOptions, res *http.Response) error {
		if opt == nil {
			return nil
		}

		opt.body = body
		return nil
	}
}

func WithRequestBodyByte(body []byte) RequestModifier {
	return func(opt *RequestOptions, res *http.Response) error {
		if opt == nil {
			return nil
		}

		opt.bodyBytes = body
		return nil
	}
}

func WithRequestBodyObject(body any) RequestModifier {
	return func(opt *RequestOptions, res *http.Response) error {
		var reqJSONBody []byte
		var err error

		if opt == nil {
			return nil
		}

		if reqJSONBody, err = json.Marshal(body); err != nil {
			return err
		}

		opt.bodyBytes = reqJSONBody
		return nil
	}
}

func WithExpectedStatusCode(status int) RequestModifier {
	return func(opt *RequestOptions, res *http.Response) error {
		if opt == nil {
			return nil
		}

		opt.status = status
		return nil
	}
}

func ExtractGenericModel(response any) RequestModifier {
	return func(opt *RequestOptions, res *http.Response) error {
		var err error
		var body []byte

		if res == nil || response == nil {
			return nil
		}

		if body, err = io.ReadAll(res.Body); err != nil {
			return err
		}

		if err = json.Unmarshal(body, &response); err != nil {
			return err
		}

		return nil
	}
}
