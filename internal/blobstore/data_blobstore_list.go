package blobstore

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nduyphuong/go-nexus-client/nexus3"
)

var _ datasource.DataSource = &BlobStoreListSource{}

func NewBlobStoreListSource() datasource.DataSource {
	return &BlobStoreListSource{}
}

type BlobStoreListSource struct {
	client *nexus3.NexusClient
}

type BlobStoreListSourceModel struct {
	Id    types.String                    `tfsdk:"id"`
	Items []*BlobStoreListSourceItemModel `tfsdk:"items"`
}

type BlobStoreListSourceItemModel struct {
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

func (d *BlobStoreListSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blobstore_list"
}

func (d *BlobStoreListSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get a list with all Blob Stores.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Used to identify data source at nexus",
				Computed:    true,
			},
			"items": schema.ListNestedAttribute{
				MarkdownDescription: "A List of all Blob Stores",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Blobstore name",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The type of current blob store",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *BlobStoreListSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BlobStoreListSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state, newState BlobStoreListSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	blobStoreList, err := d.client.BlobStore.List()
	if err != nil {
		resp.Diagnostics.AddError("Get blobStore list failed", err.Error())
		return
	}
	for _, item := range blobStoreList {
		newState.Items = append(newState.Items, &BlobStoreListSourceItemModel{
			Name: types.StringValue(item.Name),
			Type: types.StringValue(item.Type),
		})
	}

	tflog.Trace(ctx, "read a BlobStoreList data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
