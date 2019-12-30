package pages 

import (
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/layout"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
	"github.com/raedahgroup/godcr/gio/pages/common"
)

type (
	OverviewPage struct {
		multiWallet *helper.MultiWallet
		syncer *common.Syncer

		err error
		totalBalance string

		transactions []*dcrlibwallet.Transaction
	}
)

func NewOverviewPage() *OverviewPage {
	return &OverviewPage{}
}

func (o *OverviewPage) BeforeRender(syncer *common.Syncer, multiWallet *helper.MultiWallet) {
	o.syncer = syncer 
	o.multiWallet = multiWallet

	o.totalBalance, o.err = o.multiWallet.TotalBalance()
}

func (o *OverviewPage) Render(ctx *layout.Context, refreshWindowFunc func()) {
	widgets.NewLabel("Overview").
		SetColor(helper.BlackColor).
		SetWeight(text.Bold).
		SetSize(6).
		Draw(ctx)

	inset := layout.Inset{
		Top: unit.Dp(45),
		Left: unit.Dp(15),
	}
	inset.Layout(ctx, func(){
		widgets.NewLabel(o.totalBalance).
			SetColor(helper.BlackColor).
			SetWeight(text.Bold).
			SetSize(7).
			Draw(ctx)
	})

	inset = layout.Inset{
		Top: unit.Dp(80),
		Left: unit.Dp(15),
	}
	inset.Layout(ctx, func(){
		widgets.NewLabel("Current Total Balance").
			SetSize(5).
			SetColor(helper.GrayColor).
			Draw(ctx)
	})

	var nextTopInset float32  = 105

	inset = layout.Inset{
		Top: unit.Dp(nextTopInset),
	}
	inset.Layout(ctx, func(){
		if len(o.transactions) == 0 {
			o.drawNoTransactionsCard(ctx)
			nextTopInset += 85 
		} else {
			o.drawRecentTransactionsCard(ctx)
			nextTopInset += 300
		}
	})

	inset = layout.Inset{
		Top: unit.Dp(nextTopInset),
	}
	inset.Layout(ctx, func(){
		o.syncer.Render(ctx)
	})
}

func (o *OverviewPage) drawNoTransactionsCard(ctx *layout.Context) {
	helper.PaintArea(ctx, helper.WhiteColor, ctx.Constraints.Width.Max, 80)

	inset := layout.UniformInset(unit.Dp(15))
	inset.Layout(ctx, func(){
		widgets.NewLabel("Recent Transactions").
			SetSize(5).
			SetColor(helper.BlackColor).
			SetWeight(text.Bold).
			Draw(ctx)

		inset := layout.Inset{
			Top: unit.Dp(40),
		}
		inset.Layout(ctx, func(){
			widgets.NewLabel("No transactions yet").
				SetSize(4).
				SetColor(helper.GrayColor).
				SetWeight(text.Bold).
				Draw(ctx)
		})
	})
}

func (o *OverviewPage) drawRecentTransactionsCard(ctx *layout.Context) {

}