package application

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application/command"
	"github.com/adamaso/wallet-service/internal/application/query"
	"github.com/adamaso/wallet-service/internal/domain"
	"github.com/adamaso/wallet-service/internal/infrastructure/projection"
)

type Application struct {
	walletv1.UnimplementedWalletServiceServer
	Commands Commands
	Queries  Queries
}

var _ walletv1.WalletServiceServer = (*Application)(nil)

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

func NewApplication(repo domain.WalletRepository, store projection.Repository) *Application {
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

func (a *Application) CreateWallet(ctx context.Context, req *walletv1.CreateWalletRequest) (*walletv1.CreateWalletResponse, error) {
	return a.Commands.CreateWallet.Handle(ctx, req)
}

func (a *Application) Deposit(ctx context.Context, req *walletv1.DepositRequest) (*emptypb.Empty, error) {
	return a.Commands.Deposit.Handle(ctx, req)
}

func (a *Application) Withdraw(ctx context.Context, req *walletv1.WithdrawRequest) (*emptypb.Empty, error) {
	return a.Commands.Withdraw.Handle(ctx, req)
}

func (a *Application) Transfer(ctx context.Context, req *walletv1.TransferRequest) (*emptypb.Empty, error) {
	return a.Commands.Transfer.Handle(ctx, req)
}

func (a *Application) FreezeWallet(ctx context.Context, req *walletv1.FreezeWalletRequest) (*emptypb.Empty, error) {
	return a.Commands.FreezeWallet.Handle(ctx, req)
}

func (a *Application) GetWallet(ctx context.Context, req *walletv1.GetWalletRequest) (*walletv1.WalletResponse, error) {
	return a.Queries.GetWallet.Handle(ctx, req)
}

func (a *Application) ListWalletsByOwner(ctx context.Context, req *walletv1.ListWalletsByOwnerRequest) (*walletv1.ListWalletsByOwnerResponse, error) {
	return a.Queries.ListWalletsByOwner.Handle(ctx, req)
}
