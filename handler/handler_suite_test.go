package handler_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/SAP/gologger"

	"testing"
)

func TestHandler(t *testing.T) {
	gologger.DefaultLogger = gologger.NewNativeLogger(GinkgoWriter, GinkgoWriter)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}
