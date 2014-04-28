package hcpclient_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHcp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hcp Suite")
}
