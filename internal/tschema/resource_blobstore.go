package tschema

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

var (
	RsoftQuota = schema.SingleNestedBlock{
		MarkdownDescription: "Soft quota of the blobstore",
		Attributes: map[string]schema.Attribute{
			"limit": schema.Int64Attribute{
				Description: "The limit in Bytes. Minimum value is 1000000",

				Required: true,
			},
			"type": schema.StringAttribute{
				Description: "The type to use such as spaceRemainingQuota, or spaceUsedQuota",
				Required:    true,
			},
		},
	}
)
