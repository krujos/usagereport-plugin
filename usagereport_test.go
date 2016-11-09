package main

import (
	"errors"

	"github.com/krujos/usagereport-plugin/apihelper"
	"github.com/krujos/usagereport-plugin/apihelper/fakes"
	"io/ioutil"

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

	Describe("get single org errors", func() {
		It("should return an error if cf curl /v2/organizations fails", func() {
			fakeAPI.GetOrgReturns(apihelper.Organization{}, errors.New("Bad Things"))
			_, err := cmd.getOrg("test")
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("when showing output", func() {
		var orgs []org

		BeforeEach(func() {
			orgs = []org{
				org{
					name:        "test-org",
					memoryQuota: 4096,
					spaces: []space{space{
						name: "test-space",
						apps: []app{
							app{ram: 128, instances: 2, running: true},
							app{ram: 128, instances: 2, running: false},
						},
					},
					},
				},
			}
		})

		It("should output in expected default human readable format", func() {
		})

		It("should output in csv format", func() {
			expectedOutput, err := ioutil.ReadFile("fixtures/result.csv")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(cmd.printOrgsCSV(orgs)).To(Equal(string(expectedOutput)))
		})
	})

	Describe("get orgs errors", func() {

		It("should return an error if cf curl /v2/organizations fails", func() {
			fakeAPI.GetOrgsReturns(nil, errors.New("Bad Things"))
			_, err := cmd.getOrgs()
			Expect(err).ToNot(BeNil())
		})

		Context("good org bad other thigns", func() {
			BeforeEach(func() {
				fakeAPI.GetOrgsReturns([]apihelper.Organization{apihelper.Organization{}}, nil)
			})

			It("should return an error if cf curl /v2/organizations/{guid}/memory_usage fails", func() {
				fakeAPI.GetOrgMemoryUsageReturns(0, errors.New("Bad Things"))
				_, err := cmd.getOrgs()
				Expect(err).ToNot(BeNil())
			})

			It("sholud return an error if cf curl to the quota url fails", func() {
				fakeAPI.GetQuotaMemoryLimitReturns(0, errors.New("Bad Things"))
				_, err := cmd.getOrgs()
				Expect(err).ToNot(BeNil())
			})

			It("should return an error if cf curl to get org spaces fails", func() {
				fakeAPI.GetOrgSpacesReturns(nil, errors.New("Bad Things"))
				_, err := cmd.getOrgs()
				Expect(err).ToNot(BeNil())
				Expect(fakeAPI.GetOrgSpacesCallCount()).To(Equal(1))
			})

			It("Should return an error if cf curl to get the apps in a space fails", func() {
				fakeAPI.GetOrgSpacesReturns(
					[]apihelper.Space{apihelper.Space{AppsURL: "/v2/apps"}}, nil)
				fakeAPI.GetSpaceAppsReturns(nil, errors.New("Bad Things"))
				_, err := cmd.getOrgs()
				Expect(err).ToNot(BeNil())
				Expect(fakeAPI.GetSpaceAppsCallCount()).To(Equal(1))
			})
		})

	})

	Describe("Get org composes the values correctly", func() {
		org := apihelper.Organization{
			URL:      "/v2/organizations/1234",
			QuotaURL: "/v2/quotas/2345",
		}

		BeforeEach(func() {
			fakeAPI.GetOrgsReturns([]apihelper.Organization{org}, nil)
		})

		It("should return two one org using 1 mb of 2 mb quota", func() {
			fakeAPI.GetOrgMemoryUsageReturns(float64(1), nil)
			fakeAPI.GetQuotaMemoryLimitReturns(float64(2), nil)
			orgs, err := cmd.getOrgs()
			Expect(err).To(BeNil())
			Expect(len(orgs)).To(Equal(1))
			org := orgs[0]
			Expect(org.memoryQuota).To(Equal(2))
			Expect(org.memoryUsage).To(Equal(1))
		})

		It("Should return an org with 1 space", func() {
			fakeAPI.GetOrgSpacesReturns(
				[]apihelper.Space{apihelper.Space{}, apihelper.Space{}}, nil)
			orgs, _ := cmd.getOrgs()
			Expect(len(orgs[0].spaces)).To(Equal(2))
		})

		It("Should not choke on an org with no spaces", func() {
			fakeAPI.GetOrgSpacesReturns(
				[]apihelper.Space{}, nil)
			orgs, _ := cmd.getOrgs()
			Expect(len(orgs[0].spaces)).To(Equal(0))
		})

		It("Should return two apps from a space", func() {
			fakeAPI.GetOrgSpacesReturns(
				[]apihelper.Space{apihelper.Space{}}, nil)

			fakeAPI.GetSpaceAppsReturns(
				[]apihelper.App{
					apihelper.App{},
					apihelper.App{},
					apihelper.App{},
				},
				nil)
			orgs, _ := cmd.getOrgs()
			org := orgs[0]
			space := org.spaces[0]
			apps := space.apps
			Expect(len(apps)).To(Equal(3))
		})

		It("Should mark the first app as running, the second as stopped", func() {
			fakeAPI.GetOrgSpacesReturns(
				[]apihelper.Space{apihelper.Space{}}, nil)

			fakeAPI.GetSpaceAppsReturns(
				[]apihelper.App{
					apihelper.App{Running: true},
					apihelper.App{Running: false},
				},
				nil)

			orgs, _ := cmd.getOrgs()
			org := orgs[0]
			space := org.spaces[0]
			apps := space.apps
			Expect(apps[0].running).To(BeTrue())
			Expect(apps[1].running).To(BeFalse())
		})
	})
})
