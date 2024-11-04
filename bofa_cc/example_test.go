package bofa_cc

import (
	"fmt"
)

func ExampleParseStatement() {
	statement := `
Account# 4400 1234 5678 1234
September 12 - October 11, 2024

Account Summary/Payment Information
Previous Balance $1,905.57
Payments and Other Credits -$1,977.21
Purchases and Adjustments $1,121.54
Fees Charged $0.00
Interest Charged $0.00
New Balance Total $1,049.90

Transactions

TransactionDate PostingDate Description ReferenceNumber AccountNumber Amount Total

Payments and Other Credits
09/28 09/30 PAYMENT - THANK YOU 0027 1234 -1,905.57
09/30 10/02 THE HOME DEPOT #1111 TOWN STATE 1579 1234 -55.42
10/08 10/09 COSTCO WHSE #1111 TOWN STATE 6307 1234 -16.22
TOTAL PAYMENTS AND OTHER CREDITS FOR THIS PERIOD -$1,977.21


Purchases and Adjustments
09/13 09/16 ERERE RERE COUNTY SCHOOL FDFDDF-DDFDFD DF 0881 1234 42.75
09/14 09/16 MY HEALTH RERERTDFDF 4912 1234 34.18
09/15 09/16 Subway 12345 SDRE ER 0067 1234 25.48
09/15 09/16 B'S PRODUCE TOWN CITY STATE 9139 1234 4.00
09/19 09/20 WAL-MART #1111, TOWN, STATE 5210 1234 0.98
09/19 09/20 COSTCO WHSE #1111 TOWN STATE 6299 1234 274.90
09/20 09/23 TST*WATERPARK - KIOSK 1 TOWN STATE 3524 1234 14.90
09/20 09/23 TST*WATERPARK - KIOSK 1 TOWN STATE 3557 1234 6.40
09/21 09/23 METRO 111-TOWN N TOWN STATE 5679 1234 46.54
09/22 09/23 Google 122X232 111-2222222 BC 7059 1234 94.23
09/25 09/26 COSTCO WHSE #1111 TOWN STATE 8119 1234 93.49
09/27 09/30 HOMEDEPOT.COM 111-111-1111 BC 8383 1234 54.92
09/30 10/01 LOWES #01878* TOWN STATE 8740 1234 14.73
09/30 10/02 THE HOME DEPOT #3644 TOWN STATE 2309 1234 64.70
10/01 10/02 WHOLEFDS CAR 1111 TOWN STATE 5423 1234 60.10
10/02 10/03 COSTCO WHSE #1206 TOWN STATE 2910 1234 142.13
10/03 10/04 HELLOMONKEY STUDIOS HTTPSWWW.CODECA 6774 1234 27.04
10/05 10/07 MY CHURCH EWWEW WEWEWE 5336 1234 10.00
10/06 10/07 DUNKIN #111111 TOWN STATE 3379 1234 3.64
10/11 10/11 SP HAIR HTTPSWWW.HAIR 0637 1234 106.43
TOTAL PURCHASES AND ADJUSTMENTS FOR THIS PERIOD $1,121.54


Interest Charged
10/11 10/11 INTEREST CHARGED ON PURCHASES 0.00
10/11 10/11 INTEREST CHARGED ON BALANCE TRANSFERS 0.00
10/11 10/11 INTEREST CHARGED ON DIR DEP&CHK CASHADV 0.00
10/11 10/11 INTEREST CHARGED ON BANK CASH ADVANCES 0.00
TOTAL INTEREST CHARGED FOR THIS PERIOD $0.00
	`

	s, err := ParseStatement(statement)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, tx := range s.Transactions {
		fmt.Printf("%+v\n", tx)
	}

	// Output:
	// {TransactionDate:2024-09-28 00:00:00 +0000 UTC PostingDate:2024-09-30 00:00:00 +0000 UTC Description:PAYMENT - THANK YOU ReferenceNumber:0027 AccountNumber:1234 Amount:-1905.57 Category:Payments and Other Credits}
	// {TransactionDate:2024-09-30 00:00:00 +0000 UTC PostingDate:2024-10-02 00:00:00 +0000 UTC Description:THE HOME DEPOT #1111 TOWN STATE ReferenceNumber:1579 AccountNumber:1234 Amount:-55.42 Category:Payments and Other Credits}
	// {TransactionDate:2024-10-08 00:00:00 +0000 UTC PostingDate:2024-10-09 00:00:00 +0000 UTC Description:COSTCO WHSE #1111 TOWN STATE ReferenceNumber:6307 AccountNumber:1234 Amount:-16.22 Category:Payments and Other Credits}
	// {TransactionDate:2024-09-13 00:00:00 +0000 UTC PostingDate:2024-09-16 00:00:00 +0000 UTC Description:ERERE RERE COUNTY SCHOOL FDFDDF-DDFDFD DF ReferenceNumber:0881 AccountNumber:1234 Amount:42.75 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-14 00:00:00 +0000 UTC PostingDate:2024-09-16 00:00:00 +0000 UTC Description:MY HEALTH RERERTDFDF ReferenceNumber:4912 AccountNumber:1234 Amount:34.18 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-15 00:00:00 +0000 UTC PostingDate:2024-09-16 00:00:00 +0000 UTC Description:Subway 12345 SDRE ER ReferenceNumber:0067 AccountNumber:1234 Amount:25.48 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-15 00:00:00 +0000 UTC PostingDate:2024-09-16 00:00:00 +0000 UTC Description:B'S PRODUCE TOWN CITY STATE ReferenceNumber:9139 AccountNumber:1234 Amount:4 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-19 00:00:00 +0000 UTC PostingDate:2024-09-20 00:00:00 +0000 UTC Description:WAL-MART #1111, TOWN, STATE ReferenceNumber:5210 AccountNumber:1234 Amount:0.98 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-19 00:00:00 +0000 UTC PostingDate:2024-09-20 00:00:00 +0000 UTC Description:COSTCO WHSE #1111 TOWN STATE ReferenceNumber:6299 AccountNumber:1234 Amount:274.9 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-20 00:00:00 +0000 UTC PostingDate:2024-09-23 00:00:00 +0000 UTC Description:TST*WATERPARK - KIOSK 1 TOWN STATE ReferenceNumber:3524 AccountNumber:1234 Amount:14.9 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-20 00:00:00 +0000 UTC PostingDate:2024-09-23 00:00:00 +0000 UTC Description:TST*WATERPARK - KIOSK 1 TOWN STATE ReferenceNumber:3557 AccountNumber:1234 Amount:6.4 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-21 00:00:00 +0000 UTC PostingDate:2024-09-23 00:00:00 +0000 UTC Description:METRO 111-TOWN N TOWN STATE ReferenceNumber:5679 AccountNumber:1234 Amount:46.54 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-22 00:00:00 +0000 UTC PostingDate:2024-09-23 00:00:00 +0000 UTC Description:Google 122X232 111-2222222 BC ReferenceNumber:7059 AccountNumber:1234 Amount:94.23 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-25 00:00:00 +0000 UTC PostingDate:2024-09-26 00:00:00 +0000 UTC Description:COSTCO WHSE #1111 TOWN STATE ReferenceNumber:8119 AccountNumber:1234 Amount:93.49 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-27 00:00:00 +0000 UTC PostingDate:2024-09-30 00:00:00 +0000 UTC Description:HOMEDEPOT.COM 111-111-1111 BC ReferenceNumber:8383 AccountNumber:1234 Amount:54.92 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-30 00:00:00 +0000 UTC PostingDate:2024-10-01 00:00:00 +0000 UTC Description:LOWES #01878* TOWN STATE ReferenceNumber:8740 AccountNumber:1234 Amount:14.73 Category:Purchases and Adjustments}
	// {TransactionDate:2024-09-30 00:00:00 +0000 UTC PostingDate:2024-10-02 00:00:00 +0000 UTC Description:THE HOME DEPOT #3644 TOWN STATE ReferenceNumber:2309 AccountNumber:1234 Amount:64.7 Category:Purchases and Adjustments}
	// {TransactionDate:2024-10-01 00:00:00 +0000 UTC PostingDate:2024-10-02 00:00:00 +0000 UTC Description:WHOLEFDS CAR 1111 TOWN STATE ReferenceNumber:5423 AccountNumber:1234 Amount:60.1 Category:Purchases and Adjustments}
	// {TransactionDate:2024-10-02 00:00:00 +0000 UTC PostingDate:2024-10-03 00:00:00 +0000 UTC Description:COSTCO WHSE #1206 TOWN STATE ReferenceNumber:2910 AccountNumber:1234 Amount:142.13 Category:Purchases and Adjustments}
	// {TransactionDate:2024-10-03 00:00:00 +0000 UTC PostingDate:2024-10-04 00:00:00 +0000 UTC Description:HELLOMONKEY STUDIOS HTTPSWWW.CODECA ReferenceNumber:6774 AccountNumber:1234 Amount:27.04 Category:Purchases and Adjustments}
	// {TransactionDate:2024-10-05 00:00:00 +0000 UTC PostingDate:2024-10-07 00:00:00 +0000 UTC Description:MY CHURCH EWWEW WEWEWE ReferenceNumber:5336 AccountNumber:1234 Amount:10 Category:Purchases and Adjustments}
	// {TransactionDate:2024-10-06 00:00:00 +0000 UTC PostingDate:2024-10-07 00:00:00 +0000 UTC Description:DUNKIN #111111 TOWN STATE ReferenceNumber:3379 AccountNumber:1234 Amount:3.64 Category:Purchases and Adjustments}
	// {TransactionDate:2024-10-11 00:00:00 +0000 UTC PostingDate:2024-10-11 00:00:00 +0000 UTC Description:SP HAIR HTTPSWWW.HAIR ReferenceNumber:0637 AccountNumber:1234 Amount:106.43 Category:Purchases and Adjustments}
}
