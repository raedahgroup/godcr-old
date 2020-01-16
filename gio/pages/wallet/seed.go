package wallet

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"

	"github.com/decred/dcrd/hdkeychain"
	"github.com/decred/dcrwallet/walletseed"
)

type screen int

type (
	reminderScreen struct {
		checkboxes     []*widgets.Checkbox
		labels         [][]*widgets.Label
		viewSeedButton *widgets.Button
	}

	seedScreen struct {
		err                    error
		wordsString            string
		words                  []string
		wordColumns            [][]string
		goToVerifyScreenButton *widgets.Button
	}

	verifyScreen struct {
		options            []*widgets.Selectable
		verifyButton       *widgets.Button
		failedVerification bool
	}

	SeedPage struct {
		createWalletPage *CreateWalletPage
		currentScreen    screen

		backButton          *widgets.Button
		backToWalletsButton *widgets.ClickableLabel

		reminderScreen *reminderScreen
		seedScreen     *seedScreen
		verifyScreen   *verifyScreen

		isPrepared bool
	}
)

const (
	REMINDER_SCREEN screen = iota 
	SEED_SCREEN 
	VERIFY_SCREEN 
	SUCCESS_SCREEN
)

func NewSeedPage(createWalletPage *CreateWalletPage) *SeedPage {
	s := &SeedPage{
		createWalletPage: createWalletPage,
		isPrepared:       false,
		currentScreen:   REMINDER_SCREEN,
	}

	s.backButton = widgets.NewButton("", widgets.NavigationArrowBackIcon).SetBackgroundColor(helper.BackgroundColor).SetColor(helper.BlackColor).MakeRound()
	s.backToWalletsButton = widgets.NewClickableLabel("Back to wallets").SetAlignment(widgets.AlignMiddle).SetSize(5).SetColor(helper.DecredLightBlueColor).SetWeight(text.Bold)

	return s
}

func (s *SeedPage) prepare(seed string) {
	// TODO remove block before if 
	seedS, _ := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	seed = walletseed.EncodeMnemonic(seedS)

	if !s.isPrepared {
		s.prepareInformationScreenWidgets(seed)
	}
}

func (s *SeedPage) prepareInformationScreenWidgets(seed string) {
	// reminder screen widgets
	numOfCheckboxes := 5
	s.reminderScreen = &reminderScreen{
		viewSeedButton: widgets.NewButton("View seed phrase", nil),
		checkboxes:     make([]*widgets.Checkbox, numOfCheckboxes),
		labels:         make([][]*widgets.Label, numOfCheckboxes),
	}

	for index := range s.reminderScreen.checkboxes {
		s.reminderScreen.checkboxes[index] = widgets.NewCheckbox()
	}

	s.reminderScreen.labels[0] = []*widgets.Label{
		widgets.NewLabel("The 33-word seed phrase is").SetSize(4),
		widgets.NewLabel("EXTREMELY IMPORTANT.").SetSize(4),
	}

	s.reminderScreen.labels[1] = []*widgets.Label{
		widgets.NewLabel("Seed phrase iss the only way to").SetSize(4),
		widgets.NewLabel("restore your wallet.").SetSize(4),
	}

	s.reminderScreen.labels[2] = []*widgets.Label{
		widgets.NewLabel("It is recommended to store your seed").SetSize(4),
		widgets.NewLabel("phrase in a physical format (e.g.").SetSize(4),
		widgets.NewLabel("write down on a paper).").SetSize(4),
	}

	s.reminderScreen.labels[3] = []*widgets.Label{
		widgets.NewLabel("It is highly discouraged to store your").SetSize(4),
		widgets.NewLabel("seed phrase in any digital format").SetSize(4),
		widgets.NewLabel("(e.g. screenshot).").SetSize(4),
	}

	s.reminderScreen.labels[4] = []*widgets.Label{
		widgets.NewLabel("Anyone with your seed phrase can").SetSize(4),
		widgets.NewLabel("steal your funds. DO NOT show it to").SetSize(4),
		widgets.NewLabel("anyone.").SetSize(4),
	}

	// seed screen widgets
	s.seedScreen = &seedScreen{
		goToVerifyScreenButton: widgets.NewButton("I have written down all 33 words", nil).SetBackgroundColor(helper.DecredLightBlueColor),
	}

	s.seedScreen.wordColumns = make([][]string, 3)
	s.seedScreen.wordsString = seed

	s.seedScreen.words = strings.Split(seed, " ")
	maxWordCountPerColumn := int(math.Ceil(float64(len(s.seedScreen.words)) / 3.0))
	s.seedScreen.wordColumns[0] = s.seedScreen.words[:maxWordCountPerColumn]
	s.seedScreen.wordColumns[1] = s.seedScreen.words[maxWordCountPerColumn : maxWordCountPerColumn*2]
	s.seedScreen.wordColumns[2] = s.seedScreen.words[maxWordCountPerColumn*2:]

	// verify screen widgets
	s.prepareVerifyScreenData()
	s.isPrepared = true
}

