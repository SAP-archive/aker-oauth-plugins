package authorization

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/satori/go.uuid"

	"github.infra.hana.ondemand.com/I061150/aker/logging"
	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugin/handler"
	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugin/token"
)

const HeaderOAuthAccessToken = "X-Aker-OAuth-Token-Access-Token"
const HeaderOAuthRefreshToken = "X-Aker-OAuth-Token-Refresh-Token"
const HeaderOAuthTokenType = "X-Aker-OAuth-Token-Type"
const HeaderOAuthTokenExpiry = "X-Aker-OAuth-Token-Expiry"
const HeaderOAuthInfoUserID = "X-Aker-OAuth-Info-User-ID"
const HeaderOAuthInfoUserName = "X-Aker-OAuth-Info-User-Name"
const HeaderOAuthInfoScopes = "X-Aker-OAuth-Info-User-Scopes"

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

	provider := token.NewProvider(tokenConfig, cfg.OAuth.SkipSSLValidation)
	decoder := token.NewDecoder()

	store := sessions.NewCookieStore(
		[]byte(cfg.Session.AuthenticationKey),
		[]byte(cfg.Session.EncryptionKey))

	return NewHandler(provider, decoder, store, cfg.OAuth.RequiredScopes), nil
}

func NewHandler(
	tokenProvider token.Provider,
	tokenDecoder token.Decoder,
	cookieStore sessions.Store,
	requiredScopes []string) http.Handler {

	return context.ClearHandler(&authorizationHandler{
		provider:       tokenProvider,
		decoder:        tokenDecoder,
		store:          cookieStore,
		requiredScopes: requiredScopes,
	})
}

type authorizationHandler struct {
	provider       token.Provider
	decoder        token.Decoder
	store          sessions.Store
	requiredScopes []string
}

func (h *authorizationHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	session, err := h.store.Get(req, handler.SessionName)
	if err != nil {
		logging.Errorf("Failed to retrieve session: %s", err.Error())
	}

	token, ok := h.getToken(session)
	if !ok || !token.Valid() {
		delete(session.Values, "token")

		session.Values["url"] = req.URL.String()

		state := uuid.NewV4().String()
		session.Values["state"] = state

		session.Save(req, resp)

		loginURL := h.provider.LoginURL(state)
		http.Redirect(resp, req, loginURL, http.StatusFound)
		return
	}

	req.Header.Set(HeaderOAuthAccessToken, token.AccessToken)
	req.Header.Set(HeaderOAuthRefreshToken, token.RefreshToken)
	req.Header.Set(HeaderOAuthTokenType, token.TokenType)
	req.Header.Set(HeaderOAuthTokenExpiry, strconv.FormatInt(token.Expiry.UnixNano(), 10))

	info, err := h.decoder.Decode(token)
	if err != nil {
		logging.Errorf("Failed to decode token: %s", err.Error())
		http.Error(resp, "Could not process user's token.", http.StatusInternalServerError)
		return
	}
	req.Header.Set(HeaderOAuthInfoUserID, info.UserID)
	req.Header.Set(HeaderOAuthInfoUserName, info.UserName)
	req.Header[HeaderOAuthInfoScopes] = info.Scope

	if !h.containsAllRequiredScopes(info.Scope) {
		http.Error(resp, "You don't have all required OAuth scopes!", http.StatusForbidden)
		return
	}
}

func (h *authorizationHandler) getToken(session *sessions.Session) (*oauth2.Token, bool) {
	raw, ok := session.Values["token"]
	if !ok {
		return nil, false
	}
	tokenString, ok := raw.(string)
	if !ok {
		logging.Errorf("Token in session is not of type string!")
		return nil, false
	}
	var token oauth2.Token
	if err := json.Unmarshal([]byte(tokenString), &token); err != nil {
		logging.Errorf("Could not unmarshal token from session: %s", err.Error())
		return nil, false
	}
	return &token, true
}

func (h *authorizationHandler) containsAllRequiredScopes(scopes []string) bool {
	for _, requiredScope := range h.requiredScopes {
		hasScope := false
		for _, scope := range scopes {
			if scope == requiredScope {
				hasScope = true
			}
		}
		if !hasScope {
			return false
		}
	}
	return true
}
