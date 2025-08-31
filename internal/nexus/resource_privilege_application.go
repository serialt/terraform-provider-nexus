package nexus

import (
	"context"
	"fmt"

	"github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/blobstore"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/serialt/terraform-provider-nexus/internal/model"
)

// ResourcePrivilegeApplication defines the resource implementation.
type ResourcePrivilegeApplication struct {
	client *nexus3.NexusClient
}

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ResourcePrivilegeApplication{}
	_ resource.ResourceWithImportState = &ResourcePrivilegeApplication{}
)

func NewResourcePrivilegeApplication() resource.Resource {
	return &ResourcePrivilegeApplication{}
}

func (r *ResourcePrivilegeApplication) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_privilege_application"
}

func (r *ResourcePrivilegeApplication) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Example resource source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Used to identify resource at nexus",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Privilege application name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "description",
				Optional:    true,
			},
			"actions": schema.ListAttribute{
				Description: "List of the names of action",
				Computed:    true,
				ElementType: types.StringType,
			},
			"domain": schema.StringAttribute{
				Description: "The total size of the blobstore in Bytes",
				Optional:    true,
			},
		},
	}
}
func (r *ResourcePrivilegeApplication) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourcePrivilegeApplication) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	tflog.Debug(ctx, "Create Privilege Application resource")
	var plan model.PrivilegeApplication

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
	fmt.Println(plan.Name.ValueString())
	state, err := BlobstoreFileGetState(r.client, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get blob file from nexus failed", err.Error())
		return
	}
	tflog.Debug(ctx, "created a resource of blob file")

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourcePrivilegeApplication) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.BlobStoreFileModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state, err := BlobstoreFileGetState(r.client, state.Id.String())
	if err != nil {
		resp.Diagnostics.AddError("Get blob file data from nexus failed", err.Error())
		return

	}

	tflog.Trace(ctx, "read a blobStoreFile data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (r *ResourcePrivilegeApplication) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.BlobStoreFileModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	bPath := plan.Name.ValueString()
	if !plan.Path.IsNull() && plan.Path.ValueString() != "" {
		bPath = plan.Path.ValueString()
	}

	item, err := r.client.BlobStore.File.Get(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get blob file data from nexus failed", err.Error())
		return
	}
	item.Path = bPath

	if plan.SoftQuota != nil {
		item.SoftQuota = &blobstore.SoftQuota{
			Type:  plan.SoftQuota.Type.ValueString(),
			Limit: plan.SoftQuota.Limit.ValueInt64() * 1024 * 1024,
		}
	}
	fmt.Println(plan.Id.ValueString())

	err = r.client.BlobStore.File.Update(plan.Name.ValueString(), item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating blobstore file",
			"Could not update, unexpected error: "+err.Error(),
		)
		return
	}

	plan, err = BlobstoreFileGetState(r.client, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get blob file data msg from nexus failed", err.Error())
		return

	}
	tflog.Trace(ctx, "update a blobStoreFile data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ResourcePrivilegeApplication) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.BlobStoreFileModel
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

func (r *ResourcePrivilegeApplication) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
