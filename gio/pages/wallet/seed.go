package wallet 

import (
	"image"
	"math"
	"strings"
	"strconv"

	"gioui.org/text"
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

		goToVerifyScreenButton *widgets.Button

		seedWords 		  string
		seedColumns       [][]string
		err				  error
	}
)

func NewSeedPage(multiWallet *dcrlibwallet.MultiWallet, createWalletPage *CreateWalletPage) *SeedPage {
	s := &SeedPage{
		multiWallet      : multiWallet,
		createWalletPage : createWalletPage,
		currentScreen    : "reminderScreen",
	}

	s.prepareInformationScreenWidgets()
	s.seedWords, s.err = helper.GenerateSeedWords()
	if s.err == nil {
		s.seedColumns = make([][]string, 3)
		
		allWords := strings.Split(s.seedWords, " ")
		maxWordCountPerColumn := int(math.Ceil(float64(len(allWords)) / 3.0))
		s.seedColumns[0] = allWords[:maxWordCountPerColumn] 
		s.seedColumns[1] = allWords[maxWordCountPerColumn : maxWordCountPerColumn*2]
		s.seedColumns[2] = allWords[maxWordCountPerColumn*2:]
	}

	return s
}

func (s *SeedPage) prepareInformationScreenWidgets() {
	numOfCheckboxes := 5
	
	s.viewSeedButton = widgets.NewButton("View seed phrase", nil)
	s.goToVerifyScreenButton = widgets.NewButton("I have written down all 33 words", nil).SetBackgroundColor(helper.DecredLightBlueColor)
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
		Top: unit.Dp(30),
		Left: unit.Dp(helper.StandaloneScreenPadding),
		Right: unit.Dp(helper.StandaloneScreenPadding),
	}
	inset.Layout(ctx, func(){
		s.drawReminderItems(ctx, refreshWindowFunc)
	})

	inset = layout.Inset{
		Top: unit.Dp(410),
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
						innerTopInset += 23
					}
				}),
			)
		})
		outerTopInset += (27 * len(s.labels[cindex])) + 7
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
		Top: unit.Dp(30),
		Left: unit.Dp(helper.StandaloneScreenPadding),
		Right: unit.Dp(helper.StandaloneScreenPadding),
	}
	inset.Layout(ctx, func(){
		widgets.NewLabel("Write down all 33 words in the correct order.").
			SetSize(5).
			Draw(ctx)
	})	

	seedCardHeight := ctx.Constraints.Height.Max - 175

	inset = layout.Inset{
		Top: unit.Dp(55),
		Left: unit.Dp(helper.StandaloneScreenPadding),
		Right: unit.Dp(helper.StandaloneScreenPadding),
	}
	inset.Layout(ctx, func(){
		bounds := image.Point{
			X: ctx.Constraints.Width.Max,
			Y: seedCardHeight,
		}
		helper.PaintArea(ctx, helper.WhiteColor, bounds)

		layout.Stack{}.Layout(ctx, 
			layout.Expanded(func(){
				inset := layout.Inset{
					Top: unit.Dp(15),
					Left: unit.Dp(15),
					Right: unit.Dp(15),
				}
				inset.Layout(ctx, func(){
					currentItem := 1
					layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
						layout.Rigid(func(){
							inset := layout.Inset{
								Left: unit.Dp(5),
							}
							inset.Layout(ctx, func(){
								drawColumn(ctx, s.seedColumns[0], &currentItem)
							})
						}),
						layout.Rigid(func(){
							inset := layout.Inset{
								Left: unit.Dp(70),
							}
							inset.Layout(ctx, func(){
								drawColumn(ctx, s.seedColumns[1], &currentItem)
							})
						}),
						layout.Flexed(1, func(){
							inset := layout.Inset{
								Left: unit.Dp(65),
							}
							inset.Layout(ctx, func(){
								drawColumn(ctx, s.seedColumns[2], &currentItem)
							})
						}),
					)
				})	
			}),
		)
	})

	inset = layout.Inset{
		Top: unit.Dp(float32(seedCardHeight + 60)),
	}
	inset.Layout(ctx, func(){
		bounds := image.Point{
			X: ctx.Constraints.Width.Max,
			Y: 200,
		}
		helper.PaintArea(ctx, helper.WhiteColor, bounds)

		inset := layout.UniformInset(unit.Dp(20))
		inset.Layout(ctx, func(){
			widgets.NewLabel("You will be asked to enter the seed phrase on the next screen").SetSize(5).Draw(ctx)
		})

		inset = layout.Inset{
			Top: unit.Dp(45),
			Left: unit.Dp(20),
			Right: unit.Dp(20),
		}
		inset.Layout(ctx, func(){
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max 
			ctx.Constraints.Height.Min = 50
			s.goToVerifyScreenButton.Draw(ctx, func(){

			})
		})
	})
}

func drawColumn(ctx *layout.Context, words []string, currentItem *int) {
	topInset := 0
	for i := range words {
		inset := layout.Inset{
			Top: unit.Dp(float32(topInset)),
		}
		inset.Layout(ctx, func(){
			widgets.NewLabel(strconv.Itoa(*currentItem) + ".) " + words[i]).
				SetWeight(text.Bold).
				SetSize(5).
				Draw(ctx)
		})
		topInset += 26
		*currentItem++
	}
}
