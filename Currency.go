package main

import (
	"errors"
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
	newAmount := new(Amount)

	if textAmount == "" {
		return nil, errors.New("empty text amount")
	}

	textAmount = strings.ReplaceAll(textAmount, " ", "")

	var parts []string
	if strings.Contains(textAmount, ".") {
		parts = strings.Split(textAmount, ".")
	} else if strings.Contains(textAmount, ",") {
		parts = strings.Split(textAmount, ",")
	} else {
		parts = append(parts, textAmount)
	}

	switch len(parts) {
	case 1:
		nominal, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		newAmount.SetNominal(int64(nominal))

	case 2:
		nominal, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}

		cent, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		newAmount.SetNominal(int64(nominal))
		newAmount.SetCent(int64(cent))

	default:
		return nil, errors.New("Invalid amount, failed to parse: " + textAmount)
	}

	return newAmount, nil
}

func ParseNewAmountWithCurrency(textAmount string, currency Currency) (*Amount, error) {

	tickerSpaceIndex := strings.LastIndex(textAmount, " ")
	onlyAmount, ticker := textAmount[:tickerSpaceIndex], textAmount[tickerSpaceIndex:]

	newAmount, err := ParseNewAmount(onlyAmount)
	if err != nil {
		return nil, err
	}

	if currency.Ticker != "" {
		newAmount.Currency = currency
	} else if trimmedTicker := strings.TrimSpace(ticker); trimmedTicker != "" {
		newAmount.Currency = Currency{Ticker: trimmedTicker}
	}

	return newAmount, nil
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
