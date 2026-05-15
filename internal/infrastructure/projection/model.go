package projection

import (
	"time"
)

type WalletView struct {
	ID        string
	OwnerID   string
	Balance   float64
	Currency  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
