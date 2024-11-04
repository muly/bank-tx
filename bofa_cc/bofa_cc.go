// package bofa_cc provides the parsing functions to process the bofa credit card statements
package bofa_cc

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/muly/bank-tx/util"
)

// TODO: interest lines are not included in the parsing logic.

type Transaction struct {
	TransactionDate time.Time
	PostingDate     time.Time
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
	var accountNumber string
	var periodStartDate, periodEndDate time.Time
	var beginBalance, endBalance, totalPayments, totalPurchases, totalFees, totalInterest float64
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
		if strings.Contains(line, "-") && strings.HasSuffix(line, "2024") { // TODO: need to replace this with a regex
			periodStartDate, periodEndDate, err = parseStatementPeriod(line)
			if err != nil {
				return nil, err
			}

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
		if strings.HasPrefix(line, "Fees Charged") && strings.Contains(line, "$") {
			totalFees, _ = util.ParseFloat(strings.TrimPrefix(line, "Fees Charged "))
			if err != nil {
				fmt.Println("totalFees parse error:", err)
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

		if line == "Payments and Other Credits" ||
			line == "Purchases and Adjustments" ||
			line == "Interest Charged" ||
			line == "Fees" {
			inCategory = line
			continue // Skip the header line
		}

		if line == "TransactionDate PostingDate Description ReferenceNumber AccountNumber Amount Total" ||
			line == "Transactions" ||
			line == "Account Summary/Payment Information" {
			continue // Skip the transaction header line
		}

		if strings.HasPrefix(line, "TOTAL PAYMENTS AND OTHER CREDITS FOR THIS PERIOD") ||
			strings.HasPrefix(line, "TOTAL PURCHASES AND ADJUSTMENTS FOR THIS PERIOD") ||
			strings.HasPrefix(line, "TOTAL INTEREST CHARGED FOR THIS PERIOD") ||
			strings.HasPrefix(line, "TOTAL FEES FOR THIS PERIOD") {
			continue // Skip the transaction subtotal line
		}

		transaction, err := ParseTransaction(line, periodStartDate, periodEndDate)
		if err != nil {
			return nil, err
		}
		if transaction == nil {
			log.Printf("unprocessed line: %v", line)
			continue // unwanted line
		}

		transaction.Category = inCategory

		transactions = append(transactions, *transaction)
	}
	// fmt.Println(beginBalance, endBalance, totalPayments, totalPurchases, totalFees, totalInterest)

	// Validate balances
	if !validateSummaryBalance(beginBalance, totalPayments, totalPurchases, totalFees, totalInterest, endBalance) {
		return nil, fmt.Errorf("summary balance validation failed: beginBalance %v, totalPayments %v, totalPurchases %v, totalFees %v, totalInterest %v, endBalance %v",
			beginBalance, totalPayments, totalPurchases, totalFees, totalInterest, endBalance)
	}

	total := calculateTotal(transactions)
	if !validateTxBalance(beginBalance, endBalance, total) {
		return nil, fmt.Errorf("tx balance validation failed. calculated end balance %v, expected end balance %v", util.RoundToTwoDecimal(beginBalance+total), endBalance)
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
func ParseTransaction(line string, startPeriod, endPeriod time.Time) (*Transaction, error) {
	// txRegex := regexp.MustCompile(`(?m)^(\d{2}/\d{2})\s+(\d{2}/\d{2})\s+(.+?)\s+(\d+)\s+(-?\$?[\d,]+\.\d{2})$`)
	// txRegex := regexp.MustCompile(`^(\d{2}/\d{2}) (\d{2}/\d{2}) (.*?) (\d{4}) (\d{4}) (-?\$?\d+\.\d{2})$`)
	var txRegex = regexp.MustCompile(`(?m)^(\d{2}/\d{2})\s+(\d{2}/\d{2})\s+(.+?)\s+(\d{4})\s+(\d{4})\s+(-?\$?[\d,]+\.\d{2})$`)
	interestRegex := regexp.MustCompile(`^(\d{1,2}/\d{1,2})\s+(\d{1,2}/\d{1,2})\s+(.+?)\s+(\d+\.\d{2})$`)

	transaction := Transaction{}

	if matches := txRegex.FindStringSubmatch(line); matches != nil {
		transactionDate := matches[1]
		postingDate := matches[2]
		description := matches[3]
		referenceNumber := matches[4]
		accountNumber := matches[5]
		amountStr := matches[6]

		amountStr = strings.ReplaceAll(amountStr, "$", "")
		amountStr = strings.ReplaceAll(amountStr, ",", "")

		amount, err := util.ParseFloat(amountStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse amount: %s, error: %v", amountStr, err)
		}

		transaction = Transaction{
			Description:     description,
			ReferenceNumber: referenceNumber,
			AccountNumber:   accountNumber,
			Amount:          amount,
		}

		transaction.TransactionDate, err = addYearToDate(transactionDate, startPeriod, endPeriod)
		if err != nil {
			return nil, fmt.Errorf("error adding year to transaction date: %v", err)
		}
		transaction.PostingDate, err = addYearToDate(postingDate, startPeriod, endPeriod)
		if err != nil {
			return nil, fmt.Errorf("error adding year to posting date: %v", err)
		}

		return &transaction, nil
	}

	if matches := interestRegex.FindStringSubmatch(line); matches != nil {
		transactionDate := matches[1]
		postingDate := matches[2]
		description := matches[3]
		amountStr := matches[4]

		amountStr = strings.ReplaceAll(amountStr, "$", "")
		amountStr = strings.ReplaceAll(amountStr, ",", "")

		amount, err := util.ParseFloat(amountStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse amount: %s, error: %v", amountStr, err)
		}

		transaction = Transaction{
			Description: description,
			Amount:      amount,
		}

		transaction.TransactionDate, err = addYearToDate(transactionDate, startPeriod, endPeriod)
		if err != nil {
			return nil, fmt.Errorf("error adding year to transaction date: %v", err)
		}
		transaction.PostingDate, err = addYearToDate(postingDate, startPeriod, endPeriod)
		if err != nil {
			return nil, fmt.Errorf("error adding year to posting date: %v", err)
		}

		return &transaction, nil
	}

	return nil, nil
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
			tx.TransactionDate.Format("01/02/2006"),
			tx.PostingDate.Format("01/02/2006"),
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

func parseStatementPeriod(periodStr string) (time.Time, time.Time, error) {
	// Split the period string to get start and end dates
	parts := strings.Split(periodStr, " - ")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid statement period format")
	}

	// Define the year format for the end date part
	yearFormat := "January 2, 2006"
	// Define the format for the start date part, which lacks a year
	startDateFormat := "January 2"

	// Parse the end date (includes the year)
	endPeriod, err := time.Parse(yearFormat, parts[1])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to parse end period: %v", err)
	}

	// Parse the start date (no year)
	startDateNoYear, err := time.Parse(startDateFormat, parts[0])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to parse start period: %v", err)
	}

	// Determine the year for the start date
	startYear := endPeriod.Year()
	if startDateNoYear.Month() > endPeriod.Month() {
		// If start month is after end month, assume the start date is in the previous year
		startYear = endPeriod.Year() - 1
	}

	// Create startPeriod with the derived year
	startPeriod := time.Date(startYear, startDateNoYear.Month(), startDateNoYear.Day(), 0, 0, 0, 0, time.UTC)

	return startPeriod, endPeriod, nil
}

// addYearToDate adds the correct year to a given month-day date based on the statement period
func addYearToDate(dateStr string, startPeriod, endPeriod time.Time) (time.Time, error) {
	parsedDate, err := time.Parse("01/02", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	// Set the year based on start and end periods
	year := endPeriod.Year()

	// For cross-year periods, dates before the end month should use startPeriod's year
	if startPeriod.Year() != endPeriod.Year() {
		if parsedDate.Month() > endPeriod.Month() {
			year = startPeriod.Year()
		}
	}

	return parsedDate.AddDate(year, 0, 0), nil
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
func validateSummaryBalance(beginBalance, totalPayments, totalPurchases, totalFees, totalInterest, endBalance float64) bool {
	calculatedEndBalance := beginBalance + totalPayments + totalPurchases + totalFees + totalInterest
	return util.RoundToOneDecimal(calculatedEndBalance) == util.RoundToOneDecimal(endBalance)
}

// func main() {
// 	data, err := os.ReadFile("sample.txt")
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		return
// 	}

// 	statement, err := ParseStatement(string(data))
// 	if err != nil {
// 		fmt.Println("Error parsing statement:", err)
// 		return
// 	}

// 	if err := SaveTransactionsToCSV(*statement); err != nil {
// 		fmt.Println("Error writing to CSV:", err)
// 		return
// 	}

// 	fmt.Println("Statement parsed and saved successfully.")
// }
