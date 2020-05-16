package account

import (
	"context"
	"net/http"
	"strings"

	"qiniupkg.com/x/rpc.v7"
)

type AccountService interface {
	Me(ctx context.Context) (*UsersResult, error)
}

type Config struct {
	Host      string
	Transport http.RoundTripper
	Dev       bool
}

type Client struct {
	config Config
	client rpc.Client
}

var _ AccountService = new(Client)

func New(cfg Config) *Client {
	cfg.Host = cleanHost(cfg.Host)
	if cfg.Transport == nil {
		cfg.Transport = http.DefaultTransport
	}
	p := &Client{config: cfg}

	p.client = rpc.Client{&http.Client{Transport: cfg.Transport}}
	return p
}

func cleanHost(host string) string {
	for strings.HasSuffix(host, "/") {
		host = strings.TrimSuffix(host, "/")
	}

	if !strings.HasPrefix(host, "http") {
		host = "http://" + host
	}
	return host
}
