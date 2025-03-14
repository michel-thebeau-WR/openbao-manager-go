package bao_server

type listener_unix struct {
	address      string
	socket_mode  string
	socket_user  string
	socket_group string
}

// tls_prefer_server_cipher_suites is deprecated, but kept in here for preserving input data
type tls struct {
	disable                        bool
	cert_file                      string
	key_file                       string
	min_version                    string
	max_version                    string
	cipher_suites                  string
	prefer_server_cipher_suites    string
	require_and_verify_client_cert string
	client_ca_file                 string
	disable_client_certs           string
}

type acme struct {
	ca_directory           string
	cache_path             string
	ca_root                string
	eab_key_id             string
	eab_mac_key            string
	email                  string
	domains                string
	disable_http_challenge bool
	disable_alpn_challenge bool
}

type status_code struct {
	code    string
	headers map[string][]string
}

type listener_tcp struct {
	address                                   string
	cluster_address                           string
	http_idle_timeout                         string
	http_read_header_timeout                  string
	http_read_timeout                         string
	http_write_timeout                        string
	max_request_size                          int64
	max_request_duration                      string
	proxy_protocol_behavior                   string
	proxy_protocol_authorized_addrs           []string
	tls                                       tls
	acme                                      acme
	unauthenticated_metrics_access            bool
	unauthenticated_pprof_access              bool
	unauthenticated_in_flight_requests_access bool
	custom_response_headers                   []status_code
}
