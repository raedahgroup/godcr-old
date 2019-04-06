package terminalprompt

import "testing"

func Test_getTextInput(t *testing.T) {
	tests := []struct {
		name    string
		prompt  string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := getTextInput(test.prompt)
			if (err != nil) != test.wantErr {
				t.Errorf("getTextInput() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("getTextInput() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_getPasswordInput(t *testing.T) {
	tests := []struct {
		name    string
		prompt  string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := getPasswordInput(test.prompt)
			if (err != nil) != test.wantErr {
				t.Errorf("getPasswordInput() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("getPasswordInput() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_setTerminalEcho(t *testing.T) {
	tests := []struct {
		name    string
		on      bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := setTerminalEcho(test.on); (err != nil) != test.wantErr {
				t.Errorf("setTerminalEcho() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
