package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nduyphuong/go-nexus-client/nexus3"
	"github.com/nduyphuong/go-nexus-client/nexus3/pkg/client"
	"github.com/serialt/terraform-provider-nexus/internal/blobstore"
)

var _ provider.Provider = &NexusProvider{}
var _ provider.ProviderWithFunctions = &NexusProvider{}

type NexusProvider struct {
}

type NexusProviderModel struct {
	Insecure types.Bool   `tfsdk:"insecure"`
	Password types.String `tfsdk:"password"`
	URL      types.String `tfsdk:"url"`
	Username types.String `tfsdk:"username"`
}

func (p *NexusProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nexus"
}

func (p *NexusProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "nexus url.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "nexus username.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"insecure": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (p *NexusProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config NexusProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}
	url := os.Getenv("NEXUS_URL")
	username := os.Getenv("NEXUS_USERNAME")
	password := os.Getenv("NEXUS_PASSWORD")
	insecureC := os.Getenv("NEXUS_INSECURE_SKIP_VERIFY")
	insecure := false
	if insecureC == "true" {
		insecure = true
	}
	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	}
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}
	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}
	client := nexus3.NewClient(client.Config{
		Insecure: insecure,
		Password: password,
		URL:      url,
		Username: username,
	})

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *NexusProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// NewExampleResource,
		blobstore.NewResourceBlobstoreFile,
	}
}

func (p *NexusProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		blobstore.NewBlobStoreFileSource,
		blobstore.NewBlobStoreListSource,
		blobstore.NewBlobStoreGroupSource,
		// blobstore.NewBlobStoreFileSource,
	}
}

func (p *NexusProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// NewExampleFunction,
	}
}

func New() provider.Provider {
	return &NexusProvider{}
}
