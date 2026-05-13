package application

import (
	"github.com/adamaso/wallet-service/internal/application/command"
	"github.com/adamaso/wallet-service/internal/application/query"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/projection"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateWallet *command.CreateWalletHandler
	Deposit      *command.DepositHandler
	Withdraw     *command.WithdrawHandler
	Transfer     *command.TransferHandler
	FreezeWallet *command.FreezeWalletHandler
}
type Queries struct {
	GetWallet          *query.GetWalletHandler
	ListWalletsByOwner *query.ListWalletsByOwnerHandler
}

func NewApplication(repo domain.WalletRepository, store projection.WalletStore) *Application {
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