func (s *SeedPage) prepareVerifyScreenData() {
	s.verifyScreen = &verifyScreen{
		verifyButton: widgets.NewButton("Verify", nil),
	}
	s.verifyScreen.options = make([]*widgets.Selectable, len(s.seedScreen.words))

	allSeeds := dcrlibwallet.PGPWordList()
	rand.Seed(time.Now().Unix())

	for i := range s.verifyScreen.options {
		// TODO run in separate goroutines
		selectableItems := make([]string, 3)
		selectableItems[0] = s.seedScreen.words[i]

		for k := 1; k < 3; k++ {
			for {
				randomWord := getRandomWord(allSeeds)
				if randomWord != selectableItems[k-1] {
					selectableItems[k] = randomWord
					break
				}
			}
		}

		// shuffle items
		rand.Shuffle(len(selectableItems), func(i, j int) {
			selectableItems[i], selectableItems[j] = selectableItems[j], selectableItems[i]
		})
		s.verifyScreen.options[i] = widgets.NewSelectable(selectableItems)
	}
}

func getRandomWord(words []string) string {
	return words[rand.Intn(len(words))]
}


func (s *SeedPage) getWidgets(ctx *layout.Context, changePageFunc func(string)) []func() {
	switch s.currentScreen {
	case SEED_SCREEN:
		return s.getSeedPhraseScreenWidgets(ctx, changePageFunc)
	case VERIFY_SCREEN:
		return s.getVerifySeedPhraseScreenWidgets(ctx, changePageFunc)
	case SUCCESS_SCREEN:
		return s.getSuccessScreenWidgets(ctx, changePageFunc)
	default: // REMINDER_SCREEN
		return s.getReminderScreenWidgets(ctx, changePageFunc)
	}
}

func (s *SeedPage) getReminderScreenWidgets(ctx *layout.Context, changePageFunc func(string)) []func() {
	return []func(){
		// header row
		func() {
			drawHeader(ctx, func() {
			}, func() {
				widgets.NewLabel("Keep in mind").
					SetWeight(text.Bold).
					SetSize(6).
					Draw(ctx)
			})
		},

		// reminder items section
		func() {
			drawBody(ctx, nil, func(){
				list := layout.List{Axis: layout.Vertical}
				list.Layout(ctx, len(s.reminderScreen.checkboxes), func(i int){
					inset := layout.Inset{
						Top: unit.Dp(10),
					}
					inset.Layout(ctx, func(){
						layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
							layout.Rigid(func(){
								inset := layout.Inset{
									Top: unit.Dp(6),
								}
								inset.Layout(ctx, func() {
									ctx.Constraints.Width.Min = 38
									ctx.Constraints.Height.Min = 40
									s.reminderScreen.checkboxes[i].Draw(ctx)
								})
							}),
							layout.Flexed(1, func(){
								list := layout.List{Axis: layout.Vertical}
								list.Layout(ctx, len(s.reminderScreen.labels[i]), func(j int){
									inset := layout.Inset{
										Top: unit.Dp(3),
										Left: unit.Dp(30),
									}
									inset.Layout(ctx, func(){
										s.reminderScreen.labels[i][j].Draw(ctx)
									})
								})
							}),
						)
					})
				})
			})
		},

		// footer row 
		func() {
			drawFooter(ctx, 120,  func() {
				ctx.Constraints.Height.Min = 50
		
				bgcol := helper.DecredLightBlueColor
				if !s.hasCheckedAllReminders() {
					bgcol = helper.GrayColor
				}
				
				s.reminderScreen.viewSeedButton.SetBackgroundColor(bgcol).Draw(ctx, func() {
					if s.hasCheckedAllReminders() {
						s.currentScreen = SEED_SCREEN
					}
				})
			})
		},
	}
}

func (s *SeedPage) hasCheckedAllReminders() bool {
	for i := range s.reminderScreen.checkboxes {
		if !s.reminderScreen.checkboxes[i].IsChecked() {
			return false
		}
	}
	return true
}

