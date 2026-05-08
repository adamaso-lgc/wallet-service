package domain

const (
	EventWalletCreated    EventType = "WalletCreated"
	EventWalletDeposited  EventType = "WalletDeposited"
	EventMoneyWithdrawn   EventType = "MoneyWithdrawn"
	EventMoneyTransferred EventType = "MoneyTransferred"
	EventWalletFrozen     EventType = "WalletFrozen"
)

type WalletCreatedEvent struct {
	BaseEvent
	OwnerID        string  `json:"owner_id"`
	Currency       string  `json:"currency"`
	InitialBalance float64 `json:"initial_balance"`
}

type MoneyDepositedEvent struct {
	BaseEvent
	Amount       float64 `json:"amount"`
	BalanceAfter float64 `json:"balance_after"`
	Reference    string  `json:"reference"`
}

type MoneyWithdrawnEvent struct {
	BaseEvent
	Amount       float64 `json:"amount"`
	BalanceAfter float64 `json:"balance_after"`
	Reference    string  `json:"reference"`
}

type MoneyTransferredEvent struct {
	BaseEvent
	Amount         float64 `json:"amount"`
	BalanceAfter   float64 `json:"balance_after"`
	CounterpartyID string  `json:"counterparty_id"`
	Direction      string  `json:"direction"`
	Reference      string  `json:"reference"`
}

type WalletFrozenEvent struct {
	BaseEvent
	Reference string `json:"reference"`
}
