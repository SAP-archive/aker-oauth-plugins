package authorization_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.infra.hana.ondemand.com/I061150/aker/logging"

	"testing"
)

func TestAuthorization(t *testing.T) {
	logging.DefaultLogger = logging.NewNativeLogger(GinkgoWriter, GinkgoWriter)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Authorization Suite")
}