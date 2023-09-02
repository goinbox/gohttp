package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/goinbox/golog"
	"github.com/goinbox/pcontext"
)

type Client struct {
	config *Config
	client *http.Client
}

type Request struct {
	Method     string
	Url        string
	Body       []byte
	Ip         string
	ExtHeaders map[string]string

	*http.Request
}

type Response struct {
	T        time.Duration
	Contents []byte

	*http.Response
}

func NewClient(config *Config) *Client {
	c := &Client{
		config: config,
	}

	c.client = &http.Client{
		Timeout: config.Timeout,

		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   config.Timeout,
				KeepAlive: config.KeepAliveTime,
			}).DialContext,
			DisableKeepAlives:   config.DisableKeepAlives,
			MaxIdleConns:        config.MaxIdleConns,
			MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
			IdleConnTimeout:     config.IdleConnTimeout,
		},
	}

	return c
}

func (c *Client) Do(ctx pcontext.Context, req *Request, retry int) (*Response, error) {
	resp, err := c.request(ctx, req, retry)
	fields := []*golog.Field{
		{
			Key:   "t",
			Value: resp.T,
		},
	}
	if resp.Response != nil {
		defer func() { _ = resp.Body.Close() }()
		fields = append(fields, &golog.Field{
			Key:   "StatusCode",
			Value: resp.StatusCode,
		})
	}

	logger := ctx.Logger()
	if err != nil {
		fields = append(fields, &golog.Field{
			Key:   "Error",
			Value: err,
		})
		logger.Error("client request error", fields...)
		return nil, err
	}

	resp.Contents, err = io.ReadAll(resp.Body)
	if err != nil {
		fields = append(fields, &golog.Field{
			Key:   "Error",
			Value: err,
		})
		logger.Error("read response body error", fields...)
		return nil, err
	}

	if c.config.LogResponseBody {
		fields = append(fields, &golog.Field{
			Key:   "ResponseBody",
			Value: string(resp.Contents),
		})
	}

	logger.Info("response", fields...)

	return resp, nil
}

func (c *Client) request(ctx pcontext.Context, req *Request, retry int) (*Response, error) {
	fields := []*golog.Field{
		{
			Key:   "Method",
			Value: req.Method,
		},
		{
			Key:   "Host",
			Value: req.Host,
		},
		{
			Key:   "url",
			Value: req.URL.String(),
		},
	}
	if c.config.LogRequestBody {
		fields = append(fields, &golog.Field{
			Key:   "RequestBody",
			Value: string(req.Body),
		})
	}

	ctx.Logger().Info("request", fields...)

	start := time.Now()
	resp, err := c.client.Do(req.Request)
	t := time.Since(start)
	if err != nil {
		for i := 0; i < retry; i++ {
			req, _ = NewRequest(req.Method, req.Url, req.Body, req.Ip, req.ExtHeaders)
			start = time.Now()
			resp, err = c.client.Do(req.Request)
			t = time.Since(start)
			if err == nil && resp.StatusCode == 200 {
				break
			}
			if resp != nil {
				_ = resp.Body.Close()
			}
		}
	}

	return &Response{
		T:        t,
		Contents: nil,
		Response: resp,
	}, err
}

func NewRequest(method string, url string, body []byte, ip string, extHeaders map[string]string) (*Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Host = req.URL.Host

	if ip != "" {
		s := strings.Split(req.URL.Host, ":")
		s[0] = ip
		req.URL.Host = strings.Join(s, ":")
	}

	if extHeaders != nil {
		for k, v := range extHeaders {
			req.Header.Set(k, v)
		}
	}

	return &Request{
		Method:     method,
		Url:        url,
		Body:       body,
		Ip:         ip,
		ExtHeaders: extHeaders,

		Request: req,
	}, nil
}

func MakeRequestBodyUrlEncoded(params map[string]interface{}) []byte {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, fmt.Sprint(value))
	}

	return []byte(values.Encode())
}
