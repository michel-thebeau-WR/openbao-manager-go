package baoConfig

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type URL struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type MonitorConfig struct {
	DNSnames map[string]URL `yaml:"IncludeInCluster"`
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
