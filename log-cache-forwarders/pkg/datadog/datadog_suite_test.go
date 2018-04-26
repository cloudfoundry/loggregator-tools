package datadog_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDatadog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Datadog Suite")
}
