package config

/*
* Inspired from https://github.com/ffuf/ffuf/blob/master/pkg/ffuf/optionsparser.go
 */

import (
	"errors"
	"net/url"
	"time"
)

type ConfigOptions struct {
	General GeneralOptions
	HTTP    HTTPOptions
	Output  OutputOptions
}
type GeneralOptions struct {
	Verbose bool
	Debug   bool
	Thread  int
}

type HTTPOptions struct {
	ProxyURL      string
	Timeout       time.Duration
	SkipSSLVerify bool
	URL           string
}

type OutputOptions struct {
	OutputDirectory string
	OutputFile      string
}

func NewConfigOptions() *ConfigOptions {
	c := &ConfigOptions{}
	c.General.Verbose = false
	c.General.Debug = false
	c.General.Thread = 30
	c.HTTP.ProxyURL = ""
	c.HTTP.Timeout = 10
	c.HTTP.URL = ""
	c.HTTP.SkipSSLVerify = true
	c.Output.OutputDirectory = ""
	c.Output.OutputFile = ""
	return c
}

func ConfigFromOptions(parseOpts *ConfigOptions) (*Config, error) {
	conf := NewConfig()

	if len(parseOpts.HTTP.URL) == 0 {
		return &conf, errors.New("-u flag is required")
	}

	_, err := url.ParseRequestURI(parseOpts.HTTP.URL)
	if err != nil {
		return &conf, errors.New("incorrect url format")
	}

	conf.URL = parseOpts.HTTP.URL

	if len(parseOpts.HTTP.ProxyURL) > 0 {
		_, err := url.ParseRequestURI(parseOpts.HTTP.ProxyURL)
		if err != nil {
			return &conf, errors.New("incorrect proxy url format")
		}
		conf.ProxyURL = parseOpts.HTTP.ProxyURL
	}

	conf.Thread = parseOpts.General.Thread
	conf.OutputFile = parseOpts.Output.OutputFile
	conf.OutputDirectory = parseOpts.Output.OutputDirectory
	conf.Timeout = parseOpts.HTTP.Timeout
	conf.Verbose = parseOpts.General.Verbose
	conf.Debug = parseOpts.General.Debug

	return &conf, nil
}
