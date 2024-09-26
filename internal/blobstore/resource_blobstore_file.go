package blobstore

// import (
// 	"context"
// 	"fmt"

// 	"github.com/hashicorp/terraform-plugin-framework/resource"
// 	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
// 	"github.com/hashicorp/terraform-plugin-log/tflog"
// 	"github.com/nduyphuong/go-nexus-client/nexus3"
// 	"github.com/nduyphuong/go-nexus-client/nexus3/schema/blobstore"
// )

// // ResourceBlobstoreFile defines the resource implementation.
// type ResourceBlobstoreFile struct {
// 	client *nexus3.NexusClient
// }

// // Ensure provider defined types fully satisfy framework interfaces.
// var (
// 	_ resource.Resource                = &ResourceBlobstoreFile{}
// 	_ resource.ResourceWithImportState = &ResourceBlobstoreFile{}
// )

// func NewResourceBlobstoreFile() resource.Resource {
// 	return &ResourceBlobstoreFile{}
// }

// func (r *ResourceBlobstoreFile) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
// 	resp.TypeName = req.ProviderTypeName + "_blobstore_file"
// }

// func (r *ResourceBlobstoreFile) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
// 	resp.Schema = schema.Schema{
// 		MarkdownDescription: "Example data source",

// 		Attributes: map[string]schema.Attribute{
// 			"id": schema.StringAttribute{
// 				Description: "Used to identify data source at nexus",
// 				Computed:    true,
// 			},
// 			"name": schema.StringAttribute{
// 				Description: "Blobstore name",
// 				Required:    true,
// 			},
// 			"path": schema.StringAttribute{
// 				Description: "The path to the blobstore contents. This can be an absolute path to anywhere on the system nxrm has access to or it can be a path relative to the sonatype-work directory",
// 				Optional:    true,
// 			},
// 			"available_space_in_bytes": schema.Int64Attribute{
// 				Computed:    true,
// 				Description: "Available space in Bytes",
// 			},
// 			"total_size_in_bytes": schema.Int64Attribute{
// 				Computed:    true,
// 				Description: "The total size of the blobstore in Bytes",
// 			},
// 			"blob_count": schema.Int64Attribute{
// 				Computed:    true,
// 				Description: "Count of blobs",
// 			},
// 			"soft_quota": schema.SingleNestedAttribute{
// 				MarkdownDescription: "Soft quota of the blobstore",
// 				Computed:            true,
// 				Attributes: map[string]schema.Attribute{
// 					"limit": schema.Int64Attribute{
// 						Description: "The limit in Bytes. Minimum value is 1000000",
// 						Computed:    true,
// 					},
// 					"type": schema.StringAttribute{
// 						Description: "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
// 						Computed:    true,
// 					},
// 				},
// 			},
// 		},
// 	}
// }
// func (r *ResourceBlobstoreFile) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
// 	if req.ProviderData == nil {
// 		return
// 	}
// 	client, ok := req.ProviderData.(*nexus3.NexusClient)
// 	if !ok {
// 		resp.Diagnostics.AddError(
// 			"Unexpected Resource Configure Type",
// 			fmt.Sprintf("Expected Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
// 		)
// 		return
// 	}
// 	r.client = client
// }

// func (r *ResourceBlobstoreFile) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

// 	tflog.Debug(ctx, "Create BlobStore File resource")

// 	var data, newData BlobStoreFileSourceModel

// 	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
// 	if resp.Diagnostics.HasError() {
// 		return
// 	}
// 	if data.Path.IsNull() {
// 		data.Path = data.Name
// 	}

// 	err := r.client.BlobStore.File.Create(&blobstore.File{
// 		Name: data.Name.ValueString(),
// 		Path: data.Path.ValueString(),
// 		SoftQuota: &blobstore.SoftQuota{
// 			Type:  data.SoftQuota.Type.ValueString(),
// 			Limit: data.SoftQuota.Limit.ValueInt64(),
// 		},
// 	})
// 	if err != nil {
// 		resp.Diagnostics.AddError("Get projects msg from harbor failed", err.Error())
// 		return
// 	}

// 	tflog.Debug(ctx, "created a resource")
// 	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
// }

// func (r *ResourceBlobstoreFile) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
// }

// func (r *ResourceBlobstoreFile) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
// }

// func (r *ResourceBlobstoreFile) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
// }

// func (r *ResourceBlobstoreFile) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
