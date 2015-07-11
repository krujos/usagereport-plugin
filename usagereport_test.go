package main

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
	var cmd *UsageReportCmd
	var fakeCliConnection *fakes.FakeCliConnection

	BeforeEach(func() {
		fakeCliConnection = &fakes.FakeCliConnection{}
		cmd = &UsageReportCmd{}
	})

	Describe("get orgs", func() {
		var orgsJSON []string

		BeforeEach(func() {
			orgsJSON = slurp("test-data/orgs.json")
		})

		It("should return two orgs", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			orgs, _ := cmd.getOrgs(fakeCliConnection)
			Expect(len(orgs)).To(Equal(2))
		})

		It("does something intellegent when cf curl fails", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				nil, errors.New("bad things"))
			_, err := cmd.getOrgs(fakeCliConnection)
			Expect(err).ToNot(BeNil())
		})

		It("populates the url", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(orgsJSON, nil)
			orgs, _ := cmd.getOrgs(fakeCliConnection)
			org := orgs[0]
			Expect(org.url).To(Equal("/v2/organizations/b1a23fd6-ac8d-4304-a3b4-815745417acd"))
		})

	})

	Describe("get quota memory limit", func() {
		var quotaJSON []string

		BeforeEach(func() {
			quotaJSON = slurp("test-data/quota.json")
		})

		It("should return an error when it can't fetch the memory limit", func() {
			_, err := cmd.getQuotaMemoryLimit(fakeCliConnection, "/v2/somequota")
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				nil, errors.New("Bad Things"))
			Expect(err).ToNot(BeNil())
		})

		It("should reutrn 10240 as the memory limit", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				quotaJSON, nil)
			limit, _ := cmd.getQuotaMemoryLimit(fakeCliConnection, "/v2/quotas/")
			Expect(limit).To(Equal(float64(10240)))
		})
	})

	Describe("it gets the org memory usage", func() {
		var org organization
		var usageJSON []string

		BeforeEach(func() {
			usageJSON = slurp("test-data/memory_usage.json")
		})

		It("should return an error when it can't fetch the orgs memory usage", func() {
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil,
				errors.New("Bad things"))
			_, err := cmd.getOrgMemoryUsage(fakeCliConnection, org)
			Expect(err).ToNot(BeNil())
		})

		It("Shoudl return the memory usage", func() {
			org.url = "/v2/organizations/1234/"
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(usageJSON, nil)
			usage, _ := cmd.getOrgMemoryUsage(fakeCliConnection, org)
			Expect(usage).To(Equal(float64(512)))
		})
	})

})
