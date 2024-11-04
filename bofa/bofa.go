package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/muly/bank-tx/util"
)

// TODO: interest lines are not included in the parsing logic.

type Transaction struct {
	TransactionDate string
	PostingDate     string
	Description     string
	ReferenceNumber string
	AccountNumber   string
	Amount          float64
	Category        string
}

type Statement struct {
	AccountNumber    string
	PeriodStartDate  time.Time
	PeriodEndDate    time.Time
	BeginningBalance float64
	EndingBalance    float64
	Transactions     []Transaction
}

// ParseStatement parses the statement
func ParseStatement(data string) (*Statement, error) {
	lines := strings.Split(data, "\n")

	var transactions []Transaction
	var accountNumber, periodStr, year string
	var beginBalance, endBalance, totalPayments, totalPurchases, totalInterest float64
	inCategory := ""
	var err error


	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if line == "" {
			continue
		}

		// Parse Account Number
		if strings.HasPrefix(line, "Account#") {
			accountNumber = strings.TrimSpace(strings.Split(line, "#")[1])
			continue
		}

		// Parse Statement Period
		if strings.Contains(line, "-") && strings.HasSuffix(line, "2024") {
			periodStr = line
			// var err error
			// _, year, err = parseStatementPeriod(periodStr)
			// if err != nil {
			// 	return nil, fmt.Errorf("invalid period format: %w", err)
			// }
			continue
		}

		// Parse balance information
		if strings.HasPrefix(line, "Previous Balance") {
			beginBalance, err = util.ParseFloat(strings.TrimPrefix(line, "Previous Balance "))
			if err != nil {
				fmt.Println("beginBalance parse error:", err)
			}
			continue
		}

		if strings.HasPrefix(line, "New Balance Total") {
			endBalance, _ = util.ParseFloat(strings.TrimPrefix(line, "New Balance Total "))
			if err != nil {
				fmt.Println("endBalance parse error:", err)
			}
			continue
		}
		if strings.HasPrefix(line, "Payments and Other Credits") && strings.Contains(line, "-") {
			totalPayments, _ = util.ParseFloat(strings.TrimPrefix(line, "Payments and Other Credits "))
			if err != nil {
				fmt.Println("totalPayments parse error:", err)
			}
			continue
		}
		if strings.HasPrefix(line, "Purchases and Adjustments") && strings.Contains(line, "$") {
			totalPurchases, _ = util.ParseFloat(strings.TrimPrefix(line, "Purchases and Adjustments "))
			if err != nil {
				fmt.Println("totalPurchases parse error:", err)
			}
			continue
		}
		if strings.HasPrefix(line, "Interest Charged") && strings.Contains(line, "$") {
			totalInterest, _ = util.ParseFloat(strings.TrimPrefix(line, "Interest Charged "))
			if err != nil {
				fmt.Println("totalInterest parse error:", err)
			}
			continue
		}

		// Identify and categorize transaction sections
		if line == "Payments and Other Credits" || line == "Purchases and Adjustments" || line == "Interest Charged" {
			inCategory = line
			continue // Skip the header line
		}

		if line == "TransactionDate PostingDate Description ReferenceNumber AccountNumber Amount Total" ||
			line == "Transactions" ||
			line == "Account Summary/Payment Information" {
			continue // Skip the transaction header line
		}

		transaction, err := ParseTransaction(line)
		if err != nil {
			return nil, err
		}
		if transaction == nil {
			continue // unwanted line
		}

		transaction.TransactionDate = transaction.TransactionDate + "/" + year
		transaction.PostingDate = transaction.PostingDate + "/" + year
		transaction.Category = inCategory

		transactions = append(transactions, *transaction)
	}
	// fmt.Println(beginBalance, endBalance, totalPayments, totalPurchases, totalInterest)

	// Validate balances
	if !validateSummaryBalance(beginBalance, endBalance, totalPayments, totalPurchases, totalInterest) {
		return nil, fmt.Errorf("summary balance validation failed")
	}

	total := calculateTotal(transactions)
	if !validateTxBalance(beginBalance, endBalance, total) {
		return nil, fmt.Errorf("tx balance validation failed. calculated end balance %v, expected end balance %v", util.RoundToTwoDecimal(beginBalance+total), endBalance)
	}

	// Parse statement period for the CSV filename
	periodStartDate, periodEndDate, err := parseStatementPeriod(periodStr)
	if err != nil {
		return nil, err
	}

	statement := Statement{
		AccountNumber:    accountNumber,
		PeriodStartDate:  periodStartDate,
		PeriodEndDate:    periodEndDate,
		BeginningBalance: beginBalance,
		EndingBalance:    endBalance,
		Transactions:     transactions,
	}

	return &statement, nil
}

