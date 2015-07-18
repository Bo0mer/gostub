package acceptance_test

import (
	. "github.com/momchil-atanasov/gostub/acceptance"
	"github.com/momchil-atanasov/gostub/acceptance/acceptance_stubs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TypeStruct", func() {
	var stub *acceptance_stubs.StructSupportStub

	BeforeEach(func() {
		stub = new(acceptance_stubs.StructSupportStub)
	})

	It("stub is assignable to interface", func() {
		_, assignable := interface{}(stub).(StructSupport)
		Ω(assignable).Should(BeTrue())
	})
})
