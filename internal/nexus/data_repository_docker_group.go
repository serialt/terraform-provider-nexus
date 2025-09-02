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

var _ datasource.DataSource = &RepositoryDockerGroupDatasource{}

func NewRepositoryDockerGroupDatasource() datasource.DataSource {
	return &RepositoryDockerGroupDatasource{}
}

type RepositoryDockerGroupDatasource struct {
	client *nexus3.NexusClient
}

func (d *RepositoryDockerGroupDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_docker_group"
}

func (d *RepositoryDockerGroupDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to get an existing docker repository.",
		MarkdownDescription: "Use this data source to get an existing docker repository.",
		Blocks: map[string]schema.Block{
			"storage": tschema.DSStorage,
			"group":   tschema.DSGroup,
			"docker":  tschema.DSDocker,
		},
		Attributes: map[string]schema.Attribute{
			"id":     tschema.DataSourceID,
			"name":   tschema.DataSourceName,
			"online": tschema.DataSourceOnline,
		},
	}
}

func (d *RepositoryDockerGroupDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RepositoryDockerGroupDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state model.RepositoryDockerGroupModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	if state.Name.IsUnknown() {
		resp.Diagnostics.AddError("Get docker hosted datasource failed", "name is unknown")
	}
	state, err := RepositoryDockerGroupGetState(d.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get docker hosted datasource failed", err.Error())
	}
	tflog.Trace(ctx, "read a docker hosted data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func RepositoryDockerGroupGetState(client *nexus3.NexusClient, name string) (data model.RepositoryDockerGroupModel, err error) {
	if name == "" {
		err = errors.New("name is nil")
		return
	}
	repo, err := client.Repository.Docker.Hosted.Get(name)
	if err != nil {
		return
	}
	data = model.RepositoryDockerGroupModel{
		Id:     types.StringValue(repo.Name),
		Name:   types.StringValue(repo.Name),
		Online: types.BoolValue(repo.Online),
		Storage: &model.StorageModel{
			BlobStoreName:               types.StringValue(repo.Storage.BlobStoreName),
			StrictContentTypeValidation: types.BoolValue(repo.Storage.StrictContentTypeValidation),
		},
		Docker: &model.DockerModel{
			ForceBasicAuth: types.BoolValue(repo.Docker.ForceBasicAuth),
			HttpPort:       types.Int64Value(int64(GetValue(repo.Docker.HTTPPort))),
			HttpsPort:      types.Int64Value(int64(GetValue(repo.Docker.HTTPSPort))),
			V1Enabled:      types.BoolValue(repo.Docker.V1Enabled),
			Subdomain:      types.StringValue(*repo.Docker.Subdomain),
		},
	}
	return
}
