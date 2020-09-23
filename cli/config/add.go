// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package config

import (
	"os"

	"github.com/arduino/arduino-cli/cli/errorcodes"
	"github.com/arduino/arduino-cli/cli/feedback"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initAddCommand() *cobra.Command {
	addCommand := &cobra.Command{
		Use:     "add",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.MinimumNArgs(2),
		Run:     runAddCommand,
	}
	return addCommand
}

func runAddCommand(cmd *cobra.Command, args []string) {
	// We're assuming the config file already exists, we should probably
	// create one in case it doesn't
	settingsKey := args[0]

	if settingsKey != "board_manager.additional_urls" {
		feedback.Errorf("Settings key %v is not a list.", settingsKey)
		os.Exit(errorcodes.ErrGeneric)
	}

	configFlag := cmd.InheritedFlags().Lookup("config-file")
	configFile := ""
	if configFlag != nil {
		configFile = configFlag.Value.String()
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
		err := viper.MergeInConfig()
		if err != nil {
			feedback.Errorf("Error reading config file: %s", configFile)
			os.Exit(errorcodes.ErrBadArgument)
		}
	}

	savedValues := viper.GetStringSlice(settingsKey)
	savedValues = append(savedValues, args[1:]...)
	viper.Set(settingsKey, savedValues)
	if err := viper.WriteConfig(); err != nil {
		feedback.Errorf("Error saving config file %s", viper.ConfigFileUsed())
		os.Exit(errorcodes.ErrGeneric)
	}

	feedback.Printf("New value: %v", viper.GetStringSlice(settingsKey))
}
