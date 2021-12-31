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

	singleDigit := Amount{}
	singleDigit.SetCent(1)
	if singleDigit.Nominal != 0 || singleDigit.Cent != 1 || singleDigit.Sum() != 1 {
		t.Errorf("Invalid Sum and after setting single digit Cent: %v", singleDigit.Sum())
	}

	sumBoth := Amount{}
	sumBoth.SetNominal(100)
	sumBoth.SetCent(1)
	if sumBoth.Sum() != 10001 {
		t.Errorf("Invalid Sum and after setting Nominal & Cent: %v", sumBoth.Sum())
	}
}

func TestAmount_Sum_Negative(t *testing.T) {
	negativeCases := []struct {
		nominal     int64
		cent        int64
		expectedSum int64
	}{
		{nominal: 1, cent: 95, expectedSum: 195},
		{nominal: -1, cent: 95, expectedSum: 195},
		{nominal: 1, cent: -95, expectedSum: 195},
		{nominal: -1, cent: -95, expectedSum: 195},
	}

	for _, test := range negativeCases {
		t.Run(fmt.Sprintf("Negative case: %v", test), func(t *testing.T) {
			neg := Amount{}
			neg.SetNominal(test.nominal)
			neg.SetCent(test.cent)

			if neg.Sum() != test.expectedSum {
				t.Errorf("Expected %v, got %v", test.expectedSum, neg.Sum())
				t.Log(neg)
			}
		})
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
func TestAmount_ToTextWithSingleCent(t *testing.T) {
	amount := Amount{
		Cent: 5,
	}

	if amount.ToText() != "0.05" {
		t.Errorf("Invalid format: %s instead of %s", amount.ToText(), "0.05")
	}

	amountNom := Amount{
		Nominal: 2,
		Cent:    5,
	}

	if amountNom.ToText() != "2.05" {
		t.Errorf("Invalid format: %s instead of %s", amountNom.ToText(), "2.05")
	}

	amountOnlyCents := Amount{}
	amountOnlyCents.SetCent(205)
	if amountOnlyCents.ToText() != "2.05" {
		t.Errorf("Invalid format: %s instead of %s", amountOnlyCents.ToText(), "2.05")
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

func TestParseNewAmount(t *testing.T) {

	oneNominalAmount := Amount{}
	oneNominalAmount.SetNominal(1)

	oneDollarOneCent := Amount{}
	oneDollarOneCent.SetCent(101)

	oneDollarTenCent := Amount{}
	oneDollarTenCent.SetCent(110)

	ninetyNine := Amount{}
	ninetyNine.SetCent(99)

	testCases := []struct {
		amount       string
		expected     Amount
		errorMessage string
	}{
		{"1", oneNominalAmount, ""},
		{"01", oneNominalAmount, ""},
		{" 1", oneNominalAmount, ""},
		{"1 CAD", Amount{}, "strconv.Atoi: parsing \"1CAD\": invalid syntax"},
		{"1.01", oneDollarOneCent, ""},
		{" 1.01", oneDollarOneCent, ""},
		{" 1.01 ", oneDollarOneCent, ""},
		{" 1. 01 ", oneDollarOneCent, ""},
		{" 1.10 ", oneDollarTenCent, ""},
		{" 1,10 ", oneDollarTenCent, ""},
		{" 0,99 ", ninetyNine, ""},
		{" 0.99 ", ninetyNine, ""},
	}

	for _, test := range testCases {
		t.Run("Test amount: "+test.amount, func(t *testing.T) {
			result, err := ParseNewAmount(test.amount)

			if err != nil {
				errorMessage := err.Error()
				if errorMessage != test.errorMessage {
					t.Errorf("Unexpected error message: %s", errorMessage)
				} else {
					t.Logf("Expected error message: %s", errorMessage)
					return
				}

				if result == nil {
					t.Errorf("nil result, stopping test.")
					return
				}
			}

			if result.Sum() != test.expected.Sum() {
				t.Errorf("Unexpected Sum: %v instead of %v", result.Sum(), test.expected.Sum())
			}

			if result.Nominal != test.expected.Nominal {
				t.Errorf("Unexpected Nominal: %v instead of %v", result.Nominal, test.expected.Nominal)
			}

			if result.Cent != test.expected.Cent {
				t.Errorf("Unexpected Cent: %v instead of %v", result.Cent, test.expected.Cent)
			}
		})
	}
}

func TestParseNewAmountWithCurrency(t *testing.T) {
	result, err := ParseNewAmountWithCurrency("1 CAD", Currency{})
	if result == nil || result.Nominal != 1 || result.Currency.Ticker != "CAD" {
		t.Errorf("Invalid result: %v, expected %s", result, "1 CAD")
		t.Log(err)
	}

	resultCurr, err := ParseNewAmountWithCurrency("1 ", Currency{Ticker: "CAD"})
	if resultCurr == nil || resultCurr.Nominal != 1 || resultCurr.Currency.Ticker != "CAD" {
		t.Errorf("Invalid result: %v, expected %s", result, "1 CAD")
		t.Log(err)
	}

	resultNoTicker, err := ParseNewAmountWithCurrency("1 ", Currency{})
	if resultNoTicker == nil || resultNoTicker.Nominal != 1 {
		t.Errorf("Invalid result: %v, expected %s", resultNoTicker, "1")
		t.Log(err)
	}
}
