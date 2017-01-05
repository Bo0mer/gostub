package acceptance_test

import (
	. "github.com/mokiat/gostub/acceptance"
	"github.com/mokiat/gostub/acceptance/acceptance_stubs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExternalEmbeddedInterface", func() {
	var stub *acceptance_stubs.ExternalEmbeddedInterfaceSupportStub

	BeforeEach(func() {
		stub = new(acceptance_stubs.ExternalEmbeddedInterfaceSupportStub)
	})

	It("stub is assignable to interface", func() {
		_, assignable := interface{}(stub).(ExternalEmbeddedInterfaceSupport)
		Ω(assignable).Should(BeTrue())
	})
})
