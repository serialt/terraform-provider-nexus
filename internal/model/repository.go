package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RepositoryAptProxyModel struct {
	Id      types.String    `tfsdk:"id"`
	Name    types.String    `tfsdk:"name"`
	Online  types.Bool      `tfsdk:"online"`
	Flat    types.Bool      `tfsdk:"flat"`
	Cleanup []*CleanupModel `tfsdk:"cleanup"`

	Storage       *StorageModel       `tfsdk:"storage"`
	Distribution  types.String        `tfsdk:"distribution"`
	RoutingRule   types.String        `tfsdk:"routing_rule"`
	HttpClient    *HttpClientModel    `tfsdk:"http_client"`
	NegativeCache *NegativeCacheModel `tfsdk:"negative_cache"`
	Proxy         *ProxyModel         `tfsdk:"proxy"`
}

type StorageModel struct {
	BlobStoreName               types.String `tfsdk:"blob_store_name"`
	StrictContentTypeValidation types.Bool   `tfsdk:"strict_content_type_validation"`
}
type NegativeCacheModel struct {
	Enabled types.Bool  `tfsdk:"enabled"`
	TTL     types.Int64 `tfsdk:"ttl"`
}
type HttpClientModel struct {
	Authentication *HttpClientAuthenticationModel `tfsdk:"authentication"`
	AutoBlock      types.Bool                     `tfsdk:"auto_block"`
	Blocked        types.Bool                     `tfsdk:"blocked"`
	Connection     *HttpClientConnectionModel     `tfsdk:"connection"`
}

type ProxyModel struct {
	ContentMaxAge  types.Int64  `tfsdk:"content_max_age"`
	MetadataMaxAge types.Int64  `tfsdk:"metadata_max_age"`
	RemoteURL      types.String `tfsdk:"remote_url"`
}
type HttpClientAuthenticationModel struct {
	NtlmDomain types.String `tfsdk:"ntlm_domain"`
	NtlmHost   types.String `tfsdk:"ntlm_host"`
	Password   types.String `tfsdk:"password"`
	Type       types.String `tfsdk:"type"`
	Username   types.String `tfsdk:"username"`
}
type HttpClientConnectionModel struct {
	EnableCircularRedirects types.Bool   `tfsdk:"enable_circular_redirects"`
	EnableCookies           types.Bool   `tfsdk:"enable_cookies"`
	Retries                 types.Int64  `tfsdk:"retries"`
	Timeout                 types.Int64  `tfsdk:"timeout"`
	UseTrustStore           types.Bool   `tfsdk:"use_trust_store"`
	UserAgentSuffix         types.String `tfsdk:"user_agent_suffix"`
}

type CleanupModel struct {
	PolicyNames []types.String `tfsdk:"policy_names"`
}

type RepositoryAptHostedModel struct {
	Id           types.String      `tfsdk:"id"`
	Name         types.String      `tfsdk:"name"`
	Online       types.Bool        `tfsdk:"online"`
	Cleanup      []*CleanupModel   `tfsdk:"cleanup"`
	Component    []*ComponentModel `tfsdk:"component"`
	Storage      *StorageModelV2   `tfsdk:"storage"`
	Distribution types.String      `tfsdk:"distribution"`
	Signing      *SigningModel     `tfsdk:"signing"`
}

type StorageModelV2 struct {
	BlobStoreName               types.String `tfsdk:"blob_store_name"`
	StrictContentTypeValidation types.Bool   `tfsdk:"strict_content_type_validation"`
	WritePolicy                 types.String `tfsdk:"write_policy"`
}

type SigningModel struct {
	Keypair    []types.String `tfsdk:"keypair"`
	Passphrase []types.String `tfsdk:"passphrase"`
}

type RepositoryDockerHostedModel struct {
	Id        types.String    `tfsdk:"id"`
	Name      types.String    `tfsdk:"name"`
	Online    types.Bool      `tfsdk:"online"`
	Cleanup   CleanupModel    `tfsdk:"cleanup"`
	Component ComponentModel  `tfsdk:"component"`
	Storage   *StorageModelV2 `tfsdk:"storage"`
}

type DockerModel struct {
	ForceBasicAuth types.Bool   `tfsdk:"force_basic_auth"`
	HttpPort       types.Int64  `tfsdk:"http_port"`
	HttpsPort      types.Int64  `tfsdk:"https_port"`
	V1Enabled      types.Bool   `tfsdk:"v1_enabled"`
	Subdomain      types.String `tfsdk:"subdomain"`
}

type RepositoryDockerProxyModel struct {
	Id            types.String     `tfsdk:"id"`
	Name          types.String     `tfsdk:"name"`
	Online        types.Bool       `tfsdk:"online"`
	Cleanup       CleanupModel     `tfsdk:"cleanup"`
	HttpClient    *HttpClientModel `tfsdk:"http_client"`
	NegativeCache *NegativeCache   `tfsdk:"negative_cache"`
	Proxy         *ProxyModel      `tfsdk:"proxy"`
	Component     *ComponentModel  `tfsdk:"component"`
	RoutingRule   types.String     `tfsdk:"routing_rule"`
	Storage       *StorageModel    `tfsdk:"storage"`
}

type NegativeCache struct {
	Enabled types.String `tfsdk:"enabled"`
	Ttl     types.Int64  `tfsdk:"ttl"`
}

type RepositoryDockerGroupModel struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Online types.Bool   `tfsdk:"online"`

	Cleanup       *CleanupModel    `tfsdk:"cleanup"`
	HttpClient    *HttpClientModel `tfsdk:"http_client"`
	NegativeCache *NegativeCache   `tfsdk:"negative_cache"`
	Proxy         *ProxyModel      `tfsdk:"proxy"`
	Component     *ComponentModel  `tfsdk:"component"`
	RoutingRule   types.String     `tfsdk:"routing_rule"`
	Storage       *StorageModel    `tfsdk:"storage"`
}
