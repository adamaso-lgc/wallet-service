package command

import (
	"context"
	"fmt"

	"github.com/adamaso/wallet-service/internal/domain"
)

type Deposit struct {
	WalletID  string
	Amount    float64
	Reference string
}

type DepositHandler struct {
	repo domain.WalletRepository
}

func NewDepositHandler(repo domain.WalletRepository) *DepositHandler {
	return &DepositHandler{repo: repo}
}

func (h *DepositHandler) Handle(ctx context.Context, cmd Deposit) error {
	wallet, err := h.repo.Get(ctx, cmd.WalletID)
	if err != nil {
		return fmt.Errorf("get wallet: %w", err)
	}
	if err := wallet.Deposit(cmd.Amount, cmd.Reference); err != nil {
		return fmt.Errorf("deposit: %w", err)
	}
	if err := h.repo.Save(ctx, wallet); err != nil {
		return fmt.Errorf("save wallet: %w", err)
	}
	return nil
}
