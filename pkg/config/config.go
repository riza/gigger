package config

import (
	"time"
)

type Config struct {
	Thread        int
	ProxyURL      string
	SkipSSLVerify bool
	Timeout       time.Duration
	URL           string
	Verbose       bool
}

func NewConfig() Config {
	var conf Config
	conf.ProxyURL = ""
	conf.Timeout = 60 * time.Second
	conf.SkipSSLVerify = true
	conf.Thread = 30
	conf.URL = ""
	conf.Verbose = false
	return conf
}
