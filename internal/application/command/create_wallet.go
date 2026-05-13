package command

import (
	"context"
	"fmt"

	"github.com/adamaso/wallet-service/internal/domain"
)

type CreateWallet struct {
	OwnerID        string
	Currency       string
	InitialBalance float64
}

type CreateWalletResult struct {
	WalletID string
}

type CreateWalletHandler struct {
	repo domain.WalletRepository
}

func NewCreateWalletHandler(repo domain.WalletRepository) *CreateWalletHandler {
	return &CreateWalletHandler{repo: repo}
}

func (h *CreateWalletHandler) Handle(ctx context.Context, cmd CreateWallet) (CreateWalletResult, error) {
	wallet, err := domain.NewWallet(cmd.OwnerID, cmd.Currency, cmd.InitialBalance)
	if err != nil {
		return CreateWalletResult{}, fmt.Errorf("create wallet: %w", err)
	}
	if err := h.repo.Save(ctx, wallet); err != nil {
		return CreateWalletResult{}, fmt.Errorf("save wallet: %w", err)
	}
	return CreateWalletResult{WalletID: wallet.GetID()}, nil
}
