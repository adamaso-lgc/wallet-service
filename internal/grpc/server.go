package grpcserver

import (
	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application"
)

type Server struct {
	walletv1.UnimplementedWalletServiceServer
	app *application.Application
}

// Compile-time check that Server fully implements the generated interface.
var _ walletv1.WalletServiceServer = (*Server)(nil)

func NewServer(app *application.Application) *Server {
	return &Server{app: app}
}
