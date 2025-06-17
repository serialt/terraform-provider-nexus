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

var _ datasource.DataSource = &RepositoryGoGroupDatasource{}

func NewRepositoryGoGroupDatasource() datasource.DataSource {
	return &RepositoryGoGroupDatasource{}
}

type RepositoryGoGroupDatasource struct {
	client *nexus3.NexusClient
}

func (d *RepositoryGoGroupDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_go_proxy"
}

func (d *RepositoryGoGroupDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to get an existing docker repository.",
		MarkdownDescription: "Use this data source to get an existing docker repository.",
		Blocks: map[string]schema.Block{
			"cleanup":        tschema.DSCleanUp,
			"storage":        tschema.DSStorage,
			"http_client":    tschema.DSHttpClient,
			"negative_cache": tschema.DSNegativeCache,
			"proxy":          tschema.DSProxy,
		},
		Attributes: map[string]schema.Attribute{
			"id":     tschema.DataSourceID,
			"name":   tschema.DataSourceName,
			"online": tschema.DataSourceOnline,
			"routing_rule": schema.StringAttribute{
				Description:         "The name of the routing rule assigned to this repository",
				MarkdownDescription: "The name of the routing rule assigned to this repository",
				Computed:            true,
			},
		},
	}
}

func (d *RepositoryGoGroupDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RepositoryGoGroupDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state model.RepositoryGoGroupModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	if state.Name.IsUnknown() {
		resp.Diagnostics.AddError("Get go proxy datasource failed", "name is unknown")
	}
	state, err := d.getState(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get go proxy datasource failed", err.Error())
	}
	tflog.Trace(ctx, "read a go proxydata source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *RepositoryGoGroupDatasource) getState(name string) (data model.RepositoryGoGroupModel, err error) {

	if name == "" {
		err = errors.New("name is nil")
		return
	}

	repo, err := d.client.Repository.Go.Group.Get(name)
	if err != nil {
		return
	}
	data = model.RepositoryGoGroupModel{
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

func ListStringToType(list []string) (typeList []types.String) {
	for _, item := range list {
		typeList = append(typeList, types.StringValue(item))
	}
	return
}
