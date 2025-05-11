package main

type TransactionType int

const (
	Debit TransactionType = iota
	Credit
)

type Transaction struct {
	Id     int
	Date   string
	Amount float64
	Type   TransactionType
}
