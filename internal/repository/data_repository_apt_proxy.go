package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nduyphuong/go-nexus-client/nexus3"
)

var _ datasource.DataSource = &RepositoryAptProxyDatasource{}

func NewRepositoryAptProxyDatasource() datasource.DataSource {
	return &RepositoryAptProxyDatasource{}
}

type RepositoryAptProxyDatasource struct {
	client *nexus3.NexusClient
}

type RepositoryAptProxySourceModel struct {
	Id                    types.String        `tfsdk:"id"`
	Name                  types.String        `tfsdk:"name"`
	Online                types.Bool          `tfsdk:"online"`
	Flat                  types.Bool          `tfsdk:"flat"`
	Cleanup               CleanupModel        `tfsdk:"cleanup"`
	Path                  types.String        `tfsdk:"path"`
	BlobCount             types.Int64         `tfsdk:"blob_count"`
	AvailableSpaceInBytes types.Int64         `tfsdk:"available_space_in_bytes"`
	TotalSizeInBytes      types.Int64         `tfsdk:"total_size_in_bytes"`
	SoftQuota             *SoftQuotaModel     `tfsdk:"soft_quota"`
	RoutingRule           types.String        `tfsdk:"routing_rule"`
	HttpClient            *HttpClientModel    `tfsdk:"http_client"`
	NegativeCache         *NegativeCacheModel `tfsdk:"negative_cache"`
	Proxy                 ProxyModel          `tfsdk:"proxy"`
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
	UserAgentSuffix         types.String `tfsdk:"user_agent_suffix "`
}

func (d *RepositoryAptProxyDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_Repository_file"
}

