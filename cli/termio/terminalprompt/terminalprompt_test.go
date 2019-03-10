package terminalprompt

import "testing"

func Test_skipEOFError(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		err     error
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := skipEOFError(test.value, test.err)
			if (err != nil) != test.wantErr {
				t.Errorf("skipEOFError() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("skipEOFError() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRequestInput(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		validate ValidatorFunction
		want     string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := RequestInput(test.message, test.validate)
			if (err != nil) != test.wantErr {
				t.Errorf("RequestInput() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("RequestInput() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRequestNumberInput(t *testing.T) {
	tests := []struct {
		name         string
		message      string
		defaultValue []int
		wantNumber   int
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotNumber, err := RequestNumberInput(test.message, test.defaultValue...)
			if (err != nil) != test.wantErr {
				t.Errorf("RequestNumberInput() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if gotNumber != test.wantNumber {
				t.Errorf("RequestNumberInput() = %v, want %v", gotNumber, test.wantNumber)
			}
		})
	}
}

func TestRequestInputSecure(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		validate ValidatorFunction
		want     string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := RequestInputSecure(test.message, test.validate)
			if (err != nil) != test.wantErr {
				t.Errorf("RequestInputSecure() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("RequestInputSecure() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRequestSelection(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		options  []string
		validate ValidatorFunction
		want     string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := RequestSelection(test.message, test.options, test.validate)
			if (err != nil) != test.wantErr {
				t.Errorf("RequestSelection() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("RequestSelection() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestRequestYesNoConfirmation(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		defaultOption string
		want          bool
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := RequestYesNoConfirmation(test.message, test.defaultOption)
			if (err != nil) != test.wantErr {
				t.Errorf("RequestYesNoConfirmation() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("RequestYesNoConfirmation() = %v, want %v", got, test.want)
			}
		})
	}
}
