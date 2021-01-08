package task

import (
	"crypto/tls"

	"github.com/riza/gigger/pkg/config"
	"github.com/valyala/fasthttp"
)

type Task struct {
	config *config.Config
	client *fasthttp.Client
}

func NewTask(conf *config.Config) (*Task, error) {
	t := &Task{}
	t.config = conf
	t.client = &fasthttp.Client{
		MaxIdleConnDuration: conf.Timeout,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: conf.SkipSSLVerify,
		},
	}

	return t, nil
}
