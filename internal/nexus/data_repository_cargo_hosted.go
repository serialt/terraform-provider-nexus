package nexus

import (
	"context"
	"errors"
	"fmt"

	"github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/serialt/terraform-provider-nexus/internal/model"
	"github.com/serialt/terraform-provider-nexus/internal/tschema"
)

var _ datasource.DataSource = &RepositoryCargoHostedDatasource{}

func NewRepositoryCargoHostedDatasource() datasource.DataSource {
	return &RepositoryCargoHostedDatasource{}
}

type RepositoryCargoHostedDatasource struct {
	client *nexus3.NexusClient
}

func (d *RepositoryCargoHostedDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_cargo_hosted"
}

func (d *RepositoryCargoHostedDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to get an existing cargo hosted repository.",
		MarkdownDescription: "Use this data source to get an existing cargo hosted repository.",
		Blocks: map[string]schema.Block{
			"cleanup":   tschema.DSCleanUp,
			"component": tschema.DSComponent,
			"storage":   tschema.DSHostedStorage,
		},
		Attributes: map[string]schema.Attribute{
			"id":     tschema.ResourceID,
			"name":   tschema.DataSourceName,
			"online": tschema.DataSourceOnline,
		},
	}
}

func (d *RepositoryCargoHostedDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RepositoryCargoHostedDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state model.RepositoryCargoHostedModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	if state.Name.IsUnknown() {
		resp.Diagnostics.AddError("Get cargo hosted datasource failed", "name is unknown")
	}
	state, err := RepositoryCargoHostedGetState(d.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get cargo hosted datasource failed", err.Error())
	}
	tflog.Trace(ctx, "read a cargo hosted data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func RepositoryCargoHostedGetState(client *nexus3.NexusClient, name string) (data model.RepositoryCargoHostedModel, err error) {

	if name == "" {
		err = errors.New("name is nil")
		return
	}

	repo, err := client.Repository.Cargo.Hosted.Get(name)
	if err != nil {
		return
	}
	data = model.RepositoryCargoHostedModel{
		Id:     types.StringValue(repo.Name),
		Name:   types.StringValue(repo.Name),
		Online: types.BoolValue(repo.Online),
		Storage: &model.StorageModel{
			BlobStoreName:               types.StringValue(repo.Storage.BlobStoreName),
			StrictContentTypeValidation: types.BoolValue(repo.Storage.StrictContentTypeValidation),
		},
		Component: &model.ComponentModel{
			ProprietaryComponents: types.BoolValue(repo.Component.ProprietaryComponents),
		},
	}
	if repo.Cleanup != nil {
		policyNames := []attr.Value{}
		for _, item := range repo.Cleanup.PolicyNames {
			policyNames = append(policyNames, types.StringValue(item))
		}
		policyNamesTfsdk, _ := types.ListValue(types.StringType, policyNames)
		data.Cleanup = &model.CleanupModel{
			PolicyNames: policyNamesTfsdk,
		}
	}
	return
}
