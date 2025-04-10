//
// Copyright (c) 2025 Wind River Systems, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//

package baoCommands

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	openbao "github.com/openbao/openbao/api/v2"
	"github.com/spf13/cobra"
)

var optFileStr string
var secretShares int
var secretThreshold int

func initializeOpenbao(dnshost string, opts *openbao.InitRequest) error {
	slog.Debug(fmt.Sprintf("Attempting the initialize openbao on host %v", dnshost))
	newConfig, err := globalConfig.NewOpenbaoConfig(dnshost)
	if err != nil {
		return fmt.Errorf("error in creating new config for openbao: %v", err)
	}

	slog.Debug("Creating Openbao client for API access")
	newClient, err := openbao.NewClient(newConfig)
	if err != nil {
		return fmt.Errorf("error in creating new client for openbao: %v", err)
	}

	slog.Debug("Running /sys/init")
	responce, err := newClient.Sys().Init(opts)
	if err != nil {
		return fmt.Errorf("error during call to openbao init: %v", err)
	}

	slog.Debug("/sys/init complete")
	err = globalConfig.ParseInitResponse(dnshost, responce)
	if err != nil {
		return fmt.Errorf("error during parsing init responce: %v", err)
	}

	return nil
}

var initCmd = &cobra.Command{
	Use:   "init DNSHost",
	Short: "Initialize openbao",
	Long: `Initialize openbao server using the monitor configurations.
The key shards returned from the initResponse will be stored in the monitor
configurations.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		slog.Debug("Action: init")
		fileGiven := cmd.Flags().Lookup("file").Changed
		secretSharesFlag := cmd.Flags().Lookup("secret-shares").Changed
		secretThresholdFlag := cmd.Flags().Lookup("secret-threshold").Changed

		secretNumGiven := secretSharesFlag && secretThresholdFlag
		if fileGiven == secretNumGiven {
			slog.Error("Command \"init\" failed due to insufficent init options")
			fmt.Fprintf(os.Stderr, "The options for openbao init must be set by one of:\n")
			fmt.Fprintf(os.Stderr, "utilizing an option file using --file, or\n")
			fmt.Fprintf(os.Stderr, "--secret-shares and -- secret-threshold\n")
			return
		}

		var opts openbao.InitRequest
		if fileGiven {
			optFileReader, err := os.ReadFile(optFileStr)
			if err != nil {
				slog.Error("Command \"init\" failed due to invalid init option file")
				fmt.Fprintf(os.Stderr, "Unable to open init option file %v: %v\n", optFileStr, err)
				return
			}
			err = json.Unmarshal(optFileReader, &opts)
			if err != nil {
				slog.Error("Command \"init\" failed due to error in parsing init option file")
				fmt.Fprintf(os.Stderr, "Unable to parse JSON file %v: %v\n", optFileStr, err)
				return
			}
		}
		if secretNumGiven {
			if secretShares == 0 {
				slog.Error("Command \"init\" failed due to insufficent number of secret-shares")
				fmt.Fprintf(os.Stderr, "The field secret-shares cannot be 0\n")
				return
			}
			if secretThreshold == 0 {
				slog.Error("Command \"init\" failed due to insufficent number of secret-threshold")
				fmt.Fprintf(os.Stderr, "The field secret-threshold cannot be 0\n")
				return
			}
			if secretShares < secretThreshold {
				slog.Error("Command \"init\" failed due to secret-threshold being greater than secret-shares")
				fmt.Fprintf(os.Stderr, "The field secret-threshold cannot be greater than secret-shares\n")
				return
			}
			opts.SecretShares = secretShares
			opts.SecretThreshold = secretThreshold
		}
		slog.Debug(fmt.Sprintf("Parsing init option successful. Attempting to run openbao init on host %v", args[0]))
		err := initializeOpenbao(args[0], &opts)
		if err != nil {
			slog.Error("Command \"init\" failed due to error in openbao init")
			fmt.Fprintf(os.Stderr, "Openbao init failed with error: %v\n", err)
			return
		}
		slog.Debug("Openbao init successful")
		fmt.Print("Openbao init complete\n")
	},
}

func init() {
	initCmd.Flags().StringVarP(&optFileStr, "file", "f", "", "A JSON file containing the options for openbao init")
	initCmd.Flags().IntVar(&secretShares, "secret-shares", 0, "The number of shares to split the root key into.")
	initCmd.Flags().IntVar(&secretThreshold, "secret-threshold", 0, "The number of shares required to reconstruct the root key.")
	RootCmd.AddCommand(initCmd)
}
