package handler_test

import (
	. "github.com/SAP/aker-oauth-plugins/handler"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	Describe("ParseConfig", func() {

		var configData []byte

		var cfg *Config
		var err error

		JustBeforeEach(func() {
			cfg, err = ParseConfig(configData)
		})

		Context("when the configuration is invalid", func() {

			BeforeEach(func() {
				configData = []byte("broken configuration")
			})

			It("should return an error", func() {
				Ω(cfg).Should(BeNil())
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("when the authentication key has wrong length", func() {

			BeforeEach(func() {
				configData = []byte(`
---
session:
  authentication_key: "too_short_to_be_valid"
`)
			})

			It("should return an error", func() {
				Ω(cfg).Should(BeNil())
				Ω(err).Should(HaveOccurred())
				_, ok := err.(*InvalidKeyLengthError)
				Ω(ok).Should(BeTrue())
			})
		})

		Context("when the encryption key has wrong length", func() {

			BeforeEach(func() {
				configData = []byte(`
---
session:
  authentication_key: "this_key_should_be_32_bytes_long"
  encryption_key: "too_short_to_be_valid"
`)
			})

			It("should return an error", func() {
				Ω(cfg).Should(BeNil())
				Ω(err).Should(HaveOccurred())
				_, ok := err.(*InvalidKeyLengthError)
				Ω(ok).Should(BeTrue())
			})
		})

		Context("when the configuration is valid", func() {

			BeforeEach(func() {
				configData = []byte(`
---
oauth:
  client_id: id
  client_secret: secret
  skip_ssl_validation: true
  authorization_url: http://auth.com
  token_url: http://token.com
  redirect_url: http://redirect.com
  required_scopes: ["first_required", "second_required"]
  optional_scopes: ["optional"]
session:
  authentication_key: this_key_should_be_32_bytes_long
  encryption_key: this_key_should_be_32_bytes_long
`)
			})

			It("should return no error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should return a Config object with valid OAuthConfig", func() {
				Ω(cfg.OAuth.ClientID).Should(Equal("id"))
				Ω(cfg.OAuth.ClientSecret).Should(Equal("secret"))
				Ω(cfg.OAuth.SkipSSLValidation).Should(BeTrue())
				Ω(cfg.OAuth.AuthorizationURL).Should(Equal("http://auth.com"))
				Ω(cfg.OAuth.TokenURL).Should(Equal("http://token.com"))
				Ω(cfg.OAuth.RedirectURL).Should(Equal("http://redirect.com"))
				Ω(cfg.OAuth.RequiredScopes).Should(Equal([]string{"first_required", "second_required"}))
				Ω(cfg.OAuth.OptionalScopes).Should(Equal([]string{"optional"}))

			})

			It("should return a Config object with valid SessionConfig", func() {
				Ω(cfg.Session.AuthenticationKey).Should(Equal("this_key_should_be_32_bytes_long"))
				Ω(cfg.Session.EncryptionKey).Should(Equal("this_key_should_be_32_bytes_long"))
			})
		})
	})
})
