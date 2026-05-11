package domain

import "fmt"

type Money struct {
	amount   float64
	currency string
}

func NewMoney(amount float64, currency string) (Money, error) {
	if amount <= 0 {
		return Money{}, ErrInvalidAmount
	}
	if currency == "" {
		return Money{}, ErrCurrencyRequired
	}
	return Money{amount, currency}, nil
}

func (m Money) Amount() float64  { return m.amount }
func (m Money) Currency() string { return m.currency }

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("%w: %s vs %s", ErrCurrencyMismatch, m.currency, other.currency)
	}
	return Money{m.amount + other.amount, m.currency}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("%w: %s vs %s", ErrCurrencyMismatch, m.currency, other.currency)
	}
	if m.amount < other.amount {
		return Money{}, ErrInsufficientFunds
	}
	return Money{m.amount - other.amount, m.currency}, nil
}
