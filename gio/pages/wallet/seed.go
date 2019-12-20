package wallet 

import (
	"image"

	"gioui.org/unit"
	"gioui.org/layout"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
)

type (
	informationScreen struct {

	}

	SeedPage struct {
		multiWallet       *dcrlibwallet.MultiWallet 
		createWalletPage  *CreateWalletPage
		currentScreen     string
		informationScreen *informationScreen

		checkboxes        []*widgets.Checkbox 
		labels            [][]*widgets.Label
		viewSeedButton    *widgets.Button
	}
)

func NewSeedPage(multiWallet *dcrlibwallet.MultiWallet, createWalletPage *CreateWalletPage) *SeedPage {
	s := &SeedPage{
		multiWallet      : multiWallet,
		createWalletPage : createWalletPage,
		currentScreen    : "reminderScreen",
	}

	s.prepareInformationScreenWidgets()

	return s
}

func (s *SeedPage) prepareInformationScreenWidgets() {
	numOfCheckboxes := 5
	
	s.viewSeedButton = widgets.NewButton("View seed phrase", nil)
	s.checkboxes = make([]*widgets.Checkbox, numOfCheckboxes)
	s.labels     = make([][]*widgets.Label, numOfCheckboxes)

	for index := range s.checkboxes {
		s.checkboxes[index] = widgets.NewCheckbox()
	}

	s.labels[0] = []*widgets.Label{
		widgets.NewLabel("The 33-word seed phrase is").SetSize(4),
		widgets.NewLabel("EXTREMELY IMPORTANT.").SetSize(4),
	} 

	s.labels[1] = []*widgets.Label{
		widgets.NewLabel("Seed phrase iss the only way to").SetSize(4),
		widgets.NewLabel("restore your wallet.").SetSize(4),
	}

	s.labels[2] = []*widgets.Label{
		widgets.NewLabel("It is recommended to store your seed").SetSize(4),
		widgets.NewLabel("phrase in a physical format (e.g.").SetSize(4),
		widgets.NewLabel("write down on a paper).").SetSize(4),
	}

	s.labels[3] = []*widgets.Label{
		widgets.NewLabel("It is highly discouraged to store your").SetSize(4),
		widgets.NewLabel("seed phrase in any digital format").SetSize(4),
		widgets.NewLabel("(e.g. screenshot).").SetSize(4),
	}

	s.labels[4] = []*widgets.Label{
		widgets.NewLabel("Anyone with your seed phrase can").SetSize(4),
		widgets.NewLabel("steal your funds. DO NOT show it to").SetSize(4),
		widgets.NewLabel("anyone.").SetSize(4),
	}
}

func (s *SeedPage) render(ctx *layout.Context, refreshWindowFunc func()) {
	if s.currentScreen == "reminderScreen" {
		s.drawReminderScreen(ctx, refreshWindowFunc)
	} else if s.currentScreen == "seedPhraseScreen" {
		s.drawSeedPhraseScreen(ctx, refreshWindowFunc)
	}
}

func (s *SeedPage) drawReminderScreen(ctx *layout.Context, refreshWindowFunc func()) {
	inset := layout.Inset{
		Left: unit.Dp(20),
		Right: unit.Dp(20),
	}
	inset.Layout(ctx, func(){
		s.drawReminderItems(ctx, refreshWindowFunc)
	})

	inset = layout.Inset{
		Top: unit.Dp(350),
	}
	inset.Layout(ctx, func(){
		bounds := image.Point{
			X: 700,
			Y: 400,
		}
		helper.PaintArea(ctx, helper.WhiteColor, bounds)

		inset := layout.UniformInset(unit.Dp(20))
		inset.Layout(ctx, func(){
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
			ctx.Constraints.Height.Min = 50

			bgcol := helper.DecredLightBlueColor 
			if !s.hasCheckedAllReminders() {
				bgcol = helper.GrayColor
			}

			s.viewSeedButton.
				SetBackgroundColor(bgcol).
				Draw(ctx, func(){
					if s.hasCheckedAllReminders() {
						s.currentScreen = "seedPhraseScreen"
						refreshWindowFunc()
					}
				})
		})
	})
}

func (s *SeedPage) drawReminderItems(ctx *layout.Context, refreshWindowFunc func()) {
	outerTopInset := 0
	for index := range s.checkboxes {
		cindex := index

		inset := layout.Inset{
			Top: unit.Dp(float32(outerTopInset)),
			Left: unit.Dp(35),
		}
		inset.Layout(ctx, func(){
			layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
				layout.Rigid(func(){
					inset := layout.Inset{
						Top: unit.Dp(10),
					}
					inset.Layout(ctx, func(){
						s.checkboxes[cindex].Draw(ctx)
					})
				}),
				layout.Rigid(func(){
					innerTopInset := 0
					for i := range s.labels[cindex] {
						inset := layout.Inset{
							Top: unit.Dp(float32(innerTopInset)),
							Left: unit.Dp(60),
						}
						inset.Layout(ctx, func(){
							s.labels[cindex][i].Draw(ctx)
						})
						innerTopInset += 20
					}
				}),
			)
		})
		outerTopInset += (22 * len(s.labels[cindex])) + 7
	}
}

func (s *SeedPage) hasCheckedAllReminders() bool {
	for i := range s.checkboxes {
		if !s.checkboxes[i].IsChecked() {
			return false
		}
	}

	return true
}

func (s *SeedPage) drawSeedPhraseScreen(ctx *layout.Context, refreshWindowFunc func()) {
	inset := layout.Inset{
		Left: unit.Dp(20),
		Right: unit.Dp(20),
	}
	inset.Layout(ctx, func(){
		widgets.NewLabel("Write down all 33 words in the correct order.").
			SetSize(5).
			Draw(ctx)
	})
}