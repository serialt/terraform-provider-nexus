package nexus

import (
	"context"
	"errors"
	"fmt"

	"github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/serialt/terraform-provider-nexus/internal/model"
	"github.com/serialt/terraform-provider-nexus/internal/tschema"
)

var _ datasource.DataSource = &RepositoryCargoGroupDatasource{}

func NewRepositoryCargoGroupDatasource() datasource.DataSource {
	return &RepositoryCargoGroupDatasource{}
}

type RepositoryCargoGroupDatasource struct {
	client *nexus3.NexusClient
}

func (d *RepositoryCargoGroupDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_cargo_group"
}

func (d *RepositoryCargoGroupDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to get an existing cargo group repository.",
		MarkdownDescription: "Use this data source to get an existing cargo group repository.",
		Blocks: map[string]schema.Block{
			"storage": tschema.DSStorage,
			"group":   tschema.DSGroup,
		},
		Attributes: map[string]schema.Attribute{
			"id":     tschema.DataSourceID,
			"name":   tschema.DataSourceName,
			"online": tschema.DataSourceOnline,
		},
	}
}

func (d *RepositoryCargoGroupDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RepositoryCargoGroupDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state model.RepositoryCargoGroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.Name.IsUnknown() {
		resp.Diagnostics.AddError("Get cargo group datasource failed", "name is unknown")
	}
	state, err := RepositoryCargoGroupGetState(d.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get cargo group datasource failed", err.Error())
	}
	tflog.Trace(ctx, "read a cargo group data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func RepositoryCargoGroupGetState(client *nexus3.NexusClient, name string) (data model.RepositoryCargoGroupModel, err error) {
	if name == "" {
		err = errors.New("name is nil")
		return
	}
	repo, err := client.Repository.Cargo.Group.Get(name)
	if err != nil {
		return
	}
	data = model.RepositoryCargoGroupModel{
		Id:     types.StringValue(repo.Name),
		Name:   types.StringValue(repo.Name),
		Online: types.BoolValue(repo.Online),
		Storage: &model.StorageModel{
			BlobStoreName:               types.StringValue(repo.Storage.BlobStoreName),
			StrictContentTypeValidation: types.BoolValue(repo.Storage.StrictContentTypeValidation),
		},
	}
	if len(repo.MemberNames) > 0 {
		data.Group = &model.GroupModel{
			MemberNames: ListStringToType(repo.MemberNames),
		}
	}

	return
}
