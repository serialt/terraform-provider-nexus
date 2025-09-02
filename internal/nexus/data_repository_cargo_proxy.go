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

var _ datasource.DataSource = &RepositoryCargoProxyDatasource{}

func NewRepositoryCargoProxyDatasource() datasource.DataSource {
	return &RepositoryCargoProxyDatasource{}
}

type RepositoryCargoProxyDatasource struct {
	client *nexus3.NexusClient
}

func (d *RepositoryCargoProxyDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_cargo_hosted"
}

func (d *RepositoryCargoProxyDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to get an existing cargo proxy repository.",
		MarkdownDescription: "Use this data source to get an existing cargo proxy repository.",
		Blocks: map[string]schema.Block{
			"cleanup":        tschema.DSCleanUp,
			"storage":        tschema.DSStorage,
			"http_client":    tschema.DSHttpClient,
			"negative_cache": tschema.DSNegativeCache,
			"proxy":          tschema.DSProxy,
		},
		Attributes: map[string]schema.Attribute{
			"id":           tschema.DataSourceID,
			"name":         tschema.DataSourceName,
			"online":       tschema.DataSourceOnline,
			"routing_rule": tschema.DataSourceRoutingRule,
			"cargo_version": schema.StringAttribute{
				Description:         "Cargo protocol version",
				MarkdownDescription: "Cargo protocol version",
				Computed:            true,
			},
			"query_cache_item_max_age": schema.Int64Attribute{
				Description:         "How long to cache query results from the proxied repository (in seconds)",
				MarkdownDescription: "How long to cache query results from the proxied repository (in seconds)",
				Computed:            true,
			},
		},
	}
}

func (d *RepositoryCargoProxyDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RepositoryCargoProxyDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state model.RepositoryCargoProxyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	if state.Name.IsUnknown() {
		resp.Diagnostics.AddError("Get docker hosted datasource failed", "name is unknown")
	}
	state, err := d.getState(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get docker hosted datasource failed", err.Error())
	}
	tflog.Trace(ctx, "read a docker hosted data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *RepositoryCargoProxyDatasource) getState(name string) (data model.RepositoryCargoProxyModel, err error) {

	if name == "" {
		err = errors.New("name is nil")
		return
	}

	repo, err := d.client.Repository.Cargo.Proxy.Get(name)
	if err != nil {
		return
	}
	data = model.RepositoryCargoProxyModel{
		Id:     types.StringValue(repo.Name),
		Name:   types.StringValue(repo.Name),
		Online: types.BoolValue(repo.Online),
		NegativeCache: &model.NegativeCacheModel{
			Enabled: types.BoolValue(repo.NegativeCache.Enabled),
			TTL:     types.Int64Value(int64(repo.NegativeCache.TTL)),
		},
		Proxy: &model.ProxyModel{
			ContentMaxAge:  types.Int64Value(int64(repo.Proxy.ContentMaxAge)),
			MetadataMaxAge: types.Int64Value(int64(repo.Proxy.MetadataMaxAge)),
			RemoteURL:      types.StringValue(repo.Proxy.RemoteURL),
		},
		Storage: &model.StorageModel{
			BlobStoreName:               types.StringValue(repo.BlobStoreName),
			StrictContentTypeValidation: types.BoolValue(repo.StrictContentTypeValidation),
		},
		HttpClient: &model.HttpClientModel{
			AutoBlock:      types.BoolValue(repo.HTTPClient.AutoBlock),
			Blocked:        types.BoolValue(repo.HTTPClient.Blocked),
			Authentication: &model.HttpClientAuthenticationModel{},
		},
	}
	if repo.HTTPClient.Authentication != nil {
		data.HttpClient = &model.HttpClientModel{
			Authentication: &model.HttpClientAuthenticationModel{
				NtlmDomain: types.StringValue(repo.HTTPClient.Authentication.NTLMDomain),
				NtlmHost:   types.StringValue(repo.HTTPClient.Authentication.NTLMHost),
				Type:       types.StringValue(string(repo.HTTPClient.Authentication.Type)),
				Username:   types.StringValue(repo.HTTPClient.Authentication.Username),
				Password:   types.StringValue(repo.HTTPClient.Authentication.Password),
			},
		}
	}
	if repo.HTTPClient.Connection != nil {
		data.HttpClient.Connection = &model.HttpClientConnectionModel{
			EnableCircularRedirects: types.BoolValue(GetValue(repo.HTTPClient.Connection.EnableCircularRedirects)),
			EnableCookies:           types.BoolValue(GetValue(repo.HTTPClient.Connection.EnableCookies)),
			Retries:                 types.Int64Value(int64(GetValue(repo.HTTPClient.Connection.Retries))),
			Timeout:                 types.Int64Value(int64(GetValue(repo.HTTPClient.Connection.Timeout))),
			UseTrustStore:           types.BoolValue(GetValue(repo.HTTPClient.Connection.UseTrustStore)),
			UserAgentSuffix:         types.StringValue(repo.HTTPClient.Connection.UserAgentSuffix),
		}
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
