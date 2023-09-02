package httpclient

import (
	"github.com/goinbox/golog"
	"github.com/goinbox/pcontext"

	"net/http"
	"testing"
	"time"
)

var ctx pcontext.Context
var client *Client

func init() {
	w, _ := golog.NewFileWriter("/dev/stdout", 0)
	logger := golog.NewSimpleLogger(w, golog.NewSimpleFormater())
	ctx = pcontext.NewSimpleContext(logger)

	config := NewConfig()
	config.Timeout = time.Second * 1
	client = NewClient(config)
}

func TestClientGet(t *testing.T) {
	extHeaders := map[string]string{
		"GO-CLIENT-1": "gobox-httpclient-1",
		"GO-CLIENT-2": "gobox-httpclient-2",
	}
	req, _ := NewRequest(http.MethodGet, "http://www.vmubt.com/test.php?a=1&b=2", nil, "127.0.0.1", extHeaders)

	resp, err := client.Do(ctx, req, 1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(resp.Contents), resp.T.String())
	}
}

func TestClientPost(t *testing.T) {
	extHeaders := map[string]string{
		"GO-CLIENT-1":  "gobox-httpclient-1",
		"GO-CLIENT-2":  "gobox-httpclient-2",
		"Content-Type": "application/x-www-form-urlencoded;charset=utf-8",
	}
	params := map[string]interface{}{
		"a": 1,
		"b": "bb",
		"c": "测试post",
	}
	req, _ := NewRequest(http.MethodPost, "http://www.vmubt.com/test.php", MakeRequestBodyUrlEncoded(params), "127.0.0.1", extHeaders)

	resp, err := client.Do(ctx, req, 1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(resp.Contents), resp.T.String())
	}
}
