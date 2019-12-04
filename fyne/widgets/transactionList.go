package widgets

// import (
// 	"image/color"

// 	"fyne.io/fyne"
// 	"fyne.io/fyne/canvas"
// 	"fyne.io/fyne/theme"
// 	"fyne.io/fyne/widget"
// )

// type TransactionLister struct {
// 	Amount   string
// 	Decimals string
// 	Status   string
// 	Icon     fyne.Resource
// }

// type TransactionList struct {
// 	amount   string
// 	decimals string
// 	status   string
// 	icon     fyne.Resource
// 	size     fyne.Size
// 	position fyne.Position

// 	OnTapped func()
// }

// type transactionListRenderer struct {
// 	icon     *canvas.Image
// 	label    *canvas.Text
// 	subLabel *canvas.Text
// 	status   *canvas.Text

// 	objects         []fyne.CanvasObject
// 	transactionList *TransactionList
// }

// func (t *TransactionList) CreateRenderer() fyne.WidgetRenderer {
// 	var objects []fyne.CanvasObject
// 	var img *canvas.Image
// 	if t.icon != nil {
// 		img = canvas.NewImageFromResource(t.icon)
// 		img = canvas.NewImageFromResource(t.icon)
// 		objects = append(objects, img)
// 	}
// 	text := canvas.NewText(t.amount, color.Black)
// 	text.Alignment = fyne.TextAlignLeading
// 	subText := canvas.NewText(t.decimals, color.Black)
// 	subText.Alignment = fyne.TextAlignLeading
// 	subText.TextSize = 10
// 	status := canvas.NewText(t.status, color.Gray{})
// 	status.TextSize = 12
// 	status.Alignment = fyne.TextAlignTrailing
// 	objects = append(objects, text)
// 	objects = append(objects, subText)
// 	objects = append(objects, status)
// 	return &transactionListRenderer{icon: img, label: text, subLabel: subText, status: status, objects: objects, transactionList: t}
// }

// func (t *transactionListRenderer) Layout(size fyne.Size) {
// 	inner := size.Subtract(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
// 	inner = inner.Subtract(fyne.NewSize(24, 24))
// 	statusSize := size.Subtract(fyne.NewSize(theme.Padding()*10, theme.Padding()*8))
// 	height := t.transactionList.size.Height

// 	t.label.Resize(inner)
// 	t.label.Move(fyne.NewPos(theme.Padding()*2+30, theme.Padding()))
// 	inner = inner.Subtract(fyne.NewSize(12, 12))
// 	t.subLabel.Resize(inner)
// 	t.subLabel.Move(fyne.NewPos(t.label.MinSize().Width+39, t.label.Position().Y+5))
// 	t.status.Resize(statusSize)
// 	t.status.Move(fyne.NewPos(0, height/2))

// 	if t.icon != nil {
// 		t.icon.Resize(fyne.NewSize(height, height))
// 		t.icon.Move(fyne.NewPos(0, height/2))
// 	}
// }

// func (t *transactionListRenderer) MinSize() (size fyne.Size) {
// 	baseSize := t.label.MinSize()
// 	newSize := baseSize.Add(fyne.NewSize(250, 0))
// 	return newSize.Add(fyne.NewSize(theme.Padding()*4, theme.Padding()*2))
// }

// func (t *transactionListRenderer) Refresh() {
// 	// todo
// }

// func (t *transactionListRenderer) ApplyTheme() {
// 	// todo
// }

// func (t *transactionListRenderer) BackgroundColor() color.Color {
// 	return color.Transparent
// }

// func (t *transactionListRenderer) Objects() (objects []fyne.CanvasObject) {
// 	return t.objects
// }

// func (t *transactionListRenderer) Destroy() {
// 	// todo
// }

// func (t *TransactionList) Size() fyne.Size {
// 	return t.size
// }

// func (t *TransactionList) Resize(size fyne.Size) {
// 	t.size = size
// 	widget.Renderer(t).Layout(size)
// }

// func (t *TransactionList) Position() fyne.Position {
// 	return t.position
// }

// func (t *TransactionList) Move(position fyne.Position) {
// 	t.position = position
// 	widget.Renderer(t).Refresh()
// }

// func (t *TransactionList) MinSize() fyne.Size {
// 	if widget.Renderer(t) == nil {
// 		return fyne.NewSize(0, 0)
// 	}
// 	return widget.Renderer(t).MinSize()
// }

// func (t *TransactionList) Show() {
// }

// func (t *TransactionList) Hide() {
// }

// func (t *TransactionList) Visible() bool {
// 	return true
// }

// func NewTransactionList(lister TransactionLister, onTap func()) *TransactionList {
// 	transactionList := &TransactionList{
// 		amount:   lister.Amount,
// 		decimals: lister.Decimals,
// 		status:   lister.Status,
// 		icon:     lister.Icon,
// 		OnTapped: onTap,
// 	}
// 	widget.Renderer(transactionList).Layout(transactionList.MinSize())
// 	return transactionList
// }
