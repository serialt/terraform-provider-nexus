package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RepositoryHostBaseModel struct {
	Id        types.String    `tfsdk:"id"`
	Name      types.String    `tfsdk:"name"`
	Online    types.Bool      `tfsdk:"online"`
	Cleanup   *CleanupModel   `tfsdk:"cleanup"`
	Component *ComponentModel `tfsdk:"component"`
	Storage   *StorageModel   `tfsdk:"storage"`
}
type RepositoryGroupBaseModel struct {
	Id      types.String  `tfsdk:"id"`
	Name    types.String  `tfsdk:"name"`
	Online  types.Bool    `tfsdk:"online"`
	Group   *GroupModel   `tfsdk:"group"`
	Storage *StorageModel `tfsdk:"storage"`
}

type RepositoryProxyBaseModel struct {
	Id                types.String        `tfsdk:"id"`
	Name              types.String        `tfsdk:"name"`
	Online            types.Bool          `tfsdk:"online"`
	Cleanup           *CleanupModel       `tfsdk:"cleanup"`
	HttpClient        *HttpClientModel    `tfsdk:"http_client"`
	NegativeCache     *NegativeCacheModel `tfsdk:"negative_cache"`
	Proxy             *ProxyModel         `tfsdk:"proxy"`
	RoutingRule       types.String        `tfsdk:"routing_rule"`
	Storage           *StorageModel       `tfsdk:"storage"`
	RewritePackagUrls types.Bool          `tfsdk:"rewrite_package_urls"`
}

