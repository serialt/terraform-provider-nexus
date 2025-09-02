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

var _ datasource.DataSource = &RepositoryAptProxyDatasource{}

func NewRepositoryAptProxyDatasource() datasource.DataSource {
	return &RepositoryAptProxyDatasource{}
}

type RepositoryAptProxyDatasource struct {
	client *nexus3.NexusClient
}

func (d *RepositoryAptProxyDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_apt_proxy"
}

func (d *RepositoryAptProxyDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Use this data source to get an existing apt repository.",
		MarkdownDescription: "Use this data source to get an existing apt repository.",
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
			"flat": schema.BoolAttribute{
				Description:         "Distribution to fetch",
				MarkdownDescription: "Distribution to fetch",
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

func (d *RepositoryAptProxyDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RepositoryAptProxyDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state model.RepositoryAptProxyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	if state.Name.IsUnknown() {
		resp.Diagnostics.AddError("Get apt proxy datasource failed", "name is unknown")
	}

	state, err := RepositoryAptProxyGetState(d.client, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Get apt proxy datasource failed", err.Error())
	}
	tflog.Trace(ctx, "read a apt proxy data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func RepositoryAptProxyGetState(client *nexus3.NexusClient, name string) (data model.RepositoryAptProxyModel, err error) {

	if name == "" {
		err = errors.New("name is nil")
		return
	}

	repo, err := client.Repository.Apt.Proxy.Get(name)
	if err != nil {
		return
	}

	data = model.RepositoryAptProxyModel{
		Id:           types.StringValue(repo.Name),
		Name:         types.StringValue(repo.Name),
		Online:       types.BoolValue(repo.Online),
		Flat:         types.BoolValue(repo.Apt.Flat),
		Distribution: types.StringValue(repo.Apt.Distribution),

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
		// Cleanup: &CleanupModel{
		// 	PolicyNames: []types.String{types.StringValue("")},
		// },
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

	return
}

func stringListValue(inputs []string) types.List {
	ret, _ := types.ListValueFrom(
		context.TODO(),
		types.StringType,
		inputs,
	)
	return ret
}

// 定义一个泛型函数，接受任何类型的指针并返回其值
func GetValue[T any](ptr *T) T {
	if ptr != nil {
		return *ptr // 解引用指针以获取值
	}
	var zero T // 如果指针为nil，则返回类型的零值
	return zero
}
