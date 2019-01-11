package config

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	flags "github.com/jessevdk/go-flags"
)

const defaultConfigFilename = "godcr.conf"
var AppConfigFilePath = filepath.Join(defaultAppDataDir, defaultConfigFilename)

// createConfigFile create the configuration file in AppConfigFilePath using the default values
func createConfigFile() (successful bool) {
	configFile, err := os.Create(AppConfigFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error in creating config file: %s\n", err.Error())
			return
		}
		err = os.Mkdir(defaultAppDataDir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in creating config file directory: %s\n", err.Error())
			return
		}
		// we were unable to create the file because the dir was not found.
		// we shall attempt to recreate the file now that we have successfully created the dir
		configFile, err = os.Create(AppConfigFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in creating config file: %s\n", err.Error())
			return
		}
	}
	defer configFile.Close()

	tmpl := template.New("config")

	tmpl, err = tmpl.Parse(configTextTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error preparing default config file content: %s", err.Error())
		return
	}

	err = tmpl.Execute(configFile, defaultFileOptions())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error saving default configuration to file: %s\n", err.Error())
		return
	}

	fmt.Println("Config file created with default values at", AppConfigFilePath)
	return true
}

func parseConfigFile(parser *flags.Parser) error {
	if (parser.Options & flags.IgnoreUnknown) != flags.None {
		options := parser.Options
		parser.Options = flags.None
		defer func() { parser.Options = options }()
	}
	err := flags.NewIniParser(parser).ParseFile(AppConfigFilePath)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return fmt.Errorf("Error parsing configuration file: %v", err.Error())
		}
		return err
	}
	return nil
}

