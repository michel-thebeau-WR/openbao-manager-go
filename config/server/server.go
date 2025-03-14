package bao_server

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type user_lockout struct {
	ul_type               string
	lockout_threshold     string
	lockout_duration      string
	lockout_counter_reset string
	disable_lockout       bool
}

type service_reg_kubernetes struct {
	namespace string
	pod_name  string
}

type server_config struct {
	// must be one of: nil, storage_file, storage_inmem, storage_raft, storage_postgresql
	storage    any
	ha_storage any

	// must be one of: nil, listener_unix, or listener_tcp
	listener any

	// must be one of: nil, telemetry_statsite, telemetry_statsd, telemetry_dogstatsd,
	// telemetry_circonus, or telemetry_prometheus
	telemetry any

	// must be one of: nil, seal_alicloudkms, seal_awskms, seal_azurekeyvault, seal_gcpckms,
	// seal_ocikms, seal_pkcs11, seal_transit
	seal any

	service_registration                service_reg_kubernetes
	user_lockout                        user_lockout
	cluster_name                        string
	cache_size                          string
	disable_cache                       bool
	plugin_directory                    string
	plugin_file_uid                     int64
	plugin_file_permissions             string
	default_lease_ttl                   string
	max_lease_ttl                       string
	default_max_request_duration        string
	detect_deadlocks                    string
	raw_storage_endpoint                string
	introspection_endpoint              bool
	ui                                  bool
	pid_file                            string
	enable_response_header_hostname     bool
	enable_response_header_raft_node_id bool
	log_level                           string
	log_format                          string
	log_file                            string
	log_rotate_duration                 string
	log_rotate_bytes                    int64
	imprecise_lease_role_tracking       bool

	// following fields are for HA configurations
	api_addr           string
	cluster_addr       string
	disable_clustering bool
}

// read_yaml reads from a YAML file named file_name, and parses the contents into
func (s server_config) read_yaml(file_name string) {
	data, err := os.ReadFile(file_name)
	if err != nil {
		log.Fatalf("DEBUG Unable to read file %v. Reason: %v", file_name, err)
	}

	err = yaml.Unmarshal(data, &s)
	if err != nil {
		log.Fatalf("DEBUG Issue with parsing YAML file %v. Reason: %v", file_name, err)
	}
}

func (s server_config) write_yaml(file_name string) {
	data, err := yaml.Marshal(&s)
	if err != nil {
		log.Fatalf("DEBUG Issue with parsing configuration into YAML, Reason: %v", err)
	}

	err = os.WriteFile(file_name, data, 0644)
	if err != nil {
		log.Fatalf("DEBUG Unable to write YAML data to file %v. Reason: %v", file_name, err)
	}
}