func (s *SeedPage) getSeedPhraseScreenWidgets(ctx *layout.Context, changePageFunc func(string)) []func() {
	return []func(){
		// header section 
		func(){
			drawHeader(ctx, func(){
				s.backButton.Draw(ctx, func() {
					s.currentScreen = REMINDER_SCREEN
				})
			}, func(){
				widgets.NewLabel("Write down seed phrase").SetWeight(text.Bold).SetSize(6).Draw(ctx)
				inset := layout.Inset{
					Top: unit.Dp(23),
				}
				inset.Layout(ctx, func() {
					widgets.NewLabel("Step 1/2").SetSize(4).Draw(ctx)
				})
			})
		},

		// body section 
		func(){
			currentNumber := 1
			columnWidths :=  ctx.Constraints.Width.Max / 3 

			drawCardBody(ctx,
				widgets.NewLabel("Write down all 33 words in the correct order.").SetSize(5),
				func(){
					layout.UniformInset(unit.Dp(10)).Layout(ctx, func(){
						layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
							layout.Rigid(func(){
								ctx.Constraints.Width.Min = columnWidths
								drawWordColumn(ctx, s.seedScreen.wordColumns[0], &currentNumber)
							}),
							layout.Rigid(func(){
								ctx.Constraints.Width.Min = columnWidths
								drawWordColumn(ctx, s.seedScreen.wordColumns[1], &currentNumber)
							}),
							layout.Rigid(func(){
								ctx.Constraints.Width.Min = columnWidths
								drawWordColumn(ctx, s.seedScreen.wordColumns[2], &currentNumber)
							}),
						)
					})
			})
		},

		// footer section 
		func() {
			drawFooter(ctx, 25, func(){
				inset := layout.UniformInset(unit.Dp(-10))
				inset.Left = unit.Dp(20)
				inset.Layout(ctx, func(){
					widgets.NewLabel("You will be asked to enter the seed phrase on the next screen").SetSize(5).Draw(ctx)
				})

				inset = layout.Inset{
					Top: unit.Dp(12),
				}
				inset.Layout(ctx, func(){
					ctx.Constraints.Height.Min = 50
					s.seedScreen.goToVerifyScreenButton.Draw(ctx, func() {
						s.currentScreen = VERIFY_SCREEN
					})
				})
			})
		},
	}
}

func drawWordColumn(ctx *layout.Context, words []string, currentNumber *int) {
	list := layout.List{Axis: layout.Vertical}
	list.Layout(ctx, len(words), func(i int){
		inset := layout.Inset{
			Top: unit.Dp(7),
		}
		inset.Layout(ctx, func(){
			widgets.NewLabel(strconv.Itoa(*currentNumber) + ") " + words[i]).
				SetWeight(text.Bold).
				SetSize(5).
				Draw(ctx)
		})
		*currentNumber++
	})
}

