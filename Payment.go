package main

import "time"

type Payment struct {
	Date     time.Time
	Year     string
	Month    string
	Name     string
	Category string
	Cashback float32
	Debit    float32
	Credit   float32
}
