package tschema

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	DSsoftQuota = schema.SingleNestedBlock{
		MarkdownDescription: "Soft quota of the blobstore",
		Attributes: map[string]schema.Attribute{
			"limit": schema.Int64Attribute{
				Description: "The limit in Bytes. Minimum value is 1000000",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
				Computed:    true,
			},
		}}

	DSNegativeCache = schema.SingleNestedBlock{
		Description:         "Configuration of the negative cache handling",
		MarkdownDescription: "Configuration of the negative cache handling",
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description:         "Whether to cache responses for content not present in the proxied repository",
				MarkdownDescription: "Whether to cache responses for content not present in the proxied repository",
				Computed:            true,
			},
			"ttl": schema.Int64Attribute{
				Description:         "How long to cache the fact that a file was not found in the repository (in minutes)",
				MarkdownDescription: "How long to cache the fact that a file was not found in the repository (in minutes)",
				Computed:            true,
			},
		},
	}
	DSStorage = schema.SingleNestedBlock{
		Description:         "The storage configuration of the repository",
		MarkdownDescription: "The storage configuration of the repository",
		Attributes: map[string]schema.Attribute{
			"blob_store_name": schema.StringAttribute{
				Description:         "Blob store used to store repository contents",
				MarkdownDescription: "Blob store used to store repository contents",
				Computed:            true,
			},
			"strict_content_type_validation": schema.BoolAttribute{
				Description:         "Whether to validate uploaded content's MIME type appropriate for the repository format",
				MarkdownDescription: "Whether to validate uploaded content's MIME type appropriate for the repository format",
				Computed:            true,
			},
		},
	}
)
