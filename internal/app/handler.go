package app

import (
	"context"
	"fmt"

	"github.com/adamaso/wallet-service/internal/domain"
)

// WalletCommandHandler handles all wallet use cases.
// It is the only entry point into the domain from the outside world.
type WalletCommandHandler struct {
	repo domain.WalletRepository
}

func NewWalletCommandHandler(repo domain.WalletRepository) *WalletCommandHandler {
	return &WalletCommandHandler{repo: repo}
}

func (h *WalletCommandHandler) CreateWallet(ctx context.Context, cmd CreateWalletCommand) (CreateWalletResult, error) {
	wallet, err := domain.NewWallet(cmd.OwnerID, cmd.Currency, cmd.InitialBalance)
	if err != nil {
		return CreateWalletResult{}, fmt.Errorf("create wallet: %w", err)
	}
	if err := h.repo.Save(ctx, wallet); err != nil {
		return CreateWalletResult{}, fmt.Errorf("save wallet: %w", err)
	}
	return CreateWalletResult{WalletID: wallet.GetID()}, nil
}

func (h *WalletCommandHandler) Deposit(ctx context.Context, cmd DepositCommand) error {
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

func (h *WalletCommandHandler) Withdraw(ctx context.Context, cmd WithdrawCommand) error {
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

// Transfer debits the source wallet and credits the destination wallet.
// The two saves are intentionally not atomic — see TransferCommand for details.
func (h *WalletCommandHandler) Transfer(ctx context.Context, cmd TransferCommand) error {
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

	if err := h.repo.Save(ctx, source); err != nil {
		return fmt.Errorf("save source wallet: %w", err)
	}
	if err := h.repo.Save(ctx, destination); err != nil {
		return fmt.Errorf("save destination wallet: %w", err)
	}
	return nil
}

func (h *WalletCommandHandler) FreezeWallet(ctx context.Context, cmd FreezeWalletCommand) error {
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
