package security

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/gio/helper"
	"github.com/raedahgroup/godcr/gio/widgets"
	"github.com/raedahgroup/godcr/gio/widgets/editor"
)

type (
	colors struct {
		passwordStrengthFillColor color.RGBA
		confirmPasswordBorderColor color.RGBA
		confirmPasswordFocusBorderColor color.RGBA
		cancelLabelColor color.RGBA
		createButtonBackgroundColor color.RGBA
	}

	passwordTab struct {
		passwordInput             *editor.Input
		confirmPasswordInput      *editor.Input
		passwordStrengthIndicator *widgets.ProgressBar
		passwordStrength          float64

		errStr string
	}

	pinTab struct {
		pinInput        *editor.Input
		confirmPinInput *editor.Input
		pinStrength     int
	}

	PinAndPasswordWidget struct {
		tabContainer *widgets.TabContainer
		passwordTab  *passwordTab
		pinTab       *pinTab

		currentTab string
		IsCreating bool

		createButton *widgets.Button
		cancelLabel  *widgets.ClickableLabel

		colors colors

		cancelFunc func()
		createFunc func()
	}
)

func NewPinAndPasswordWidget(cancelFunc, createFunc func()) *PinAndPasswordWidget {
	return &PinAndPasswordWidget{
		currentTab:   "password",
		tabContainer: widgets.NewTabContainer().AddTab("Password").AddTab("PIN"),
		createButton: widgets.NewButton("Create", nil),
		cancelLabel:  widgets.NewClickableLabel("Cancel").SetSize(4).SetWeight(text.Bold).SetColor(helper.DecredLightBlueColor),
		cancelFunc:   cancelFunc,
		createFunc:   createFunc,
		passwordTab: &passwordTab{
			passwordInput:             editor.NewInput("Spending Password").SetMask("*"),
			confirmPasswordInput:      editor.NewInput("Confirm Spending Password").SetMask("*"),
			passwordStrength:          0,
			passwordStrengthIndicator: widgets.NewProgressBar().SetHeight(6),
		},
		pinTab: &pinTab{
			pinInput:        editor.NewInput("Pin").SetMask("*").Numeric(),
			confirmPinInput: editor.NewInput("Confirm Pin").SetMask("*").Numeric(),
		},
		colors: colors{},
	}
}

func (p *PinAndPasswordWidget) Reset() {
	p.passwordTab.passwordInput.Clear()
	p.passwordTab.confirmPasswordInput.Clear()
}

func (p *PinAndPasswordWidget) Value() string {
	if p.currentTab == "password" {
		return p.passwordTab.passwordInput.Text()
	}

	return p.pinTab.pinInput.Text()
}

func (p *PinAndPasswordWidget) Render(ctx *layout.Context) { 
	// perform these actions in separate goroutines as changes happen
	go func(){
		// update password strength 
		p.passwordTab.passwordStrength = (dcrlibwallet.ShannonEntropy(p.passwordTab.passwordInput.Text()) / 4) * 100

		p.validate()
		
		// watch for changes and update colors
		p.updateColors()
	}()

	layout.Stack{Alignment: layout.NW}.Layout(ctx,
		layout.Expanded(func() {
			widgets.NewLabel("Create a Spending Password", 5).
				SetWeight(text.Bold).
				Draw(ctx)
		}),
		layout.Stacked(func() {
			inset := layout.Inset{
				Top: unit.Dp(25),
			}
			inset.Layout(ctx, func() {
				p.tabContainer.Draw(ctx, p.passwordRenderFunc, p.pinRenderFunc)
			})
		}),
	)
}


