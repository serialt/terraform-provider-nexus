package blobstore

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nduyphuong/go-nexus-client/nexus3"
)

var _ datasource.DataSource = &BlobStoreS3Source{}

func NewBlobStoreS3Source() datasource.DataSource {
	return &BlobStoreS3Source{}
}

type BlobStoreS3Source struct {
	client *nexus3.NexusClient
}

type BlobStoreS3SourceModel struct {
	Id                    types.String              `tfsdk:"id"`
	Name                  types.String              `tfsdk:"name"`
	BlobCount             types.Int64               `tfsdk:"blob_count"`
	AvailableSpaceInBytes types.Int64               `tfsdk:"available_space_in_bytes"`
	TotalSizeInBytes      types.Int64               `tfsdk:"total_size_in_bytes"`
	SoftQuota             *SoftQuotaModel           `tfsdk:"soft_quota"`
	BucketConfiguration   *bucketConfigurationModel `tfsdk:"bucket_configuration"`
}

type advancedBucketConnectionModel struct {
	Endpoint       types.String `tfsdk:"endpoint"`
	ForcePathStyle types.Bool   `tfsdk:"force_path_style"`
	SignerType     types.String `tfsdk:"signer_type"`
}

type bucketModel struct {
	Expiration types.Int64  `tfsdk:"expiration"`
	Name       types.String `tfsdk:"name"`
	Prefix     types.String `tfsdk:"prefix"`
	region     types.String `tfsdk:"region"`
}

type bucketSecurityModel struct {
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	Role            types.String `tfsdk:"role"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	SessionToken    types.String `tfsdk:"session_token"`
}

type encryptionModel struct {
	EncryptionKey  types.String `tfsdk:"encryption_key"`
	EncryptionType types.String `tfsdk:"encryption_type"`
}
type bucketConfigurationModel struct {
	AdvancedBucketConnection *advancedBucketConnectionModel `tfsdk:"advanced_bucket_connection"`
	Bucket                   *bucketModel                   `tfsdk:"bucket"`
	BucketSecurity           *bucketSecurityModel           `tfsdk:"bucket_security"`
	Encryption               *encryptionModel               `tfsdk:"encryption"`
}

func (d *BlobStoreS3Source) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blobstore_file"
}

func (d *BlobStoreS3Source) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Example data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Used to identify data source at nexus",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Blobstore name",
				Required:    true,
			},
			"total_size_in_bytes": schema.Int64Attribute{
				Computed:    true,
				Description: "The total size of the blobstore in Bytes",
			},
			"blob_count": schema.Int64Attribute{
				Computed:    true,
				Description: "Count of blobs",
			},
			"soft_quota": schema.SingleNestedAttribute{
				MarkdownDescription: "Soft quota of the blobstore",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"limit": schema.Int64Attribute{
						Description: "The limit in Bytes. Minimum value is 1000000",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
						Computed:    true,
					},
				},
			},
			"bucket_configuration": schema.SingleNestedAttribute{
				MarkdownDescription: "The S3 bucket configuration.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"advanced_bucket_connection": schema.SingleNestedAttribute{
						Description: "Additional connection configurations",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"endpoint": schema.StringAttribute{
								Description: "A custom endpoint URL for third party object stores using the S3 API.",
								Computed:    true,
							},
							"force_path_style": schema.BoolAttribute{
								Description: "Setting this flag will result in path-style access being used for all requests.",
								Computed:    true,
							},
							"signer_type": schema.StringAttribute{
								Description: "An API signature version which may be required for third party object stores using the S3 API.",
								Computed:    true,
							},
						},
					},
					"bucket": schema.SingleNestedAttribute{
						Description: "The S3 bucket configuration",
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"region": schema.StringAttribute{
								Description: "The AWS region to create a new S3 bucket in or an existing S3 bucket's region",
								Computed:    true,
							},
							"name": schema.StringAttribute{
								Description: "The name of the S3 bucket",
								Computed:    true,
							},
							"prefix": schema.StringAttribute{
								Description: "The S3 blob store (i.e S3 object) key prefix",
								Computed:    true,
							},
							"expiration": schema.Int64Attribute{
								Description: "How many days until deleted blobs are finally removed from the S3 bucket (-1 to disable)",
								Computed:    true,
							},
						},
					},
					"bucket_security": schema.SingleNestedAttribute{
						Description: "Additional security configurations",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"access_key_id": schema.StringAttribute{
								Description: "An IAM access key ID for granting access to the S3 bucket",
								Computed:    true,
							},
							"secret_access_key": schema.StringAttribute{
								Description: "The secret access key associated with the specified IAM access key ID",
								Computed:    true,
								Sensitive:   true,
							},
							"role": schema.StringAttribute{
								Description: "An IAM role to assume in order to access the S3 bucket",
								Computed:    true,
							},
							"session_token": schema.StringAttribute{
								Description: "An AWS STS session token associated with temporary security credentials which grant access to the S3 bucket",
								Computed:    true,
								Sensitive:   true,
							},
						},
					},
					"encryption": schema.SingleNestedAttribute{
						Description: "Additional bucket encryption configurations",
						Computed:    true,

						Attributes: map[string]schema.Attribute{
							"encryption_key": schema.StringAttribute{
								Description: "The encryption key.",
								Computed:    true,
							},
							"encryption_type": schema.StringAttribute{
								Description: "The type of S3 server side encryption to use.",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

func (d *BlobStoreS3Source) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *BlobStoreS3Source) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state, newState BlobStoreS3SourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// BlobStoreS3, err := d.client.BlobStore.File.Get(state.Name.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError("Get projects msg from harbor failed", err.Error())
	// 	return
	// }

	// var genericBlobstoreInformation blobstore.Generic
	// genericBlobstores, err := d.client.BlobStore.List()
	// if err != nil {
	// 	resp.Diagnostics.AddError("Get projects msg from harbor failed", err.Error())
	// 	return
	// }
	// for _, generic := range genericBlobstores {
	// 	if generic.Name == BlobStoreS3.Name {
	// 		genericBlobstoreInformation = generic
	// 	}
	// }

	// newState = BlobStoreS3SourceModel{
	// 	Id:                    types.StringValue(BlobStoreS3.Name),
	// 	Name:                  types.StringValue(BlobStoreS3.Name),
	// 	Path:                  types.StringValue(BlobStoreS3.Path),
	// 	AvailableSpaceInBytes: types.Int64Value(int64(genericBlobstoreInformation.AvailableSpaceInBytes)),
	// 	TotalSizeInBytes:      types.Int64Value(int64(genericBlobstoreInformation.TotalSizeInBytes)),
	// 	BlobCount:             types.Int64Value(int64(genericBlobstoreInformation.BlobCount)),
	// 	SoftQuota: SoftQuota{
	// 		Limit: types.Int64Value(BlobStoreS3.SoftQuota.Limit),
	// 		Type:  types.StringValue(BlobStoreS3.SoftQuota.Type),
	// 	},
	// }
	tflog.Trace(ctx, "read a BlobStoreS3 data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
