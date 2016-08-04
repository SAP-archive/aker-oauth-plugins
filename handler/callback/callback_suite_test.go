package callback_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.infra.hana.ondemand.com/I061150/aker/logging"

	"testing"
)

func TestCallback(t *testing.T) {
	logging.DefaultLogger = logging.NewNativeLogger(GinkgoWriter, GinkgoWriter)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Callback Suite")
}
