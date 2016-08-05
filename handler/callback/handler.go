package callback

import (
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugins/handler"
	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugins/token"
	"github.infra.hana.ondemand.com/cloudfoundry/aker/logging"
)

func HandlerFromConfig(cfg *handler.Config) (http.Handler, error) {
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

	return NewHandler(
		token.NewProvider(tokenConfig, cfg.OAuth.SkipSSLValidation),
		sessions.NewCookieStore(
			[]byte(cfg.Session.AuthenticationKey),
			[]byte(cfg.Session.EncryptionKey)),
	), nil
}

func NewHandler(tokenProvider token.Provider, sessionStore sessions.Store) http.Handler {
	return context.ClearHandler(&callbackHandler{
		provider: tokenProvider,
		store:    sessionStore,
	})
}

type callbackHandler struct {
	provider token.Provider
	store    sessions.Store
}

func (h *callbackHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	session, _ := h.store.Get(req, "aker-oauth-authorization")

	errorParam := req.FormValue("error")
	if errorParam != "" {
		session.Values = nil
		session.Save(req, resp)
		if errorParam == "invalid_scope" {
			http.Error(resp, "You do not have the required authorization!", http.StatusForbidden)
		} else {
			logging.Errorf("UAA returned an error authorization grant response '%s'.", errorParam)
			http.Error(resp, "We could not authorize you due to an internal error!", http.StatusInternalServerError)
		}
		return
	}

	state := req.FormValue("state")
	if state == "" {
		session.Values = nil
		session.Save(req, resp)
		http.Error(resp, "Missing state parameter!", http.StatusBadRequest)
		return
	}

	oauthCode := req.FormValue("code")
	if oauthCode == "" {
		session.Values = nil
		session.Save(req, resp)
		http.Error(resp, "Missing code parameter!", http.StatusBadRequest)
		return
	}

	originalURL, ok := h.getOriginalURL(session)
	if !ok {
		session.Values = nil
		session.Save(req, resp)
		http.Error(resp, "Missing redirect URL!", http.StatusBadRequest)
		return
	}
	delete(session.Values, "url")

	expectedState, ok := h.getState(session)
	if !ok {
		session.Values = nil
		session.Save(req, resp)
		http.Error(resp, "Missing state parameter in session!", http.StatusBadRequest)
		return
	}

	if state != expectedState {
		session.Values = nil
		session.Save(req, resp)
		http.Error(resp, "Invalid state parameter!", http.StatusBadRequest)
		return
	}
	delete(session.Values, "state")

	oauthToken, err := h.provider.RequestToken(oauthCode)
	if err != nil {
		session.Values = nil
		session.Save(req, resp)
		logging.Errorf("Could not retrieve Token from UAA due to '%s'.", err)
		http.Error(resp, "Could not retrieve token", http.StatusInternalServerError)
		return
	}

	byteToken, err := json.Marshal(oauthToken)
	if err != nil {
		panic(err)
	}
	session.Values["token"] = string(byteToken)
	session.Save(req, resp)

	http.Redirect(resp, req, originalURL, http.StatusFound)
}

func (h *callbackHandler) getOriginalURL(session *sessions.Session) (string, bool) {
	raw, ok := session.Values["url"]
	if !ok {
		return "", false
	}
	url, ok := raw.(string)
	if !ok {
		return "", false
	}
	return url, true
}

func (h *callbackHandler) getState(session *sessions.Session) (string, bool) {
	raw, ok := session.Values["state"]
	if !ok {
		return "", false
	}
	state, ok := raw.(string)
	if !ok {
		return "", false
	}
	return state, true
}
