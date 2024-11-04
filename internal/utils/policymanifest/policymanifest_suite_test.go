package polman_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPolicymanifest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PolicyManifest Suite")
}
