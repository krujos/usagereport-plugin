package apihelper

import (
	"bufio"
	"errors"
	"os"

	"github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func slurp(filename string) []string {
	var b []string
	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		b = append(b, scanner.Text())
	}
	return b
}

var _ = Describe("UsageReport", func() {
	var api CFAPIHelper
	var fakeCliConnection *fakes.FakeCliConnection

	BeforeEach(func() {
		fakeCliConnection = &fakes.FakeCliConnection{}
		api = New(fakeCliConnection)
	})

	Describe("Get orgs", func() {
		var orgsJSON []string

		BeforeEach(func() {
			orgsJSON = slurp("test-data/orgs.json")
		})

		It("should return two orgs", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			orgs, _ := api.GetOrgs()
			Expect(len(orgs)).To(Equal(2))
		})

		It("does something intellegent when cf curl fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				nil, errors.New("bad things"))
			_, err := api.GetOrgs()
			Expect(err).ToNot(BeNil())
		})

		It("populates the url", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			orgs, _ := api.GetOrgs()
			org := orgs[0]
			Expect(org.URL).To(Equal("/v2/organizations/b1a23fd6-ac8d-4304-a3b4-815745417acd"))
		})

		It("calls /v2/orgs", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			api.GetOrgs()
			args := fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)
			Expect(args[1]).To(Equal("/v2/organizations"))
		})

	})

	Describe("paged org output", func() {
		var orgsPage1 []string

		BeforeEach(func() {
			orgsPage1 = slurp("test-data/paged-orgs-page-1.json")
		})

		It("deals with paged output", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsPage1, nil)
			api.GetOrgs()
			args := fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)
			Expect(args[1]).To(Equal("/v2/organizations"))
			Ω(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(2))
		})

		It("Should have 100 orgs", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsPage1, nil)
			orgs, _ := api.GetOrgs()
			args := fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(1)
			Expect(args[1]).To(Equal("/v2/organizations?page=2"))
			Ω(orgs).To(HaveLen(100))
		})
	})

	Describe("Get quota memory limit", func() {
		var quotaJSON []string

		BeforeEach(func() {
			quotaJSON = slurp("test-data/quota.json")
		})

		It("should return an error when it can't fetch the memory limit", func() {
			_, err := api.GetQuotaMemoryLimit("/v2/somequota")
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				nil, errors.New("Bad Things"))
			Expect(err).ToNot(BeNil())
		})

		It("should reutrn 10240 as the memory limit", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				quotaJSON, nil)
			limit, _ := api.GetQuotaMemoryLimit("/v2/quotas/")
			Expect(limit).To(Equal(float64(10240)))
		})
	})

	Describe("it Gets the org memory usage", func() {
		var org Organization
		var usageJSON []string

		BeforeEach(func() {
			usageJSON = slurp("test-data/memory_usage.json")
		})

		It("should return an error when it can't fetch the orgs memory usage", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil,
				errors.New("Bad things"))
			_, err := api.GetOrgMemoryUsage(org)
			Expect(err).ToNot(BeNil())
		})

		It("should return the memory usage", func() {
			org.URL = "/v2/organizations/1234"
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(usageJSON, nil)
			usage, _ := api.GetOrgMemoryUsage(org)
			Expect(usage).To(Equal(float64(512)))
		})
	})

	Describe("get spaces", func() {
		var spacesJSON []string

		BeforeEach(func() {
			spacesJSON = slurp("test-data/spaces.json")
		})

		It("should error when the the spaces url fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("Bad Things"))
			_, err := api.GetOrgSpaces("/v2/organizations/12345/spaces")
			Expect(err).ToNot(BeNil())
		})

		It("should return two spaces", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(spacesJSON, nil)
			spaces, _ := api.GetOrgSpaces("/v2/organizations/12345/spaces")
			Expect(len(spaces)).To(Equal(2))
		})

		It("should have name jdk-space", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(spacesJSON, nil)
			spaces, _ := api.GetOrgSpaces("/v2/organizations/12345/spaces")
			Expect(spaces[0].Name).To(Equal("jdk-space"))
			Expect(spaces[0].AppsURL).To(Equal("/v2/spaces/81c310ed-d258-48d7-a57a-6522d93a4217/apps"))
		})
	})

	Describe("get apps", func() {
		var appsJSON []string

		BeforeEach(func() {
			appsJSON = slurp("test-data/apps.json")
		})

		It("should return an error when the apps url fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("Bad Things"))
			_, err := api.GetSpaceApps("/v2/whateverapps")
			Expect(err).ToNot(BeNil())
		})

		It("should return one app with 1 instance and 1024 mb of ram", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(appsJSON, nil)
			apps, _ := api.GetSpaceApps("/v2/whateverapps")
			Expect(len(apps)).To(Equal(1))
			Expect(apps[0].Instances).To(Equal(float64(1)))
			Expect(apps[0].RAM).To(Equal(float64(1024)))
			Expect(apps[0].Running).To(BeTrue())
		})
	})

	// TODO need tests for no spaces and no apps in org.

	Describe("error calling CF API", func() {
		var errorJSON []string

		BeforeEach(func() {
			errorJSON = slurp("test-data/not-authenticated.json")
		})

		It("should return an error when CF API call fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(errorJSON, nil)
			_, err := api.GetOrgs()
			Expect(err.Error()).To(Equal("Error calling CF API: Authentication error"))
		})
	})
})
