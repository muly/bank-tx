package td

import (
	"fmt"
)

func ExampleParseStatement() {
	statement := `
Statement Period: Mar 21 2023-Apr 20 2023
TD Convenience Checking
Some Name
Account # 123-4567890

ACCOUNT SUMMARY
Beginning Balance 10,750.35
Electronic Deposits 16,918.44
Electronic Payments 17,537.72
Ending Balance 10,131.07
 
DAILY ACCOUNT ACTIVITY
Electronic Deposits
POSTING DATE DESCRIPTION AMOUNT
03/22 TD ZELLE RECEIVED, erer434ree r34rere re5rerer4344re 418.00
04/02 ACH DEPOSIT, rere erer ererereL 6,377.30
04/08 ACH DEPOSIT, fererer  erereer dfdferr 6,377.30
04/18 ACH DEPOSIT, rere rer4tr rtrtrrtr 3,745.84
Subtotal: 16,918.44
Electronic Payments
POSTING DATE DESCRIPTION AMOUNT
03/28 ELECTRONIC PMT-WEB, BKOFAM CK WEBXFR TRANSFER ****234533 1,000.00
03/28 TD BILL PAY SERV, BANK OF AMERICA ONLINE PMT TDB****34454454 4,223.27
04/11 TD BILL PAY SERV, BANK OF AMERICA ONLINE PMT TDB****34454454 500.00
04/11 TD BILL PAY SERV, BANK OF AMERICA ONLINE PMT TDB****34454454 4,313.34
04/11 ELECTRONIC PMT-WEB, EEERERERERE MTG PAYMENTS ****311244 6,208.46
04/12 ELECTRONIC PMT-WEB, ERERERE CARD RER PAYMNT ****63101563793 292.65
04/12 ELECTRONIC PMT-WEB, REREREER CK WEBXFR TRANSFER ****062167 1,000.00
Subtotal: 17,537.72
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
	// {Category:Electronic Deposits PostingDate:2023-03-22 00:00:00 +0000 UTC Description:TD ZELLE RECEIVED, erer434ree r34rere re5rerer4344re Amount:418}
	// {Category:Electronic Deposits PostingDate:2023-04-02 00:00:00 +0000 UTC Description:ACH DEPOSIT, rere erer ererereL Amount:6377.3}
	// {Category:Electronic Deposits PostingDate:2023-04-08 00:00:00 +0000 UTC Description:ACH DEPOSIT, fererer  erereer dfdferr Amount:6377.3}
	// {Category:Electronic Deposits PostingDate:2023-04-18 00:00:00 +0000 UTC Description:ACH DEPOSIT, rere rer4tr rtrtrrtr Amount:3745.84}
	// {Category:Electronic Payments PostingDate:2023-03-28 00:00:00 +0000 UTC Description:ELECTRONIC PMT-WEB, BKOFAM CK WEBXFR TRANSFER ****234533 Amount:1000}
	// {Category:Electronic Payments PostingDate:2023-03-28 00:00:00 +0000 UTC Description:TD BILL PAY SERV, BANK OF AMERICA ONLINE PMT TDB****34454454 Amount:4223.27}
	// {Category:Electronic Payments PostingDate:2023-04-11 00:00:00 +0000 UTC Description:TD BILL PAY SERV, BANK OF AMERICA ONLINE PMT TDB****34454454 Amount:500}
	// {Category:Electronic Payments PostingDate:2023-04-11 00:00:00 +0000 UTC Description:TD BILL PAY SERV, BANK OF AMERICA ONLINE PMT TDB****34454454 Amount:4313.34}
	// {Category:Electronic Payments PostingDate:2023-04-11 00:00:00 +0000 UTC Description:ELECTRONIC PMT-WEB, EEERERERERE MTG PAYMENTS ****311244 Amount:6208.46}
	// {Category:Electronic Payments PostingDate:2023-04-12 00:00:00 +0000 UTC Description:ELECTRONIC PMT-WEB, ERERERE CARD RER PAYMNT ****63101563793 Amount:292.65}
	// {Category:Electronic Payments PostingDate:2023-04-12 00:00:00 +0000 UTC Description:ELECTRONIC PMT-WEB, REREREER CK WEBXFR TRANSFER ****062167 Amount:1000}
}
