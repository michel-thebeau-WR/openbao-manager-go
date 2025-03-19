package baoConfig

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"

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

type MonitorConfig struct {
	// A map value listing all DNS names
	// Key: Domain name
	// Value: URL. consisting of host address and port number
	DNSnames map[string]URL `yaml:"IncludeInCluster"`

	// A map value listing all authentication tokens
	// Key: release id
	// Value: Token. consisting of lease duration and the token key
	Tokens map[string]Token `yaml:"Tokens"`
}

func (s *MonitorConfig) ReadYAMLMonitorConfig(in io.Reader) error {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return fmt.Errorf("unable to read Host DNS config data from input. Error message: %v", err)
	}

	err = yaml.Unmarshal(data, s)
	if err != nil {
		return fmt.Errorf("unable to unmarshal Host DNS config YAML data. Error message: %v", err)
	}

	// Validate YAML input for DNSnames
	for domain_name, url := range s.DNSnames {
		// If either Host or Port number is empty, then the domain entry is invalid
		if url.Host == "" || url.Port == 0 {
			return fmt.Errorf("the domain entry %v in config InvludeInCluster is invalid", domain_name)
		}
	}

	// Validate YAML input for Tokens
	rootExists := false
	r, _ := regexp.Compile("[sbr][.][a-zA-Z0-9]{24,}")
	for releaseID, token := range s.Tokens {
		if token.Duration == 0 {
			// There can only be one root token
			if rootExists {
				return fmt.Errorf("there are two or more root tokens listed")
			} else {
				rootExists = true
			}
		}
		// Token key should have s, b, or r as the first character, and . as the second.
		// The body of the token (key[2:]) should be 24 characters or more
		if !r.MatchString(token.Key) {
			return fmt.Errorf("the token with release id %v has wrong key format", releaseID)
		}
	}

	return nil
}

func (s MonitorConfig) WriteYAMLMonitorConfig(out io.Writer) error {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("unable to marshal Host DNS config data to YAML. Error message: %v", err)
	}

	_, err = out.Write(data)
	if err != nil {
		return fmt.Errorf("unable to write marshaled Host DNS config YAML data. Error message: %v", err)
	}

	return nil
}
