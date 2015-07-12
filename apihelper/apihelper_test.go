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
	var api *APIHelper
	var fakeCliConnection *fakes.FakeCliConnection

	BeforeEach(func() {
		fakeCliConnection = &fakes.FakeCliConnection{}
		api = &APIHelper{}
	})

	Describe("Get orgs", func() {
		var orgsJSON []string

		BeforeEach(func() {
			orgsJSON = slurp("test-data/orgs.json")
		})

		It("should return two orgs", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			orgs, _ := api.GetOrgs(fakeCliConnection)
			Expect(len(orgs)).To(Equal(2))
		})

		It("does something intellegent when cf curl fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				nil, errors.New("bad things"))
			_, err := api.GetOrgs(fakeCliConnection)
			Expect(err).ToNot(BeNil())
		})

		It("populates the url", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			orgs, _ := api.GetOrgs(fakeCliConnection)
			org := orgs[0]
			Expect(org.url).To(Equal("/v2/organizations/b1a23fd6-ac8d-4304-a3b4-815745417acd"))
		})

		It("calls /v2/orgs", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			api.GetOrgs(fakeCliConnection)
			args := fakeCliConnection.CliCommandWithoutTerminalOutputArgsForCall(0)
			Expect(args[1]).To(Equal("/v2/organizations"))
		})

	})

	Describe("Get quota memory limit", func() {
		var quotaJSON []string

		BeforeEach(func() {
			quotaJSON = slurp("test-data/quota.json")
		})

		It("should return an error when it can't fetch the memory limit", func() {
			_, err := api.GetQuotaMemoryLimit(fakeCliConnection, "/v2/somequota")
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				nil, errors.New("Bad Things"))
			Expect(err).ToNot(BeNil())
		})

		It("should reutrn 10240 as the memory limit", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				quotaJSON, nil)
			limit, _ := api.GetQuotaMemoryLimit(fakeCliConnection, "/v2/quotas/")
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
			_, err := api.GetOrgMemoryUsage(fakeCliConnection, org)
			Expect(err).ToNot(BeNil())
		})

		It("Shoudl return the memory usage", func() {
			org.url = "/v2/organizations/1234/"
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(usageJSON, nil)
			usage, _ := api.GetOrgMemoryUsage(fakeCliConnection, org)
			Expect(usage).To(Equal(float64(512)))
		})
	})

})
