package crawler

import (
	"log"
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

	// TODO: create cookiejar
	cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

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

	log.Printf("starting worker %d", w.id)

	for {
		select {
		case req := <-w.reqCh:
			if err := w.query(req); err != nil {
				// TODO: stop and fail or continue?
			}
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
	err = creq.Process(w.reqCh, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
