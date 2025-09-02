package nexus

import (
	"context"
	"fmt"

	"github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	actions := []string{}
	_ = plan.Actions.ElementsAs(ctx, actions, true)
	var appActions []security.SecurityPrivilegeApplicationActions
	for _, v := range actions {
		appActions = append(appActions, security.SecurityPrivilegeApplicationActions(v))
	}
	privilegeApp := security.PrivilegeApplication{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Actions:     appActions,
		Domain:      plan.Domain.ValueString(),
	}
	err := r.client.Security.Privilege.Application.Create(privilegeApp)
	if err != nil {
		resp.Diagnostics.AddError("Create Privilege failed", err.Error())
		return
	}

	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ResourcePrivilegeApplication) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state model.PrivilegeApplication

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	privilege, err := r.client.Security.Privilege.Get(state.Id.String())
	if err != nil {
		resp.Diagnostics.AddError("Get Privilege failed", err.Error())
		return
	}

	actions := []attr.Value{}
	for _, v := range privilege.Actions {
		actions = append(actions, types.StringValue(v))
	}
	actionTfsdk, _ := types.ListValue(types.StringType, actions)
	state = model.PrivilegeApplication{
		Id:          types.StringValue(privilege.Name),
		Name:        types.StringValue(privilege.Name),
		Description: types.StringValue(privilege.Description),
		Actions:     actionTfsdk,
		Domain:      types.StringValue(privilege.Domain),
	}

	tflog.Trace(ctx, "read a blobStoreFile data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ResourcePrivilegeApplication) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan model.PrivilegeApplication
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	actions := []string{}
	_ = plan.Actions.ElementsAs(ctx, actions, true)
	var appActions []security.SecurityPrivilegeApplicationActions
	for _, v := range actions {
		appActions = append(appActions, security.SecurityPrivilegeApplicationActions(v))
	}
	privilegeApp := security.PrivilegeApplication{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Actions:     appActions,
		Domain:      plan.Domain.ValueString(),
	}

	err := r.client.Security.Privilege.Application.Update(plan.Name.String(), privilegeApp)
	if err != nil {
		resp.Diagnostics.AddError("update Privilege Apllication failed", err.Error())
		return
	}

	tflog.Trace(ctx, "update a Privilige application data")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ResourcePrivilegeApplication) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state model.PrivilegeApplication
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Security.Privilege.Delete(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Privilege application",
			"Could not delete, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *ResourcePrivilegeApplication) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
