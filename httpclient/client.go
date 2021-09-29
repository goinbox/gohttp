package httpclient

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/gomisc"

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
	config  *Config
	logger  golog.Logger
	traceId []byte

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

func NewClient(config *Config, logger golog.Logger) *Client {
	c := &Client{
		config:  config,
		traceId: []byte("-"),
	}

	if logger == nil {
		c.logger = new(golog.NoopLogger)
	} else {
		c.logger = logger
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

func (c *Client) SetTraceId(traceId []byte) *Client {
	c.traceId = traceId

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
			if resp != nil {
				_ = resp.Body.Close()
			}
		}
	}

	if resp != nil {
		defer resp.Body.Close()
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
		c.logger.Error(c.fmtLog(bytes.Join(msg, []byte("\t"))))
		return nil, err
	}
	msg = append(msg, []byte("StatusCode:"+strconv.Itoa(resp.StatusCode)))
	_ = c.logger.Log(c.config.LogLevel, c.fmtLog(bytes.Join(msg, []byte("\t"))))

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

func (c *Client) fmtLog(msg []byte) []byte {
	return gomisc.AppendBytes(
		c.traceId, []byte("\t"),
		[]byte("[HttpClient]\t"),
		msg,
	)
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
