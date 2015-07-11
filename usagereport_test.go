package main

import (
	"bufio"
	"errors"
	"os"

	"github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
			file, _ := os.Open("test-data/orgs.json")
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				orgsJSON = append(orgsJSON, scanner.Text())
			}
		})

		AfterEach(func() {
			orgsJSON = nil
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

	})
	Describe("get quota memory limit", func() {
		It("should return an error when it can't fetch the memory limit", func() {
			_, err := cmd.getQuotaMemoryLimit(fakeCliConnection, "/v2/somequota")
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(
				nil, errors.New("Bad Things"))
			Expect(err).ToNot(BeNil())
		})
	})
})
