package callback_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	"golang.org/x/oauth2"

	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugin/handler"
	. "github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugin/handler/callback"
	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugin/session/sessionfakes"
	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugin/token/tokenfakes"

	"github.com/gorilla/sessions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	const unsetResponseCode = -1
	const originalURL = "http://some/resource/somewhere"
	const oauthState = "SOME_STATE_VALUE"
	const oauthCode = "SOME_OAUTH_CODE"

	var oauthToken *oauth2.Token
	var tokenProvider *tokenfakes.FakeProvider

	var sessionStore *sessionfakes.FakeStore
	var session *sessions.Session

	var callbackHandler http.Handler

	var request *http.Request
	var response *httptest.ResponseRecorder

	BeforeEach(func() {
		sessionStore = new(sessionfakes.FakeStore)
		session = sessions.NewSession(sessionStore, handler.SessionName)
		session.Values["url"] = originalURL
		session.Values["state"] = oauthState
		sessionStore.GetReturns(session, nil)

		oauthToken = &oauth2.Token{
			AccessToken:  "SomeAccessToken",
			RefreshToken: "SomeRefreshToken",
			TokenType:    "bearer",
			Expiry:       time.Now().Add(time.Hour),
		}
		tokenProvider = new(tokenfakes.FakeProvider)
		tokenProvider.RequestTokenReturns(oauthToken, nil)

		callbackHandler = NewHandler(tokenProvider, sessionStore)
		Ω(callbackHandler).ShouldNot(BeNil())

		var err error
		request, err = http.NewRequest("GET", "http://some/resource", nil)
		Ω(err).ShouldNot(HaveOccurred())
		request.Form = make(map[string][]string)
		request.Form.Add("state", oauthState)
		request.Form.Add("code", oauthCode)

		response = httptest.NewRecorder()
		response.Code = unsetResponseCode
	})

	JustBeforeEach(func() {
		callbackHandler.ServeHTTP(response, request)
	})

	Context("when everything works as expected", func() {
		It("should get a session", func() {
			Ω(sessionStore.GetCallCount()).Should(Equal(1))
			argRequest, argName := sessionStore.GetArgsForCall(0)
			Ω(argRequest).Should(Equal(request))
			Ω(argName).Should(Equal(handler.SessionName))
		})

		It("should remove state parameter from session", func() {
			Ω(session.Values).ShouldNot(HaveKey("state"))
		})

		It("should remove original url from session", func() {
			Ω(session.Values).ShouldNot(HaveKey("url"))
		})

		It("should request a token", func() {
			Ω(tokenProvider.RequestTokenCallCount()).Should(Equal(1))
			argCode := tokenProvider.RequestTokenArgsForCall(0)
			Ω(argCode).Should(Equal(oauthCode))
		})

		It("should store token in session", func() {
			tokenObj, exists := session.Values["token"]
			Ω(exists).Should(BeTrue())
			tokenString, isString := tokenObj.(string)
			Ω(isString).Should(BeTrue())
			var token oauth2.Token
			err := json.Unmarshal([]byte(tokenString), &token)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(token).Should(Equal(*oauthToken))
		})

		It("should save the session", func() {
			Ω(sessionStore.SaveCallCount()).Should(Equal(1))
			argRequest, argResponse, argSession := sessionStore.SaveArgsForCall(0)
			Ω(argRequest).Should(Equal(request))
			Ω(argResponse).Should(Equal(response))
			Ω(argSession).Should(Equal(session))
		})

		It("should redirect to original address", func() {
			Ω(response.Code).Should(Equal(http.StatusFound))
			Ω(response.Header().Get("Location")).Should(Equal(originalURL))
		})
	})

	itShouldDeleteSession := func() {
		It("should delete session", func() {
			Ω(session.Values).Should(BeNil())
			Ω(sessionStore.SaveCallCount()).Should(Equal(1))
			argRequest, argResponse, argSession := sessionStore.SaveArgsForCall(0)
			Ω(argRequest).Should(Equal(request))
			Ω(argResponse).Should(Equal(response))
			Ω(argSession).Should(Equal(session))
		})
	}

	itShouldReturnBadRequest := func() {
		It("should return bad request", func() {
			Ω(response.Code).Should(Equal(http.StatusBadRequest))
		})
	}

	itShouldReturnInternalServerError := func() {
		It("should return bad request", func() {
			Ω(response.Code).Should(Equal(http.StatusInternalServerError))
		})
	}

	itShouldReturnForbidden := func() {
		It("should return forbidden", func() {
			Ω(response.Code).Should(Equal(http.StatusForbidden))
		})
	}

	Context("when UAA responds with invalid_scope error param", func() {
		BeforeEach(func() {
			request.Form.Add("error", "invalid_scope")
		})

		itShouldDeleteSession()

		itShouldReturnForbidden()
	})

	Context("when UAA responds with error param other than invalid_scope", func() {
		BeforeEach(func() {
			request.Form.Add("error", "unauthorized_client")
		})

		itShouldDeleteSession()

		itShouldReturnInternalServerError()
	})

	Context("when state parameter in request is missing", func() {
		BeforeEach(func() {
			request.Form.Del("state")
		})

		itShouldDeleteSession()

		itShouldReturnBadRequest()
	})

	Context("when a oauth code is missing in the request", func() {
		BeforeEach(func() {
			request.Form.Del("code")
		})

		itShouldDeleteSession()

		itShouldReturnBadRequest()
	})

	Context("when original URL parameter is missing from session", func() {
		BeforeEach(func() {
			delete(session.Values, "url")
		})

		itShouldDeleteSession()

		itShouldReturnBadRequest()
	})

	Context("when original URL parameter in session is not string", func() {
		BeforeEach(func() {
			session.Values["url"] = 1
		})

		itShouldDeleteSession()

		itShouldReturnBadRequest()
	})

	Context("when state parameter is missing from session", func() {
		BeforeEach(func() {
			delete(session.Values, "state")
		})

		itShouldDeleteSession()

		itShouldReturnBadRequest()
	})

	Context("when state parameter in session is not string", func() {
		BeforeEach(func() {
			session.Values["state"] = 1
		})

		itShouldDeleteSession()

		itShouldReturnBadRequest()
	})

	Context("when state parameter in request does not equal the session one", func() {
		BeforeEach(func() {
			session.Values["state"] = "some_other_value"
		})

		itShouldDeleteSession()

		itShouldReturnBadRequest()
	})

	Context("when provider fails to retrieve token", func() {
		BeforeEach(func() {
			tokenProvider.RequestTokenReturns(nil, errors.New("Could not get token"))
		})

		itShouldDeleteSession()

		itShouldReturnInternalServerError()
	})

})
