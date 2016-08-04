package token

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

//go:generate counterfeiter . Provider

type Provider interface {
	RequestToken(code string) (*oauth2.Token, error)
	LoginURL(state string) string
}

type provider struct {
	config          oauth2.Config
	exchangeContext context.Context
}

func NewProvider(config oauth2.Config, skipSSLValidation bool) Provider {
	client := &http.Client{Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipSSLValidation,
		},
	}}
	return &provider{
		config: config,
		exchangeContext: context.WithValue(
			oauth2.NoContext,
			oauth2.HTTPClient,
			client,
		),
	}
}

func (p *provider) RequestToken(code string) (*oauth2.Token, error) {
	return p.config.Exchange(p.exchangeContext, code)
}

func (p *provider) LoginURL(state string) string {
	return p.config.AuthCodeURL(state)
}
