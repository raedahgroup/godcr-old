package widgets

const progressBarHeight = 25

func (window *Window) AddProgressBar(progress *int, max int)  {
	window.Row(progressBarHeight).Dynamic(1)
	window.Progress(progress, max, false)
}
