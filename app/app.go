package app

import (
	"context"

	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/wallet"
)

const DisplayName = "GoDCR"

type UserInterface interface {
	// DisplayPreLaunchError produces a minimal interface to
	// display errors that occur before the main app interface is loaded.
	DisplayPreLaunchError(errorMessage string)

	// LaunchApp loads main app interface for use.
	LaunchApp(ctx context.Context, cfg *config.Config, wallet wallet.Wallet)
}