func (s *SeedPage) getVerifySeedPhraseScreenWidgets(ctx *layout.Context, changePageFunc func(string)) []func() {
	return []func(){
		func(){
			drawHeader(ctx, func(){
				s.backButton.Draw(ctx, func(){
					s.currentScreen = SEED_SCREEN 
					s.verifyScreen.failedVerification = false
				})
			}, func(){
				widgets.NewLabel("Verify seed phrase").SetWeight(text.Bold).SetSize(6).Draw(ctx)
				inset := layout.Inset{
					Top: unit.Dp(23),
				}
				inset.Layout(ctx, func() {
					widgets.NewLabel("Step 2/2").SetSize(4).Draw(ctx)
				})
			})
		},
		func(){
			drawBody(ctx, 
				widgets.NewLabel("Select the correct words to verify.").SetSize(5),
				func() {
					topInset := float32(0) 
					if s.verifyScreen.failedVerification {
						inset := layout.Inset{
							Top: unit.Dp(topInset),
						}
						inset.Layout(ctx, func(){
							helper.PaintArea(ctx, helper.DangerColor, ctx.Constraints.Width.Max, 30)
							ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
							widgets.NewLabel("Failed to verify. Please go through every word and try again").
								SetSize(5).
								SetColor(helper.WhiteColor).
								SetAlignment(widgets.AlignMiddle).
								Draw(ctx)
						})
						topInset += 33
					}

					inset := layout.Inset{
						Top: unit.Dp(topInset),
					}
					inset.Layout(ctx, func(){
						list := layout.List{Axis: layout.Vertical}
						list.Layout(ctx, len(s.verifyScreen.options), func(i int){
							ctx.Constraints.Height.Min = 53
							helper.PaintArea(ctx, helper.WhiteColor, ctx.Constraints.Width.Max, 75)
						
							inset := layout.Inset{
								Top: unit.Dp(5),
								Left: unit.Dp(7),
							}
							inset.Layout(ctx, func(){
								txt := strconv.Itoa(i+1) + ")"
								lbl := widgets.NewLabel(txt).SetSize(4).SetWeight(text.Bold).SetColor(helper.GrayColor)
								selectedText := s.verifyScreen.options[i].Selected()
								if selectedText == "" {
									txt += " ..."
								} else {
									txt += " " + selectedText
								}
								lbl.SetText(txt).Draw(ctx)
							})

							inset = layout.Inset{
								Top:  unit.Dp(27),
								Left: unit.Dp(10),
							}
							inset.Layout(ctx, func() {
								ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
								ctx.Constraints.Height.Max = 40
								s.verifyScreen.options[i].Draw(ctx)
							})
						})
					})
				})
		},
		func(){
			ctx.Constraints.Height.Min = 50

			bgcol := helper.DecredLightBlueColor
			if !s.hasSelectedAllVerifyWords() {
				bgcol = helper.GrayColor
			}

			drawFooter(ctx, 10, func(){
				s.verifyScreen.verifyButton.
					SetBackgroundColor(bgcol).
					Draw(ctx, func() {
						if s.hasSelectedAllVerifyWords() {
							s.verifyScreen.failedVerification = false
							if s.doVerify() {
								s.currentScreen = SUCCESS_SCREEN
							} else {
								s.verifyScreen.failedVerification = true
							}
						}
					})
				})
		},
	}
}

func (s *SeedPage) hasSelectedAllVerifyWords() bool {
	for i := range s.verifyScreen.options {
		if s.verifyScreen.options[i].Selected() == "" {
			return false
		}
	}
	return true
}

func (s *SeedPage) doVerify() bool {
	for i := range s.verifyScreen.options {
		if s.verifyScreen.options[i].Selected() != s.seedScreen.words[i] {
			return false
		}
	}

	return true
}

func (s *SeedPage) getSuccessScreenWidgets(ctx *layout.Context, changePageFunc func(string)) []func() {
	return []func(){
		func(){
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
			layout.Align(layout.Center).Layout(ctx, func() {
				inset := layout.Inset{
					Top: unit.Dp(80),
				}
				inset.Layout(ctx, func() {
					ctx.Constraints.Width.Min = 50
					widgets.NewCheckbox().SetSize(80).MakeAsIcon().Draw(ctx)
				})
			})
		},
		func(){
			inset := layout.Inset{
				Top: unit.Dp(20),
			}
			inset.Layout(ctx, func() {
				widgets.NewLabel("Your seed phrase backup is").
					SetSize(6).
					SetWeight(text.Bold).
					SetAlignment(widgets.AlignMiddle).
					SetColor(helper.BlackColor).
					Draw(ctx)
			})
		},
		func(){
			inset := layout.Inset{
				Top: unit.Dp(5),
			}
			inset.Layout(ctx, func() {
				widgets.NewLabel("verified").
					SetSize(6).
					SetWeight(text.Bold).
					SetAlignment(widgets.AlignMiddle).
					SetColor(helper.BlackColor).
					Draw(ctx)
			})
		},
		func(){
			inset := layout.Inset{
				Top: unit.Dp(30),
			}
			inset.Layout(ctx, func() {
				widgets.NewLabel("Be sure to store your seed phrase backup in a").
					SetSize(5).
					SetAlignment(widgets.AlignMiddle).
					SetColor(helper.BlackColor).
					Draw(ctx)
			})
		},
		func(){
			inset := layout.Inset{
				Top: unit.Dp(5),
			}
			inset.Layout(ctx, func() {
				widgets.NewLabel("secure location.").
					SetSize(5).
					SetAlignment(widgets.AlignMiddle).
					SetColor(helper.BlackColor).
					Draw(ctx)
			})
		},
		func() {
			inset := layout.Inset{
				Top: unit.Dp(200),
			}
			inset.Layout(ctx, func() {
				s.backToWalletsButton.SetWidth(ctx.Constraints.Width.Max).Draw(ctx, func() {
					changePageFunc("overview")
				})
			})
		},
	}
}
