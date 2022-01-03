package main

import (
	"time"
)

type Payment struct {
	Stamp    time.Time
	Date     string
	Year     string
	Month    string
	Name     string
	Category string
	Cashback *Amount
	Debit    *Amount
	Credit   *Amount
}

const (
	Unknown   = 0
	Tangerine = 1
	CIBC      = 2
	RBC       = 3
)

const (
	DEBIT  = "DEBIT"
	CREDIT = "CREDIT"
)

func (payment Payment) GetPaymentAmount() string {

	if payment.Debit != nil {
		return payment.Debit.ToText()
	} else if payment.Credit != nil {
		return payment.Credit.ToText()
	} else {
		return ""
	}
}
