package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestUsagereportPlugin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UsagereportPlugin Suite")
}
