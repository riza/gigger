package main

import (
	"flag"
	"os"
	"time"

	"github.com/riza/gigger/pkg/config"
	"github.com/riza/gigger/pkg/gigger"
	"github.com/riza/gigger/pkg/task"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ParseFlags(opts *config.ConfigOptions) *config.ConfigOptions {
	flag.StringVar(&opts.HTTP.URL, "u", opts.HTTP.URL, "Target URL")
	flag.StringVar(&opts.HTTP.ProxyURL, "x", opts.HTTP.ProxyURL, "HTTP Proxy URL")
	flag.DurationVar(&opts.HTTP.Timeout, "timeout", opts.HTTP.Timeout, "HTTP request timeout in seconds.")
	flag.BoolVar(&opts.HTTP.SkipSSLVerify, "ssl", opts.HTTP.SkipSSLVerify, "todo")

	flag.StringVar(&opts.Output.OutputDirectory, "od", opts.Output.OutputDirectory, "Directory path to store matched results to.")
	flag.StringVar(&opts.Output.OutputFile, "o", opts.Output.OutputFile, "Write output to file")

	flag.BoolVar(&opts.General.Verbose, "v", opts.General.Verbose, "Verbose output, printing full URL and redirect location (if any) with the results.")
	flag.BoolVar(&opts.General.Debug, "d", opts.General.Debug, "You can catch errors more easily in debug mode.")
	flag.Parse()
	return opts
}

func main() {
	var err error

	var opts *config.ConfigOptions
	opts = config.NewConfigOptions()
	opts = ParseFlags(opts)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC822})

	conf, err := config.ConfigFromOptions(opts)
	if err != nil {
		log.Error().Msgf("[Config]: %s\n", err)
		os.Exit(1)
	}

	t, err := task.NewTask(conf)
	if err != nil {
		log.Error().Msgf("[Task]: %s\n", err)
		os.Exit(1)
	}

	p, err := gigger.NewPool(conf, t)
	if err != nil {
		log.Error().Msgf("[Pool]: %s\n", err)
		os.Exit(1)
	}

	p.Run()
}
