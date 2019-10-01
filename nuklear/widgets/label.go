package widgets

import (
	"image/color"
	"strings"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	f "golang.org/x/image/font"
)

const (
	LeftCenterAlign = "LC"
	CenterAlign     = "CC"
)

type fontFace f.Face

// AddLabel adds a single line label to the window. The label added does not wrap.
func (window *Window) AddLabel(text string, align label.Align) {
	window.Row(window.SingleLineLabelHeight()).Dynamic(1)
	window.Label(text, align)
}

// AddLabel adds a single line label to the window. The label added does not wrap.
func (window *Window) AddLabelFixedWidth(text string, align label.Align, width int) {
	window.Row(window.SingleLineLabelHeight()).Static(width)
	window.Label(text, align)
}

// AddLabelWithFont adds a single line label to the window. The label added does not wrap.
func (window *Window) AddLabelWithFont(text string, align label.Align, font fontFace) {
	window.UseFontAndResetToPrevious(font, func() {
		window.Row(window.SingleLineLabelHeight()).Dynamic(1)
		window.Label(text, align)
	})
}

func (window *Window) AddColoredLabel(text string, color color.RGBA, align label.Align) {
	window.Row(window.SingleLineLabelHeight()).Dynamic(1)
	window.LabelColored(text, align, color)
}

// AddWrappedLabel adds a label to the window.
// The label added wraps it's text and assumes the height required to display all it's text.
func (window *Window) AddWrappedLabel(text string, align label.Align) {
	singleLineHeight := window.SingleLineLabelHeight()
	lines := window.WrapLabelText(text, window.Font())

	for _, line := range lines {
		window.Row(singleLineHeight).Dynamic(1)
		window.Label(line, align)
	}
}

// AddWrappedLabel adds a label to the window.
// The label added wraps it's text and assumes the height required to display all it's text.
func (window *Window) AddWrappedLabelWithColor(text string, align label.Align, color color.RGBA) {
	singleLineHeight := window.SingleLineLabelHeight()
	lines := window.WrapLabelText(text, window.Font())

	for _, line := range lines {
		window.Row(singleLineHeight).Dynamic(1)
		window.LabelColored(line, align, color)
	}
}

func (window *Window) AddWrappedLabelWithFont(text string, align label.Align, font fontFace) {
	singleLineHeight := window.SingleLineLabelHeight()
	lines := window.WrapLabelText(text, window.Font())

	window.UseFontAndResetToPrevious(font, func() {
		for _, line := range lines {
			window.Row(singleLineHeight).Dynamic(1)
			window.Label(line, align)
		}
	})
}

func (window *Window) WrapLabelText(text string, font fontFace) (wrappedLines []string) {
	textWidth := nucular.FontWidth(font, text)
	maxWidth := window.Bounds.W - window.Master().Style().GroupWindow.Padding.X

	words := strings.Split(text, " ")
	wordsCountPerLine := maxWidth * len(words) / textWidth
	if textWidth < maxWidth {
		wordsCountPerLine = len(words)
	}

	wordCountForCurrentLine := wordsCountPerLine
	for {
		currentLine := strings.Join(words[:wordCountForCurrentLine], " ")
		if nucular.FontWidth(font, currentLine) > maxWidth {
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

func (window *Window) AddLabels(labels ...*LabelTableCell) {
	widths := make([]int, len(labels))
	for i, labelCell := range labels {
		widths[i] = window.LabelWidth(labelCell.text)
	}
	window.AddLabelsWithWidths(widths, labels...)
}

func (window *Window) AddLabelsWithWidths(widths []int, labels ...*LabelTableCell) {
	window.Row(window.SingleLineLabelHeight()).Static(widths...)
	window.AddLabelsToCurrentRow(labels...)
}

func (window *Window) AddLabelsToCurrentRow(labels ...*LabelTableCell) {
	for _, labelCell := range labels {
		if labelCell == nil {
			// need to fill this column with empty space so the next cell is added to the next column instead of this column
			window.Spacing(1)
		} else {
			labelCell.Render(window)
		}
	}
}

func (window *Window) AddLinkLabelsToCurrentRow(labels ...*LinkLabelCell) {
	for _, labelCell := range labels {
		if labelCell == nil {
			// need to fill this column with empty space so the next cell is added to the next column instead of this column
			window.Spacing(1)
		} else {
			labelCell.Render(window)
		}
	}
}

func (window *Window) LabelWidth(text string) int {
	return nucular.FontWidth(window.Font(), text) + 8 // add 8 to text width to avoid text being cut off in label
}

func (window *Window) SingleLineLabelHeight() int {
	singleLineHeight := nucular.FontHeight(window.Font()) + 1
	if singleLineHeight < 20 {
		singleLineHeight = 20 // seems labels will not be drawn if row height is less than 20
	}
	return singleLineHeight
}

type LinkLabelCell struct {
	text      string
	selected  *bool
	clickFunc func()
}

func NewLinkLabelCellCell(text string, clickFunc func()) *LinkLabelCell {
	selected := false

	return &LinkLabelCell{
		text:      text,
		selected:  &selected,
		clickFunc: clickFunc,
	}
}

func (link *LinkLabelCell) Render(window *Window) {
	if window.SelectableLabel(link.text, "LC", link.selected) {
		link.clickFunc()
	}
}

func (link *LinkLabelCell) MinWidth(window *Window) int {
	return window.LabelWidth(link.text)
}
