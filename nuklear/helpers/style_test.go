package helpers

import (
	"reflect"
	"testing"

	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
	"golang.org/x/image/font"
)

func TestInitFonts(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := InitFonts(); (err != nil) != test.wantErr {
				t.Errorf("InitFonts() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func Test_getFont(t *testing.T) {
	tests := []struct {
		name     string
		fontSize int
		DPI      int
		fontData []byte
		want     font.Face
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := getFont(test.fontSize, test.DPI, test.fontData)
			if (err != nil) != test.wantErr {
				t.Errorf("getFont() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("getFont() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestSetFont(t *testing.T) {
	tests := []struct {
		name   string
		window *nucular.Window
		font   font.Face
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SetFont(test.window, test.font)
		})
	}
}

func TestGetStyle(t *testing.T) {
	tests := []struct {
		name string
		want *nstyle.Style
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := GetStyle(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetStyle() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestSetNavStyle(t *testing.T) {
	tests := []struct {
		name   string
		window nucular.MasterWindow
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SetNavStyle(test.window)
		})
	}
}

func TestSetPageStyle(t *testing.T) {
	tests := []struct {
		name   string
		window nucular.MasterWindow
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			SetPageStyle(test.window)
		})
	}
}
