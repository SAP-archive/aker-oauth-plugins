package authorization_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"golang.org/x/oauth2"

	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugins/handler"
	. "github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugins/handler/authorization"
	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugins/session/sessionfakes"
	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugins/token"
	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugins/token/tokenfakes"

	"github.com/gorilla/sessions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	const loginURL = "http://login.here"
	const unsetResponseCode = -1

	var oauthToken *oauth2.Token
	var oauthTokenInfo token.Info

	var tokenProvider *tokenfakes.FakeProvider

	var tokenDecoder *tokenfakes.FakeDecoder

	var sessionStore *sessionfakes.FakeStore
	var session *sessions.Session

	var authHandler http.Handler

	var request *http.Request
	var response *httptest.ResponseRecorder

	tokenToString := func(token *oauth2.Token) string {
		bytes, err := json.Marshal(token)
		Ω(err).ShouldNot(HaveOccurred())
		return string(bytes)
	}

	BeforeEach(func() {
		oauthToken = &oauth2.Token{
			AccessToken:  "SomeAccessToken",
			RefreshToken: "SomeRefreshToken",
			TokenType:    "bearer",
			Expiry:       time.Now().Add(time.Hour),
		}

		sessionStore = new(sessionfakes.FakeStore)
		session = sessions.NewSession(sessionStore, handler.SessionName)
		session.Values["token"] = tokenToString(oauthToken)
		sessionStore.GetReturns(session, nil)

		tokenProvider = new(tokenfakes.FakeProvider)
		tokenProvider.LoginURLReturns(loginURL)

		oauthTokenInfo = token.Info{
			UserID:   "012345",
			UserName: "Patrick",
			Scope:    []string{"logs", "messages"},
		}
		tokenDecoder = new(tokenfakes.FakeDecoder)
		tokenDecoder.DecodeReturns(oauthTokenInfo, nil)

		authHandler = NewHandler(tokenProvider, tokenDecoder, sessionStore, []string{"logs", "messages"})
		Ω(authHandler).ShouldNot(BeNil())

		var err error
		request, err = http.NewRequest("GET", "http://some/resource", nil)
		Ω(err).ShouldNot(HaveOccurred())
		response = httptest.NewRecorder()
		response.Code = unsetResponseCode
	})

	JustBeforeEach(func() {
		authHandler.ServeHTTP(response, request)
	})

	Context("when everything goes as planned", func() {
		It("should get a session", func() {
			Ω(sessionStore.GetCallCount()).Should(Equal(1))
			argRequest, argName := sessionStore.GetArgsForCall(0)
			Ω(argRequest).Should(Equal(request))
			Ω(argName).Should(Equal(handler.SessionName))
		})

		It("should save the token inside the headers", func() {
			Ω(request.Header.Get(HeaderOAuthAccessToken)).Should(Equal(oauthToken.AccessToken))
			Ω(request.Header.Get(HeaderOAuthRefreshToken)).Should(Equal(oauthToken.RefreshToken))
			Ω(request.Header.Get(HeaderOAuthTokenType)).Should(Equal(oauthToken.TokenType))
			Ω(request.Header.Get(HeaderOAuthTokenExpiry)).Should(Equal(strconv.FormatInt(oauthToken.Expiry.UnixNano(), 10)))
		})

		It("should decode the token", func() {
			Ω(tokenDecoder.DecodeCallCount()).Should(Equal(1))
			argToken := tokenDecoder.DecodeArgsForCall(0)
			Ω(argToken).Should(Equal(oauthToken))
		})

		It("should save token info inside the headers", func() {
			Ω(request.Header.Get(HeaderOAuthInfoUserID)).Should(Equal(oauthTokenInfo.UserID))
			Ω(request.Header.Get(HeaderOAuthInfoUserName)).Should(Equal(oauthTokenInfo.UserName))
			Ω(request.Header[HeaderOAuthInfoScopes]).Should(Equal(oauthTokenInfo.Scope))
		})

		It("should not have written output", func() {
			Ω(response.Code).Should(Equal(unsetResponseCode))
			Ω(response.Body.Len()).Should(Equal(0))
		})
	})

	itShouldRedirectToLogin := func() {
		It("should remove token from session", func() {
			Ω(session.Values).ShouldNot(HaveKey("token"))
		})

		It("should store the request URL in the session", func() {
			Ω(session.Values["url"]).Should(Equal(request.URL.String()))
		})

		It("should store some random state parameter in the session", func() {
			state, exists := session.Values["state"]
			Ω(exists).Should(BeTrue())
			_, isString := state.(string)
			Ω(isString).Should(BeTrue())
		})

		It("should save the session", func() {
			Ω(sessionStore.SaveCallCount()).Should(Equal(1))
			argRequest, argResponse, argSession := sessionStore.SaveArgsForCall(0)
			Ω(argRequest).Should(Equal(request))
			Ω(argResponse).Should(Equal(response))
			Ω(argSession).Should(Equal(session))
		})

		It("should get login URL from provider", func() {
			Ω(tokenProvider.LoginURLCallCount()).Should(Equal(1))
			argState := tokenProvider.LoginURLArgsForCall(0)
			Ω(argState).Should(Equal(session.Values["state"]))
		})

		It("should redirect to login URL", func() {
			Ω(response.Code).Should(Equal(http.StatusFound))
			Ω(response.Header().Get("Location")).Should(Equal(loginURL))
		})
	}

	Context("when session does not contain token", func() {
		BeforeEach(func() {
			delete(session.Values, "token")
		})

		itShouldRedirectToLogin()
	})

	Context("when token in session is not string", func() {
		BeforeEach(func() {
			session.Values["token"] = 1
		})

		itShouldRedirectToLogin()
	})

	Context("when token in session is not a valid json", func() {
		BeforeEach(func() {
			session.Values["token"] = "{"
		})

		itShouldRedirectToLogin()
	})

	Context("when token is expired", func() {
		BeforeEach(func() {
			oauthToken.Expiry = time.Now().Add(-time.Hour)
			session.Values["token"] = tokenToString(oauthToken)
		})

		itShouldRedirectToLogin()
	})

	Context("when token info decoding fails", func() {
		BeforeEach(func() {
			tokenDecoder.DecodeReturns(token.Info{}, errors.New("Could not parse token!"))
		})

		It("should not save token info inside headers", func() {
			Ω(request.Header).ShouldNot(HaveKey(HeaderOAuthInfoUserID))
			Ω(request.Header).ShouldNot(HaveKey(HeaderOAuthInfoUserName))
			Ω(request.Header).ShouldNot(HaveKey(HeaderOAuthInfoScopes))
		})

		It("should return internal server error", func() {
			Ω(response.Code).Should(Equal(http.StatusInternalServerError))
		})
	})

	Context("when token does not contain all required scopes", func() {
		BeforeEach(func() {
			oauthTokenInfo.Scope = []string{"messages"}
			tokenDecoder.DecodeReturns(oauthTokenInfo, nil)
		})

		It("should return forbidden", func() {
			Ω(response.Code).Should(Equal(http.StatusForbidden))
		})
	})
})
