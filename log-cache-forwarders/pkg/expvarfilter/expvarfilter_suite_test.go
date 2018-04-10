package expvarfilter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestExpvarfilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Expvarfilter Suite")
}
