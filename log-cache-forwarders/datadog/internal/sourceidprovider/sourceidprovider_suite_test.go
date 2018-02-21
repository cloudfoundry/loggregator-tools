package sourceidprovider_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSourceidprovider(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sourceidprovider Suite")
}
