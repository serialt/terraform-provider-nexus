package repository

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nduyphuong/go-nexus-client/nexus3"
)

// var _ datasource.DataSource = &RepositoryDockerHostedDatasource{}

// func NewRepositoryDockerHostedDatasource() datasource.DataSource {
// 	return &RepositoryDockerHostedDatasource{}
// }

type RepositoryDockerHostedDatasource struct {
	client *nexus3.NexusClient
}

type RepositoryDockerHostedSourceModel struct {
	Id        types.String            `tfsdk:"id"`
	Name      types.String            `tfsdk:"name"`
	Online    types.Bool              `tfsdk:"online"`
	Cleanup   CleanupModel            `tfsdk:"cleanup"`
	Component ComponentModel          `tfsdk:"component"`
	Storage   *StorageDataSourceModel `tfsdk:"storage"`
}

type DockerModel struct {
	ForceBasicAuth types.Bool   `tfsdk:"force_basic_auth"`
	HttpPort       types.Int64  `tfsdk:"http_port"`
	HttpsPort      types.Int64  `tfsdk:"https_port"`
	V1Enabled      types.Bool   `tfsdk:"v1_enabled"`
	Subdomain      types.String `tfsdk:"subdomain"`
}

func (d *RepositoryDockerHostedDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_docker_hosted"
}

func (d *RepositoryDockerHostedDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"cleanup": schema.SingleNestedAttribute{
				MarkdownDescription: "Cleanup policies",
				Description:         "Cleanup policies",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"policy_names": schema.SetAttribute{
						Description:         "List of policy names",
						MarkdownDescription: "List of policy names",
						Computed:            true,
						ElementType:         types.StringType,
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
					"write_policy": schema.StringAttribute{
						Description:         "Controls if deployments of and updates to assets are allowed",
						MarkdownDescription: "Controls if deployments of and updates to assets are allowed",
						Computed:            true,
					},
					"latest_policy": schema.BoolAttribute{
						Description:         "Whether to allow redeploying the 'latest' tag but defer to the Deployment Policy for all other tags. Only usable with write_policy \"ALLOW_ONCE\"",
						MarkdownDescription: "Whether to allow redeploying the 'latest' tag but defer to the Deployment Policy for all other tags. Only usable with write_policy \"ALLOW_ONCE\"",
						Computed:            true,
					},
				},
			},
			"distribution": schema.StringAttribute{
				Description:         "Distribution to fetch",
				MarkdownDescription: "Distribution to fetch",
				Computed:            true,
			},
			"signing": schema.SingleNestedAttribute{
				Description:         "Contains signing data of repositores",
				MarkdownDescription: "Contains signing data of repositores",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"keypair": schema.StringAttribute{
						Description: `PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor)
							If passphrase is unset, the keypair cannot be read from the nexus api.
							When reading the resource, the keypair will be read from the previous state,
							so external changes won't be detected in this case.`,
						MarkdownDescription: `PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor)
							If passphrase is unset, the keypair cannot be read from the nexus api.
							When reading the resource, the keypair will be read from the previous state,
							so external changes won't be detected in this case.`,
						Computed:  true,
						Sensitive: true,
					},
					"passphrase": schema.StringAttribute{
						Description: `Passphrase to access PGP signing key.
							This value cannot be read from the nexus api.
							When reading the resource, the value will be read from the previous state,
							so external changes won't be detected.`,
						MarkdownDescription: `Passphrase to access PGP signing key.
							This value cannot be read from the nexus api.
							When reading the resource, the value will be read from the previous state,
							so external changes won't be detected.`,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},
	}
}

// func (d *RepositoryDockerHostedDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
// 	if req.ProviderData == nil {
// 		return
// 	}
// 	client, ok := req.ProviderData.(*nexus3.NexusClient)
// 	if !ok {
// 		resp.Diagnostics.AddError(
// 			"Unexpected Data Source Configure Type",
// 			fmt.Sprintf("Expected *nexus3.NexusClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
// 		)
// 		return
// 	}
// 	d.client = client
// }

// func (d *RepositoryDockerHostedDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

// 	// var state, newState RepositoryDockerHostedSourceModel

// 	// tflog.Trace(ctx, "read a RepositoryDockerHosted data source")
// 	// resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
// }

// func (d *RepositoryDockerHostedDatasource) getState(name string)
