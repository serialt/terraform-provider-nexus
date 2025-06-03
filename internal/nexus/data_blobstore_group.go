package nexus

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nduyphuong/go-nexus-client/nexus3"
	"github.com/serialt/terraform-provider-nexus/internal/model"
)

var _ datasource.DataSource = &BlobStoreGroupSource{}

func NewBlobStoreGroupSource() datasource.DataSource {
	return &BlobStoreGroupSource{}
}

type BlobStoreGroupSource struct {
	client *nexus3.NexusClient
}

func (d *BlobStoreGroupSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blobstore_group"
}

func (d *BlobStoreGroupSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `~> PRO Feature

	Use this data source to get details of an existing Nexus Group blobstore.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Used to identify data source at nexus",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Blobstore name",
				Required:    true,
			},
			"available_space_in_bytes": schema.Int64Attribute{
				Computed:    true,
				Description: "Available space in Bytes",
			},
			"total_size_in_bytes": schema.Int64Attribute{
				Computed:    true,
				Description: "The total size of the blobstore in Bytes",
			},
			"blob_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Count of blobs",
			},
			"fill_policy": schema.StringAttribute{
				MarkdownDescription: "The policy how to fill the members. Possible values: `roundRobin` or `writeToFirst`",
				Computed:            true,
			},
			"members": schema.ListAttribute{
				Description: "List of the names of blob stores that are members of this group",
				Computed:    true,
				ElementType: types.StringType,
			},
			"soft_quota": schema.SingleNestedAttribute{
				MarkdownDescription: "Soft quota of the blobstore",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"limit": schema.Int64Attribute{
						Description: "The limit in Bytes. Minimum value is 1000000",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
						Computed:    true,
					},
				},
			},
		},
	}

}

func (d *BlobStoreGroupSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BlobStoreGroupSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state, newState model.BlobStoreGroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, err := BlobstoreGroupGetState(d.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get blobStore group failed", err.Error())
		return
	}

	tflog.Trace(ctx, "read a BlobStoreGroup data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
