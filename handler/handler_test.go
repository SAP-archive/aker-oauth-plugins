package handler_test

import (
	"net/http"
	"net/http/httptest"

	. "github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugins/handler"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HeaderRenameHandler", func() {
	var handler http.Handler
	var request *http.Request
	var response *httptest.ResponseRecorder

	BeforeEach(func() {
		handler = HeaderRenameHandler("X-Goauth", "X-Aker")
		Ω(handler).ShouldNot(BeNil())

		var err error
		request, err = http.NewRequest("GET", "http://some/resource", nil)
		Ω(err).ShouldNot(HaveOccurred())
		request.Header.Set("X-Goauth-Test-Header", "test-header-value")
		request.Header.Set("X-Not-Goauth-Test-Header", "test-other-value")
		response = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		handler.ServeHTTP(response, request)
	})

	It("should rename headers with specified prefix", func() {
		Ω(request.Header).ShouldNot(HaveKey("X-Goauth-Test-Header"))
		Ω(request.Header).Should(HaveKey("X-Aker-Test-Header"))
		Ω(request.Header.Get("X-Aker-Test-Header")).Should(Equal("test-header-value"))
	})

	It("should not rename headers with other prefix", func() {
		Ω(request.Header).Should(HaveKey("X-Not-Goauth-Test-Header"))
		Ω(request.Header.Get("X-Not-Goauth-Test-Header")).Should(Equal("test-other-value"))
	})
})
