package viridian_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestViridian(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Viridian Suite")
}
