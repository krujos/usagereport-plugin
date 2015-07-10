package main

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Report", func() {

	It("should fail", func() {
		Fail("Yep")
	})
})
