//
// Copyright (c) 2025 Wind River Systems, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//

package baoConfig

import (
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"

	"github.com/go-yaml/yaml"
	openbao "github.com/openbao/openbao/api/v2"
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

	// A string of path to the PEM-encoded CA cert file to use to verify
	// Openbao's server SSL certificate
	// Leave this empty if using the default CA cert file location
	CACert string `yaml:"CACert"`

	// ClientCert is the path to the certificate for Vault communication
	ClientCert string `yaml:"ClientCert"`

	// ClientKey is the path to the private key for Vault communication
	ClientKey string `yaml:"ClientKey"`

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

	// Validate YAML input for CACert
	err = configInstance.validateCACert()
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

// Create a new openbao config based on the monitor config
func (configInstance MonitorConfig) NewOpenbaoConfig(dnshost string) (*openbao.Config, error) {
	defConfig := openbao.DefaultConfig()

	// Check if DefaultConfig has issues
	if defConfig.Error != nil {
		return defConfig, fmt.Errorf("issue found in openbao default config: %v", defConfig.Error)
	}
	slog.Debug("No issues found in retrieving openbao default config.")

	// Check if there is a domain name listed under IncludeInCluster
	dnsAddr, ok := configInstance.DNSnames[dnshost]
	if !ok {
		return defConfig, fmt.Errorf("unable to find %v under the list of available DNS names", dnshost)
	}

	// Set the DNS address as the address to openbao
	defConfig.Address = strings.Join([]string{"https://", dnsAddr.Host, ":", strconv.Itoa(dnsAddr.Port)}, "")

	slog.Debug(fmt.Sprintf("Openbao address set to %v", defConfig.Address))

	// Apply CACert entry to openbao config
	var newTLSconfig openbao.TLSConfig
	slog.Debug("Applying the following cert configs:")
	slog.Debug(fmt.Sprintf("CACert: %v", configInstance.CACert))
	slog.Debug(fmt.Sprintf("ClientCert: %v", configInstance.ClientCert))
	slog.Debug(fmt.Sprintf("ClientKey: %v", configInstance.ClientKey))

	newTLSconfig.CACert = configInstance.CACert
	newTLSconfig.ClientCert = configInstance.ClientCert
	newTLSconfig.ClientKey = configInstance.ClientKey

	// This does nothing if newTLSconfig is empty
	err := defConfig.ConfigureTLS(&newTLSconfig)
	if err != nil {
		return defConfig, fmt.Errorf("error with configuring TLS for openbao: %v", err)
	}

	slog.Debug("Configuring TLS successful")
	return defConfig, nil
}

// Parse the new keys from the init responce into the monitor config
func (configInstance *MonitorConfig) ParseInitResponse(dnshost string, responce *openbao.InitResponse) error {
	slog.Debug("Parsing response from /sys/init to monitor configs")

	keyShardheader := strings.Join([]string{"key", "shard", dnshost}, "-")

	// Parse in the root token
	if _, ok := configInstance.Tokens["root_token"]; ok {
		return fmt.Errorf("an entry of the root token was already found")
	}
	if configInstance.Tokens == nil {
		configInstance.Tokens = make(map[string]Token)
	}
	configInstance.Tokens["root_token"] = Token{
		Duration: 0,
		Key:      responce.RootToken,
	}

	// Parse in the key shards for unseal
	for i := range len(responce.Keys) {
		keyShardName := strings.Join([]string{keyShardheader, strconv.Itoa(i)}, "-")
		if _, ok := configInstance.UnsealKeyShards[keyShardName]; ok {
			return fmt.Errorf("an entry of %v was already found under UnsealKeyShards", keyShardName)
		}
		if configInstance.UnsealKeyShards == nil {
			configInstance.UnsealKeyShards = make(map[string]KeyShards)
		}
		configInstance.UnsealKeyShards[keyShardName] = KeyShards{
			Key:       responce.Keys[i],
			KeyBase64: responce.KeysB64[i],
		}
	}

	// Parse in the recovery key shards
	for i := range len(responce.RecoveryKeys) {
		keyShardName := strings.Join([]string{keyShardheader, "recovery", strconv.Itoa(i)}, "-")
		if _, ok := configInstance.UnsealKeyShards[keyShardName]; ok {
			return fmt.Errorf("an entry of %v was already found under UnsealKeyShards", keyShardName)
		}
		if configInstance.UnsealKeyShards == nil {
			configInstance.UnsealKeyShards = make(map[string]KeyShards)
		}
		configInstance.UnsealKeyShards[keyShardName] = KeyShards{
			Key:       responce.RecoveryKeys[i],
			KeyBase64: responce.RecoveryKeysB64[i],
		}
	}

	return nil
}
