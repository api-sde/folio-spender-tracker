package main

import (
	"strconv"
	"strings"
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

func ParseNewAmount(textAmount string) (*Amount, error) {
	return ParseNewAmountWithCurrency(textAmount, Currency{})
}

func ParseNewAmountWithCurrency(textAmount string, currency Currency) (*Amount, error) {
	newAmount := new(Amount)

	if currency.Ticker != "" {
		newAmount.Currency = currency
	}

	if textAmount != "" {
		parts := strings.Split(textAmount, ".")

		if len(parts) == 1 {
			nominal, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, err
			}

			newAmount.SetNominal(int64(nominal))
		}

	}
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
