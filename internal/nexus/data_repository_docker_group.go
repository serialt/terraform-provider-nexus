package nexus

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nduyphuong/go-nexus-client/nexus3"
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
			"negative_cache": tschema.DSNegativeCache,
			"storage":        tschema.DSStorage,
			"http_client":    tschema.DSHttpClient,
			"proxy":          tschema.DSProxy,
			"cleanup":        tschema.DSCleanUp,
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Used to identify data source at nexus",
				MarkdownDescription: "Used to identify data source at nexus",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "A unique identifier for this repository",
				MarkdownDescription: "A unique identifier for this repository",
				Required:            true,
			},
			"online": schema.BoolAttribute{
				Description:         "Whether this repository accepts incoming requests",
				MarkdownDescription: "Whether this repository accepts incoming requests",
				Computed:            true,
			},
			"distribution": schema.StringAttribute{
				Description:         "Distribution to fetch",
				MarkdownDescription: "Distribution to fetch",
				Computed:            true,
			},
			"routing_rule": schema.StringAttribute{
				Description:         "The name of the routing rule assigned to this repository",
				MarkdownDescription: "The name of the routing rule assigned to this repository",
				Computed:            true,
			},
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
		Component: &model.ComponentModel{
			ProprietaryComponents: types.BoolValue(repo.Component.ProprietaryComponents),
		},
	}
	if repo.Cleanup != nil {
		var plicyNames []types.String
		for _, item := range repo.Cleanup.PolicyNames {
			plicyNames = append(plicyNames, types.StringValue(item))
		}
		data.Cleanup = &model.CleanupModel{
			PolicyNames: plicyNames,
		}
	}
	return
}
