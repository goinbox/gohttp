package httpclient

import (
	"github.com/goinbox/golog"

	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	config *Config
	logger golog.ILogger

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

func NewClient(config *Config, logger golog.ILogger) *Client {
	c := &Client{
		config: config,
	}

	if logger == nil {
		c.logger = new(golog.NoopLogger)
	} else {
		c.logger = logger
	}

	c.client = &http.Client{
		Timeout: config.Timeout,

		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   config.Timeout,
				KeepAlive: config.KeepAliveTime,
			}).Dial,
		},
	}

	return c
}

func (c *Client) Do(req *Request, retry int) (*Response, error) {
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
		}
	}

	msg := [][]byte{
		[]byte("Method:" + req.Method),
		[]byte("Host: " + req.Host),
		[]byte("Url:" + req.URL.String()),
		[]byte("Time:" + t.String()),
	}
	if err != nil {
		if resp != nil {
			msg = append(msg, []byte("StatusCode:"+strconv.Itoa(resp.StatusCode)))
		}
		msg = append(msg, []byte("ErrMsg:"+err.Error()))
		c.logger.Error(bytes.Join(msg, []byte("\t")))
		return nil, err
	}
	msg = append(msg, []byte("StatusCode:"+strconv.Itoa(resp.StatusCode)))
	c.logger.Log(c.config.LogLevel, bytes.Join(msg, []byte("\t")))

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		T:        t,
		Contents: contents,
		Response: resp,
	}, nil
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
