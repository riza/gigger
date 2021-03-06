<img align="right" width="200" src="https://github.com/riza/gigger/blob/master/res/pickaxe.png?raw=true" />

# Gigger
> Git folder digger, I'm sure it's worthwhile stuff.

[![Build Status](https://github.com/riza/gigger/workflows/.github/workflows/test.yml/badge.svg)](https://github.com/riza/gigger/)  [![GitHub version](https://badge.fury.io/gh/riza%2Fgigger.svg)](https://github.com/riza/gigger/releases) [![Go Report Card](https://goreportcard.com/badge/github.com/riza/gigger)](https://goreportcard.com/report/github.com/riza/gigger) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/riza/gigger) [![codecov](https://codecov.io/gh/riza/gigger/branch/master/graph/badge.svg)](https://codecov.io/gh/riza/gigger)


## Installation

- Download a prebuilt binary from [releases page](https://github.com/riza/gigger/releases/latest).

  _or_
- If you have recent go compiler installed: `GO11MODULES=on go get -u github.com/riza/gigger`
  
## Usage

```
█▀▀ █ █▀▀ █▀▀ █▀▀ █▀█
█▄█ █ █▄█ █▄█ ██▄ █▀▄

Usage of gigger:
  -ssl
        Disable SSL verification (default true)
  -t int
        Concurrent process count (default 30)
  -timeout duration
        HTTP request timeout in seconds. (default 10ns)
  -u string
        Target URL
  -v    Verbose output, printing full URL and redirect location (if any) with the results.
  -x string
        HTTP Proxy URL
```

## p.s.

We also need a new icon, anyone who can make me a pickaxe icon can reach me on [Twitter](https://twitter.com/rizasabuncu). :)

## TODO

- [ ] Output Provider (html,json,csv)
- [ ] Usage docs
- [ ] Comments for godoc
- [ ] Codecov with Github Actions
- [ ] Multi url scan
- [ ] Index of detection
- [x] Download real files (entries, name)

## License

gigger is released under MIT license. See [LICENSE](https://github.com/riza/gigger/blob/master/LICENSE).
