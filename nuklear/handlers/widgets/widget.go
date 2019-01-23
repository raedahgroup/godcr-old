package widgets

import (
	"fmt"

	"github.com/aarzilli/nucular"
)

type Widget interface {
	BeforeRender(window *nucular.Window)
	Render(finishHandler func())
	AfterRender()
}

var widgets = map[string]Widget{
	"passphrase": NewPassphraseWidget(),
	"loading":    NewLoadingWidget(),
}

func RegisterWidget(key string, widget Widget) error {
	if _, ok := widgets[key]; ok {
		return fmt.Errorf("Widget %s is already registered", key)
	}
	widgets[key] = widget
	return nil
}

func GetWidget(key string) (Widget, error) {
	if widget, ok := widgets[key]; ok {
		return widget, nil
	}

	return nil, fmt.Errorf("Widget :%s is not registered. Forgotten import?", key)
}

func Run(key string, window *nucular.Window, finishHandler func()) error {
	widget, err := GetWidget(key) // this will panic if widget is no registered
	if err != nil {
		return err
	}
	widget.BeforeRender(window)
	widget.Render(finishHandler)
	widget.AfterRender()

	return nil
}
