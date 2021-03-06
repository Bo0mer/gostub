package acceptance_test

import (
	. "github.com/mokiat/gostub/acceptance"
	"github.com/mokiat/gostub/acceptance/acceptance_stubs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AliasedRefSupport", func() {
	var stub *acceptance_stubs.AliasedRefSupportStub

	BeforeEach(func() {
		stub = new(acceptance_stubs.AliasedRefSupportStub)
	})

	It("stub is assignable to interface", func() {
		_, assignable := interface{}(stub).(AliasedRefSupport)
		Ω(assignable).Should(BeTrue())
	})
})
