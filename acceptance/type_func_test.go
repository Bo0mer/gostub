package acceptance_test

import (
	. "github.com/mokiat/gostub/acceptance"
	"github.com/mokiat/gostub/acceptance/acceptance_stubs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TypeFunc", func() {
	var stub *acceptance_stubs.FuncSupportStub

	BeforeEach(func() {
		stub = new(acceptance_stubs.FuncSupportStub)
	})

	It("stub is assignable to interface", func() {
		_, assignable := interface{}(stub).(FuncSupport)
		Ω(assignable).Should(BeTrue())
	})
})
