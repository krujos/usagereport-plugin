package apihelper_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestApihelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Apihelper Suite")
}
