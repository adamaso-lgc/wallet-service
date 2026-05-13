package command

import (
	"context"
	"fmt"

	"github.com/adamaso/wallet-service/internal/domain"
)

type Withdraw struct {
	WalletID  string
	Amount    float64
	Reference string
}

type WithdrawHandler struct {
	repo domain.WalletRepository
}

func NewWithdrawHandler(repo domain.WalletRepository) *WithdrawHandler {
	return &WithdrawHandler{repo: repo}
}

func (h *WithdrawHandler) Handle(ctx context.Context, cmd Withdraw) error {
	wallet, err := h.repo.Get(ctx, cmd.WalletID)
	if err != nil {
		return fmt.Errorf("get wallet: %w", err)
	}
	if err := wallet.Withdraw(cmd.Amount, cmd.Reference); err != nil {
		return fmt.Errorf("withdraw: %w", err)
	}
	if err := h.repo.Save(ctx, wallet); err != nil {
		return fmt.Errorf("save wallet: %w", err)
	}
	return nil
}
