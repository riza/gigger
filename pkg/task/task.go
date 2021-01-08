package task

import (
	"crypto/tls"

	"github.com/riza/gigger/pkg/config"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type Task struct {
	config *config.Config
	client *fasthttp.Client
}

func NewTask(conf *config.Config) (*Task, error) {
	t := &Task{}
	t.config = conf
	t.client = &fasthttp.Client{
		Dial:                fasthttp.Dial,
		MaxIdleConnDuration: conf.Timeout,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: conf.SkipSSLVerify,
		},
	}

	if len(t.config.ProxyURL) >= 1 {
		t.client.Dial = fasthttpproxy.FasthttpHTTPDialerTimeout(conf.ProxyURL, conf.Timeout)
	}

	return t, nil
}
