package application

// Commands are plain data structs — they carry intent, no logic.

type CreateWalletCommand struct {
	OwnerID        string
	Currency       string
	InitialBalance float64
}

type CreateWalletResult struct {
	WalletID string
}

type DepositCommand struct {
	WalletID  string
	Amount    float64
	Reference string
}

type WithdrawCommand struct {
	WalletID  string
	Amount    float64
	Reference string
}

// TransferCommand moves funds between two wallets atomically via SaveAll.
type TransferCommand struct {
	SourceWalletID      string
	DestinationWalletID string
	Amount              float64
	Reference           string
}

type FreezeWalletCommand struct {
	WalletID  string
	Reference string
}
