package main

import "time"

type Payment struct {
	Stamp    time.Time
	Date     string
	Year     string
	Month    string
	Name     string
	Category string
	Cashback float32
	Debit    float32
	Credit   float32
}

const (
	Unknown   = 0
	Tangerine = 1
	CIBC
	RBC
)
