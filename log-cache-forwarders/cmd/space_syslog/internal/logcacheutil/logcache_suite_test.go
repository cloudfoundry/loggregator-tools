package logcacheutil_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogcache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logcache Suite")
}
