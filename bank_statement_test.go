package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	statement "github.com/wenxin-liu/bank-statement"
	"time"
)

var _ = Describe("New entry from row", func() {
	It("should error when there are not enough fields", func() {
		row := []string{}
		_, err := statement.NewEntryFromRow(row)
		Expect(err).To(HaveOccurred())
	})

	It("should error when the date is invalid", func() {
		row := []string{"adate", "atime", "desc", "pi", "po", "bal"}
		_, err := statement.NewEntryFromRow(row)
		Expect(err).To(HaveOccurred())
	})

	It("should return a StatementEntry with input values", func() {
		row := []string{"21 November 2019", "POS", "UBER TRIP", "-", "£5.77", "£3,624.57"}
		statementEntry, err := statement.NewEntryFromRow(row)

		const shortForm = "2 January 2006"
		t, err := time.Parse(shortForm, "21 November 2019")
		Expect(err).NotTo(HaveOccurred())

		expectedStatementEntry := statement.StatementEntry{t, "UBER TRIP", 577}
		Expect(statementEntry).To(Equal(expectedStatementEntry))
	})

	It("should error when PaidOut is invalid", func() {
		row := []string{"21 November 2019", "atime", "desc", "pi", "po", "bal"}
		_, err := statement.NewEntryFromRow(row)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("Filter statement for outgoing payments", func() {
	It("should work", func() {
		statementEntries := []statement.StatementEntry{
			statement.StatementEntry{PaidOut: 0},
			statement.StatementEntry{PaidOut: 1},
			statement.StatementEntry{PaidOut: 100},
			statement.StatementEntry{PaidOut: 0},
		}

		filtered := statement.FilterStatementForOutgoingPayments(statementEntries)

		Expect(filtered).To(HaveLen(2))
	})
})

var _ = Describe("Filter statement for uber trips", func() {
	It("should work", func() {
		statementEntries := []statement.StatementEntry{
			statement.StatementEntry{Description: "breakfast"},
			statement.StatementEntry{Description: "UBER TRIP"},
			statement.StatementEntry{Description: "lunch UBER EATS"},
			statement.StatementEntry{Description: "UBER TRIP"},
			statement.StatementEntry{Description: "dinner"},
		}

		filtered := statement.FilterStatementForUberTrips(statementEntries)

		Expect(filtered).To(HaveLen(2))
	})
})

var _ = Describe("Computing the sum of statements", func() {
	It("should work", func() {
		statementEntries := []statement.StatementEntry{
			statement.StatementEntry{PaidOut: 0},
			statement.StatementEntry{PaidOut: 1},
			statement.StatementEntry{PaidOut: 100},
			statement.StatementEntry{PaidOut: 0},
		}

		total := statement.ComputeSumOfStatement(statementEntries)

		Expect(total).To(Equal(101))
	})
})
