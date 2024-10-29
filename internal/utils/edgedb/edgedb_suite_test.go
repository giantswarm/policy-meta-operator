package edgedbutils_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEdgedb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Edgedb Suite")
}