func (p *PinAndPasswordWidget) updateColors() {
	// password strength fill color
	if p.passwordTab.passwordStrength > 70 {
		p.colors.passwordStrengthFillColor = helper.DecredGreenColor
	} else {
		p.colors.passwordStrengthFillColor = helper.DecredOrangeColor
	}

	// confirm password border colors
	if !p.bothPasswordsMatch() && p.passwordTab.confirmPasswordInput.Len() > 0 {
		p.colors.confirmPasswordBorderColor = helper.DangerColor
		p.colors.confirmPasswordFocusBorderColor = helper.DangerColor
	} else {
		p.colors.confirmPasswordBorderColor = helper.GrayColor
		p.colors.confirmPasswordFocusBorderColor = helper.DecredLightBlueColor
	}

	// cancel label color 
	if p.IsCreating {
		p.colors.cancelLabelColor = helper.GrayColor
	} else {
		p.colors.cancelLabelColor = helper.DecredLightBlueColor
	}

	// create button 
	if p.IsCreating {
		p.colors.createButtonBackgroundColor = helper.GrayColor
	} else {
		if p.bothPasswordsMatch() && p.passwordTab.confirmPasswordInput.Len() > 0 {
			p.colors.createButtonBackgroundColor = helper.DecredLightBlueColor
		} else {
			p.colors.createButtonBackgroundColor = helper.GrayColor
		}
	}
}

func (p *PinAndPasswordWidget) bothPasswordsMatch() bool {
	if p.passwordTab.confirmPasswordInput.Text() == p.passwordTab.passwordInput.Text() { 
		return true
	}

	return false
}

func (p *PinAndPasswordWidget) passwordRenderFunc(ctx *layout.Context) {
	p.currentTab = "password"

	widgets := []func(){
		// password row
		func(){
			p.passwordTab.passwordInput.Draw(ctx)
		},
		
		// password strength row
		func(){
			helper.Inset(ctx, 10, 0, 0, 0, func(){
				p.passwordTab.passwordStrengthIndicator.
					SetProgressColor(p.colors.passwordStrengthFillColor).
					Draw(ctx, &p.passwordTab.passwordStrength)
			})
		},

		// confirm password row
		func() {
			helper.Inset(ctx, 10, 0, 0, 0, func(){
				p.passwordTab.confirmPasswordInput.
					SetBorderColor(p.colors.confirmPasswordBorderColor).
					SetFocusedBorderColor(p.colors.confirmPasswordFocusBorderColor).
					Draw(ctx)
			})
		},

		// buttons and error text row
		func() {
			helper.Inset(ctx, 10, 0, 0, 0, func(){
				layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
					layout.Flexed(2, func(){
						widgets.NewLabel(p.passwordTab.errStr).SetColor(helper.DangerColor).Draw(ctx)
					}), 
					layout.Rigid(func(){
						helper.Inset(ctx, 10, 0, 0, 10, func(){
							p.cancelLabel.SetColor(p.colors.cancelLabelColor).Draw(ctx, func() {
								if !p.IsCreating {
									p.cancelFunc()
								}
							})
						})
					}),
					layout.Rigid(func(){
						var txt string 
						if p.IsCreating {
							txt = "Creating..."
						} else {
							txt = "Create"
						}

						p.createButton.SetText(txt).
							SetBackgroundColor(p.colors.createButtonBackgroundColor).
							Draw(ctx, func() {
							p.validateAndCreate()
						})
					}),
				)
			})
		},
	}

	list := &layout.List{
		Axis: layout.Vertical,
	}

	list.Layout(ctx, len(widgets), func(i int){
		layout.UniformInset(unit.Dp(0)).Layout(ctx, widgets[i])
	})
}

func (p *PinAndPasswordWidget) validate() bool {
	if p.passwordTab.passwordInput.Text() == "" { 
		return false
	}
	
	
	if !p.bothPasswordsMatch() {
		p.passwordTab.errStr = "Both passwords do not match"
		return false
	}

	p.passwordTab.errStr = ""

	return true
}

func (p *PinAndPasswordWidget) validateAndCreate() {
	if p.IsCreating {
		return
	}
	
	if p.passwordTab.passwordInput.Text() == "" {
		p.passwordTab.errStr = "Please enter your desired password"
		return
	}

	if !p.validate() {
		return
	}

	p.passwordTab.errStr = ""
	p.createFunc()
}

func (p *PinAndPasswordWidget) pinRenderFunc(ctx *layout.Context) {
	p.currentTab = "pin"
}
