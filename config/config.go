//
// Copyright (c) 2025 Wind River Systems, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//

package baoConfig

import (
	"fmt"
	"io"

	"github.com/go-yaml/yaml"
)

type URL struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Token struct {
	Duration int    `yaml:"duration"`
	Key      string `yaml:"key"`
}

type KeyShards struct {
	Key       string `yaml:"key"`
	KeyBase64 string `yaml:"key_base64"`
}

type MonitorConfig struct {
	// A map value listing all DNS names
	// Key: Domain name
	// Value: URL. consisting of host address and port number
	DNSnames map[string]URL `yaml:"IncludeInCluster"`

	// A map value listing all authentication tokens
	// Key: release id
	// Value: Token. consisting of lease duration and the token key
	Tokens map[string]Token `yaml:"Tokens"`

	// A map value listing all key shards for unseal
	// Key: shard name
	// Value: The shard key and the base64 encoded version of that key
	UnsealKeyShards map[string]KeyShards `yaml:"UnsealKeyShards"`

	// The default path of the log file
	LogPath string `yaml:"logPath"`

	// The default log level
	// Available log levels: DEBUG, INFO, WARN and ERROR
	LogLevel string `yaml:"logLevel"`
}

func (configInstance *MonitorConfig) ReadYAMLMonitorConfig(in io.Reader) error {
	data, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf(
			"unable to read Host DNS config data from input. Error message: %v", err)
	}

	err = yaml.Unmarshal(data, configInstance)
	if err != nil {
		return fmt.Errorf(
			"unable to unmarshal Host DNS config YAML data. Error message: %v", err)
	}

	// Validate YAML input for DNSnames
	err = configInstance.validateDNS()
	if err != nil {
		return err
	}

	// Validate YAML input for Tokens
	err = configInstance.validateTokens()
	if err != nil {
		return err
	}

	// Validate YAML input for unseal key shards
	err = configInstance.validateKeyShards()
	if err != nil {
		return err
	}

	// Validate YAML input for log configs
	err = configInstance.validateLogConfig()
	if err != nil {
		return err
	}

	return nil
}

func (configInstance MonitorConfig) WriteYAMLMonitorConfig(out io.Writer) error {
	data, err := yaml.Marshal(configInstance)
	if err != nil {
		return fmt.Errorf(
			"unable to marshal Host DNS config data to YAML. Error message: %v", err)
	}

	_, err = out.Write(data)
	if err != nil {
		return fmt.Errorf(
			"unable to write marshaled Host DNS config YAML data. Error message: %v", err)
	}

	return nil
}
