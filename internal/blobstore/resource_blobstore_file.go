package blobstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nduyphuong/go-nexus-client/nexus3"
	"github.com/nduyphuong/go-nexus-client/nexus3/schema/blobstore"
)

// ResourceBlobstoreFile defines the resource implementation.
type ResourceBlobstoreFile struct {
	client *nexus3.NexusClient
}

type BlobStoreFileReourceModel struct {
	Id                    types.String    `tfsdk:"id"`
	Name                  types.String    `tfsdk:"name"`
	Path                  types.String    `tfsdk:"path"`
	BlobCount             types.Int64     `tfsdk:"blob_count"`
	AvailableSpaceInBytes types.Int64     `tfsdk:"available_space_in_bytes"`
	TotalSizeInBytes      types.Int64     `tfsdk:"total_size_in_bytes"`
	SoftQuota             *SoftQuotaModel `tfsdk:"soft_quota"`
}

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ResourceBlobstoreFile{}
	_ resource.ResourceWithImportState = &ResourceBlobstoreFile{}
)

func NewResourceBlobstoreFile() resource.Resource {
	return &ResourceBlobstoreFile{}
}

func (r *ResourceBlobstoreFile) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blobstore_file"
}

func (r *ResourceBlobstoreFile) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Example resource source",

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
				Optional:    true,
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
				Optional:            true,
				// Computed:            true,
				Attributes: map[string]schema.Attribute{
					"limit": schema.Int64Attribute{
						Description: "The limit in Bytes. Minimum value is 1000000",
						Optional:    true,
						// Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
						Optional:    true,
						// Computed:    true,
					},
				},
			},
		},
	}
}
func (r *ResourceBlobstoreFile) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*nexus3.NexusClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *ResourceBlobstoreFile) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	tflog.Debug(ctx, "Create BlobStore File resource")
	var plan BlobStoreFileReourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	bPath := plan.Name.ValueString()
	if !plan.Path.IsNull() && plan.Path.ValueString() != "" {
		bPath = plan.Path.ValueString()
	}

	bFile := blobstore.File{
		Name: plan.Name.ValueString(),
		Path: bPath,
	}
	if plan.SoftQuota != nil {
		bFile.SoftQuota = &blobstore.SoftQuota{
			Type:  plan.SoftQuota.Type.ValueString(),
			Limit: plan.SoftQuota.Limit.ValueInt64() * 1024 * 1024,
		}
	}
	err := r.client.BlobStore.File.Create(&bFile)
	if err != nil {
		resp.Diagnostics.AddError("Get projects msg from harbor failed", err.Error())
		return
	}

	blobStoreFile, err := r.client.BlobStore.File.Get(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get projects msg from harbor failed", err.Error())
		return
	}

	var genericBlobstoreInformation blobstore.Generic
	genericBlobstores, err := r.client.BlobStore.List()
	if err != nil {
		resp.Diagnostics.AddError("Get projects msg from harbor failed", err.Error())
		return
	}
	for _, generic := range genericBlobstores {
		if generic.Name == blobStoreFile.Name {
			genericBlobstoreInformation = generic
		}
	}

	plan.Id = types.StringValue(blobStoreFile.Name)
	plan.Path = types.StringValue(bFile.Path)
	plan.AvailableSpaceInBytes = types.Int64Value(int64(genericBlobstoreInformation.AvailableSpaceInBytes))
	plan.TotalSizeInBytes = types.Int64Value(int64(genericBlobstoreInformation.TotalSizeInBytes))
	plan.BlobCount = types.Int64Value(int64(genericBlobstoreInformation.BlobCount))

	tflog.Trace(ctx, "read a blobStoreFile data source")

	tflog.Debug(ctx, "created a resource")

	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceBlobstoreFile) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BlobStoreFileReourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state, err := r.getState(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get blob file data msg from nexus failed", err.Error())
		return

	}

	tflog.Trace(ctx, "read a blobStoreFile data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (r *ResourceBlobstoreFile) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BlobStoreFileReourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	bPath := plan.Name.ValueString()
	if !plan.Path.IsNull() && plan.Path.ValueString() != "" {
		bPath = plan.Path.ValueString()
	}

	bFile := blobstore.File{
		Path: bPath,
	}
	if plan.SoftQuota != nil {
		bFile.SoftQuota = &blobstore.SoftQuota{
			Type:  plan.SoftQuota.Type.ValueString(),
			Limit: plan.SoftQuota.Limit.ValueInt64() * 1024 * 1024,
		}
	}
	err := r.client.BlobStore.File.Update(plan.Id.ValueString(), &bFile)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating blobstore file",
			"Could not update, unexpected error: "+err.Error(),
		)
		return
	}

	plan, err = r.getState(plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get blob file data msg from nexus failed", err.Error())
		return

	}
	tflog.Trace(ctx, "update a blobStoreFile data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ResourceBlobstoreFile) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state BlobStoreFileReourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.BlobStore.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting blobstore file",
			"Could not delete, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *ResourceBlobstoreFile) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ResourceBlobstoreFile) getState(name string) (data BlobStoreFileReourceModel, err error) {
	if name == "" {
		err = errors.New("name is nil")
		return
	}
	blobStoreFile, err := r.client.BlobStore.File.Get(name)
	if err != nil {
		return
	}

	var genericBlobstoreInformation blobstore.Generic
	genericBlobstores, err := r.client.BlobStore.List()
	if err != nil {
		return
	}
	for _, generic := range genericBlobstores {
		if generic.Name == blobStoreFile.Name {
			genericBlobstoreInformation = generic
		}
	}

	data = BlobStoreFileReourceModel{
		Id:                    types.StringValue(blobStoreFile.Name),
		Name:                  types.StringValue(blobStoreFile.Name),
		Path:                  types.StringValue(blobStoreFile.Path),
		AvailableSpaceInBytes: types.Int64Value(int64(genericBlobstoreInformation.AvailableSpaceInBytes)),
		TotalSizeInBytes:      types.Int64Value(int64(genericBlobstoreInformation.TotalSizeInBytes)),
		BlobCount:             types.Int64Value(int64(genericBlobstoreInformation.BlobCount)),
	}
	if blobStoreFile.SoftQuota != nil {
		data.SoftQuota = &SoftQuotaModel{
			Limit: types.Int64Value(blobStoreFile.SoftQuota.Limit / (1024 * 1024)),
			Type:  types.StringValue(blobStoreFile.SoftQuota.Type),
		}
	}

	return

}
