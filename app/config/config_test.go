package config

import (
	"reflect"
	"testing"
)

func Test_defaultConfig(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{
			name: "default config",
			want: Config{
				ConfFileOptions: ConfFileOptions{
					AppDataDir:    DefaultAppDataDir,
					WalletRPCCert: defaultRPCCertFile,
					HTTPHost:      defaultHTTPHost,
					HTTPPort:      defaultHTTPPort,
					DebugLevel:    defaultLogLevel,
				},
				CommandLineOptions: CommandLineOptions{
					InterfaceMode: "cli",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := defaultConfig(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("defaultConfig() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_defaultCommandLineOptions(t *testing.T) {
	tests := []struct {
		name string
		want CommandLineOptions
	}{
		{
			name: "default command line options",
			want: CommandLineOptions{
				InterfaceMode: "cli",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := defaultCommandLineOptions(); !reflect.DeepEqual(got, test.want) {
				t.Errorf("defaultCommandLineOptions() = %v, want %v", got, test.want)
			}
		})
	}
}

/**
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *Config
		want1   []string
		wantErr bool
	}{
		{
			name: "load config",
			args: []string{"test"},
			want: &Config{
				ConfFileOptions: ConfFileOptions{
					AppDataDir:    DefaultAppDataDir,
					WalletRPCCert: defaultRPCCertFile,
					HTTPHost:      defaultHTTPHost,
					HTTPPort:      defaultHTTPPort,
					DebugLevel:    defaultLogLevel,
				},
				CommandLineOptions: CommandLineOptions{
					InterfaceMode: "cli",
				},
			},
			want1:   []string{},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.args != nil {
				os.Args = test.args
			}

			got, got1, err := LoadConfig()
			if (err != nil) != test.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("LoadConfig() got = %v, want %v", got, test.want)
			}
			if !reflect.DeepEqual(got1, test.want1) {
				t.Errorf("LoadConfig() got1 = %v, want %v", got1, test.want1)
			}
		})
	}
}**/

func Test_hasConfigFileOption(t *testing.T) {
	tests := []struct {
		name            string
		commandLineArgs []string
		want            bool
	}{
		{
			name:            "valid config file option",
			commandLineArgs: []string{"--httpport=:9001"},
			want:            true,
		},
		{
			name:            "invalid config file option",
			commandLineArgs: []string{"--mode=cli"},
			want:            false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := hasConfigFileOption(test.commandLineArgs); got != test.want {
				t.Errorf("hasConfigFileOption() = %v, want %v", got, test.want)
			}
		})
	}
}
