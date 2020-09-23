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
	"strings"

	"github.com/arduino/arduino-cli/cli/errorcodes"
	"github.com/arduino/arduino-cli/cli/feedback"
	"github.com/arduino/go-paths-helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func initRemoveCommand() *cobra.Command {
	removeCommand := &cobra.Command{
		Use:     "remove",
		Short:   "",
		Long:    "",
		Example: "",
		Args:    cobra.ExactArgs(1),
		Run:     runRemoveCommand,
	}
	return removeCommand
}

func runRemoveCommand(cmd *cobra.Command, args []string) {
	settings := viper.AllSettings()

	// This is ugly but we don't care since this is just a simple proof of concept
	// also we gotta do it like this because Viper doesn't a way to unset keys
	k := strings.Split(args[0], ".")
	switch len(k) {
	case 1:
		delete(settings, k[0])
	case 2:
		delete(settings[k[0]].(map[string]interface{}), k[1])
		if len(settings[k[0]].(map[string]interface{})) == 0 {
			delete(settings, k[0])
		}
	case 3:
		delete(settings[k[0]].(map[string]interface{})[k[1]].(map[string]interface{}), k[2])
		if len(settings[k[0]].(map[string]interface{})[k[1]].(map[string]interface{})) == 0 {
			delete(settings, k[1])
		}
		if len(settings[k[0]].(map[string]interface{})) == 0 {
			delete(settings, k[0])
		}
	}

	bs, err := yaml.Marshal(settings)
	if err != nil {
		feedback.Errorf("unable to marshal configs to YAML: %v", err)
		os.Exit(errorcodes.ErrGeneric)
	}

	configFile := paths.New(viper.ConfigFileUsed())
	if err = configFile.WriteFile(bs); err != nil {
		feedback.Errorf("unable to update configs: %v", err)
		os.Exit(errorcodes.ErrGeneric)
	}
}
