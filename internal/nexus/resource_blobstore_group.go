package nexus

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nduyphuong/go-nexus-client/nexus3"
	"github.com/nduyphuong/go-nexus-client/nexus3/schema/blobstore"
	"github.com/samber/lo"
	"github.com/serialt/terraform-provider-nexus/internal/model"
	"github.com/serialt/terraform-provider-nexus/internal/tschema"
)

// ResourceBlobstoreGroup defines the resource implementation.
type ResourceBlobstoreGroup struct {
	client *nexus3.NexusClient
}

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ResourceBlobstoreGroup{}
	_ resource.ResourceWithImportState = &ResourceBlobstoreGroup{}
)

func NewResourceBlobstoreGroup() resource.Resource {
	return &ResourceBlobstoreGroup{}
}

func (r *ResourceBlobstoreGroup) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blobstore_group"
}

func (r *ResourceBlobstoreGroup) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Example resource source",
		Blocks: map[string]schema.Block{
			"soft_quota": tschema.RsoftQuota,
		},
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
			"fill_policy": schema.StringAttribute{
				Description: "The policy how to fill the members. Possible values: `roundRobin` or `writeToFirst`",
				Required:    true,
				Validators:  []validator.String{stringvalidator.OneOf("roundRobin", "writeToFirst")},
			},
			"members": schema.ListAttribute{
				Description: "List of the names of blob stores that are members of this group",
				Required:    true,
				ElementType: types.StringType,
				// Validators:  []validator.List{listvalidator.SizeAtLeast(1)}, // 校验,至少一个
			},
		},
	}
}
func (r *ResourceBlobstoreGroup) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceBlobstoreGroup) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	tflog.Debug(ctx, "Create BlobStore File resource")
	var plan model.BlobStoreGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	members := lo.Map(plan.Members, func(x types.String, index int) string {
		return x.ValueString()
	})
	bGroup := blobstore.Group{
		Name:       plan.Name.ValueString(),
		Members:    members,
		FillPolicy: plan.FillPolicy.ValueString(),
	}

	if plan.SoftQuota != nil {
		bGroup.SoftQuota = &blobstore.SoftQuota{
			Type:  plan.SoftQuota.Type.ValueString(),
			Limit: plan.SoftQuota.Limit.ValueInt64() * 1024 * 1024,
		}
	}
	err := r.client.BlobStore.Group.Create(&bGroup)
	if err != nil {
		resp.Diagnostics.AddError("Get projects msg from harbor failed", err.Error())
		return
	}

	state, err := BlobstoreGroupGetState(r.client, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get state failed", err.Error())
		return
	}
	tflog.Debug(ctx, "created a resource")

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceBlobstoreGroup) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.BlobStoreGroupModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state, err := BlobstoreGroupGetState(r.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get blob file data msg from nexus failed", err.Error())
		return

	}

	tflog.Trace(ctx, "read a BlobstoreGroup data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

}

func (r *ResourceBlobstoreGroup) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.BlobStoreGroupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	bGroup, err := r.client.BlobStore.Group.Get(plan.Id.String())
	if err != nil {
		resp.Diagnostics.AddError("Get blob group data msg from nexus failed", err.Error())
		return
	}

	bGroup.FillPolicy = plan.FillPolicy.ValueString()
	bGroup.Members = lo.Map(plan.Members, func(x types.String, index int) string {
		return x.ValueString()
	})

	if plan.SoftQuota != nil {
		bGroup.SoftQuota = &blobstore.SoftQuota{
			Type:  plan.SoftQuota.Type.ValueString(),
			Limit: plan.SoftQuota.Limit.ValueInt64() * 1024 * 1024,
		}
	}
	err = r.client.BlobStore.Group.Update(plan.Id.ValueString(), bGroup)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating blobstore file",
			"Could not update, unexpected error: "+err.Error(),
		)
		return
	}

	state, err := BlobstoreGroupGetState(r.client, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get blob file data msg from nexus failed", err.Error())
		return

	}
	tflog.Trace(ctx, "update a BlobstoreGroup data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ResourceBlobstoreGroup) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.BlobStoreGroupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.BlobStore.Group.Delete(state.Id.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting blobstore group",
			"Could not delete, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *ResourceBlobstoreGroup) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func BlobstoreGroupGetState(client *nexus3.NexusClient, name string) (data model.BlobStoreGroupModel, err error) {
	if name == "" {
		err = errors.New("name is nil")
		return
	}
	BlobstoreGroup, err := client.BlobStore.Group.Get(name)
	if err != nil {
		return
	}

	var genericBlobstoreInformation blobstore.Generic
	genericBlobstores, err := client.BlobStore.List()
	if err != nil {
		return
	}

	for _, generic := range genericBlobstores {
		if generic.Name == BlobstoreGroup.Name {
			genericBlobstoreInformation = generic
		}
	}

	members := lo.Map(BlobstoreGroup.Members, func(x string, index int) types.String {
		return types.StringValue(x)
	})
	data = model.BlobStoreGroupModel{
		Id:                    types.StringValue(BlobstoreGroup.Name),
		Name:                  types.StringValue(BlobstoreGroup.Name),
		AvailableSpaceInBytes: types.Int64Value(int64(genericBlobstoreInformation.AvailableSpaceInBytes)),
		TotalSizeInBytes:      types.Int64Value(int64(genericBlobstoreInformation.TotalSizeInBytes)),
		BlobCount:             types.Int64Value(int64(genericBlobstoreInformation.BlobCount)),
		Members:               members,
		FillPolicy:            types.StringValue(BlobstoreGroup.FillPolicy),
	}
	if BlobstoreGroup.SoftQuota != nil {
		data.SoftQuota = &model.SoftQuotaModel{
			Limit: types.Int64Value(BlobstoreGroup.SoftQuota.Limit / (1024 * 1024)),
			Type:  types.StringValue(BlobstoreGroup.SoftQuota.Type),
		}
	}

	return

}
