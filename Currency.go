package main

import (
	"strconv"
)

type Currency struct {
	Ticker string
	Name   string
}

type Amount struct {
	Currency Currency
	sum      int64 // Hundreds of Nominal + Cent for calculations
	Nominal  int64 // To do: Check for the actual economic definition, seems better than Main
	Cent     int64
}

func (amount Amount) Sum() int64 {
	return amount.sum
}

func (amount *Amount) updateSum() {
	amount.sum = amount.Nominal*100 + amount.Cent
}

func (amount *Amount) SetNominal(nominal int64) int64 {
	amount.Nominal = nominal
	amount.updateSum()

	return amount.Nominal
}

func (amount *Amount) SetCent(cent int64) int64 {
	if cent > 99 {
		centRemainder := cent % 100
		amount.Nominal += (cent - centRemainder) / 100
		cent = centRemainder
	}

	amount.Cent = cent
	amount.updateSum()

	return amount.Cent
}

func (amount Amount) ToText() string {
	return strconv.FormatInt(amount.Nominal, 10) + "." + strconv.FormatInt(amount.Cent, 10)
}

func (amount Amount) ToTextCurrency() string {
	if amount.Currency.Ticker == "" {
		return amount.ToText()
	}

	return amount.ToText() + " " + amount.Currency.Ticker
}
