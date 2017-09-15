package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHttpDrain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HttpDrain Suite")
}
