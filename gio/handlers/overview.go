package handlers

import (
	"gioui.org/layout"

	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
)

type OverviewHandler struct {
}

func NewOverviewHandler() *OverviewHandler {
	return &OverviewHandler{}
}

func (o *OverviewHandler) BeforeRender(walletMiddleware app.WalletMiddleware, config *config.Settings) {

}

func (o *OverviewHandler) Render(ctx *layout.Context, refreshWindowFunc func()) {

}
