package bao_config

type seal struct {
	disabled string
}

type kms struct {
	region     string
	access_key string
	secret_key string
	kms_key_id string
	seal
}

type seal_alicloudkms struct {
	domain string
	kms
}

type seal_awskms struct {
	session_tocken string
	endpoint       string
	kms
}

type seal_azurekeyvault struct {
	tenant_id     string
	client_id     string
	client_secret string
	environment   string
	vault_name    string
	key_name      string
	resource      string
	seal
}

type seal_gcpckms struct {
	credentials string
	project     string
	region      string
	key_ring    string
	crypto_key  string
	seal
}

type seal_ocikms struct {
	key_id            string
	crypto_endpoint   string
	manage_endpoint   string
	auth_type_api_key bool
	seal
}

type seal_pkcs11 struct {
	lib               string
	slot              string
	token_label       string
	pin               string
	key_label         string
	default_key_label string
	key_id            string
	mechanism         string
	seal
}

type seal_transit struct {
	address         string
	token           string
	key_name        string
	mount_path      string
	namespace       string
	disable_renewal string
	tls_ca_cert     string
	tls_client_cert string
	tls_client_key  string
	tls_server_name string
	tls_skip_verify bool
	seal
}
