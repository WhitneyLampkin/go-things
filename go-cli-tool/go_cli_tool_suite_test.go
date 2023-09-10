package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Entrypoint for Ginkgo
func TestGoCliTool(t *testing.T) {
	// Glue code that connects Ginkgo to Gomega
	RegisterFailHandler(Fail)
	// Starts Ginkgo test suite
	RunSpecs(t, "GoCliTool Suite")
}
