package command

import (
	"context"
	"fmt"

	"github.com/adamaso/wallet-service/internal/domain"
)

type FreezeWallet struct {
	WalletID  string
	Reference string
}

type FreezeWalletHandler struct {
	repo domain.WalletRepository
}

func NewFreezeWalletHandler(repo domain.WalletRepository) *FreezeWalletHandler {
	return &FreezeWalletHandler{repo: repo}
}

func (h *FreezeWalletHandler) Handle(ctx context.Context, cmd FreezeWallet) error {
	wallet, err := h.repo.Get(ctx, cmd.WalletID)
	if err != nil {
		return fmt.Errorf("get wallet: %w", err)
	}
	if err := wallet.Freeze(cmd.Reference); err != nil {
		return fmt.Errorf("freeze: %w", err)
	}
	if err := h.repo.Save(ctx, wallet); err != nil {
		return fmt.Errorf("save wallet: %w", err)
	}
	return nil
}
