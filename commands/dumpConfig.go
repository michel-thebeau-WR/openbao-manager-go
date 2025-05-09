//
// Copyright (c) 2025 Wind River Systems, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//

package baoCommands

import (
	"fmt"
	"log/slog"
	"os"

	baoConfig "github.com/michel-thebeau-WR/openbao-manager-go/baomon/config"
	"github.com/spf13/cobra"
)

var dumpConfigReadCmd = &cobra.Command{
	Use:   "read readFile",
	Short: "Read config from a YAML file",
	Long:  "Read baomon configuration from a specified YAML file, and prints to stdout",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug(fmt.Sprintf("Action: dumpConfig read. File to read: %v", args[0]))
		var S baoConfig.MonitorConfig
		openedFile, err := os.Open(args[0])
		if err != nil {
			slog.Error(fmt.Sprintf("Unable to open file %v. Error message: %v\n", args[0], err))
			fmt.Fprint(os.Stderr, "Error with reading YAML file. Check the logs for details.\n")
			return
		}
		defer openedFile.Close()

		err = S.ReadYAMLMonitorConfig(openedFile)
		if err != nil {
			slog.Error(fmt.Sprintf("Error with parsing read yaml data: %v\n", err))
			fmt.Fprint(os.Stderr, "Error with reading YAML file. Check the logs for details.\n")
			return
		}
		fmt.Printf("Result: \n%#v\n", S)
	},
}

var dumpConfigWriteCmd = &cobra.Command{
	Use:   "write readFile readWrite",
	Short: "Write config to another YAML file",
	Long:  "Copy a config file from the first file to the second file, using baoConfig's write method",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug(fmt.Sprintf("Action: dumpConfig write. File to read: %v. File to write: %v", args[0], args[1]))
		var S baoConfig.MonitorConfig
		openedFile, err := os.Open(args[0])
		if err != nil {
			slog.Error(fmt.Sprintf("Unable to open file %v. Error message: %v\n", args[0], err))
			fmt.Fprint(os.Stderr, "Error with writing YAML file. Check the logs for details.\n")
			return
		}
		defer openedFile.Close()

		err = S.ReadYAMLMonitorConfig(openedFile)
		if err != nil {
			slog.Error(fmt.Sprintf("Error with parsing read yaml data: %v\n", err))
			fmt.Fprint(os.Stderr, "Error with writing YAML file. Check the logs for details.\n")
			return
		}

		writeFile, err := os.Create(args[1])
		if err != nil {
			slog.Error(fmt.Sprintf("Unable to open %v for writing\n", args[1]))
			fmt.Fprint(os.Stderr, "Error with writing YAML file. Check the logs for details.\n")
			return
		}
		defer writeFile.Close()

		err = S.WriteYAMLMonitorConfig(writeFile)
		if err != nil {
			slog.Error(fmt.Sprintf("Error with writing parsed yaml data to file: %v\n", err))
			fmt.Fprint(os.Stderr, "Error with writing YAML file. Check the logs for details.\n")
			return
		}
		fmt.Print("Write Complete\n")
	},
}

var dumpConfigCmd = &cobra.Command{
	Use:   "dumpConfig",
	Short: "Dev command for read/write YAML config files",
	Long:  `A dev command for interacting with the baoConfig package using a YAML file.`,
}

func init() {
	dumpConfigCmd.AddCommand(dumpConfigReadCmd)
	dumpConfigCmd.AddCommand(dumpConfigWriteCmd)
	RootCmd.AddCommand(dumpConfigCmd)
}
