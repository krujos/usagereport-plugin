package main

import (
	"errors"

	"github.com/krujos/usagereport-plugin/apihelper"
	"github.com/krujos/usagereport-plugin/apihelper/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Usagereport", func() {
	var fakeAPI *fakes.FakeCFAPIHelper
	var cmd *UsageReportCmd

	BeforeEach(func() {
		fakeAPI = &fakes.FakeCFAPIHelper{}
		cmd = &UsageReportCmd{apiHelper: fakeAPI}
	})

	Describe("get org errors", func() {

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

	Describe("Get org composes the values correctly", func() {
		org := apihelper.Organization{
			URL:      "/v2/orginzations/1234",
			QuotaURL: "/v2/quotas/2345",
		}

		It("should return two one org using 1 mb of 2 mb quota", func() {
			fakeAPI.GetOrgsReturns([]apihelper.Organization{org}, nil)
			fakeAPI.GetOrgMemoryUsageReturns(float64(1), nil)
			fakeAPI.GetQuotaMemoryLimitReturns(float64(2), nil)
			orgs, err := cmd.getOrgs()
			Expect(err).To(BeNil())
			Expect(len(orgs)).To(Equal(1))
			org := orgs[0]
			Expect(org.memoryQuota).To(Equal(float64(2)))
			Expect(org.memoryUsage).To(Equal(float64(1)))
		})
	})
})
