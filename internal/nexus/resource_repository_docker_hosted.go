package nexus

import (
	"context"
	"fmt"

	"github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/repository"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/serialt/terraform-provider-nexus/internal/model"
	"github.com/serialt/terraform-provider-nexus/internal/tschema"
	"github.com/spf13/cast"
)

// ResourceRepositoryDockerHostedModel defines the resource implementation.
type ResourceRepositoryDockerHostedModel struct {
	client *nexus3.NexusClient
}

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ResourceRepositoryDockerHostedModel{}
	_ resource.ResourceWithImportState = &ResourceRepositoryDockerHostedModel{}
)

func NewResourceRepositoryDockerHostedResource() resource.Resource {
	return &ResourceRepositoryDockerHostedModel{}
}

func (r *ResourceRepositoryDockerHostedModel) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_docker_hosted"
}

func (r *ResourceRepositoryDockerHostedModel) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Example resource source",
		Blocks: map[string]schema.Block{
			"cleanup":   tschema.ResourceCleanUp,
			"component": tschema.ResourceComponent,
			"storage":   tschema.ResourceDockerHostedStorage,
			"docker":    tschema.ResourceDocker,
		},
		Attributes: map[string]schema.Attribute{
			"id":     tschema.ResourceID,
			"name":   tschema.ResourceName,
			"online": tschema.ResourceOnline,
		},
	}
}
func (r *ResourceRepositoryDockerHostedModel) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ResourceRepositoryDockerHostedModel) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	tflog.Debug(ctx, "Create Privilege Application resource")
	var plan model.RepositoryDockerHostedModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// r.client.Repository.Docker.Hosted.Create(repo repository.DockerHostedRepository)
	repo := repository.DockerHostedRepository{
		Name:   plan.Name.ValueString(),
		Online: plan.Online.ValueBool(),
		Storage: repository.DockerHostedStorage{
			BlobStoreName:               plan.Storage.BlobStoreName.ValueString(),
			StrictContentTypeValidation: plan.Storage.StrictContentTypeValidation.ValueBool(),
			WritePolicy:                 repository.StorageWritePolicy(plan.Storage.WritePolicy.ValueString()),
		},
	}

	if plan.Docker != nil {
		repo.Docker = repository.Docker{
			ForceBasicAuth: plan.Docker.ForceBasicAuth.ValueBool(),
			V1Enabled:      plan.Docker.V1Enabled.ValueBool(),
		}
		if !plan.Docker.HttpPort.IsNull() && plan.Docker.HttpPort.ValueInt64() > 0 {
			httpPort := cast.ToInt(plan.Docker.HttpPort.ValueInt64())
			repo.Docker.HTTPPort = &httpPort
		}
		if !plan.Docker.HttpsPort.IsNull() && plan.Docker.HttpsPort.ValueInt64() > 0 {
			httpsPort := cast.ToInt(plan.Docker.HttpsPort.ValueInt64())
			repo.Docker.HTTPPort = &httpsPort
		}
	}
	if !plan.Storage.LatestPolicy.IsNull() {
		repo.Storage.LatestPolicy = plan.Storage.LatestPolicy.ValueBoolPointer()
	}

	err := r.client.Repository.Docker.Hosted.Create(repo)
	if err != nil {
		resp.Diagnostics.AddError("Create Privilege failed", err.Error())
		return
	}
	plan.Id = plan.Name

	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceRepositoryDockerHostedModel) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.RepositoryDockerHostedModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	repo, err := r.client.Repository.Docker.Hosted.Get(state.Id.ValueString())
	if err != nil {
		return
	}
	state = model.RepositoryDockerHostedModel{
		Id:     types.StringValue(repo.Name),
		Name:   types.StringValue(repo.Name),
		Online: types.BoolValue(repo.Online),
		Storage: &model.DockerHostedStorageModel{
			BlobStoreName:               types.StringValue(repo.Storage.BlobStoreName),
			StrictContentTypeValidation: types.BoolValue(repo.Storage.StrictContentTypeValidation),
			WritePolicy:                 types.StringValue(string(repo.Storage.WritePolicy)),
			LatestPolicy:                types.BoolPointerValue(repo.Storage.LatestPolicy),
		},
	}

	tflog.Trace(ctx, "read a blobStoreFile data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ResourceRepositoryDockerHostedModel) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.RepositoryDockerHostedModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// r.client.Repository.Docker.Hosted.Create(repo repository.DockerHostedRepository)
	repo := repository.DockerHostedRepository{
		Name:   plan.Name.ValueString(),
		Online: plan.Online.ValueBool(),
		Storage: repository.DockerHostedStorage{
			BlobStoreName:               plan.Storage.BlobStoreName.ValueString(),
			StrictContentTypeValidation: plan.Storage.StrictContentTypeValidation.ValueBool(),
			WritePolicy:                 repository.StorageWritePolicy(plan.Storage.WritePolicy.ValueString()),
		},
	}

	err := r.client.Repository.Docker.Hosted.Update(plan.Id.ValueString(), repo)
	if err != nil {
		resp.Diagnostics.AddError("update Privilege Apllication failed", err.Error())
		return
	}

	tflog.Trace(ctx, "update a Privilige application data")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ResourceRepositoryDockerHostedModel) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.RepositoryDockerHostedModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Repository.Docker.Hosted.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Privilege application",
			"Could not delete, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *ResourceRepositoryDockerHostedModel) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
