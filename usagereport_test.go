package main

import (
	"errors"

	. "github.com/krujos/usagereport-plugin/fakes"
	. "github.com/onsi/ginkgo"
)

//. "github.com/onsi/gomega"

var _ = Describe("Usagereport", func() {
	Describe("get org errors", func() {
		var fakeAPI *FakeCFApiHelper
		var cmd *UsageReportCmd

		BeforeEach(func() {
			fakeAPI = FakeCFApiHelper{}
			cmd = UsageReportCmd{apiHelper: fakeAPI}
		})

		It("should return an error if cf curl /v2/organizations fails", func() {
			fakeApi.getOrgsReturns(nil, errors.New("Bad Things"))
			_, err := cmd.getOrgs()
			Expect(err).To(BeNil())
		})
	})
})
