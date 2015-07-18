package acceptance_test

import (
	. "github.com/momchil-atanasov/gostub/acceptance"
	"github.com/momchil-atanasov/gostub/acceptance/acceptance_stubs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmbeddedEmbeddedInterface", func() {
	var stub *acceptance_stubs.EmbeddedEmbeddedInterfaceSupportStub

	BeforeEach(func() {
		stub = new(acceptance_stubs.EmbeddedEmbeddedInterfaceSupportStub)
	})

	It("stub is assignable to interface", func() {
		_, assignable := interface{}(stub).(EmbeddedEmbeddedInterfaceSupport)
		Ω(assignable).Should(BeTrue())
	})
})
