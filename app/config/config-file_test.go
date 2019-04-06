package config

import (
	"reflect"
	"testing"
)

func Test_defaultFileOptions(t *testing.T) {
	tests := []struct {
		name string
		want ConfFileOptions
	}{
		{
			name: "default file options",
			want: ConfFileOptions{
				AppDataDir:    DefaultAppDataDir,
				WalletRPCCert: defaultRPCCertFile,
				HTTPHost:      defaultHTTPHost,
				HTTPPort:      defaultHTTPPort,
				DebugLevel:    defaultLogLevel,
			},
		},
	}
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
	}{
		{
			name:           "create config file",
			wantSuccessful: true,
		},
	}
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
	}{
		{
			name: "read config file",
			wantConfig: ConfFileOptions{
				AppDataDir:    DefaultAppDataDir,
				WalletRPCCert: defaultRPCCertFile,
				HTTPHost:      defaultHTTPHost,
				HTTPPort:      defaultHTTPPort,
				DebugLevel:    defaultLogLevel,
			},
			wantErr: false,
		},
	}
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
	}{
		{
			name: "update config file",
			updateConfig: func(config *ConfFileOptions) {
				config.HTTPHost = defaultHTTPHost
			},
			wantErr: false,
		},
	}
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
	}{
		{
			name: "save config to file",
			config: ConfFileOptions{
				AppDataDir:    DefaultAppDataDir,
				WalletRPCCert: defaultRPCCertFile,
				HTTPHost:      defaultHTTPHost,
				HTTPPort:      defaultHTTPPort,
				DebugLevel:    defaultLogLevel,
			},
			wantErr: false,
		},
	}
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
		{
			name: "config file options",
			wantOptions: []string{
				"--appdata",
				"--walletrpcserver",
				"--walletrpccert",
				"--nowalletrpctls",
				"--httphost",
				"--httpport",
				"--debuglevel",
				"--wallets",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if gotOptions := configFileOptions(); !reflect.DeepEqual(gotOptions, test.wantOptions) {
				t.Errorf("configFileOptions() = %v, want %v", gotOptions, test.wantOptions)
			}
		})
	}
}
