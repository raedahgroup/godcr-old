package widgets

import (
	"image"
	"image/color"
	"strings"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
	"github.com/raedahgroup/godcr/nuklear/styles"
)

const (
	alertWidgetWidth = 300
	lineHeight       = 20
	displayDuration  = 4      // 4 seconds
	alertWindowFlag  = 524288 // this displays the popup as a tooltip so the UI does not block
)

type Alert struct {
	isErrorAlert     bool
	text             string
	lines            []string
	closeAlertWindow bool
}

func NewAlertWidget(text string, isErrorAlert bool, window *Window) {
	a := &Alert{
		text:             text,
		isErrorAlert:     isErrorAlert,
		closeAlertWindow: false,
	}
	a.lines = a.splitToLines(window)
	windowBounds := window.Window.Bounds

	bounds := rect.Rect{
		W: alertWidgetWidth,
		H: lineHeight*len(a.lines) + 80, // 80 is to allow for vertical spacing between window border and this widget
		X: windowBounds.W - 80,          // 80 is to allow for horizontal spacing between window border and this widget
		Y: 10,
	}

	flags := nucular.WindowClosable | nucular.WindowDynamic | nucular.WindowNonmodal | alertWindowFlag
	window.Window.Master().PopupOpen("", flags, bounds, true, a.popup)
	time.AfterFunc(time.Second*displayDuration, func() {
		a.closeAlertWindow = true
		window.Window.Master().Changed()
	})
}

func (a *Alert) splitToLines(window *Window) (wrappedLines []string) {
	font := window.Font()
	textWidth := nucular.FontWidth(font, a.text)

	words := strings.Split(a.text, " ")
	wordsCountPerLine := alertWidgetWidth * len(words) / textWidth
	if textWidth < alertWidgetWidth {
		wordsCountPerLine = len(words)
	}

	wordCountForCurrentLine := wordsCountPerLine
	for {
		currentLine := strings.Join(words[:wordCountForCurrentLine], " ")
		if nucular.FontWidth(font, currentLine) > alertWidgetWidth {
			wordCountForCurrentLine--
			continue // skip remainder of code and come back to calculating current line width using adjusted word count
		}

		wrappedLines = append(wrappedLines, currentLine)
		words = words[wordCountForCurrentLine:]

		if len(words) < wordsCountPerLine {
			wordCountForCurrentLine = len(words)
		} else {
			wordCountForCurrentLine = wordsCountPerLine
		}

		if len(words) == 0 {
			break
		}
	}

	return
}

func (a *Alert) popup(window *nucular.Window) {
	masterWindow := window.Master()

	// style popup
	var backgroundColor color.RGBA
	if a.isErrorAlert {
		backgroundColor = styles.DecredOrangeColor
	} else {
		backgroundColor = styles.DecredGreenColor
	}

	style := window.Master().Style()
	defaultTextColor := style.Text.Color
	defaultTooltipStyle := style.TooltipWindow

	style.Text.Color = styles.WhiteColor
	style.TooltipWindow.Background = backgroundColor
	style.TooltipWindow.Border = 0
	style.TooltipWindow.Padding = image.Point{7, 10}
	masterWindow.SetStyle(style)

	defer func() {
		// reset style
		style.Text.Color = defaultTextColor
		style.TooltipWindow.Background = defaultTooltipStyle.Background
		style.TooltipWindow.Padding = defaultTooltipStyle.Padding
		masterWindow.SetStyle(style)
	}()

	for index := range a.lines {
		window.Row(lineHeight).Dynamic(1)
		window.Label(a.lines[index], "LC")
	}

	if a.closeAlertWindow {
		window.Close()
		window.Master().Changed()
	}
}
