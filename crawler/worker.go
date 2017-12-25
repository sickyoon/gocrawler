package crawler

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/net/publicsuffix"
)

// Worker -
type Worker struct {
	id        int
	reqCh     chan Request
	client    *http.Client
	cookieJar *cookiejar.Jar
}

// NewWorker -
func NewWorker(id int, reqCh chan Request) (*Worker, error) {
	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	// TODO: create cookiejar
	// TODO: create client
	client := &http.Client{
		Transport: &http.Transport{
			//Proxy: nil, // TODO: proxy support
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(viper.GetInt("tcp_timeout")) * time.Second,
				KeepAlive: time.Duration(viper.GetInt("keepalive")) * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:        viper.GetInt("max_conn"),
			IdleConnTimeout:     time.Duration(viper.GetInt("idle_timout")) * time.Second,
			TLSHandshakeTimeout: time.Duration(viper.GetInt("tls_timeout")) * time.Second,
		},
		CheckRedirect: nil,
		Jar:           cookieJar,
		Timeout:       time.Duration(viper.GetInt("client_timeout")) * time.Second,
	}
	return &Worker{
		id:        id,
		reqCh:     reqCh,
		client:    client,
		cookieJar: cookieJar,
	}, nil
}

// Start -
func (w *Worker) Start() {
	for {
		select {
		case req := <-w.reqCh:
			w.query(req)
		}
	}
}

func (w *Worker) query(creq Request) error {
	req, err := http.NewRequest(creq.Method(), creq.URL(), nil)
	if err != nil {
		return err
	}
	// TODO: add headers
	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	creq.Process(w.reqCh, body)
	return nil
}
