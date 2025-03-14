package bao_config

type telemetry struct {
	usage_gauge_period                 string
	maximum_gauge_cardinality          int64
	disable_hostname                   bool
	enable_hostname_label              bool
	lease_metrics_epsilon              string
	num_lease_metrics_buckets          int64
	add_lease_metrics_namespace_labels bool
	filter_default                     bool
	prefix_filter                      []string
}

type telemetry_statsite struct {
	address string
	telemetry
}

type telemetry_statsd struct {
	address string
	telemetry
}

type telemetry_dogstatsd struct {
	address string
	tags    []string
	telemetry
}

type telemetry_circonus struct {
	api_token                     string
	api_app                       string
	api_url                       string
	submission_interval           string
	submission_url                string
	check_id                      string
	check_force_metric_activation bool
	check_instance_id             string
	check_search_tag              string
	check_display_name            string
	check_tags                    string
	broker_id                     string
	broker_select_tag             string
	telemetry
}

type telemetry_prometheus struct {
	retention_time   string
	disable_hostname bool
	telemetry
}
