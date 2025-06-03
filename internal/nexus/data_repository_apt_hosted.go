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

var _ datasource.DataSource = &RepositoryAptHostedDatasource{}

func NewRepositoryAptHostedDatasource() datasource.DataSource {
	return &RepositoryAptHostedDatasource{}
}

type RepositoryAptHostedDatasource struct {
	client *nexus3.NexusClient
}

func (d *RepositoryAptHostedDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_apt_hosted"
}

func (d *RepositoryAptHostedDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to get an existing apt repository.",
		MarkdownDescription: "Use this data source to get an existing apt repository.",
		Blocks: map[string]schema.Block{
			"component": tschema.DSComponent,
			"storage":   tschema.DSStorage,
			"signing":   tschema.DSSigning,
			"cleanup":   tschema.DSCleanUp,
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "Used to identify data source at nexus",
				MarkdownDescription: "Used to identify data source at nexus",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				Description:         "Repository name",
				MarkdownDescription: "Repository name",
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
		},
	}
}

func (d *RepositoryAptHostedDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*nexus3.NexusClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *nexus3.NexusClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *RepositoryAptHostedDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state model.RepositoryAptHostedModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	if state.Name.IsUnknown() {
		resp.Diagnostics.AddError("Get apt hosted datasource failed", "name is unknown")
	}
	state, err := d.getState(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get apt hosted datasource failed", err.Error())
	}
	tflog.Trace(ctx, "read a apt hosted data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *RepositoryAptHostedDatasource) getState(name string) (data model.RepositoryAptHostedModel, err error) {
	if name == "" {
		err = errors.New("name is nil")
		return
	}

	repo, err := d.client.Repository.Apt.Hosted.Get(name)
	if err != nil {
		return
	}

	data = model.RepositoryAptHostedModel{
		Id:     types.StringValue(repo.Name),
		Name:   types.StringValue(repo.Name),
		Online: types.BoolValue(repo.Online),
		Storage: &model.StorageModelV2{
			BlobStoreName:               types.StringValue(repo.Storage.BlobStoreName),
			StrictContentTypeValidation: types.BoolValue(repo.Storage.StrictContentTypeValidation),
			WritePolicy:                 types.StringValue(string(*repo.Storage.WritePolicy)),
		},
		Distribution: types.StringValue(repo.Apt.Distribution),
		Component: []*model.ComponentModel{{
			ProprietaryComponents: types.BoolValue(repo.Component.ProprietaryComponents),
		}},
		Signing: &model.SigningModel{
			Keypair:    []types.String{types.StringValue(repo.AptSigning.Keypair)},
			Passphrase: []types.String{types.StringValue(GetValue(repo.AptSigning.Passphrase))},
		},
	}
	if repo.Cleanup != nil {
		var plicyNames []types.String
		for _, item := range repo.Cleanup.PolicyNames {
			plicyNames = append(plicyNames, types.StringValue(item))
		}
		data.Cleanup = []*model.CleanupModel{{
			PolicyNames: plicyNames,
		}}
	}

	return
}
