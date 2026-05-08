package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type WalletStatus string

const (
	StatusActive WalletStatus = "active"
	StatusFrozen WalletStatus = "frozen"
	StatusClosed WalletStatus = "closed"
)

type Wallet struct {
	BaseAggregate
	ownerID   string
	balance   Money
	status    WalletStatus
	createdAt time.Time
}

func NewWallet(ownerID string, currency string, initialBalance float64) (*Wallet, error) {
	if ownerID == "" {
		return nil, errors.New("ownerID is required")
	}
	if _, err := NewMoney(initialBalance, currency); err != nil {
		return nil, err
	}
	w := &Wallet{}
	event := WalletCreatedEvent{
		BaseEvent: BaseEvent{
			EventID:     uuid.New().String(),
			Type:        EventWalletCreated,
			AggregateID: uuid.New().String(),
			OccurredAt:  time.Now().UTC(),
			Version:     1,
		},
		OwnerID:        ownerID,
		Currency:       currency,
		InitialBalance: initialBalance,
	}

	if err := w.Raise(w, event); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Wallet) Deposit(amount float64, reference string) error {
	if err := w.ensureActive(); err != nil {
		return err
	}
	deposit, err := NewMoney(amount, w.balance.currency)
	if err != nil {
		return err
	}

	newBalance, _ := w.balance.Add(deposit)
	event := MoneyDepositedEvent{
		BaseEvent: BaseEvent{
			EventID:     uuid.New().String(),
			Type:        EventWalletDeposited,
			AggregateID: w.id,
			OccurredAt:  time.Now().UTC(),
			Version:     w.version + 1,
		},
		Amount:       amount,
		BalanceAfter: newBalance.amount,
		Reference:    reference,
	}

	if err := w.Raise(w, event); err != nil {
		return err
	}

	return nil
}

func (w *Wallet) Withdraw(amount float64, reference string) error {
	if err := w.ensureActive(); err != nil {
		return err
	}
	withdraw, err := NewMoney(amount, w.balance.currency)
	if err != nil {
		return err
	}
	if withdraw.amount == 0 {
		return errors.New("withdraw amount must be greater than zero")
	}

	newBalance, err := w.balance.Subtract(withdraw)
	if err != nil {
		return err
	}

	event := MoneyWithdrawnEvent{
		BaseEvent: BaseEvent{
			EventID:     uuid.New().String(),
			Type:        EventMoneyWithdrawn,
			AggregateID: w.id,
			OccurredAt:  time.Now().UTC(),
			Version:     w.version + 1,
		},
		Amount:       amount,
		BalanceAfter: newBalance.amount,
		Reference:    reference,
	}
	if err := w.Raise(w, event); err != nil {
		return err
	}

	return nil
}

func (w *Wallet) DebitForTransfer(amount float64, counterpartyID string, reference string) error {
	if err := w.ensureActive(); err != nil {
		return err
	}
	debit, err := NewMoney(amount, w.balance.currency)
	if err != nil {
		return err
	}
	if debit.amount == 0 {
		return errors.New("transfer amount must be greater than zero")
	}

	newBalance, err := w.balance.Subtract(debit)
	if err != nil {
		return err
	}

	event := MoneyTransferredEvent{
		BaseEvent: BaseEvent{
			EventID:     uuid.New().String(),
			Type:        EventMoneyTransferred,
			AggregateID: w.id,
			OccurredAt:  time.Now().UTC(),
			Version:     w.version + 1,
		},
		Amount:         amount,
		BalanceAfter:   newBalance.amount,
		CounterpartyID: counterpartyID,
		Direction:      "debit",
		Reference:      reference,
	}
	return w.Raise(w, event)
}

func (w *Wallet) CreditForTransfer(amount float64, counterpartyID string, reference string) error {
	if err := w.ensureActive(); err != nil {
		return err
	}
	credit, err := NewMoney(amount, w.balance.currency)
	if err != nil {
		return err
	}
	if credit.amount == 0 {
		return errors.New("transfer amount must be greater than zero")
	}

	newBalance, _ := w.balance.Add(credit)

	event := MoneyTransferredEvent{
		BaseEvent: BaseEvent{
			EventID:     uuid.New().String(),
			Type:        EventMoneyTransferred,
			AggregateID: w.id,
			OccurredAt:  time.Now().UTC(),
			Version:     w.version + 1,
		},
		Amount:         amount,
		BalanceAfter:   newBalance.amount,
		CounterpartyID: counterpartyID,
		Direction:      "credit",
		Reference:      reference,
	}
	return w.Raise(w, event)
}

func (w *Wallet) Freeze(reference string) error {
	if err := w.ensureActive(); err != nil {
		return err
	}

	event := WalletFrozenEvent{
		BaseEvent: BaseEvent{
			EventID:     uuid.New().String(),
			Type:        EventWalletFrozen,
			AggregateID: w.id,
			OccurredAt:  time.Now().UTC(),
			Version:     w.version + 1,
		},
		Reference: reference,
	}
	return w.Raise(w, event)
}

func (w *Wallet) Apply(event Event) error {
	switch e := event.(type) {
	case WalletCreatedEvent:
		return w.onWalletCreated(e)
	case MoneyDepositedEvent:
		w.balance.amount = e.BalanceAfter
		return nil
	case MoneyWithdrawnEvent:
		w.balance.amount = e.BalanceAfter
		return nil
	case MoneyTransferredEvent:
		w.balance.amount = e.BalanceAfter
		return nil
	case WalletFrozenEvent:
		w.status = StatusFrozen
		return nil
	default:
		return errors.New("unsupported event type")
	}
}

func (w *Wallet) onWalletCreated(e WalletCreatedEvent) error {
	w.id = e.AggregateID
	w.ownerID = e.OwnerID
	w.status = StatusActive
	w.createdAt = e.OccurredAt

	money, err := NewMoney(e.InitialBalance, e.Currency)
	if err != nil {
		return err
	}
	w.balance = money

	return nil
}

func (w *Wallet) ensureActive() error {
	if w.status != StatusActive {
		return fmt.Errorf("wallet is not active (status: %s)", w.status)
	}
	return nil
}