func (d *RepositoryAptProxyDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to get an existing apt repository.",
		MarkdownDescription: "Use this data source to get an existing apt repository.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Used to identify data source at nexus",
				MarkdownDescription: "Used to identify data source at nexus",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "A unique identifier for this repository",
				MarkdownDescription: "A unique identifier for this repository",
				Required:            true,
			},
			"online": schema.BoolAttribute{
				Description:         "Whether this repository accepts incoming requests",
				MarkdownDescription: "Whether this repository accepts incoming requests",
				Computed:            true,
			},
			"flat": schema.BoolAttribute{
				Description:         "Distribution to fetch",
				MarkdownDescription: "Distribution to fetch",
				Computed:            true,
			},
			"cleanup": schema.SingleNestedAttribute{
				MarkdownDescription: "Cleanup policies",
				Description:         "Cleanup policies",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"policy_names": schema.ListAttribute{
						Description:         "List of policy names",
						MarkdownDescription: "List of policy names",
						Computed:            true,
					},
				},
			},
			"component": schema.SingleNestedAttribute{
				Description:         "Component configuration for the hosted repository",
				MarkdownDescription: "Component configuration for the hosted repository",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"proprietary_components": schema.BoolAttribute{
						Description:         "Components in this repository count as proprietary for namespace conflict attacks (requires Sonatype Nexus Firewall)",
						MarkdownDescription: "Components in this repository count as proprietary for namespace conflict attacks (requires Sonatype Nexus Firewall)",
						Computed:            true,
					},
				},
			},
			"storage": schema.SingleNestedAttribute{
				Description:         "The storage configuration of the repository",
				MarkdownDescription: "The storage configuration of the repository",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"blob_store_name": schema.StringAttribute{
						Description:         "Blob store used to store repository contents",
						MarkdownDescription: "Blob store used to store repository contents",
						Computed:            true,
					},
					"strict_content_type_validation": schema.BoolAttribute{
						Description:         "Whether to validate uploaded content's MIME type appropriate for the repository format",
						MarkdownDescription: "Whether to validate uploaded content's MIME type appropriate for the repository format",
						Computed:            true,
					},
				},
			},
			"distribution": schema.StringAttribute{
				Description:         "Distribution to fetch",
				MarkdownDescription: "Distribution to fetch",
				Computed:            true,
			},
			"routing_rule": schema.StringAttribute{
				Description:         "The name of the routing rule assigned to this repository",
				MarkdownDescription: "The name of the routing rule assigned to this repository",
				Computed:            true,
			},
			"http_client": schema.SingleNestedAttribute{
				Description:         "HTTP Client configuration for proxy repositories",
				MarkdownDescription: "HTTP Client configuration for proxy repositories",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"authentication": schema.SingleNestedAttribute{
						Description:         "Authentication configuration of the HTTP client",
						MarkdownDescription: "Authentication configuration of the HTTP client",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description:         "Authentication type. Possible values: `ntlm` or `username`",
								MarkdownDescription: "Authentication type. Possible values: `ntlm` or `username`",
								Computed:            true,
							},
							"username": schema.StringAttribute{
								Description:         "The username used by the proxy repository",
								MarkdownDescription: "The username used by the proxy repository",
								Computed:            true,
							},
							"password": schema.StringAttribute{
								Description:         "The password used by the proxy repository",
								MarkdownDescription: "The password used by the proxy repository",
								Computed:            true,
								Sensitive:           true,
							},
							"ntlm_domain": schema.StringAttribute{
								Description:         "The ntlm domain to connect",
								MarkdownDescription: "The ntlm domain to connect",
								Computed:            true,
							},
							"ntlm_host": schema.StringAttribute{
								Description:         "The ntlm host to connect",
								MarkdownDescription: "The ntlm host to connect",
								Computed:            true,
							},
							"connection": schema.SingleNestedAttribute{
								Description:         "Connection configuration of the HTTP client",
								MarkdownDescription: "Connection configuration of the HTTP client",
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"enable_circular_redirects": schema.BoolAttribute{
										Description:         "Whether to enable redirects to the same location (may be required by some servers)",
										MarkdownDescription: "Whether to enable redirects to the same location (may be required by some servers)",
										Computed:            true,
									},
									"enable_cookies": schema.BoolAttribute{
										Description:         "Whether to allow cookies to be stored and used",
										MarkdownDescription: "Whether to allow cookies to be stored and used",
										Computed:            true,
									},
									"retries": schema.Int64Attribute{
										Description:         "Total retries if the initial connection attempt suffers a timeout",
										MarkdownDescription: "Total retries if the initial connection attempt suffers a timeout",
										Computed:            true,
									},
									"timeout": schema.Int64Attribute{
										Description:         "Seconds to wait for activity before stopping and retrying the connection",
										MarkdownDescription: "Seconds to wait for activity before stopping and retrying the connection",
										Computed:            true,
									},
									"user_agent_suffix": schema.StringAttribute{
										Description:         "Custom fragment to append to User-Agent header in HTTP requests",
										MarkdownDescription: "Custom fragment to append to User-Agent header in HTTP requests",
										Computed:            true,
									},
									"use_trust_store": schema.BoolAttribute{
										Description:         "Use certificates stored in the Nexus Repository Manager truststore to connect to external systems",
										MarkdownDescription: "Use certificates stored in the Nexus Repository Manager truststore to connect to external systems",
										Computed:            true,
									},
								},
							},
						},
					},
					"auto_block": schema.BoolAttribute{
						Description:         "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
						MarkdownDescription: "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
						Computed:            true,
					},
					"blocked": schema.BoolAttribute{
						Description:         "Whether to block outbound connections on the repository",
						MarkdownDescription: "Whether to block outbound connections on the repository",
						Computed:            true,
					},
				},
			},
			"negative_cache": schema.SingleNestedAttribute{
				Description:         "Configuration of the negative cache handling",
				MarkdownDescription: "Configuration of the negative cache handling",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description:         "Whether to cache responses for content not present in the proxied repository",
						MarkdownDescription: "Whether to cache responses for content not present in the proxied repository",
						Computed:            true,
					},
					"ttl": schema.Int64Attribute{
						Description:         "How long to cache the fact that a file was not found in the repository (in minutes)",
						MarkdownDescription: "How long to cache the fact that a file was not found in the repository (in minutes)",
						Computed:            true,
					},
				},
			},
			"proxy": schema.SingleNestedAttribute{
				Description:         "Configuration for the proxy repository",
				MarkdownDescription: "Configuration for the proxy repository",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"content_max_age": schema.Int64Attribute{
						Description:         "How long (in minutes) to cache artifacts before rechecking the remote repository",
						MarkdownDescription: "How long (in minutes) to cache artifacts before rechecking the remote repository",
						Computed:            true,
					},
					"metadata_max_age": schema.Int64Attribute{
						Description:         "How long (in minutes) to cache metadata before rechecking the remote repository.",
						MarkdownDescription: "How long (in minutes) to cache metadata before rechecking the remote repository.",
						Computed:            true,
					},
					"remote_url": schema.StringAttribute{
						Description:         "Location of the remote repository being proxied",
						MarkdownDescription: "Location of the remote repository being proxied",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *RepositoryAptProxyDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*nexus3.NexusClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *nexus3.NexusClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *RepositoryAptProxyDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// var state, newState RepositoryAptProxySourceModel

	// tflog.Trace(ctx, "read a RepositoryAptProxy data source")
	// resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (d *RepositoryAptProxyDatasource) getState(ctx, name string) (data *RepositoryAptProxySourceModel, err error) {

	if name == "" {
		err = errors.New("name is nil")
		return
	}

	repo, err := d.client.Repository.Apt.Proxy.Get(name)
	if err != nil {
		return
	}

	data = &RepositoryAptProxySourceModel{
		Id:     types.StringValue(repo.Name),
		Name:   types.StringValue(repo.Name),
		Online: types.BoolValue(repo.Online),
		Cleanup: CleanupModel{
			PolicyNames: stringListValue(repo.Cleanup.PolicyNames),
		},
		Flat: types.BoolValue(repo.Apt.Flat),
		HttpClient: &HttpClientModel{
			Authentication: &HttpClientAuthenticationModel{
				NtlmDomain: types.StringValue(repo.HTTPClient.Authentication.NTLMDomain),
				NtlmHost:   types.StringValue(repo.HTTPClient.Authentication.NTLMHost),
				Type:       types.StringValue(string(repo.HTTPClient.Authentication.Type)),
				Username:   types.StringValue(repo.HTTPClient.Authentication.Username),
				Password:   types.StringValue(repo.HTTPClient.Authentication.Password),
			},
		},
		NegativeCache: &NegativeCacheModel{
			Enabled: types.BoolValue(repo.NegativeCache.Enabled),
			TTL:     types.Int64Value(int64(repo.NegativeCache.TTL)),
		},
		Proxy: ProxyModel{
			ContentMaxAge:  types.Int64Value(int64(repo.Proxy.ContentMaxAge)),
			MetadataMaxAge: types.Int64Value(int64(repo.Proxy.MetadataMaxAge)),
			RemoteURL:      types.StringValue(repo.Proxy.RemoteURL),
		},
	}

	return
}

func stringListValue(inputs []string) types.List {
	ret, _ := types.ListValueFrom(
		context.TODO(),
		types.StringType,
		inputs,
	)
	return ret
}
