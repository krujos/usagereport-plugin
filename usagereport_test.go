package main

import (
	"github.com/cloudfoundry/cli/plugin"
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
)

type testAPIHelper struct {
}

func (t *testAPIHelper) GetOrgs(cli plugin.CliConnection) ([]Organization, error) {
	return nil, nil
}

func (t *testAPIHelper) GetQuotaMemoryLimit(cli plugin.CliConnection, quotaURL string) (float64, error) {
	return 0, nil
}

func (t *testAPIHelper) GetQuotaMemoryUsage(cli plugin.CliConnection, org Organization) (float64, error) {
	return 0, nil
}

var _ = Describe("Usagereport", func() {
	Describe("get org errors", func() {

		It("should return an error if cf curl /v2/organizations fails", func() {
			Fail("NYI")
		})
	})
})
