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

var configFile string
var globalConfig baoConfig.MonitorConfig
var logWriter *os.File
var baoLogger *slog.Logger

var RootCmd = &cobra.Command{
	Use:   "baomon",
	Short: "A monitor service for managing Openbao in StarlingX",
	Long:  `A monitor service for managing Openbao servers in StarlingX systems.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Open config from file
		configReader, err := os.Open(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Error in opening config file: %v, message: %v", configFile, err)
			os.Exit(1)
		}
		defer configReader.Close()
		err = globalConfig.ReadYAMLMonitorConfig(configReader)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"Error in parsing config file: %v, message: %v", configFile, err)
			os.Exit(1)
		}

		// Set default configuration for logs if no custum configs are given
		logFile := globalConfig.LogPath
		logLevel := globalConfig.LogLevel
		if logFile == "" {
			logFile = "/workdir/openbao-monitor.log"
		}
		if logLevel == "" {
			logLevel = "INFO"
		}

		// Setup Logs
		logWriter, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in opening the log file to write. %v\n", err)
			os.Exit(1)
		}

		var LogLevel slog.Level
		LogLevel.UnmarshalText([]byte(logLevel))
		baoLogger = slog.New(slog.NewTextHandler(logWriter, &slog.HandlerOptions{
			Level: LogLevel,
		}))
		slog.SetDefault(baoLogger)
		slog.Info(fmt.Sprintf("New call to the monitor service. Logs attached to file %v", logFile))
		slog.Info(fmt.Sprintf("Set log level: %v", logLevel))
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// For now write back the configs back to the config file
		// This should be changed to writing back to configs from file only
		configWriter, err := os.OpenFile(configFile, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error with opening config file to write in the changed configs\n")
			fmt.Fprintf(os.Stderr, "Message: %v\n", err)
		}
		err = globalConfig.WriteYAMLMonitorConfig(configWriter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error with writing the changed configs\n")
			fmt.Fprintf(os.Stderr, "Message: %v\n", err)
		}

		// Close the log file
		err = logWriter.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error with closing the log file\n")
			fmt.Fprintf(os.Stderr, "Message: %v\n", err)
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Declarations for global flags
	RootCmd.PersistentFlags().StringVar(&configFile, "config",
		"/workdir/testConfig.yaml", "file path to the monitor config file")
}
