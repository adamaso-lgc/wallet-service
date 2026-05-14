package command

import (
	"context"
	"fmt"

	"github.com/adamaso/wallet-service/internal/domain"
)

// Transfer moves funds between two wallets atomically via SaveAll.
// Both wallets are saved in a single transaction — either both succeed or
// neither does, preventing partial state.
type Transfer struct {
	SourceWalletID      string
	DestinationWalletID string
	Amount              float64
	Reference           string
}

type TransferHandler struct {
	repo domain.WalletRepository
}

func NewTransferHandler(repo domain.WalletRepository) *TransferHandler {
	return &TransferHandler{repo: repo}
}

func (h *TransferHandler) Handle(ctx context.Context, cmd Transfer) error {
	if cmd.SourceWalletID == cmd.DestinationWalletID {
		return fmt.Errorf("cannot transfer to self")
	}

	source, err := h.repo.Get(ctx, cmd.SourceWalletID)
	if err != nil {
		return fmt.Errorf("get source wallet: %w", err)
	}
	destination, err := h.repo.Get(ctx, cmd.DestinationWalletID)
	if err != nil {
		return fmt.Errorf("get destination wallet: %w", err)
	}
	if err := source.DebitForTransfer(cmd.Amount, cmd.DestinationWalletID, cmd.Reference); err != nil {
		return fmt.Errorf("debit: %w", err)
	}
	if err := destination.CreditForTransfer(cmd.Amount, cmd.SourceWalletID, cmd.Reference); err != nil {
		return fmt.Errorf("credit: %w", err)
	}
	if err := h.repo.SaveAll(ctx, source, destination); err != nil {
		return fmt.Errorf("save transfer: %w", err)
	}
	return nil
}
