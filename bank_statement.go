package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

type StatementEntry struct {
	Date time.Time
	// Type	string
	Description string
	// PaidIn	int
	PaidOut int
	// Balance	int
}

func NewEntryFromRow(row []string) (StatementEntry, error) {
	if len(row) != 6 {
		err := fmt.Errorf("wanted 6 columns, got %d", len(row))
		return StatementEntry{}, err
	}

	var se StatementEntry
	const shortForm = "2 January 2006"
	t, err := time.Parse(shortForm, row[0])
	if err != nil {
		return se, err
	}
	se.Date = t
	se.Description = row[2]
	se.PaidOut, err = parseMoney(row[4])
	if err != nil {
		return se, err
	}
	return se, nil
}

func parseMoney(money string) (int, error) {
	if money == "-" {
		return 0, nil
	} else {
		money = strings.ReplaceAll(money, "£", "")
		money = strings.ReplaceAll(money, ",", "")
		money, err := strconv.ParseFloat(money, 64)
		if err != nil {
			return 0, err
		}
		money = money * 100
		money1 := int(money)
		return money1, nil
	}
}

func FilterStatementForOutgoingPayments(payments []StatementEntry) []StatementEntry {
	statementEntries := []StatementEntry{}
	for _, entry := range payments {
		if entry.PaidOut != 0 {
			statementEntries = append(statementEntries, entry)
		}
	}
	return statementEntries
}

func FilterStatementForUberTrips(payments []StatementEntry) []StatementEntry {
	uberTrips := []StatementEntry{}
	for _, entry := range payments {
		if entry.Description == "UBER TRIP" {
			uberTrips = append(uberTrips, entry)
		}
	}
	return uberTrips
}

func ComputeSumOfStatement(payments []StatementEntry) int {
	var total int
	for _, entry := range payments {
		total = total + entry.PaidOut
	}
	return total
}

func main() {
	content, err := ioutil.ReadFile("bankstatement")
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("File contents: %s", content)
	reader := csv.NewReader(bytes.NewReader(content))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	entries := []StatementEntry{}

	for rowIndex, row := range rows {
		if rowIndex == 0 {
			continue
		}

		entry, err := NewEntryFromRow(row)

		if err != nil {
			log.Fatal(err)
		}

		entries = append(entries, entry)
	}

	outgoingPayments := FilterStatementForOutgoingPayments(entries)
	uberTrips := FilterStatementForUberTrips(outgoingPayments)
	total := ComputeSumOfStatement(uberTrips)
	fmt.Println(total / 100)
}
