package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBankStatement(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BankStatement Suite")
}