type RepositoryAptProxyModel struct {
	Id      types.String  `tfsdk:"id"`
	Name    types.String  `tfsdk:"name"`
	Online  types.Bool    `tfsdk:"online"`
	Flat    types.Bool    `tfsdk:"flat"`
	Cleanup *CleanupModel `tfsdk:"cleanup"`

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
type GroupModel struct {
	MemberNames []types.String `tfsdk:"member_names"`
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
	PolicyNames types.List `tfsdk:"policy_names"`
}

type RepositoryAptHostedModel struct {
	Id           types.String        `tfsdk:"id"`
	Name         types.String        `tfsdk:"name"`
	Online       types.Bool          `tfsdk:"online"`
	Cleanup      *CleanupModel       `tfsdk:"cleanup"`
	Component    []*ComponentModel   `tfsdk:"component"`
	Storage      *HostedStorageModel `tfsdk:"storage"`
	Distribution types.String        `tfsdk:"distribution"`
	Signing      *SigningModel       `tfsdk:"signing"`
}

type HostedStorageModel struct {
	BlobStoreName               types.String `tfsdk:"blob_store_name"`
	StrictContentTypeValidation types.Bool   `tfsdk:"strict_content_type_validation"`
	WritePolicy                 types.String `tfsdk:"write_policy"`
}

type SigningModel struct {
	Keypair    []types.String `tfsdk:"keypair"`
	Passphrase []types.String `tfsdk:"passphrase"`
}

type RepositoryDockerHostedModel struct {
	Id        types.String              `tfsdk:"id"`
	Name      types.String              `tfsdk:"name"`
	Online    types.Bool                `tfsdk:"online"`
	Cleanup   *CleanupModel             `tfsdk:"cleanup"`
	Component *ComponentModel           `tfsdk:"component"`
	Storage   *DockerHostedStorageModel `tfsdk:"storage"`
	Docker    *DockerModel              `tfsdk:"docker"`
}
type DockerHostedStorageModel struct {
	BlobStoreName               types.String `tfsdk:"blob_store_name"`
	StrictContentTypeValidation types.Bool   `tfsdk:"strict_content_type_validation"`
	WritePolicy                 types.String `tfsdk:"write_policy"`
	LatestPolicy                types.Bool   `tfsdk:"latest_policy"`
}

type DockerModel struct {
	ForceBasicAuth types.Bool   `tfsdk:"force_basic_auth"`
	HttpPort       types.Int64  `tfsdk:"http_port"`
	HttpsPort      types.Int64  `tfsdk:"https_port"`
	V1Enabled      types.Bool   `tfsdk:"v1_enabled"`
	Subdomain      types.String `tfsdk:"subdomain"`
}

type RepositoryDockerProxyModel struct {
	Id            types.String        `tfsdk:"id"`
	Name          types.String        `tfsdk:"name"`
	Online        types.Bool          `tfsdk:"online"`
	Cleanup       *CleanupModel       `tfsdk:"cleanup"`
	HttpClient    *HttpClientModel    `tfsdk:"http_client"`
	NegativeCache *NegativeCacheModel `tfsdk:"negative_cache"`
	Proxy         *ProxyModel         `tfsdk:"proxy"`
	Component     *ComponentModel     `tfsdk:"component"`
	RoutingRule   types.String        `tfsdk:"routing_rule"`
	Storage       *StorageModel       `tfsdk:"storage"`
}

type RepositoryDockerGroupModel struct {
	Id     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Online types.Bool   `tfsdk:"online"`

	Group   *GroupModel   `tfsdk:"group"`
	Storage *StorageModel `tfsdk:"storage"`
	Docker  *DockerModel  `tfsdk:"docker"`
}

type RepositoryGoProxyModel struct {
	Id            types.String        `tfsdk:"id"`
	Name          types.String        `tfsdk:"name"`
	Online        types.Bool          `tfsdk:"online"`
	Cleanup       *CleanupModel       `tfsdk:"cleanup"`
	HttpClient    *HttpClientModel    `tfsdk:"http_client"`
	NegativeCache *NegativeCacheModel `tfsdk:"negative_cache"`
	Proxy         *ProxyModel         `tfsdk:"proxy"`
	RoutingRule   types.String        `tfsdk:"routing_rule"`
	Storage       *StorageModel       `tfsdk:"storage"`
}

type RepositoryGoGroupModel struct {
	Id      types.String  `tfsdk:"id"`
	Name    types.String  `tfsdk:"name"`
	Online  types.Bool    `tfsdk:"online"`
	Storage *StorageModel `tfsdk:"storage"`
	Group   *GroupModel   `tfsdk:"group"`
}

type RepositoryHelmProxyModel struct {
	Id            types.String        `tfsdk:"id"`
	Name          types.String        `tfsdk:"name"`
	Online        types.Bool          `tfsdk:"online"`
	Cleanup       *CleanupModel       `tfsdk:"cleanup"`
	HttpClient    *HttpClientModel    `tfsdk:"http_client"`
	NegativeCache *NegativeCacheModel `tfsdk:"negative_cache"`
	Proxy         *ProxyModel         `tfsdk:"proxy"`
	RoutingRule   types.String        `tfsdk:"routing_rule"`
	Storage       *StorageModel       `tfsdk:"storage"`
}
type RepositoryHelmHostModel struct {
	Id        types.String    `tfsdk:"id"`
	Name      types.String    `tfsdk:"name"`
	Online    types.Bool      `tfsdk:"online"`
	Cleanup   *CleanupModel   `tfsdk:"cleanup"`
	Component *ComponentModel `tfsdk:"component"`
	Storage   *StorageModel   `tfsdk:"storage"`
}

type RepositoryListItemModel struct {
	Name   types.String `tfsdk:"name"`
	Format types.String `tfsdk:"format"`
	Type   types.String `tfsdk:"type"`
	Url    types.String `tfsdk:"url"`
}
type RepositoryListModel struct {
	Id    types.String               `tfsdk:"id"`
	Items []*RepositoryListItemModel `tfsdk:"items"`
}

type RepositoryBowerGroupModel struct {
	Id      types.String  `tfsdk:"id"`
	Name    types.String  `tfsdk:"name"`
	Online  types.Bool    `tfsdk:"online"`
	Group   *GroupModel   `tfsdk:"group"`
	Storage *StorageModel `tfsdk:"storage"`
}

type RepositoryBowerHostedModel struct {
	Id        types.String    `tfsdk:"id"`
	Name      types.String    `tfsdk:"name"`
	Online    types.Bool      `tfsdk:"online"`
	Cleanup   *CleanupModel   `tfsdk:"cleanup"`
	Component *ComponentModel `tfsdk:"component"`
	Storage   *StorageModel   `tfsdk:"storage"`
}

type RepositoryBowerProxyModel struct {
	Id                types.String        `tfsdk:"id"`
	Name              types.String        `tfsdk:"name"`
	Online            types.Bool          `tfsdk:"online"`
	Cleanup           *CleanupModel       `tfsdk:"cleanup"`
	HttpClient        *HttpClientModel    `tfsdk:"http_client"`
	NegativeCache     *NegativeCacheModel `tfsdk:"negative_cache"`
	Proxy             *ProxyModel         `tfsdk:"proxy"`
	RoutingRule       types.String        `tfsdk:"routing_rule"`
	Storage           *StorageModel       `tfsdk:"storage"`
	RewritePackagUrls types.Bool          `tfsdk:"rewrite_package_urls"`
}

type RepositoryCargoHostedModel struct {
	Id        types.String    `tfsdk:"id"`
	Name      types.String    `tfsdk:"name"`
	Online    types.Bool      `tfsdk:"online"`
	Cleanup   *CleanupModel   `tfsdk:"cleanup"`
	Component *ComponentModel `tfsdk:"component"`
	Storage   *StorageModel   `tfsdk:"storage"`
}

type RepositoryCargoProxyModel struct {
	Id                   types.String        `tfsdk:"id"`
	Name                 types.String        `tfsdk:"name"`
	Online               types.Bool          `tfsdk:"online"`
	Cleanup              *CleanupModel       `tfsdk:"cleanup"`
	HttpClient           *HttpClientModel    `tfsdk:"http_client"`
	NegativeCache        *NegativeCacheModel `tfsdk:"negative_cache"`
	Proxy                *ProxyModel         `tfsdk:"proxy"`
	RoutingRule          types.String        `tfsdk:"routing_rule"`
	Storage              *StorageModel       `tfsdk:"storage"`
	QueryCacheItemMaxAge types.Int64         `tfsdk:"query_cache_item_max_age"`
}

type RepositoryCargoGroupModel struct {
	Id      types.String  `tfsdk:"id"`
	Name    types.String  `tfsdk:"name"`
	Online  types.Bool    `tfsdk:"online"`
	Group   *GroupModel   `tfsdk:"group"`
	Storage *StorageModel `tfsdk:"storage"`
}

type RepositoryGitlfsHostedModel struct {
	Id        types.String    `tfsdk:"id"`
	Name      types.String    `tfsdk:"name"`
	Online    types.Bool      `tfsdk:"online"`
	Cleanup   *CleanupModel   `tfsdk:"cleanup"`
	Component *ComponentModel `tfsdk:"component"`
	Storage   *StorageModel   `tfsdk:"storage"`
}

type RepositoryMavenGroupModel struct {
	Id      types.String  `tfsdk:"id"`
	Name    types.String  `tfsdk:"name"`
	Online  types.Bool    `tfsdk:"online"`
	Group   *GroupModel   `tfsdk:"group"`
	Storage *StorageModel `tfsdk:"storage"`
}

type RepositoryMavenMavenModel struct {
	VersionPolicy      types.String `tfsdk:"version_policy"`
	LayoutPolicy       types.String `tfsdk:"layout_policy"`
	ContentDisposition types.String `tfsdk:"content_disposition"`
}
type RepositoryMavenProxyModel struct {
	Id            types.String                 `tfsdk:"id"`
	Name          types.String                 `tfsdk:"name"`
	Online        types.Bool                   `tfsdk:"online"`
	Cleanup       *CleanupModel                `tfsdk:"cleanup"`
	HttpClient    *HttpClientModel             `tfsdk:"http_client"`
	NegativeCache *NegativeCacheModel          `tfsdk:"negative_cache"`
	Proxy         *ProxyModel                  `tfsdk:"proxy"`
	RoutingRule   types.String                 `tfsdk:"routing_rule"`
	Storage       *StorageModel                `tfsdk:"storage"`
	Maven         []*RepositoryMavenMavenModel `tfsdk:"maven"`
}
type RepositoryMavenHostModel struct {
	Id        types.String                 `tfsdk:"id"`
	Name      types.String                 `tfsdk:"name"`
	Online    types.Bool                   `tfsdk:"online"`
	Cleanup   *CleanupModel                `tfsdk:"cleanup"`
	Component *ComponentModel              `tfsdk:"component"`
	Storage   *StorageModel                `tfsdk:"storage"`
	Maven     []*RepositoryMavenMavenModel `tfsdk:"maven"`
}

type RepositoryNPMGroupModel struct {
	RepositoryGroupBaseModel
}

type RepositoryNPMHostModel struct {
	RepositoryHostBaseModel
}

type RepositoryNPMProxyModel struct {
	RepositoryProxyBaseModel
	RemoveNonCataloged types.Bool `tfsdk:"remove_non_cataloged"`
	RemoveQuarantined  types.Bool `tfsdk:"remove_quarantined"`
}

type RepositoryNugetGroupModel struct {
	RepositoryGroupBaseModel
}

type RepositoryNugetHostModel struct {
	RepositoryHostBaseModel
}

type RepositoryNugetProxyModel struct {
	RepositoryProxyBaseModel
	NugetVersion         types.String `tfsdk:"nuget_version"`
	QueryCacheItemMaxAge types.Int64  `tfsdk:"query_cache_item_max_age"`
}
type RepositoryP2GroupModel struct {
	RepositoryGroupBaseModel
}

type RepositoryPypiGroupModel struct {
	RepositoryGroupBaseModel
}

type RepositoryPypiHostModel struct {
	RepositoryHostBaseModel
}
type RepositoryPypiProxyModel struct {
	RepositoryProxyBaseModel
}

type RepositoryRGroupModel struct {
	RepositoryGroupBaseModel
}

type RepositoryRHostModel struct {
	RepositoryHostBaseModel
}
type RepositoryRProxyModel struct {
	RepositoryProxyBaseModel
}

type RepositoryRawGroupModel struct {
	RepositoryGroupBaseModel
}
type RepositoryRawHostModel struct {
	RepositoryHostBaseModel
}
type RepositoryRawProxyModel struct {
	RepositoryProxyBaseModel
}

type RepositoryYumSigning struct {
	Keypair    types.String `tfsdk:"keypair"`
	Passphrase types.String `tfsdk:"passphrase"`
}
type RepositoryYumGroupModel struct {
	RepositoryGroupBaseModel
	YumSigning []*RepositoryYumSigning `tfsdk:"yum_signing"`
}
type RepositoryYumHostModel struct {
	RepositoryHostBaseModel
	DeployPolicy  types.String `tfsdk:"deploy_policy"`
	RepodataDepth types.String `tfsdk:"repodata_depth"`
}
type RepositoryYumProxyModel struct {
	RepositoryProxyBaseModel
	YumSigning []*RepositoryYumSigning `tfsdk:"yum_signing"`
}
