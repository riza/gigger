package task

import (
	"crypto/tls"
	"net/url"
	"os"
	"strings"

	"github.com/riza/gigger/pkg/config"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

const (
	resultsPath = "./results/"
)

type Task struct {
	conf   *config.Config
	Client *fasthttp.Client
	host   string
}

func NewTask(conf *config.Config) (*Task, error) {
	t := &Task{}
	t.conf = conf
	t.Client = &fasthttp.Client{
		Dial:                fasthttp.Dial,
		MaxIdleConnDuration: conf.Timeout,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: conf.SkipSSLVerify,
		},
	}

	if len(t.conf.ProxyURL) >= 1 {
		t.Client.Dial = fasthttpproxy.FasthttpHTTPDialerTimeout(conf.ProxyURL, conf.Timeout)
	}

	var err error
	t.host, err = t.getDomainName()
	if err != nil {
		return t, err
	}

	err = t.NewResultFolder()
	if err != nil {
		return t, err
	}

	err = t.NewTaskFolder()
	if err != nil {
		return t, err
	}

	return t, nil
}

func (t *Task) NewResultFolder() error {
	if _, err := os.Stat(resultsPath); os.IsNotExist(err) {
		err = os.Mkdir(resultsPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) NewTaskFolder() error {
	filesPath := resultsPath + t.host + "/files"
	if _, err := os.Stat(filesPath); os.IsNotExist(err) {
		err = os.MkdirAll(filesPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Task) SaveFile(name string, data string) error {
	filesPath := resultsPath + t.host + "/files/" + name
	f, err := os.Create(filesPath + name)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func (t *Task) getDomainName() (domain string, err error) {
	u, _ := url.Parse(t.conf.URL)
	if err != nil {
		return "", err
	}

	parts := strings.Split(u.Hostname(), ".")
	domain = parts[len(parts)-2] + "." + parts[len(parts)-1]
	return
}
