package blobstore

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nduyphuong/go-nexus-client/nexus3"
	"github.com/nduyphuong/go-nexus-client/nexus3/schema/blobstore"
)

var _ datasource.DataSource = &BlobStoreFileSource{}

func NewBlobStoreFileSource() datasource.DataSource {
	return &BlobStoreFileSource{}
}

type BlobStoreFileSource struct {
	client *nexus3.NexusClient
}

type BlobStoreFileSourceModel struct {
	Id                    types.String    `tfsdk:"id"`
	Name                  types.String    `tfsdk:"name"`
	Path                  types.String    `tfsdk:"path"`
	BlobCount             types.Int64     `tfsdk:"blob_count"`
	AvailableSpaceInBytes types.Int64     `tfsdk:"available_space_in_bytes"`
	TotalSizeInBytes      types.Int64     `tfsdk:"total_size_in_bytes"`
	SoftQuota             *SoftQuotaModel `tfsdk:"soft_quota"`
}

type SoftQuotaModel struct {
	Limit types.Int64  `tfsdk:"limit"`
	Type  types.String `tfsdk:"type"`
}

func (d *BlobStoreFileSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blobstore_file"
}

func (d *BlobStoreFileSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Used to identify data source at nexus",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Blobstore name",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "The path to the blobstore contents",
				Computed:    true,
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
			"soft_quota": schema.SingleNestedAttribute{
				MarkdownDescription: "Soft quota of the blobstore",
				Computed:            true,
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

func (d *BlobStoreFileSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BlobStoreFileSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state, newState BlobStoreFileSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	blobStoreFile, err := d.client.BlobStore.File.Get(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get projects msg from harbor failed", err.Error())
		return
	}

	var genericBlobstoreInformation blobstore.Generic
	genericBlobstores, err := d.client.BlobStore.List()
	if err != nil {
		resp.Diagnostics.AddError("Get projects msg from harbor failed", err.Error())
		return
	}
	for _, generic := range genericBlobstores {
		if generic.Name == blobStoreFile.Name {
			genericBlobstoreInformation = generic
		}
	}

	newState = BlobStoreFileSourceModel{
		Id:                    types.StringValue(blobStoreFile.Name),
		Name:                  types.StringValue(blobStoreFile.Name),
		Path:                  types.StringValue(blobStoreFile.Path),
		AvailableSpaceInBytes: types.Int64Value(int64(genericBlobstoreInformation.AvailableSpaceInBytes)),
		TotalSizeInBytes:      types.Int64Value(int64(genericBlobstoreInformation.TotalSizeInBytes)),
		BlobCount:             types.Int64Value(int64(genericBlobstoreInformation.BlobCount)),
	}
	if blobStoreFile.SoftQuota != nil {
		newState.SoftQuota = &SoftQuotaModel{
			Limit: types.Int64Value(blobStoreFile.SoftQuota.Limit / (1024 * 1024)),
			Type:  types.StringValue(blobStoreFile.SoftQuota.Type),
		}
	}

	tflog.Trace(ctx, "read a blobStoreFile data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
