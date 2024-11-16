package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nduyphuong/go-nexus-client/nexus3"
)

// var _ datasource.DataSource = &RepositoryDockerProxyDatasource{}

// func NewRepositoryDockerProxyDatasource() datasource.DataSource {
// 	return &RepositoryDockerProxyDatasource{}
// }

type RepositoryDockerProxyDatasource struct {
	client *nexus3.NexusClient
}

type RepositoryDockerProxySourceModel struct {
	Id            types.String            `tfsdk:"id"`
	Name          types.String            `tfsdk:"name"`
	Online        types.Bool              `tfsdk:"online"`
	Cleanup       CleanupModel            `tfsdk:"cleanup"`
	HttpClient    *HttpClientModel        `tfsdk:"http_client"`
	NegativeCache NegativeCache           `tfsdk:"negative_cache"`
	Proxy         *ProxyModel             `tfsdk:"proxy"`
	Component     ComponentModel          `tfsdk:"component"`
	RoutingRule   types.String            `tfsdk:"routing_rule"`
	Storage       *StorageDataSourceModel `tfsdk:"storage"`
}

type NegativeCache struct {
	Enabled types.String `tfsdk:"enabled"`
	Ttl     types.Int64  `tfsdk:"ttl"`
}
type DataSourceDocker struct {
	ForceBasicAuth types.Bool  `tfsdk:"force_basic_auth"`
	HttpPort       types.Int64 `tfsdk:"http_port"`
	HttpsPort      types.Int64 `tfsdk:"https_port"`
	V1Enabled      types.Bool  `tfsdk:"v1_enabled"`
}

func (d *RepositoryDockerProxyDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_docker_hosted"
}

func (d *RepositoryDockerProxyDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to get an existing docker repository.",
		MarkdownDescription: "Use this data source to get an existing docker repository.",
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
			"cleanup": schema.ListNestedAttribute{
				MarkdownDescription: "Cleanup policies",
				Description:         "Cleanup policies",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"policy_names": schema.ListAttribute{
							Description:         "List of policy names",
							MarkdownDescription: "List of policy names",
							Computed:            true,
							ElementType:         types.StringType,
						},
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
								// Sensitive:           true,
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
						},
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

func (d *RepositoryDockerProxyDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RepositoryDockerProxyDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state RepositoryDockerProxySourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	if state.Name.IsUnknown() {
		resp.Diagnostics.AddError("Get docker hosted datasource failed", "name is unknown")
	}
	state, err := d.getState(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get docker hosted datasource failed", err.Error())
	}
	tflog.Trace(ctx, "read a docker hosted data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *RepositoryDockerProxyDatasource) getState(name string) (data RepositoryDockerProxySourceModel, err error) {

	if name == "" {
		err = errors.New("name is nil")
		return
	}

	repo, err := d.client.Repository.Docker.Hosted.Get(name)
	if err != nil {
		return
	}
	data = RepositoryDockerProxySourceModel{
		Id:     types.StringValue(repo.Name),
		Name:   types.StringValue(repo.Name),
		Online: types.BoolValue(repo.Online),
		Storage: &StorageDataSourceModel{
			BlobStoreName:               types.StringValue(repo.Storage.BlobStoreName),
			StrictContentTypeValidation: types.BoolValue(repo.Storage.StrictContentTypeValidation),
		},
		Component: ComponentModel{
			ProprietaryComponents: types.BoolValue(repo.Component.ProprietaryComponents),
		},
	}
	if repo.Cleanup != nil {
		var plicyNames []types.String
		for _, item := range repo.Cleanup.PolicyNames {
			plicyNames = append(plicyNames, types.StringValue(item))
		}
		data.Cleanup = CleanupModel{
			PolicyNames: plicyNames,
		}
	}
	return
}