// ParseTransaction parses the given transaction entry
func ParseTransaction(line string) (*Transaction, error) {
	// txRegex := regexp.MustCompile(`(?m)^(\d{2}/\d{2})\s+(\d{2}/\d{2})\s+(.+?)\s+(\d+)\s+(-?\$?[\d,]+\.\d{2})$`)
	// txRegex := regexp.MustCompile(`^(\d{2}/\d{2}) (\d{2}/\d{2}) (.*?) (\d{4}) (\d{4}) (-?\$?\d+\.\d{2})$`)
	var txRegex = regexp.MustCompile(`(?m)^(\d{2}/\d{2})\s+(\d{2}/\d{2})\s+(.+?)\s+(\d{4})\s+(\d{4})\s+(-?\$?[\d,]+\.\d{2})$`)

	matches := txRegex.FindStringSubmatch(line)
	if len(matches) != 7 {
		fmt.Printf("line not processed: %s\n", line)
		return nil, nil
	}

	transactionDate := matches[1]
	postingDate := matches[2]
	description := matches[3]
	referenceNumber := matches[4]
	accountNumber := matches[5]
	amountStr := matches[6]

	amountStr = strings.ReplaceAll(amountStr, "$", "")
	amountStr = strings.ReplaceAll(amountStr, ",", "")

	amount, err := util.ParseFloat(amountStr)
	// amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %s, error: %v", amountStr, err)
	}

	transaction := Transaction{
		TransactionDate: transactionDate,
		PostingDate:     postingDate,
		Description:     description,
		ReferenceNumber: referenceNumber,
		AccountNumber:   accountNumber,
		Amount:          amount,
	}

	return &transaction, nil
}

// SaveTransactionsToCSV writes transactions to CSV
func SaveTransactionsToCSV(statement Statement) error {
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

	// Write CSV headers
	writer.Write([]string{"TransactionDate", "PostingDate", "Description", "ReferenceNumber", "AccountNumber", "Amount", "Category"})

	// Write transaction data
	for _, tx := range statement.Transactions {
		writer.Write([]string{
			tx.TransactionDate,
			tx.PostingDate,
			tx.Description,
			tx.ReferenceNumber,
			tx.AccountNumber,
			fmt.Sprintf("%.2f", tx.Amount),
			tx.Category,
		})
	}

	return nil
}

// Helper functions

// Helper function to parse the date format "September 12 - October 11, 2024"
func parseStatementPeriod(periodStr string) (time.Time, time.Time, error) {
	periodRegex := regexp.MustCompile(`(\w+ \d+) - (\w+ \d+), (\d{4})`)
	matches := periodRegex.FindStringSubmatch(periodStr)
	if len(matches) != 4 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid statement period format")
	}
	year := matches[3]
	// Format period as YYYY-MM-DD to YYYY-MM-DD
	startDate, err := time.Parse("January 2 2006", matches[1]+" "+year)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endDate, err := time.Parse("January 2 2006", matches[2]+" "+year)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return startDate, endDate, nil
}

func calculateTotal(transactions []Transaction) float64 {
	total := 0.0
	for _, tx := range transactions {
		total += tx.Amount
	}
	return total
}

func validateTxBalance(beginBalance, endBalance, total float64) bool {
	return util.RoundToOneDecimal(beginBalance+total) == util.RoundToOneDecimal(endBalance)
}

// Validate balance based on parsed totals
func validateSummaryBalance(beginBalance, endBalance, totalPayments, totalPurchases, totalInterest float64) bool {
	calculatedEndBalance := beginBalance + totalPayments + totalPurchases + totalInterest
	return util.RoundToOneDecimal(calculatedEndBalance) == util.RoundToOneDecimal(endBalance)
}

func main() {
	data, err := os.ReadFile("sample.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	statement, err := ParseStatement(string(data))
	if err != nil {
		fmt.Println("Error parsing statement:", err)
		return
	}

	if err := SaveTransactionsToCSV(*statement); err != nil {
		fmt.Println("Error writing to CSV:", err)
		return
	}

	fmt.Println("Statement parsed and saved successfully.")
}
