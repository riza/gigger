package config

import (
	"time"
)

//Config
type Config struct {
	OutputDirectory string
	OutputFile      string

	Thread        int
	ProxyURL      string
	SkipSSLVerify bool
	Timeout       time.Duration
	URL           string

	Verbose bool
	Debug   bool
}

func NewConfig() Config {
	var conf Config
	conf.ProxyURL = ""
	conf.Timeout = 60 * time.Second
	conf.SkipSSLVerify = true
	conf.Thread = 30
	conf.URL = ""
	conf.Verbose = false
	conf.Debug = false
	return conf
}
