package handler

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	"github.infra.hana.ondemand.com/cloudfoundry/goauth_handlers"
	"github.infra.hana.ondemand.com/cloudfoundry/goauth_handlers/cookie"
	"github.infra.hana.ondemand.com/cloudfoundry/goauth_handlers/session"
	"github.infra.hana.ondemand.com/cloudfoundry/goauth_handlers/token"
	"github.infra.hana.ondemand.com/cloudfoundry/gologger"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func AuthorizationHandlerFromConfig(cfg *Config) (http.Handler, error) {
	return &goauth_handlers.AuthorizationHandler{
		Handler:                HeaderRenameHandler("X-Goauth", "X-Aker"),
		Provider:               buildProviderFromConfig(cfg),
		Decoder:                token.DefaultDecoder,
		Store:                  buildStoreFromConfig(cfg),
		RequiredScopes:         cfg.OAuth.RequiredScopes,
		Logger:                 gologger.DefaultLogger,
		StoreTokenInHeaders:    true,
		StoreUserInfoInHeaders: true,
	}, nil
}

func CallbackHandlerFromConfig(cfg *Config) (http.Handler, error) {
	return &goauth_handlers.CallbackHandler{
		Provider: buildProviderFromConfig(cfg),
		Store:    buildStoreFromConfig(cfg),
		Logger:   gologger.DefaultLogger,
	}, nil
}

func HeaderRenameHandler(oldPrefix, newPrefix string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for key, value := range req.Header {
			if strings.HasPrefix(key, oldPrefix) {
				delete(req.Header, key)
				key = newPrefix + strings.TrimPrefix(key, oldPrefix)
				req.Header[key] = value
			}
		}
	})
}

func buildProviderFromConfig(cfg *Config) goauth_handlers.TokenProvider {
	tokenConfig := oauth2.Config{
		ClientID:     cfg.OAuth.ClientID,
		ClientSecret: cfg.OAuth.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.OAuth.AuthorizationURL,
			TokenURL: cfg.OAuth.TokenURL,
		},
		RedirectURL: cfg.OAuth.RedirectURL,
		Scopes:      append(cfg.OAuth.RequiredScopes, cfg.OAuth.OptionalScopes...),
	}
	return &token.Provider{
		Config:  tokenConfig,
		Context: buildNetContext(cfg.OAuth.SkipSSLValidation),
	}
}

func buildNetContext(skipSSLValidation bool) context.Context {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipSSLValidation,
			},
		},
	}
	return context.WithValue(oauth2.NoContext, oauth2.HTTPClient, client)
}

func buildStoreFromConfig(cfg *Config) session.Store {
	encryptor := cookie.NewEncryptor(
		[]byte(cfg.Session.AuthenticationKey),
		[]byte(cfg.Session.EncryptionKey),
	)
	return cookie.NewStore(encryptor, gologger.DefaultLogger)
}
