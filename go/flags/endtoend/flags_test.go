/*
Copyright 2022 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flags

// These tests ensure that we are changing flags intentionally and do not accidentally make
// changes such as removing a flag. Since there's no way to test the command-line
// flag handling portion explicitly in the unit tests we do so here.

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	helpOutput = map[string]string{
		"vtgate": `Usage of vtgate:
  --allowed_tablet_types value
	Specifies the tablet types this vtgate is allowed to route queries to
  --alsologtostderr
	log to standard error as well as files
  --buffer_drain_concurrency int
	Maximum number of requests retried simultaneously. More concurrency will increase the load on the PRIMARY vttablet when draining the buffer. (default 1)
  --buffer_implementation string
	Allowed values: healthcheck (legacy implementation), keyspace_events (default) (default keyspace_events)
  --buffer_keyspace_shards string
	If not empty, limit buffering to these entries (comma separated). Entry format: keyspace or keyspace/shard. Requires --enable_buffer=true.
  --buffer_max_failover_duration duration
	Stop buffering completely if a failover takes longer than this duration. (default 20s)
  --buffer_min_time_between_failovers duration
	Minimum time between the end of a failover and the start of the next one (tracked per shard). Faster consecutive failovers will not trigger buffering. (default 1m0s)
  --buffer_size int
	Maximum number of buffered requests in flight (across all ongoing failovers). (default 1000)
  --buffer_window duration
	Duration for how long a request should be buffered at most. (default 10s)
  --catch-sigpipe
	catch and ignore SIGPIPE on stdout and stderr if specified
  --cell string
	cell to use (default test_nj)
  --cells_to_watch string
	comma-separated list of cells for watching tablets
  --consul_auth_static_file string
	JSON File to read the topos/tokens from.
  --cpu_profile string
	deprecated: use '-pprof=cpu' instead
  --datadog-agent-host string
	host to send spans to. if empty, no tracing will be done
  --datadog-agent-port string
	port to send spans to. if empty, no tracing will be done
  --dbddl_plugin string
	controls how to handle CREATE/DROP DATABASE. use it if you are using your own database provisioning service (default fail)
  --ddl_strategy string
	Set default strategy for DDL statements. Override with @@ddl_strategy session variable (default direct)
  --default_tablet_type value
	The default tablet type to set for queries, when one is not explicitly selected (default PRIMARY)
  --disable_local_gateway
	deprecated: if specified, this process will not route any queries to local tablets in the local cell
  --discovery_high_replication_lag_minimum_serving duration
	the replication lag that is considered too high when applying the min_number_serving_vttablets threshold (default 2h0m0s)
  --discovery_low_replication_lag duration
	the replication lag that is considered low enough to be healthy (default 30s)
  --emit_stats
	If set, emit stats to push-based monitoring and stats backends
  --enable_buffer
	Enable buffering (stalling) of primary traffic during failovers.
  --enable_buffer_dry_run
	Detect and log failover events, but do not actually buffer requests.
  --enable_direct_ddl
	Allow users to submit direct DDL statements (default true)
  --enable_online_ddl
	Allow users to submit, review and control Online DDL (default true)
  --enable_set_var
	This will enable the use of MySQL's SET_VAR query hint for certain system variables instead of using reserved connections (default true)
  --enable_system_settings
	This will enable the system settings to be changed per session at the database connection level (default true)
  --foreign_key_mode string
	This is to provide how to handle foreign key constraint in create/alter table. Valid values are: allow, disallow (default allow)
  --gate_query_cache_lfu
	gate server cache algorithm. when set to true, a new cache algorithm based on a TinyLFU admission policy will be used to improve cache behavior and prevent pollution from sparse queries (default true)
  --gate_query_cache_memory int
	gate server query cache size in bytes, maximum amount of memory to be cached. vtgate analyzes every incoming query and generate a query plan, these plans are being cached in a lru cache. This config controls the capacity of the lru cache. (default 33554432)
  --gate_query_cache_size int
	gate server query cache size, maximum number of queries to be cached. vtgate analyzes every incoming query and generate a query plan, these plans are being cached in a cache. This config controls the expected amount of unique entries in the cache. (default 5000)
  --gateway_initial_tablet_timeout duration
	At startup, the tabletGateway will wait up to this duration to get at least one tablet per keyspace/shard/tablet type (default 30s)
  --grpc_auth_mode string
	Which auth plugin implementation to use (eg: static)
  --grpc_auth_mtls_allowed_substrings string
	List of substrings of at least one of the client certificate names (separated by colon).
  --grpc_auth_static_client_creds string
	when using grpc_static_auth in the server, this file provides the credentials to use to authenticate with server
  --grpc_auth_static_password_file string
	JSON File to read the users/passwords from.
  --grpc_ca string
	server CA to use for gRPC connections, requires TLS, and enforces client certificate check
  --grpc_cert string
	server certificate to use for gRPC connections, requires grpc_key, enables TLS
  --grpc_compression string
	Which protocol to use for compressing gRPC. Default: nothing. Supported: snappy
  --grpc_crl string
	path to a certificate revocation list in PEM format, client certificates will be further verified against this file during TLS handshake
  --grpc_enable_optional_tls
	enable optional TLS mode when a server accepts both TLS and plain-text connections on the same port
  --grpc_enable_tracing
	Enable GRPC tracing
  --grpc_initial_conn_window_size int
	gRPC initial connection window size
  --grpc_initial_window_size int
	gRPC initial window size
  --grpc_keepalive_time duration
	After a duration of this time, if the client doesn't see any activity, it pings the server to see if the transport is still alive. (default 10s)
  --grpc_keepalive_timeout duration
	After having pinged for keepalive check, the client waits for a duration of Timeout and if no activity is seen even after that the connection is closed. (default 10s)
  --grpc_key string
	server private key to use for gRPC connections, requires grpc_cert, enables TLS
  --grpc_max_connection_age duration
	Maximum age of a client connection before GoAway is sent. (default 2562047h47m16.854775807s)
  --grpc_max_connection_age_grace duration
	Additional grace period after grpc_max_connection_age, after which connections are forcibly closed. (default 2562047h47m16.854775807s)
  --grpc_max_message_size int
	Maximum allowed RPC message size. Larger messages will be rejected by gRPC with the error 'exceeding the max size'. (default 16777216)
  --grpc_port int
	Port to listen on for gRPC calls
  --grpc_prometheus
	Enable gRPC monitoring with Prometheus
  --grpc_server_ca string
	path to server CA in PEM format, which will be combine with server cert, return full certificate chain to clients
  --grpc_server_initial_conn_window_size int
	gRPC server initial connection window size
  --grpc_server_initial_window_size int
	gRPC server initial window size
  --grpc_server_keepalive_enforcement_policy_min_time duration
	gRPC server minimum keepalive time (default 10s)
  --grpc_server_keepalive_enforcement_policy_permit_without_stream
	gRPC server permit client keepalive pings even when there are no active streams (RPCs)
  --grpc_use_effective_callerid
	If set, and SSL is not used, will set the immediate caller id from the effective caller id's principal.
  --healthcheck_retry_delay duration
	health check retry delay (default 2ms)
  --healthcheck_timeout duration
	the health check timeout period (default 1m0s)
  --jaeger-agent-host string
	host and port to send spans to. if empty, no tracing will be done
  --keep_logs duration
	keep logs for this long (using ctime) (zero to keep forever)
  --keep_logs_by_mtime duration
	keep logs for this long (using mtime) (zero to keep forever)
  --keyspaces_to_watch value
	Specifies which keyspaces this vtgate should have access to while routing queries or accessing the vschema
  --lameduck-period duration
	keep running at least this long after SIGTERM before stopping (default 50ms)
  --legacy_replication_lag_algorithm
	use the legacy algorithm when selecting the vttablets for serving (default true)
  --lock_heartbeat_time duration
	If there is lock function used. This will keep the lock connection active by using this heartbeat (default 5s)
  --log_backtrace_at value
	when logging hits line file:N, emit a stack trace
  --log_dir string
	If non-empty, write log files in this directory
  --log_err_stacks
	log stack traces for errors
  --log_queries_to_file string
	Enable query logging to the specified file
  --log_rotate_max_size uint
	size in bytes at which logs are rotated (glog.MaxSize) (default 1887436800)
  --logtostderr
	log to standard error instead of files
  --max_memory_rows int
	Maximum number of rows that will be held in memory for intermediate results as well as the final result. (default 300000)
  --max_payload_size int
	The threshold for query payloads in bytes. A payload greater than this threshold will result in a failure to handle the query.
  --mem-profile-rate int
	deprecated: use '-pprof=mem' instead (default 524288)
  --message_stream_grace_period duration
	the amount of time to give for a vttablet to resume if it ends a message stream, usually because of a reparent. (default 30s)
  --min_number_serving_vttablets int
	the minimum number of vttablets for each replicating tablet_type (e.g. replica, rdonly) that will be continue to be used even with replication lag above discovery_low_replication_lag, but still below discovery_high_replication_lag_minimum_serving (default 2)
  --mutex-profile-fraction int
	deprecated: use '-pprof=mutex' instead
  --mysql_allow_clear_text_without_tls
	If set, the server will allow the use of a clear text password over non-SSL connections.
  --mysql_auth_server_impl string
	Which auth server implementation to use. Options: none, ldap, clientcert, static, vault. (default static)
  --mysql_auth_server_static_file string
	JSON File to read the users/passwords from.
  --mysql_auth_server_static_string string
	JSON representation of the users/passwords config.
  --mysql_auth_static_reload_interval duration
	Ticker to reload credentials
  --mysql_auth_vault_addr string
	URL to Vault server
  --mysql_auth_vault_path string
	Vault path to vtgate credentials JSON blob, e.g.: secret/data/prod/vtgatecreds
  --mysql_auth_vault_role_mountpoint string
	Vault AppRole mountpoint; can also be passed using VAULT_MOUNTPOINT environment variable (default approle)
  --mysql_auth_vault_role_secretidfile string
	Path to file containing Vault AppRole secret_id; can also be passed using VAULT_SECRETID environment variable
  --mysql_auth_vault_roleid string
	Vault AppRole id; can also be passed using VAULT_ROLEID environment variable
  --mysql_auth_vault_timeout duration
	Timeout for vault API operations (default 10s)
  --mysql_auth_vault_tls_ca string
	Path to CA PEM for validating Vault server certificate
  --mysql_auth_vault_tokenfile string
	Path to file containing Vault auth token; token can also be passed using VAULT_TOKEN environment variable
  --mysql_auth_vault_ttl duration
	How long to cache vtgate credentials from the Vault server (default 30m0s)
  --mysql_clientcert_auth_method string
	client-side authentication method to use. Supported values: mysql_clear_password, dialog. (default mysql_clear_password)
  --mysql_default_workload string
	Default session workload (OLTP, OLAP, DBA) (default OLTP)
  --mysql_ldap_auth_config_file string
	JSON File from which to read LDAP server config.
  --mysql_ldap_auth_config_string string
	JSON representation of LDAP server config.
  --mysql_ldap_auth_method string
	client-side authentication method to use. Supported values: mysql_clear_password, dialog. (default mysql_clear_password)
  --mysql_server_bind_address string
	Binds on this address when listening to MySQL binary protocol. Useful to restrict listening to 'localhost' only for instance.
  --mysql_server_flush_delay duration
	Delay after which buffered response will be flushed to the client. (default 100ms)
  --mysql_server_port int
	If set, also listen for MySQL binary protocol connections on this port. (default -1)
  --mysql_server_query_timeout duration
	mysql query timeout
  --mysql_server_read_timeout duration
	connection read timeout
  --mysql_server_require_secure_transport
	Reject insecure connections but only if mysql_server_ssl_cert and mysql_server_ssl_key are provided
  --mysql_server_socket_path string
	This option specifies the Unix socket file to use when listening for local connections. By default it will be empty and it won't listen to a unix socket
  --mysql_server_ssl_ca string
	Path to ssl CA for mysql server plugin SSL. If specified, server will require and validate client certs.
  --mysql_server_ssl_cert string
	Path to the ssl cert for mysql server plugin SSL
  --mysql_server_ssl_crl string
	Path to ssl CRL for mysql server plugin SSL
  --mysql_server_ssl_key string
	Path to ssl key for mysql server plugin SSL
  --mysql_server_ssl_server_ca string
	path to server CA in PEM format, which will be combine with server cert, return full certificate chain to clients
  --mysql_server_tls_min_version string
	Configures the minimal TLS version negotiated when SSL is enabled. Defaults to TLSv1.2. Options: TLSv1.0, TLSv1.1, TLSv1.2, TLSv1.3.
  --mysql_server_version string
	MySQL server version to advertise.
  --mysql_server_write_timeout duration
	connection write timeout
  --mysql_slow_connect_warn_threshold duration
	Warn if it takes more than the given threshold for a mysql connection to establish
  --mysql_tcp_version string
	Select tcp, tcp4, or tcp6 to control the socket type. (default tcp)
  --no_scatter
	when set to true, the planner will fail instead of producing a plan that includes scatter queries
  --normalize_queries
	Rewrite queries with bind vars. Turn this off if the app itself sends normalized queries with bind vars. (default true)
  --onclose_timeout duration
	wait no more than this for OnClose handlers before stopping (default 1ns)
  --onterm_timeout duration
	wait no more than this for OnTermSync handlers before stopping (default 10s)
  --opentsdb_uri string
	URI of opentsdb /api/put method
  --pid_file string
	If set, the process will write its pid to the named file, and delete it on graceful shutdown.
  --planner_version string
	Sets the default planner to use when the session has not changed it. Valid values are: V3, Gen4, Gen4Greedy and Gen4Fallback. Gen4Fallback tries the gen4 planner and falls back to the V3 planner if the gen4 fails. (default gen4)
  --port int
	port for the server
  --pprof string
	enable profiling
  --proxy_protocol
	Enable HAProxy PROXY protocol on MySQL listener socket
  --purge_logs_interval duration
	how often try to remove old logs (default 1h0m0s)
  --querylog-filter-tag string
	string that must be present in the query for it to be logged; if using a value as the tag, you need to disable query normalization
  --querylog-format string
	format for query logs ("text" or "json") (default text)
  --querylog-row-threshold uint
	Number of rows a query has to return or affect before being logged; not useful for streaming queries. 0 means all queries will be logged.
  --redact-debug-ui-queries
	redact full queries and bind variables from debug UI
  --remote_operation_timeout duration
	time to wait for a remote operation (default 30s)
  --retry-count int
	retry count (default 2)
  --schema_change_signal
	Enable the schema tracker; requires queryserver-config-schema-change-signal to be enabled on the underlying vttablets for this to work
  --schema_change_signal_user string
	User to be used to send down query to vttablet to retrieve schema changes
  --security_policy string
	the name of a registered security policy to use for controlling access to URLs - empty means allow all for anyone (built-in policies: deny-all, read-only)
  --service_map value
	comma separated list of services to enable (or disable if prefixed with '-') Example: grpc-vtworker
  --sql-max-length-errors int
	truncate queries in error logs to the given length (default unlimited)
  --sql-max-length-ui int
	truncate queries in debug UIs to the given length (default 512) (default 512)
  --srv_topo_cache_refresh duration
	how frequently to refresh the topology for cached entries (default 1s)
  --srv_topo_cache_ttl duration
	how long to use cached entries for topology (default 1s)
  --srv_topo_timeout duration
	topo server timeout (default 5s)
  --stats_backend string
	The name of the registered push-based monitoring/stats backend to use
  --stats_combine_dimensions string
	List of dimensions to be combined into a single "all" value in exported stats vars
  --stats_common_tags string
	Comma-separated list of common tags for the stats backend. It provides both label and values. Example: label1:value1,label2:value2
  --stats_drop_variables string
	Variables to be dropped from the list of exported variables.
  --stats_emit_period duration
	Interval between emitting stats to all registered backends (default 1m0s)
  --statsd_address string
	Address for statsd client
  --statsd_sample_rate float
	 (default 1)
  --stderrthreshold value
	logs at or above this threshold go to stderr (default 1)
  --stream_buffer_size int
	the number of bytes sent from vtgate for each stream call. It's recommended to keep this value in sync with vttablet's query-server-config-stream-buffer-size. (default 32768)
  --tablet_filters value
	Specifies a comma-separated list of 'keyspace|shard_name or keyrange' values to filter the tablets to watch
  --tablet_grpc_ca string
	the server ca to use to validate servers when connecting
  --tablet_grpc_cert string
	the cert to use to connect
  --tablet_grpc_crl string
	the server crl to use to validate server certificates when connecting
  --tablet_grpc_key string
	the key to use to connect
  --tablet_grpc_server_name string
	the server name to use to validate server certificate
  --tablet_manager_protocol string
	the protocol to use to talk to vttablet (default grpc)
  --tablet_protocol string
	how to talk to the vttablets (default grpc)
  --tablet_refresh_interval duration
	tablet refresh interval (default 1m0s)
  --tablet_refresh_known_tablets
	tablet refresh reloads the tablet address/port map from topo in case it changes (default true)
  --tablet_types_to_wait string
	wait till connected for specified tablet types during Gateway initialization
  --tablet_url_template string
	format string describing debug tablet url formatting. See the Go code for getTabletDebugURL() how to customize this. (default http://{{.GetTabletHostPort}})
  --topo_consul_lock_delay duration
	LockDelay for consul session. (default 15s)
  --topo_consul_lock_session_checks string
	List of checks for consul session. (default serfHealth)
  --topo_consul_lock_session_ttl string
	TTL for consul session.
  --topo_consul_watch_poll_duration duration
	time of the long poll for watch queries. (default 30s)
  --topo_etcd_lease_ttl int
	Lease TTL for locks and leader election. The client will use KeepAlive to keep the lease going. (default 30)
  --topo_etcd_tls_ca string
	path to the ca to use to validate the server cert when connecting to the etcd topo server
  --topo_etcd_tls_cert string
	path to the client cert to use to connect to the etcd topo server, requires topo_etcd_tls_key, enables TLS
  --topo_etcd_tls_key string
	path to the client key to use to connect to the etcd topo server, enables TLS
  --topo_global_root string
	the path of the global topology data in the global topology server
  --topo_global_server_address string
	the address of the global topology server
  --topo_implementation string
	the topology implementation to use
  --topo_k8s_context string
	The kubeconfig context to use, overrides the 'current-context' from the config
  --topo_k8s_kubeconfig string
	Path to a valid kubeconfig file. When running as a k8s pod inside the same cluster you wish to use as the topo, you may omit this and the below arguments, and Vitess is capable of auto-discovering the correct values. https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/#accessing-the-api-from-a-pod
  --topo_k8s_namespace string
	The kubernetes namespace to use for all objects. Default comes from the context or in-cluster config
  --topo_read_concurrency int
	concurrent topo reads (default 32)
  --topo_zk_auth_file string
	auth to use when connecting to the zk topo server, file contents should be <scheme>:<auth>, e.g., digest:user:pass
  --topo_zk_base_timeout duration
	zk base timeout (see zk.Connect) (default 30s)
  --topo_zk_max_concurrency int
	maximum number of pending requests to send to a Zookeeper server. (default 64)
  --topo_zk_tls_ca string
	the server ca to use to validate servers when connecting to the zk topo server
  --topo_zk_tls_cert string
	the cert to use to connect to the zk topo server, requires topo_zk_tls_key, enables TLS
  --topo_zk_tls_key string
	the key to use to connect to the zk topo server, enables TLS
  --tracer string
	tracing service to use (default noop)
  --tracing-enable-logging
	whether to enable logging in the tracing service
  --tracing-sampling-rate value
	sampling rate for the probabilistic jaeger sampler (default 0.1)
  --tracing-sampling-type value
	sampling strategy to use for jaeger. possible values are 'const', 'probabilistic', 'rateLimiting', or 'remote' (default const)
  --transaction_mode string
	SINGLE: disallow multi-db transactions, MULTI: allow multi-db transactions with best effort commit, TWOPC: allow multi-db transactions with 2pc commit (default MULTI)
  --v value
	log level for V logs
  --version
	print binary version
  --vmodule value
	comma-separated list of pattern=N settings for file-filtered logging
  --vschema_ddl_authorized_users string
	List of users authorized to execute vschema ddl operations, or '%' to allow all users.
  --vtctld_addr string
	address of a vtctld instance
  --vtgate-config-terse-errors
	prevent bind vars from escaping in returned errors
  --warn_memory_rows int
	Warning threshold for in-memory results. A row count higher than this amount will cause the VtGateWarnings.ResultsExceeded counter to be incremented. (default 30000)
  --warn_payload_size int
	The warning threshold for query payloads in bytes. A payload greater than this threshold will cause the VtGateWarnings.WarnPayloadSizeExceeded counter to be incremented.
  --warn_sharded_only
	If any features that are only available in unsharded mode are used, query execution warnings will be added to the session
`,
		"vttablet": `Usage of vttablet:
  --allowed_tablet_types value
	Specifies the tablet types this vtgate is allowed to route queries to
  --alsologtostderr
	log to standard error as well as files
  --app_idle_timeout duration
	Idle timeout for app connections (default 1m0s)
  --app_pool_size int
	Size of the connection pool for app connections (default 40)
  --azblob_backup_account_key_file string
	Path to a file containing the Azure Storage account key; if this flag is unset, the environment variable VT_AZBLOB_ACCOUNT_KEY will be used as the key itself (NOT a file path)
  --azblob_backup_account_name string
	Azure Storage Account name for backups; if this flag is unset, the environment variable VT_AZBLOB_ACCOUNT_NAME will be used
  --azblob_backup_container_name string
	Azure Blob Container Name
  --azblob_backup_parallelism int
	Azure Blob operation parallelism (requires extra memory when increased) (default 1)
  --azblob_backup_storage_root string
	Root prefix for all backup-related Azure Blobs; this should exclude both initial and trailing '/' (e.g. just 'a/b' not '/a/b/')
  --backup_engine_implementation string
	Specifies which implementation to use for creating new backups (builtin or xtrabackup). Restores will always be done with whichever engine created a given backup. (default builtin)
  --backup_storage_block_size int
	if backup_storage_compress is true, backup_storage_block_size sets the byte size for each block while compressing (default is 250000). (default 250000)
  --backup_storage_compress
	if set, the backup files will be compressed (default is true). Set to false for instance if a backup_storage_hook is specified and it compresses the data. (default true)
  --backup_storage_hook string
	if set, we send the contents of the backup files through this hook.
  --backup_storage_implementation string
	which implementation to use for the backup storage feature
  --backup_storage_number_blocks int
	if backup_storage_compress is true, backup_storage_number_blocks sets the number of blocks that can be processed, at once, before the writer blocks, during compression (default is 2). It should be equal to the number of CPUs available for compression (default 2)
  --binlog_host string
	PITR restore parameter: hostname/IP of binlog server.
  --binlog_password string
	PITR restore parameter: password of binlog server.
  --binlog_player_grpc_ca string
	the server ca to use to validate servers when connecting
  --binlog_player_grpc_cert string
	the cert to use to connect
  --binlog_player_grpc_crl string
	the server crl to use to validate server certificates when connecting
  --binlog_player_grpc_key string
	the key to use to connect
  --binlog_player_grpc_server_name string
	the server name to use to validate server certificate
  --binlog_player_protocol string
	the protocol to download binlogs from a vttablet (default grpc)
  --binlog_port int
	PITR restore parameter: port of binlog server.
  --binlog_ssl_ca string
	PITR restore parameter: Filename containing TLS CA certificate to verify binlog server TLS certificate against.
  --binlog_ssl_cert string
	PITR restore parameter: Filename containing mTLS client certificate to present to binlog server as authentication.
  --binlog_ssl_key string
	PITR restore parameter: Filename containing mTLS client private key for use in binlog server authentication.
  --binlog_ssl_server_name string
	PITR restore parameter: TLS server name (common name) to verify against for the binlog server we are connecting to (If not set: use the hostname or IP supplied in -binlog_host).
  --binlog_use_v3_resharding_mode
	True iff the binlog streamer should use V3-style sharding, which doesn't require a preset sharding key column. (default true)
  --binlog_user string
	PITR restore parameter: username of binlog server.
  --builtinbackup_mysqld_timeout duration
	how long to wait for mysqld to shutdown at the start of the backup (default 10m0s)
  --builtinbackup_progress duration
	how often to send progress updates when backing up large files (default 5s)
  --catch-sigpipe
	catch and ignore SIGPIPE on stdout and stderr if specified
  --ceph_backup_storage_config string
	Path to JSON config file for ceph backup storage (default ceph_backup_config.json)
  --client-found-rows-pool-size int
	DEPRECATED: queryserver-config-transaction-cap will be used instead.
  --consul_auth_static_file string
	JSON File to read the topos/tokens from.
  --cpu_profile string
	deprecated: use '-pprof=cpu' instead
  --datadog-agent-host string
	host to send spans to. if empty, no tracing will be done
  --datadog-agent-port string
	port to send spans to. if empty, no tracing will be done
  --db-config-allprivs-charset string
	deprecated: use db_charset (default utf8mb4)
  --db-config-allprivs-flags uint
	deprecated: use db_flags
  --db-config-allprivs-flavor string
	deprecated: use db_flavor
  --db-config-allprivs-host string
	deprecated: use db_host
  --db-config-allprivs-pass string
	db allprivs deprecated: use db_allprivs_password
  --db-config-allprivs-port int
	deprecated: use db_port
  --db-config-allprivs-server_name string
	deprecated: use db_server_name
  --db-config-allprivs-ssl-ca string
	deprecated: use db_ssl_ca
  --db-config-allprivs-ssl-ca-path string
	deprecated: use db_ssl_ca_path
  --db-config-allprivs-ssl-cert string
	deprecated: use db_ssl_cert
  --db-config-allprivs-ssl-key string
	deprecated: use db_ssl_key
  --db-config-allprivs-uname string
	deprecated: use db_allprivs_user (default vt_allprivs)
  --db-config-allprivs-unixsocket string
	deprecated: use db_socket
  --db-config-app-charset string
	deprecated: use db_charset (default utf8mb4)
  --db-config-app-flags uint
	deprecated: use db_flags
  --db-config-app-flavor string
	deprecated: use db_flavor
  --db-config-app-host string
	deprecated: use db_host
  --db-config-app-pass string
	db app deprecated: use db_app_password
  --db-config-app-port int
	deprecated: use db_port
  --db-config-app-server_name string
	deprecated: use db_server_name
  --db-config-app-ssl-ca string
	deprecated: use db_ssl_ca
  --db-config-app-ssl-ca-path string
	deprecated: use db_ssl_ca_path
  --db-config-app-ssl-cert string
	deprecated: use db_ssl_cert
  --db-config-app-ssl-key string
	deprecated: use db_ssl_key
  --db-config-app-uname string
	deprecated: use db_app_user (default vt_app)
  --db-config-app-unixsocket string
	deprecated: use db_socket
  --db-config-appdebug-charset string
	deprecated: use db_charset (default utf8mb4)
  --db-config-appdebug-flags uint
	deprecated: use db_flags
  --db-config-appdebug-flavor string
	deprecated: use db_flavor
  --db-config-appdebug-host string
	deprecated: use db_host
  --db-config-appdebug-pass string
	db appdebug deprecated: use db_appdebug_password
  --db-config-appdebug-port int
	deprecated: use db_port
  --db-config-appdebug-server_name string
	deprecated: use db_server_name
  --db-config-appdebug-ssl-ca string
	deprecated: use db_ssl_ca
  --db-config-appdebug-ssl-ca-path string
	deprecated: use db_ssl_ca_path
  --db-config-appdebug-ssl-cert string
	deprecated: use db_ssl_cert
  --db-config-appdebug-ssl-key string
	deprecated: use db_ssl_key
  --db-config-appdebug-uname string
	deprecated: use db_appdebug_user (default vt_appdebug)
  --db-config-appdebug-unixsocket string
	deprecated: use db_socket
  --db-config-dba-charset string
	deprecated: use db_charset (default utf8mb4)
  --db-config-dba-flags uint
	deprecated: use db_flags
  --db-config-dba-flavor string
	deprecated: use db_flavor
  --db-config-dba-host string
	deprecated: use db_host
  --db-config-dba-pass string
	db dba deprecated: use db_dba_password
  --db-config-dba-port int
	deprecated: use db_port
  --db-config-dba-server_name string
	deprecated: use db_server_name
  --db-config-dba-ssl-ca string
	deprecated: use db_ssl_ca
  --db-config-dba-ssl-ca-path string
	deprecated: use db_ssl_ca_path
  --db-config-dba-ssl-cert string
	deprecated: use db_ssl_cert
  --db-config-dba-ssl-key string
	deprecated: use db_ssl_key
  --db-config-dba-uname string
	deprecated: use db_dba_user (default vt_dba)
  --db-config-dba-unixsocket string
	deprecated: use db_socket
  --db-config-erepl-charset string
	deprecated: use db_charset (default utf8mb4)
  --db-config-erepl-dbname string
	deprecated: dbname does not need to be explicitly configured
  --db-config-erepl-flags uint
	deprecated: use db_flags
  --db-config-erepl-flavor string
	deprecated: use db_flavor
  --db-config-erepl-host string
	deprecated: use db_host
  --db-config-erepl-pass string
	db erepl deprecated: use db_erepl_password
  --db-config-erepl-port int
	deprecated: use db_port
  --db-config-erepl-server_name string
	deprecated: use db_server_name
  --db-config-erepl-ssl-ca string
	deprecated: use db_ssl_ca
  --db-config-erepl-ssl-ca-path string
	deprecated: use db_ssl_ca_path
  --db-config-erepl-ssl-cert string
	deprecated: use db_ssl_cert
  --db-config-erepl-ssl-key string
	deprecated: use db_ssl_key
  --db-config-erepl-uname string
	deprecated: use db_erepl_user (default vt_erepl)
  --db-config-erepl-unixsocket string
	deprecated: use db_socket
  --db-config-filtered-charset string
	deprecated: use db_charset (default utf8mb4)
  --db-config-filtered-flags uint
	deprecated: use db_flags
  --db-config-filtered-flavor string
	deprecated: use db_flavor
  --db-config-filtered-host string
	deprecated: use db_host
  --db-config-filtered-pass string
	db filtered deprecated: use db_filtered_password
  --db-config-filtered-port int
	deprecated: use db_port
  --db-config-filtered-server_name string
	deprecated: use db_server_name
  --db-config-filtered-ssl-ca string
	deprecated: use db_ssl_ca
  --db-config-filtered-ssl-ca-path string
	deprecated: use db_ssl_ca_path
  --db-config-filtered-ssl-cert string
	deprecated: use db_ssl_cert
  --db-config-filtered-ssl-key string
	deprecated: use db_ssl_key
  --db-config-filtered-uname string
	deprecated: use db_filtered_user (default vt_filtered)
  --db-config-filtered-unixsocket string
	deprecated: use db_socket
  --db-config-repl-charset string
	deprecated: use db_charset (default utf8mb4)
  --db-config-repl-flags uint
	deprecated: use db_flags
  --db-config-repl-flavor string
	deprecated: use db_flavor
  --db-config-repl-host string
	deprecated: use db_host
  --db-config-repl-pass string
	db repl deprecated: use db_repl_password
  --db-config-repl-port int
	deprecated: use db_port
  --db-config-repl-server_name string
	deprecated: use db_server_name
  --db-config-repl-ssl-ca string
	deprecated: use db_ssl_ca
  --db-config-repl-ssl-ca-path string
	deprecated: use db_ssl_ca_path
  --db-config-repl-ssl-cert string
	deprecated: use db_ssl_cert
  --db-config-repl-ssl-key string
	deprecated: use db_ssl_key
  --db-config-repl-uname string
	deprecated: use db_repl_user (default vt_repl)
  --db-config-repl-unixsocket string
	deprecated: use db_socket
  --db-credentials-file string
	db credentials file; send SIGHUP to reload this file
  --db-credentials-server string
	db credentials server type ('file' - file implementation; 'vault' - HashiCorp Vault implementation) (default file)
  --db-credentials-vault-addr string
	URL to Vault server
  --db-credentials-vault-path string
	Vault path to credentials JSON blob, e.g.: secret/data/prod/dbcreds
  --db-credentials-vault-role-mountpoint string
	Vault AppRole mountpoint; can also be passed using VAULT_MOUNTPOINT environment variable (default approle)
  --db-credentials-vault-role-secretidfile string
	Path to file containing Vault AppRole secret_id; can also be passed using VAULT_SECRETID environment variable
  --db-credentials-vault-roleid string
	Vault AppRole id; can also be passed using VAULT_ROLEID environment variable
  --db-credentials-vault-timeout duration
	Timeout for vault API operations (default 10s)
  --db-credentials-vault-tls-ca string
	Path to CA PEM for validating Vault server certificate
  --db-credentials-vault-tokenfile string
	Path to file containing Vault auth token; token can also be passed using VAULT_TOKEN environment variable
  --db-credentials-vault-ttl duration
	How long to cache DB credentials from the Vault server (default 30m0s)
  --db_allprivs_password string
	db allprivs password
  --db_allprivs_use_ssl
	Set this flag to false to make the allprivs connection to not use ssl (default true)
  --db_allprivs_user string
	db allprivs user userKey (default vt_allprivs)
  --db_app_password string
	db app password
  --db_app_use_ssl
	Set this flag to false to make the app connection to not use ssl (default true)
  --db_app_user string
	db app user userKey (default vt_app)
  --db_appdebug_password string
	db appdebug password
  --db_appdebug_use_ssl
	Set this flag to false to make the appdebug connection to not use ssl (default true)
  --db_appdebug_user string
	db appdebug user userKey (default vt_appdebug)
  --db_charset string
	Character set used for this tablet. (default utf8mb4)
  --db_conn_query_info
	enable parsing and processing of QUERY_OK info fields
  --db_connect_timeout_ms int
	connection timeout to mysqld in milliseconds (0 for no timeout)
  --db_dba_password string
	db dba password
  --db_dba_use_ssl
	Set this flag to false to make the dba connection to not use ssl (default true)
  --db_dba_user string
	db dba user userKey (default vt_dba)
  --db_erepl_password string
	db erepl password
  --db_erepl_use_ssl
	Set this flag to false to make the erepl connection to not use ssl (default true)
  --db_erepl_user string
	db erepl user userKey (default vt_erepl)
  --db_filtered_password string
	db filtered password
  --db_filtered_use_ssl
	Set this flag to false to make the filtered connection to not use ssl (default true)
  --db_filtered_user string
	db filtered user userKey (default vt_filtered)
  --db_flags uint
	Flag values as defined by MySQL.
  --db_flavor string
	Flavor overrid. Valid value is FilePos.
  --db_host string
	The host name for the tcp connection.
  --db_port int
	tcp port
  --db_repl_password string
	db repl password
  --db_repl_use_ssl
	Set this flag to false to make the repl connection to not use ssl (default true)
  --db_repl_user string
	db repl user userKey (default vt_repl)
  --db_server_name string
	server name of the DB we are connecting to.
  --db_socket string
	The unix socket to connect on. If this is specified, host and port will not be used.
  --db_ssl_ca string
	connection ssl ca
  --db_ssl_ca_path string
	connection ssl ca path
  --db_ssl_cert string
	connection ssl certificate
  --db_ssl_key string
	connection ssl key
  --db_ssl_mode value
	SSL mode to connect with. One of disabled, preferred, required, verify_ca & verify_identity.
  --db_tls_min_version string
	Configures the minimal TLS version negotiated when SSL is enabled. Defaults to TLSv1.2. Options: TLSv1.0, TLSv1.1, TLSv1.2, TLSv1.3.
  --dba_idle_timeout duration
	Idle timeout for dba connections (default 1m0s)
  --dba_pool_size int
	Size of the connection pool for dba connections (default 20)
  --degraded_threshold duration
	replication lag after which a replica is considered degraded (default 30s)
  --disable_active_reparents
	if set, do not allow active reparents. Use this to protect a cluster using external reparents.
  --discovery_high_replication_lag_minimum_serving duration
	the replication lag that is considered too high when applying the min_number_serving_vttablets threshold (default 2h0m0s)
  --discovery_low_replication_lag duration
	the replication lag that is considered low enough to be healthy (default 30s)
  --emit_stats
	If set, emit stats to push-based monitoring and stats backends
  --enable-autocommit
	This flag is deprecated. Autocommit is always allowed. (default true)
  --enable-consolidator
	Synonym to -enable_consolidator (default true)
  --enable-consolidator-replicas
	Synonym to -enable_consolidator_replicas
  --enable-lag-throttler
	Synonym to -enable_lag_throttler
  --enable-query-plan-field-caching
	Synonym to -enable_query_plan_field_caching (default true)
  --enable-tx-throttler
	Synonym to -enable_tx_throttler
  --enable_consolidator
	This option enables the query consolidator. (default true)
  --enable_consolidator_replicas
	This option enables the query consolidator only on replicas.
  --enable_hot_row_protection
	If true, incoming transactions for the same row (range) will be queued and cannot consume all txpool slots.
  --enable_hot_row_protection_dry_run
	If true, hot row protection is not enforced but logs if transactions would have been queued.
  --enable_lag_throttler
	If true, vttablet will run a throttler service, and will implicitly enable heartbeats
  --enable_query_plan_field_caching
	This option fetches & caches fields (columns) when storing query plans (default true)
  --enable_replication_reporter
	Use polling to track replication lag.
  --enable_semi_sync
	(DEPRECATED - Set the correct durability_policy instead) Enable semi-sync when configuring replication, on primary and replica tablets only (rdonly tablets will not ack).
  --enable_transaction_limit
	If true, limit on number of transactions open at the same time will be enforced for all users. User trying to open a new transaction after exhausting their limit will receive an error immediately, regardless of whether there are available slots or not.
  --enable_transaction_limit_dry_run
	If true, limit on number of transactions open at the same time will be tracked for all users, but not enforced.
  --enable_tx_throttler
	If true replication-lag-based throttling on transactions will be enabled.
  --enforce-tableacl-config
	if this flag is true, vttablet will fail to start if a valid tableacl config does not exist
  --enforce_strict_trans_tables
	If true, vttablet requires MySQL to run with STRICT_TRANS_TABLES or STRICT_ALL_TABLES on. It is recommended to not turn this flag off. Otherwise MySQL may alter your supplied values before saving them to the database. (default true)
  --file_backup_storage_root string
	root directory for the file backup storage
  --filecustomrules string
	file based custom rule path
  --filecustomrules_watch
	set up a watch on the target file and reload query rules when it changes
  --gc_check_interval duration
	Interval between garbage collection checks (default 1h0m0s)
  --gc_purge_check_interval duration
	Interval between purge discovery checks (default 1m0s)
  --gcs_backup_storage_bucket string
	Google Cloud Storage bucket to use for backups
  --gcs_backup_storage_root string
	root prefix for all backup-related object names
  --gh-ost-path string
	override default gh-ost binary full path
  --grpc_auth_mode string
	Which auth plugin implementation to use (eg: static)
  --grpc_auth_mtls_allowed_substrings string
	List of substrings of at least one of the client certificate names (separated by colon).
  --grpc_auth_static_client_creds string
	when using grpc_static_auth in the server, this file provides the credentials to use to authenticate with server
  --grpc_auth_static_password_file string
	JSON File to read the users/passwords from.
  --grpc_ca string
	server CA to use for gRPC connections, requires TLS, and enforces client certificate check
  --grpc_cert string
	server certificate to use for gRPC connections, requires grpc_key, enables TLS
  --grpc_compression string
	Which protocol to use for compressing gRPC. Default: nothing. Supported: snappy
  --grpc_crl string
	path to a certificate revocation list in PEM format, client certificates will be further verified against this file during TLS handshake
  --grpc_enable_optional_tls
	enable optional TLS mode when a server accepts both TLS and plain-text connections on the same port
  --grpc_enable_tracing
	Enable GRPC tracing
  --grpc_initial_conn_window_size int
	gRPC initial connection window size
  --grpc_initial_window_size int
	gRPC initial window size
  --grpc_keepalive_time duration
	After a duration of this time, if the client doesn't see any activity, it pings the server to see if the transport is still alive. (default 10s)
  --grpc_keepalive_timeout duration
	After having pinged for keepalive check, the client waits for a duration of Timeout and if no activity is seen even after that the connection is closed. (default 10s)
  --grpc_key string
	server private key to use for gRPC connections, requires grpc_cert, enables TLS
  --grpc_max_connection_age duration
	Maximum age of a client connection before GoAway is sent. (default 2562047h47m16.854775807s)
  --grpc_max_connection_age_grace duration
	Additional grace period after grpc_max_connection_age, after which connections are forcibly closed. (default 2562047h47m16.854775807s)
  --grpc_max_message_size int
	Maximum allowed RPC message size. Larger messages will be rejected by gRPC with the error 'exceeding the max size'. (default 16777216)
  --grpc_port int
	Port to listen on for gRPC calls
  --grpc_prometheus
	Enable gRPC monitoring with Prometheus
  --grpc_server_ca string
	path to server CA in PEM format, which will be combine with server cert, return full certificate chain to clients
  --grpc_server_initial_conn_window_size int
	gRPC server initial connection window size
  --grpc_server_initial_window_size int
	gRPC server initial window size
  --grpc_server_keepalive_enforcement_policy_min_time duration
	gRPC server minimum keepalive time (default 10s)
  --grpc_server_keepalive_enforcement_policy_permit_without_stream
	gRPC server permit client keepalive pings even when there are no active streams (RPCs)
  --health_check_interval duration
	Interval between health checks (default 20s)
  --heartbeat_enable
	If true, vttablet records (if master) or checks (if replica) the current time of a replication heartbeat in the table _vt.heartbeat. The result is used to inform the serving state of the vttablet via healthchecks.
  --heartbeat_interval duration
	How frequently to read and write replication heartbeat. (default 1s)
  --hot_row_protection_concurrent_transactions int
	Number of concurrent transactions let through to the txpool/MySQL for the same hot row. Should be > 1 to have enough 'ready' transactions in MySQL and benefit from a pipelining effect. (default 5)
  --hot_row_protection_max_global_queue_size int
	Global queue limit across all row (ranges). Useful to prevent that the queue can grow unbounded. (default 1000)
  --hot_row_protection_max_queue_size int
	Maximum number of BeginExecute RPCs which will be queued for the same row (range). (default 20)
  --init_db_name_override string
	(init parameter) override the name of the db used by vttablet. Without this flag, the db name defaults to vt_<keyspacename>
  --init_keyspace string
	(init parameter) keyspace to use for this tablet
  --init_populate_metadata
	(init parameter) populate metadata tables even if restore_from_backup is disabled. If restore_from_backup is enabled, metadata tables are always populated regardless of this flag.
  --init_shard string
	(init parameter) shard to use for this tablet
  --init_tablet_type string
	(init parameter) the tablet type to use for this tablet.
  --init_tags value
	(init parameter) comma separated list of key:value pairs used to tag the tablet
  --init_timeout duration
	(init parameter) timeout to use for the init phase. (default 1m0s)
  --jaeger-agent-host string
	host and port to send spans to. if empty, no tracing will be done
  --keep_logs duration
	keep logs for this long (using ctime) (zero to keep forever)
  --keep_logs_by_mtime duration
	keep logs for this long (using mtime) (zero to keep forever)
  --keyspaces_to_watch value
	Specifies which keyspaces this vtgate should have access to while routing queries or accessing the vschema
  --lameduck-period duration
	keep running at least this long after SIGTERM before stopping (default 50ms)
  --legacy_replication_lag_algorithm
	use the legacy algorithm when selecting the vttablets for serving (default true)
  --lock_tables_timeout duration
	How long to keep the table locked before timing out (default 1m0s)
  --log_backtrace_at value
	when logging hits line file:N, emit a stack trace
  --log_dir string
	If non-empty, write log files in this directory
  --log_err_stacks
	log stack traces for errors
  --log_queries
	Enable query logging to syslog.
  --log_queries_to_file string
	Enable query logging to the specified file
  --log_rotate_max_size uint
	size in bytes at which logs are rotated (glog.MaxSize) (default 1887436800)
  --logtostderr
	log to standard error instead of files
  --master_connect_retry duration
	Deprecated, use -replication_connect_retry (default 10s)
  --mem-profile-rate int
	deprecated: use '-pprof=mem' instead (default 524288)
  --migration_check_interval duration
	Interval between migration checks (default 1m0s)
  --min_number_serving_vttablets int
	the minimum number of vttablets for each replicating tablet_type (e.g. replica, rdonly) that will be continue to be used even with replication lag above discovery_low_replication_lag, but still below discovery_high_replication_lag_minimum_serving (default 2)
  --mutex-profile-fraction int
	deprecated: use '-pprof=mutex' instead
  --mycnf-file string
	path to my.cnf, if reading all config params from there
  --mycnf_bin_log_path string
	mysql binlog path
  --mycnf_data_dir string
	data directory for mysql
  --mycnf_error_log_path string
	mysql error log path
  --mycnf_general_log_path string
	mysql general log path
  --mycnf_innodb_data_home_dir string
	Innodb data home directory
  --mycnf_innodb_log_group_home_dir string
	Innodb log group home directory
  --mycnf_master_info_file string
	mysql master.info file
  --mycnf_mysql_port int
	port mysql is listening on
  --mycnf_pid_file string
	mysql pid file
  --mycnf_relay_log_index_path string
	mysql relay log index path
  --mycnf_relay_log_info_path string
	mysql relay log info path
  --mycnf_relay_log_path string
	mysql relay log path
  --mycnf_secure_file_priv string
	mysql path for loading secure files
  --mycnf_server_id int
	mysql server id of the server (if specified, mycnf-file will be ignored)
  --mycnf_slow_log_path string
	mysql slow query log path
  --mycnf_socket_file string
	mysql socket file
  --mycnf_tmp_dir string
	mysql tmp directory
  --mysql_auth_server_static_file string
	JSON File to read the users/passwords from.
  --mysql_auth_server_static_string string
	JSON representation of the users/passwords config.
  --mysql_auth_static_reload_interval duration
	Ticker to reload credentials
  --mysql_clientcert_auth_method string
	client-side authentication method to use. Supported values: mysql_clear_password, dialog. (default mysql_clear_password)
  --mysql_server_flush_delay duration
	Delay after which buffered response will be flushed to the client. (default 100ms)
  --mysql_server_version string
	MySQL server version to advertise.
  --mysqlctl_client_protocol string
	the protocol to use to talk to the mysqlctl server (default grpc)
  --mysqlctl_mycnf_template string
	template file to use for generating the my.cnf file during server init
  --mysqlctl_socket string
	socket file to use for remote mysqlctl actions (empty for local actions)
  --onclose_timeout duration
	wait no more than this for OnClose handlers before stopping (default 1ns)
  --onterm_timeout duration
	wait no more than this for OnTermSync handlers before stopping (default 10s)
  --opentsdb_uri string
	URI of opentsdb /api/put method
  --orc_api_password string
	(Optional) Basic auth password to authenticate with Orchestrator's HTTP API.
  --orc_api_url string
	Address of Orchestrator's HTTP API (e.g. http://host:port/api/). Leave empty to disable Orchestrator integration.
  --orc_api_user string
	(Optional) Basic auth username to authenticate with Orchestrator's HTTP API. Leave empty to disable basic auth.
  --orc_discover_interval duration
	How often to ping Orchestrator's HTTP API endpoint to tell it we exist. 0 means never.
  --orc_timeout duration
	Timeout for calls to Orchestrator's HTTP API (default 30s)
  --pid_file string
	If set, the process will write its pid to the named file, and delete it on graceful shutdown.
  --pitr_gtid_lookup_timeout duration
	PITR restore parameter: timeout for fetching gtid from timestamp. (default 1m0s)
  --pool-name-prefix string
	Deprecated
  --pool_hostname_resolve_interval duration
	if set force an update to all hostnames and reconnect if changed, defaults to 0 (disabled)
  --port int
	port for the server
  --pprof string
	enable profiling
  --pt-osc-path string
	override default pt-online-schema-change binary full path
  --publish_retry_interval duration
	how long vttablet waits to retry publishing the tablet record (default 30s)
  --purge_logs_interval duration
	how often try to remove old logs (default 1h0m0s)
  --query-log-stream-handler string
	URL handler for streaming queries log (default /debug/querylog)
  --querylog-filter-tag string
	string that must be present in the query for it to be logged; if using a value as the tag, you need to disable query normalization
  --querylog-format string
	format for query logs ("text" or "json") (default text)
  --querylog-row-threshold uint
	Number of rows a query has to return or affect before being logged; not useful for streaming queries. 0 means all queries will be logged.
  --queryserver-config-acl-exempt-acl string
	an acl that exempt from table acl checking (this acl is free to access any vitess tables).
  --queryserver-config-allowunsafe-dmls
	deprecated
  --queryserver-config-annotate-queries
	prefix queries to MySQL backend with comment indicating vtgate principal (user) and target tablet type
  --queryserver-config-enable-table-acl-dry-run
	If this flag is enabled, tabletserver will emit monitoring metrics and let the request pass regardless of table acl check results
  --queryserver-config-idle-timeout float
	query server idle timeout (in seconds), vttablet manages various mysql connection pools. This config means if a connection has not been used in given idle timeout, this connection will be removed from pool. This effectively manages number of connection objects and optimize the pool performance. (default 1800)
  --queryserver-config-max-dml-rows int
	query server max dml rows per statement, maximum number of rows allowed to return at a time for an update or delete with either 1) an equality where clauses on primary keys, or 2) a subselect statement. For update and delete statements in above two categories, vttablet will split the original query into multiple small queries based on this configuration value. 
  --queryserver-config-max-result-size int
	query server max result size, maximum number of rows allowed to return from vttablet for non-streaming queries. (default 10000)
  --queryserver-config-message-conn-pool-prefill-parallelism int
	DEPRECATED: Unused.
  --queryserver-config-message-conn-pool-size int
	DEPRECATED
  --queryserver-config-message-postpone-cap int
	query server message postpone cap is the maximum number of messages that can be postponed at any given time. Set this number to substantially lower than transaction cap, so that the transaction pool isn't exhausted by the message subsystem. (default 4)
  --queryserver-config-passthrough-dmls
	query server pass through all dml statements without rewriting
  --queryserver-config-pool-prefill-parallelism int
	query server read pool prefill parallelism, a non-zero value will prefill the pool using the specified parallism.
  --queryserver-config-pool-size int
	query server read pool size, connection pool is used by regular queries (non streaming, not in a transaction) (default 16)
  --queryserver-config-query-cache-lfu
	query server cache algorithm. when set to true, a new cache algorithm based on a TinyLFU admission policy will be used to improve cache behavior and prevent pollution from sparse queries (default true)
  --queryserver-config-query-cache-memory int
	query server query cache size in bytes, maximum amount of memory to be used for caching. vttablet analyzes every incoming query and generate a query plan, these plans are being cached in a lru cache. This config controls the capacity of the lru cache. (default 33554432)
  --queryserver-config-query-cache-size int
	query server query cache size, maximum number of queries to be cached. vttablet analyzes every incoming query and generate a query plan, these plans are being cached in a lru cache. This config controls the capacity of the lru cache. (default 5000)
  --queryserver-config-query-pool-timeout float
	query server query pool timeout (in seconds), it is how long vttablet waits for a connection from the query pool. If set to 0 (default) then the overall query timeout is used instead.
  --queryserver-config-query-pool-waiter-cap int
	query server query pool waiter limit, this is the maximum number of queries that can be queued waiting to get a connection (default 5000)
  --queryserver-config-query-timeout float
	query server query timeout (in seconds), this is the query timeout in vttablet side. If a query takes more than this timeout, it will be killed. (default 30)
  --queryserver-config-schema-change-signal
	query server schema signal, will signal connected vtgates that schema has changed whenever this is detected. VTGates will need to have -schema_change_signal enabled for this to work
  --queryserver-config-schema-change-signal-interval float
	query server schema change signal interval defines at which interval the query server shall send schema updates to vtgate. (default 5)
  --queryserver-config-schema-reload-time float
	query server schema reload time, how often vttablet reloads schemas from underlying MySQL instance in seconds. vttablet keeps table schemas in its own memory and periodically refreshes it from MySQL. This config controls the reload time. (default 1800)
  --queryserver-config-stream-buffer-size int
	query server stream buffer size, the maximum number of bytes sent from vttablet for each stream call. It's recommended to keep this value in sync with vtgate's stream_buffer_size. (default 32768)
  --queryserver-config-stream-pool-prefill-parallelism int
	query server stream pool prefill parallelism, a non-zero value will prefill the pool using the specified parallelism
  --queryserver-config-stream-pool-size int
	query server stream connection pool size, stream pool is used by stream queries: queries that return results to client in a streaming fashion (default 200)
  --queryserver-config-stream-pool-timeout float
	query server stream pool timeout (in seconds), it is how long vttablet waits for a connection from the stream pool. If set to 0 (default) then there is no timeout.
  --queryserver-config-stream-pool-waiter-cap int
	query server stream pool waiter limit, this is the maximum number of streaming queries that can be queued waiting to get a connection
  --queryserver-config-strict-table-acl
	only allow queries that pass table acl checks
  --queryserver-config-terse-errors
	prevent bind vars from escaping in client error messages
  --queryserver-config-transaction-cap int
	query server transaction cap is the maximum number of transactions allowed to happen at any given point of a time for a single vttablet. E.g. by setting transaction cap to 100, there are at most 100 transactions will be processed by a vttablet and the 101th transaction will be blocked (and fail if it cannot get connection within specified timeout) (default 20)
  --queryserver-config-transaction-prefill-parallelism int
	query server transaction prefill parallelism, a non-zero value will prefill the pool using the specified parallism.
  --queryserver-config-transaction-timeout float
	query server transaction timeout (in seconds), a transaction will be killed if it takes longer than this value (default 30)
  --queryserver-config-txpool-timeout float
	query server transaction pool timeout, it is how long vttablet waits if tx pool is full (default 1)
  --queryserver-config-txpool-waiter-cap int
	query server transaction pool waiter limit, this is the maximum number of transactions that can be queued waiting to get a connection (default 5000)
  --queryserver-config-warn-result-size int
	query server result size warning threshold, warn if number of rows returned from vttablet for non-streaming queries exceeds this
  --queryserver_enable_online_ddl
	Enable online DDL. (default true)
  --redact-debug-ui-queries
	redact full queries and bind variables from debug UI
  --relay_log_max_items int
	Maximum number of rows for VReplication target buffering. (default 5000)
  --relay_log_max_size int
	Maximum buffer size (in bytes) for VReplication target buffering. If single rows are larger than this, a single row is buffered at a time. (default 250000)
  --remote_operation_timeout duration
	time to wait for a remote operation (default 30s)
  --replication_connect_retry duration
	how long to wait in between replica reconnect attempts. Only precise to the second. (default 10s)
  --restore_concurrency int
	(init restore parameter) how many concurrent files to restore at once (default 4)
  --restore_from_backup
	(init restore parameter) will check BackupStorage for a recent backup at startup and start there
  --restore_from_backup_ts string
	(init restore parameter) if set, restore the latest backup taken at or before this timestamp. Example: '2021-04-29.133050'
  --retain_online_ddl_tables duration
	How long should vttablet keep an old migrated table before purging it (default 24h0m0s)
  --s3_backup_aws_endpoint string
	endpoint of the S3 backend (region must be provided)
  --s3_backup_aws_region string
	AWS region to use (default us-east-1)
  --s3_backup_aws_retries int
	AWS request retries (default -1)
  --s3_backup_force_path_style
	force the s3 path style
  --s3_backup_log_level string
	determine the S3 loglevel to use from LogOff, LogDebug, LogDebugWithSigning, LogDebugWithHTTPBody, LogDebugWithRequestRetries, LogDebugWithRequestErrors (default LogOff)
  --s3_backup_server_side_encryption string
	server-side encryption algorithm (e.g., AES256, aws:kms, sse_c:/path/to/key/file)
  --s3_backup_storage_bucket string
	S3 bucket to use for backups
  --s3_backup_storage_root string
	root prefix for all backup-related object names
  --s3_backup_tls_skip_verify_cert
	skip the 'certificate is valid' check for SSL connections
  --sanitize_log_messages
	Remove potentially sensitive information in tablet INFO, WARNING, and ERROR log messages such as query parameters.
  --security_policy string
	the name of a registered security policy to use for controlling access to URLs - empty means allow all for anyone (built-in policies: deny-all, read-only)
  --service_map value
	comma separated list of services to enable (or disable if prefixed with '-') Example: grpc-vtworker
  --serving_state_grace_period duration
	how long to pause after broadcasting health to vtgate, before enforcing a new serving state
  --shard_sync_retry_delay duration
	delay between retries of updates to keep the tablet and its shard record in sync (default 30s)
  --shutdown_grace_period float
	how long to wait (in seconds) for queries and transactions to complete during graceful shutdown.
  --sql-max-length-errors int
	truncate queries in error logs to the given length (default unlimited)
  --sql-max-length-ui int
	truncate queries in debug UIs to the given length (default 512) (default 512)
  --srv_topo_cache_refresh duration
	how frequently to refresh the topology for cached entries (default 1s)
  --srv_topo_cache_ttl duration
	how long to use cached entries for topology (default 1s)
  --srv_topo_timeout duration
	topo server timeout (default 5s)
  --stats_backend string
	The name of the registered push-based monitoring/stats backend to use
  --stats_combine_dimensions string
	List of dimensions to be combined into a single "all" value in exported stats vars
  --stats_common_tags string
	Comma-separated list of common tags for the stats backend. It provides both label and values. Example: label1:value1,label2:value2
  --stats_drop_variables string
	Variables to be dropped from the list of exported variables.
  --stats_emit_period duration
	Interval between emitting stats to all registered backends (default 1m0s)
  --statsd_address string
	Address for statsd client
  --statsd_sample_rate float
	 (default 1)
  --stderrthreshold value
	logs at or above this threshold go to stderr (default 1)
  --stream_health_buffer_size uint
	max streaming health entries to buffer per streaming health client (default 20)
  --table-acl-config string
	path to table access checker config file; send SIGHUP to reload this file
  --table-acl-config-reload-interval duration
	Ticker to reload ACLs. Duration flag, format e.g.: 30s. Default: do not reload
  --table_gc_lifecycle string
	States for a DROP TABLE garbage collection cycle. Default is 'hold,purge,evac,drop', use any subset ('drop' implcitly always included) (default hold,purge,evac,drop)
  --tablet-path string
	tablet alias
  --tablet_config string
	YAML file config for tablet
  --tablet_dir string
	The directory within the vtdataroot to store vttablet/mysql files. Defaults to being generated by the tablet uid.
  --tablet_filters value
	Specifies a comma-separated list of 'keyspace|shard_name or keyrange' values to filter the tablets to watch
  --tablet_grpc_ca string
	the server ca to use to validate servers when connecting
  --tablet_grpc_cert string
	the cert to use to connect
  --tablet_grpc_crl string
	the server crl to use to validate server certificates when connecting
  --tablet_grpc_key string
	the key to use to connect
  --tablet_grpc_server_name string
	the server name to use to validate server certificate
  --tablet_hostname string
	if not empty, this hostname will be assumed instead of trying to resolve it
  --tablet_manager_grpc_ca string
	the server ca to use to validate servers when connecting
  --tablet_manager_grpc_cert string
	the cert to use to connect
  --tablet_manager_grpc_concurrency int
	concurrency to use to talk to a vttablet server for performance-sensitive RPCs (like ExecuteFetchAs{Dba,AllPrivs,App}) (default 8)
  --tablet_manager_grpc_connpool_size int
	number of tablets to keep tmclient connections open to (default 100)
  --tablet_manager_grpc_crl string
	the server crl to use to validate server certificates when connecting
  --tablet_manager_grpc_key string
	the key to use to connect
  --tablet_manager_grpc_server_name string
	the server name to use to validate server certificate
  --tablet_manager_protocol string
	the protocol to use to talk to vttablet (default grpc)
  --tablet_protocol string
	how to talk to the vttablets (default grpc)
  --tablet_refresh_interval duration
	tablet refresh interval (default 1m0s)
  --tablet_refresh_known_tablets
	tablet refresh reloads the tablet address/port map from topo in case it changes (default true)
  --tablet_url_template string
	format string describing debug tablet url formatting. See the Go code for getTabletDebugURL() how to customize this. (default http://{{.GetTabletHostPort}})
  --throttle_check_as_check_self
	Should throttler/check return a throttler/check-self result (changes throttler behavior for writes)
  --throttle_metrics_query SELECT
	Override default heartbeat/lag metric. Use either SELECT (must return single row, single value) or ` + "`SHOW GLOBAL ... LIKE ...`" + ` queries. Set -throttle_metrics_threshold respectively.
  --throttle_metrics_threshold float
	Override default throttle threshold, respective to -throttle_metrics_query (default 1.7976931348623157e+308)
  --throttle_tablet_types string
	Comma separated VTTablet types to be considered by the throttler. default: 'replica'. example: 'replica,rdonly'. 'replica' aways implicitly included (default replica)
  --throttle_threshold duration
	Replication lag threshold for default lag throttling (default 1s)
  --topo_consul_lock_delay duration
	LockDelay for consul session. (default 15s)
  --topo_consul_lock_session_checks string
	List of checks for consul session. (default serfHealth)
  --topo_consul_lock_session_ttl string
	TTL for consul session.
  --topo_consul_watch_poll_duration duration
	time of the long poll for watch queries. (default 30s)
  --topo_etcd_lease_ttl int
	Lease TTL for locks and leader election. The client will use KeepAlive to keep the lease going. (default 30)
  --topo_etcd_tls_ca string
	path to the ca to use to validate the server cert when connecting to the etcd topo server
  --topo_etcd_tls_cert string
	path to the client cert to use to connect to the etcd topo server, requires topo_etcd_tls_key, enables TLS
  --topo_etcd_tls_key string
	path to the client key to use to connect to the etcd topo server, enables TLS
  --topo_global_root string
	the path of the global topology data in the global topology server
  --topo_global_server_address string
	the address of the global topology server
  --topo_implementation string
	the topology implementation to use
  --topo_k8s_context string
	The kubeconfig context to use, overrides the 'current-context' from the config
  --topo_k8s_kubeconfig string
	Path to a valid kubeconfig file. When running as a k8s pod inside the same cluster you wish to use as the topo, you may omit this and the below arguments, and Vitess is capable of auto-discovering the correct values. https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/#accessing-the-api-from-a-pod
  --topo_k8s_namespace string
	The kubernetes namespace to use for all objects. Default comes from the context or in-cluster config
  --topo_read_concurrency int
	concurrent topo reads (default 32)
  --topo_zk_auth_file string
	auth to use when connecting to the zk topo server, file contents should be <scheme>:<auth>, e.g., digest:user:pass
  --topo_zk_base_timeout duration
	zk base timeout (see zk.Connect) (default 30s)
  --topo_zk_max_concurrency int
	maximum number of pending requests to send to a Zookeeper server. (default 64)
  --topo_zk_tls_ca string
	the server ca to use to validate servers when connecting to the zk topo server
  --topo_zk_tls_cert string
	the cert to use to connect to the zk topo server, requires topo_zk_tls_key, enables TLS
  --topo_zk_tls_key string
	the key to use to connect to the zk topo server, enables TLS
  --topocustomrule_cell string
	topo cell for customrules file. (default global)
  --topocustomrule_path string
	path for customrules file. Disabled if empty.
  --tracer string
	tracing service to use (default noop)
  --tracing-enable-logging
	whether to enable logging in the tracing service
  --tracing-sampling-rate value
	sampling rate for the probabilistic jaeger sampler (default 0.1)
  --tracing-sampling-type value
	sampling strategy to use for jaeger. possible values are 'const', 'probabilistic', 'rateLimiting', or 'remote' (default const)
  --track_schema_versions
	When enabled, vttablet will store versions of schemas at each position that a DDL is applied and allow retrieval of the schema corresponding to a position
  --transaction-log-stream-handler string
	URL handler for streaming transactions log (default /debug/txlog)
  --transaction_limit_by_component
	Include CallerID.component when considering who the user is for the purpose of transaction limit.
  --transaction_limit_by_principal
	Include CallerID.principal when considering who the user is for the purpose of transaction limit. (default true)
  --transaction_limit_by_subcomponent
	Include CallerID.subcomponent when considering who the user is for the purpose of transaction limit.
  --transaction_limit_by_username
	Include VTGateCallerID.username when considering who the user is for the purpose of transaction limit. (default true)
  --transaction_limit_per_user float
	Maximum number of transactions a single user is allowed to use at any time, represented as fraction of -transaction_cap. (default 0.4)
  --transaction_shutdown_grace_period float
	DEPRECATED: use shutdown_grace_period instead.
  --twopc_abandon_age float
	time in seconds. Any unresolved transaction older than this time will be sent to the coordinator to be resolved.
  --twopc_coordinator_address string
	address of the (VTGate) process(es) that will be used to notify of abandoned transactions.
  --twopc_enable
	if the flag is on, 2pc is enabled. Other 2pc flags must be supplied.
  --tx-throttler-config string
	Synonym to -tx_throttler_config (default target_replication_lag_sec: 2
max_replication_lag_sec: 10
initial_rate: 100
max_increase: 1
emergency_decrease: 0.5
min_duration_between_increases_sec: 40
max_duration_between_increases_sec: 62
min_duration_between_decreases_sec: 20
spread_backlog_across_sec: 20
age_bad_rate_after_sec: 180
bad_rate_increase: 0.1
max_rate_approach_threshold: 0.9
)
  --tx-throttler-healthcheck-cells value
	Synonym to -tx_throttler_healthcheck_cells
  --tx_throttler_config string
	The configuration of the transaction throttler as a text formatted throttlerdata.Configuration protocol buffer message (default target_replication_lag_sec: 2
max_replication_lag_sec: 10
initial_rate: 100
max_increase: 1
emergency_decrease: 0.5
min_duration_between_increases_sec: 40
max_duration_between_increases_sec: 62
min_duration_between_decreases_sec: 20
spread_backlog_across_sec: 20
age_bad_rate_after_sec: 180
bad_rate_increase: 0.1
max_rate_approach_threshold: 0.9
)
  --tx_throttler_healthcheck_cells value
	A comma-separated list of cells. Only tabletservers running in these cells will be monitored for replication lag by the transaction throttler.
  --unhealthy_threshold duration
	replication lag after which a replica is considered unhealthy (default 2h0m0s)
  --use_super_read_only
	Set super_read_only flag when performing planned failover. (default true)
  --v value
	log level for V logs
  --version
	print binary version
  --vmodule value
	comma-separated list of pattern=N settings for file-filtered logging
  --vreplication_copy_phase_duration duration
	Duration for each copy phase loop (before running the next catchup: default 1h) (default 1h0m0s)
  --vreplication_copy_phase_max_innodb_history_list_length int
	The maximum InnoDB transaction history that can exist on a vstreamer (source) before starting another round of copying rows. This helps to limit the impact on the source tablet. (default 1000000)
  --vreplication_copy_phase_max_mysql_replication_lag int
	The maximum MySQL replication lag (in seconds) that can exist on a vstreamer (source) before starting another round of copying rows. This helps to limit the impact on the source tablet. (default 43200)
  --vreplication_experimental_flags int
	(Bitmask) of experimental features in vreplication to enable (default 1)
  --vreplication_healthcheck_retry_delay duration
	healthcheck retry delay (default 5s)
  --vreplication_healthcheck_timeout duration
	healthcheck retry delay (default 1m0s)
  --vreplication_healthcheck_topology_refresh duration
	refresh interval for re-reading the topology (default 30s)
  --vreplication_heartbeat_update_interval int
	Frequency (in seconds, default 1, max 60) at which the time_updated column of a vreplication stream when idling (default 1)
  --vreplication_replica_lag_tolerance duration
	Replica lag threshold duration: once lag is below this we switch from copy phase to the replication (streaming) phase (default 1m0s)
  --vreplication_retry_delay duration
	delay before retrying a failed binlog connection (default 5s)
  --vreplication_store_compressed_gtid
	Store compressed gtids in the pos column of _vt.vreplication
  --vreplication_tablet_type string
	comma separated list of tablet types used as a source (default PRIMARY,REPLICA)
  --vstream_dynamic_packet_size
	Enable dynamic packet sizing for VReplication. This will adjust the packet size during replication to improve performance. (default true)
  --vstream_packet_size int
	Suggested packet size for VReplication streamer. This is used only as a recommendation. The actual packet size may be more or less than this amount. (default 250000)
  --vtctld_addr string
	address of a vtctld instance
  --vtgate_protocol string
	how to talk to vtgate (default grpc)
  --vttablet_skip_buildinfo_tags string
	comma-separated list of buildinfo tags to skip from merging with -init_tags. each tag is either an exact match or a regular expression of the form '/regexp/'. (default /.*/)
  --wait_for_backup_interval duration
	(init restore parameter) if this is greater than 0, instead of starting up empty when no backups are found, keep checking at this interval for a backup to appear
  --watch_replication_stream
	When enabled, vttablet will stream the MySQL replication stream from the local server, and use it to update schema when it sees a DDL.
  --xbstream_restore_flags string
	flags to pass to xbstream command during restore. These should be space separated and will be added to the end of the command. These need to match the ones used for backup e.g. --compress / --decompress, --encrypt / --decrypt
  --xtrabackup_backup_flags string
	flags to pass to backup command. These should be space separated and will be added to the end of the command
  --xtrabackup_prepare_flags string
	flags to pass to prepare command. These should be space separated and will be added to the end of the command
  --xtrabackup_root_path string
	directory location of the xtrabackup and xbstream executables, e.g., /usr/bin
  --xtrabackup_stream_mode string
	which mode to use if streaming, valid values are tar and xbstream (default tar)
  --xtrabackup_stripe_block_size uint
	Size in bytes of each block that gets sent to a given stripe before rotating to the next stripe (default 102400)
  --xtrabackup_stripes uint
	If greater than 0, use data striping across this many destination files to parallelize data transfer and decompression
  --xtrabackup_user string
	User that xtrabackup will use to connect to the database server. This user must have all necessary privileges. For details, please refer to xtrabackup documentation.
`,
	}
)

func TestHelpOutput(t *testing.T) {
	args := []string{"--help"}
	for binary, helptext := range helpOutput {
		cmd := exec.Command(binary, args...)
		output := bytes.Buffer{}
		cmd.Stderr = &output
		err := cmd.Run()
		require.NoError(t, err)
		assert.Equal(t, helptext, output.String(), fmt.Sprintf("%s does not have the expected help output. Please update the test if you intended to change it.", binary))
	}
}
