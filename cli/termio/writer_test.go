package termio

import (
	"bytes"
	"reflect"
	"testing"
	"text/tabwriter"
)

func TestTabWriter(t *testing.T) {
	tests := []struct {
		name  string
		want  *tabwriter.Writer
		wantW string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if got := TabWriter(w); !reflect.DeepEqual(got, test.want) {
				t.Errorf("TabWriter() = %v, want %v", got, test.want)
			}
			if gotW := w.String(); gotW != test.wantW {
				t.Errorf("TabWriter() = %v, want %v", gotW, test.wantW)
			}
		})
	}
}

func TestPrintTabularResult(t *testing.T) {
	tests := []struct {
		name           string
		w              *tabwriter.Writer
		columnsHeaders []string
		rows           [][]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			PrintTabularResult(test.w, test.columnsHeaders, test.rows)
		})
	}
}

func TestPrintStringResult(t *testing.T) {
	tests := []struct {
		name   string
		output []string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			PrintStringResult(test.output...)
		})
	}
}
