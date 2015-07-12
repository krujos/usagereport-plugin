package main

import (
	"errors"

	"github.com/krujos/usagereport-plugin/apihelper"
	"github.com/krujos/usagereport-plugin/apihelper/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Usagereport", func() {
	Describe("get org errors", func() {
		var fakeAPI *fakes.FakeCFAPIHelper
		var cmd *UsageReportCmd

		BeforeEach(func() {
			fakeAPI = &fakes.FakeCFAPIHelper{}
			cmd = &UsageReportCmd{apiHelper: fakeAPI}
		})

		It("should return an error if cf curl /v2/organizations fails", func() {
			fakeAPI.GetOrgsReturns(nil, errors.New("Bad Things"))
			_, err := cmd.getOrgs()
			Expect(err).ToNot(BeNil())
		})

		It("should return an error if cf curl /v2/organizations/{guid}/memory_usage fails", func() {
			fakeAPI.GetOrgsReturns([]apihelper.Organization{apihelper.Organization{}}, nil)
			fakeAPI.GetOrgMemoryUsageReturns(0, errors.New("Bad Things"))
			_, err := cmd.getOrgs()
			Expect(err).ToNot(BeNil())
		})

		It("sholud return an error if cf curl to the quota url fails", func() {
			fakeAPI.GetOrgsReturns([]apihelper.Organization{apihelper.Organization{}}, nil)
			fakeAPI.GetOrgMemoryUsageReturns(float64(1024), nil)
			fakeAPI.GetQuotaMemoryLimitReturns(0, errors.New("Bad Things"))
			_, err := cmd.getOrgs()
			Expect(err).ToNot(BeNil())
		})
	})
})
