package main

import (
	"fmt"
	"testing"
)

func TestNewZeroCurrency(t *testing.T) {
	curr := Currency{}

	if curr.Ticker != "" || curr.Name != "" {
		t.Errorf("Expected zero currency, got %v", curr)
	}
}

func TestCadCurrency(t *testing.T) {
	cad := Currency{
		Ticker: "CAD",
		Name:   "Canadian Dollar",
	}

	if cad.Ticker != "CAD" || cad.Name != "Canadian Dollar" {
		t.Errorf("Expected CAD currency, got %v", cad)
	}
}

func TestZeroAmount_SumCurrency(t *testing.T) {
	amount := Amount{}

	if amount.sum != 0 {
		t.Errorf("Expected sum 0 got %v instead", amount.sum)
	}

	if amount.Currency.Ticker != "" {
		t.Errorf("Expected zero currency, got %v", amount.Currency)
	}
}

func TestAmount_SetNominal(t *testing.T) {
	amount := Amount{Nominal: 100}

	if amount.Nominal != 100 {
		t.Errorf("Incorrect nominal amount: %v instead of %v", amount.Nominal, 100)
	}
}

func TestAmount_SetCent(t *testing.T) {
	amount := Amount{Cent: 100}

	if amount.Cent != 100 {
		t.Errorf("Incorrect Cent amount: %v instead of %v", amount.Cent, 100)
	}
}

func TestAmount_Sum(t *testing.T) {
	amount := Amount{sum: 123}

	if amount.Sum() != 123 {
		t.Errorf("Incorrect Cent amount: %v instead of %v", amount.sum, 123)
	}

	amount.SetCent(77)
	if amount.Cent != 77 || amount.Sum() != 77 {
		t.Error("Invalid Sum and after setting Cent")
	}

	amount.SetNominal(3)
	if amount.Nominal != 3 || amount.Sum() != 377 {
		t.Errorf("Invalid Sum and after setting Nominal: got %v instead of %v", amount.Sum(), 377)
	}

	amount.SetCent(333)
	if amount.Cent != 33 || amount.Sum() != 633 {
		t.Error("Invalid Sum and after setting Cent")
	}

	amount.SetNominal(1)
	amount.SetCent(123)
	if amount.Sum() != 223 {
		t.Errorf("Invalid Sum and after setting Nominal & Cent: got %v instead of %v", amount.Sum(), 223)
	}

	newAmount := Amount{}
	newAmount.SetCent(333)
	if newAmount.Nominal != 3 || newAmount.Cent != 33 || newAmount.Sum() != 333 {
		t.Error("Invalid Sum and after setting Cent")
	}
}

func TestAmount_ToText(t *testing.T) {
	amount := Amount{
		Currency: Currency{
			Ticker: "CAD",
			Name:   "",
		},
		sum:     0,
		Nominal: 0,
		Cent:    0,
	}
	amount.SetCent(450)

	if amount.ToText() != "4.50" {
		t.Errorf("Unexpected currency text format: got %v instead of %v", amount.ToText(), "4.50")
	} else {
		fmt.Println("No ticker:")
		fmt.Println(amount.ToText())
	}

}

func TestAmount_ToTextCurrency(t *testing.T) {
	amount := Amount{
		Currency: Currency{
			Ticker: "FR",
			Name:   "",
		},
		sum:     0,
		Nominal: 0,
		Cent:    0,
	}
	amount.SetNominal(9)
	amount.SetCent(99)

	if amount.ToTextCurrency() != "9.99 FR" {
		t.Errorf("Unexpected currency text format: got %v instead of %v", amount.ToTextCurrency(), "9.99 FR")
	} else {
		fmt.Println("With ticker:")
		fmt.Println(amount.ToTextCurrency())
	}
}

func TestAmount_ToTextCurrencyMissingTicker(t *testing.T) {
	amount := Amount{}
	amount.SetNominal(9)
	amount.SetCent(99)

	if amount.ToTextCurrency() != "9.99" {
		t.Errorf("Unexpected currency text format: got %v instead of %v", amount.ToTextCurrency(), "9.99")
	} else {
		fmt.Println("With missing ticker:")
		fmt.Println(amount.ToTextCurrency())
	}
}