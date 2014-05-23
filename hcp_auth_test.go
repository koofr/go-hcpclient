package hcpclient_test

import (
	. "github.com/koofr/go-hcpclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HcpAuth", func() {
	It("should generate valid Authorization header", func() {
		Expect(AuthHeader("foo", "bar")).To(Equal("HCP Zm9v:37b51d194a7513e45b56f6524f2d51f2"))
	})
})
