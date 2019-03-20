package config

import (
	"reflect"
	"testing"
)

func Test_defaultFileOptions(t *testing.T) {
	tests := []struct {
		name string
		want ConfFileOptions
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := defaultFileOptions(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("defaultFileOptions() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_createConfigFile(t *testing.T) {
	tests := []struct {
		name           string
		wantSuccessful bool
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if gotSuccessful := createConfigFile(); gotSuccessful != test.wantSuccessful {
				t.Errorf("createConfigFile() = %v, want %v", gotSuccessful, test.wantSuccessful)
			}
		})
	}
}

func TestReadConfigFile(t *testing.T) {
	tests := []struct {
		name       string
		wantConfig ConfFileOptions
		wantErr    bool
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotConfig, err := ReadConfigFile()
			if (err != nil) != test.wantErr {
				t.Errorf("ReadConfigFile() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConfig, test.wantConfig) {
				t.Errorf("ReadConfigFile() = %v, want %v", gotConfig, test.wantConfig)
			}
		})
	}
}

func TestUpdateConfigFile(t *testing.T) {
	tests := []struct {
		name         string
		updateConfig func(config *ConfFileOptions)
		wantErr      bool
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := UpdateConfigFile(test.updateConfig); (err != nil) != test.wantErr {
				t.Errorf("UpdateConfigFile() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func Test_saveConfigToFile(t *testing.T) {
	tests := []struct {
		name    string
		config  ConfFileOptions
		wantErr bool
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := saveConfigToFile(test.config); (err != nil) != test.wantErr {
				t.Errorf("saveConfigToFile() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func Test_configFileOptions(t *testing.T) {
	tests := []struct {
		name        string
		wantOptions []string
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if gotOptions := configFileOptions(); !reflect.DeepEqual(gotOptions, test.wantOptions) {
				t.Errorf("configFileOptions() = %v, want %v", gotOptions, test.wantOptions)
			}
		})
	}
}
