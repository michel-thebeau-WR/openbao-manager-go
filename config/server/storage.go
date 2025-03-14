package bao_server

// struct representing the Filesystem storage backend for storage stanza
type storage_file struct {
	// path must not be empty
	path string
}

// struct representing the In-Memory storage backend for storage stanza
type storage_inmem struct{}

// struct representing the retry_join stanza for raft
type retry_join struct {
	leader_api_addr         string
	auto_join               string
	auto_join_scheme        string
	auto_join_port          string
	leader_tls_servername   string
	leader_ca_cert_file     string
	leader_client_cert_file string
	leader_client_key_file  string
	leader_ca_cert          string
	leader_client_cert      string
	leader_client_key       string
}

// struct representing the Integrated Storage (Raft) backend for storage stanza
type storage_raft struct {
	path                         string
	node_id                      string
	performance_multiplier       int32
	trailing_logs                int32
	snapshot_threshold           int32
	snapshot_interval            int32
	retry_join_as_non_voter      bool
	max_entry_size               int64
	max_transaction_size         int64
	autopilot_reconcile_interval string
	autopilot_update_interval    string
	retry_join                   []retry_join
}

// struct representing the PostgreSQL storage backend for storage stanza
type storage_postgresql struct {
	connection_url       string
	table                string
	max_idle_connections int64
	max_parallel         string
	ha_enabled           string
	ha_table             string
	upsert_function      string
	skip_create_table    string
}

// Set default values for storage_raft
func (s *storage_raft) set_default() {
	s.performance_multiplier = 0
	s.trailing_logs = 10000
	s.snapshot_threshold = 8192
	s.snapshot_interval = 120
	s.retry_join = make([]retry_join, 0)
	s.retry_join_as_non_voter = false
	s.max_entry_size = 1048576
	s.max_transaction_size = 8388608
	s.autopilot_reconcile_interval = "10s"
	s.autopilot_update_interval = "2s"
}

// Set default values for storage_postgresql
func (s *storage_postgresql) set_default() {
	s.table = "openbao_kv_store"
	s.max_parallel = "128"
	s.ha_enabled = "false"
	s.ha_table = "openbao_ha_locks"
	s.upsert_function = "openbao_kv_put"
	s.skip_create_table = "false"
}
