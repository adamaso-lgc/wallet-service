package domain

import (
	"errors"
	"fmt"
)

type Money struct {
	amount   float64
	currency string
}

func NewMoney(amount float64, currency string) (Money, error) {
	if amount <= 0 {
		return Money{}, errors.New("amount must be greater than zero")
	}
	if currency == "" {
		return Money{}, errors.New("currency is required")
	}
	return Money{amount, currency}, nil
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.currency, other.currency)
	}
	return Money{m.amount + other.amount, m.currency}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("currency mismatch: %s vs %s", m.currency, other.currency)
	}
	if m.amount < other.amount {
		return Money{}, errors.New("insufficient funds")
	}
	return Money{m.amount - other.amount, m.currency}, nil
}
