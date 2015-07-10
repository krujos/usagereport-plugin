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

	Describe("get orgs", func() {
		var fakeCliConnection *fakes.FakeCliConnection
		var orgsJSON []string
		var cmd *UsageReportCmd

		BeforeEach(func() {
			fakeCliConnection = &fakes.FakeCliConnection{}
			cmd = &UsageReportCmd{}

			file, _ := os.Open("orgs.json")
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
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("bad things"))
			_, err := cmd.getOrgs(fakeCliConnection)
			Expect(err).ToNot(BeNil())
		})

	})
	Describe("get quota memory limit", func() {
		It("should return an error when it can't fetch the memory limit", func() {
			Fail("NYI")
		})
	})
})
