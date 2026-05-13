package application

import (
	"github.com/adamaso/wallet-service/internal/application/command"
	"github.com/adamaso/wallet-service/internal/application/query"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/projection"
)

// Application is the composition root for the application layer.
// It groups all command and query handlers and is the single entry point
// that cmd/server wires up with concrete infrastructure implementations.
type Application struct {
	Commands Commands
	Queries  Queries
}

// Commands groups all write-side handlers.
type Commands struct {
	CreateWallet *command.CreateWalletHandler
	Deposit      *command.DepositHandler
	Withdraw     *command.WithdrawHandler
	Transfer     *command.TransferHandler
	FreezeWallet *command.FreezeWalletHandler
}

// Queries groups all read-side handlers.
type Queries struct {
	GetWallet          *query.GetWalletHandler
	ListWalletsByOwner *query.ListWalletsByOwnerHandler
}

// New wires all handlers with their dependencies and returns a ready
// Application. Called once at startup from cmd/server/main.go.
func New(repo domain.WalletRepository, store projection.WalletStore) *Application {
	return &Application{
		Commands: Commands{
			CreateWallet: command.NewCreateWalletHandler(repo),
			Deposit:      command.NewDepositHandler(repo),
			Withdraw:     command.NewWithdrawHandler(repo),
			Transfer:     command.NewTransferHandler(repo),
			FreezeWallet: command.NewFreezeWalletHandler(repo),
		},
		Queries: Queries{
			GetWallet:          query.NewGetWalletHandler(store),
			ListWalletsByOwner: query.NewListWalletsByOwnerHandler(store),
		},
	}
}
