package td

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/muly/bank-tx/util"
)

// Transaction struct to hold transaction data
type Transaction struct {
	Category    string
	PostingDate time.Time
	Description string
	Amount      float64
}

// Statement struct to hold overall statement info
type Statement struct {
	AccountNumber    string
	PeriodStartDate  time.Time
	PeriodEndDate    time.Time
	BeginningBalance float64
	EndingBalance    float64
	Transactions     []Transaction
}

// ParseStatement parses the input data into a Statement struct
func ParseStatement(data string) (*Statement, error) {
	var statement Statement
	var currentCategory string
	var year int

	lines := strings.Split(data, "\n")

	// Regular expressions to capture different data fields
	rePeriod := regexp.MustCompile(`Statement Period: (\w+ \d{2} \d{4})-(\w+ \d{2} \d{4})`)
	reAccount := regexp.MustCompile(`Account # (\d{14}|\d{3}-\d{7})`)
	reBalance := regexp.MustCompile(`(Beginning|Ending) Balance ([\d,]+\.\d{2})`)
	reCategory := regexp.MustCompile(`^(Electronic Deposits|Electronic Payments)$`)
	reTransaction := regexp.MustCompile(`^(\d{2}/\d{2})\s+(.+?)\s+([\d,]+\.\d{2})$`)
	reSubtotal := regexp.MustCompile(`^Subtotal: ([\d,]+\.\d{2})$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse statement period for dates and year
		if match := rePeriod.FindStringSubmatch(line); match != nil {
			startDate, _ := time.Parse("Jan 02 2006", match[1])
			endDate, _ := time.Parse("Jan 02 2006", match[2])
			year = startDate.Year()
			statement.PeriodStartDate = startDate
			statement.PeriodEndDate = endDate
			continue
		}

		// Parse account number
		if match := reAccount.FindStringSubmatch(line); match != nil {
			statement.AccountNumber = match[1]
			continue
		}

		// Parse beginning and ending balances
		if match := reBalance.FindStringSubmatch(line); match != nil {
			balance, _ := strconv.ParseFloat(strings.ReplaceAll(match[2], ",", ""), 64)
			if match[1] == "Beginning" {
				statement.BeginningBalance = balance
			} else {
				statement.EndingBalance = balance
			}
			continue
		}

		// Parse transaction categories
		if match := reCategory.FindStringSubmatch(line); match != nil {
			currentCategory = match[1]
			continue
		}

		// Parse transaction lines
		if match := reTransaction.FindStringSubmatch(line); match != nil {
			postingDate, _ := time.Parse("01/02", match[1])
			// Set the correct year for the posting date
			postingDate = postingDate.AddDate(year-postingDate.Year(), 0, 0)

			description := match[2]
			amount, _ := strconv.ParseFloat(strings.ReplaceAll(match[3], ",", ""), 64)

			transaction := Transaction{
				Category:    currentCategory,
				PostingDate: postingDate,
				Description: description,
				Amount:      amount,
			}
			statement.Transactions = append(statement.Transactions, transaction)
			continue
		}

		// Parse sub-totals for validation
		if match := reSubtotal.FindStringSubmatch(line); match != nil {
			subtotal, _ := strconv.ParseFloat(strings.ReplaceAll(match[1], ",", ""), 64)
			// Validate subtotal against transactions
			totalAmount := 0.0
			for _, transaction := range statement.Transactions {
				if transaction.Category == currentCategory {
					totalAmount += transaction.Amount
				}
			}
			totalAmount = util.RoundToOneDecimal(totalAmount)
			subtotal = util.RoundToOneDecimal(subtotal)
			if totalAmount != subtotal {
				return nil, fmt.Errorf("subtotal mismatch in category %s: expected %.2f, got %.2f", currentCategory, subtotal, totalAmount)
			}
			continue
		}
	}

	// Final validation for beginning and ending balances
	deposits, payments := 0.0, 0.0
	for _, transaction := range statement.Transactions {
		switch transaction.Category {
		case "Electronic Deposits":
			deposits += transaction.Amount
		case "Electronic Payments":
			payments += transaction.Amount
		}
	}
	calculatedEndingBalance := util.RoundToOneDecimal(statement.BeginningBalance + deposits - payments)
	statement.EndingBalance = util.RoundToOneDecimal(statement.EndingBalance)

	if calculatedEndingBalance != statement.EndingBalance {
		return nil, fmt.Errorf("ending balance mismatch: expected %.2f, got %.2f", calculatedEndingBalance, statement.EndingBalance)
	}

	return &statement, nil
}

// SaveTransactionsToCSV saves the transactions to a CSV file
func SaveTransactionsToCSV(statement *Statement) error {
	filename := fmt.Sprintf("%s|%s_to_%s.csv", statement.AccountNumber,
		statement.PeriodStartDate.Format("Jan_02_2006"),
		statement.PeriodEndDate.Format("Jan_02_2006"))

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write transaction headers
	writer.Write([]string{"Category", "Posting Date", "Description", "Amount"})

	// Write each transaction
	for _, t := range statement.Transactions {
		writer.Write([]string{
			t.Category,
			t.PostingDate.Format("01/02/2006"),
			t.Description,
			fmt.Sprintf("%.2f", t.Amount),
		})
	}

	fmt.Printf("Transactions saved to %s\n", filename)
	return nil
}
