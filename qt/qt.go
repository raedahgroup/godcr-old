package qt

import (
	"github.com/raedahgroup/godcr/app"
	"os"

	"github.com/therecipe/qt/widgets"
)

const (
	minWindowWidth = 600
	minWindowHeight = 400
)

func LaunchApp() {
	// needs to be called once before you can start using the QWidgets
	qtApp := widgets.NewQApplication(len(os.Args), os.Args)

	// create a window and set the title
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(minWindowWidth, minWindowHeight)
	window.SetWindowTitle(app.Name)

	// create a regular widget
	// give it a QVBoxLayout
	// and make it the central widget of the window
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	// create a line edit
	// with a custom placeholder text
	// and add it to the central widgets layout
	input := widgets.NewQLineEdit(nil)
	input.SetPlaceholderText("Write something ...")
	widget.Layout().AddWidget(input)

	// create a button
	// connect the clicked signal
	// and add it to the central widgets layout
	button := widgets.NewQPushButton2("and click me!", nil)
	button.ConnectClicked(func(bool) {
		widgets.QMessageBox_Information(nil, "OK", input.Text(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})
	widget.Layout().AddWidget(button)

	// make the window visible
	window.Show()

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	qtApp.Exec()
}
