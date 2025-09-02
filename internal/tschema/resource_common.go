package tschema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
)

var (
	ResourceName = schema.StringAttribute{
		Description:         "A unique identifier for this repository",
		MarkdownDescription: "A unique identifier for this repository",
		Required:            true,
	}
	ResourceOnline = schema.BoolAttribute{
		Description:         "Whether this repository accepts incoming requests",
		MarkdownDescription: "Whether this repository accepts incoming requests",
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
	}
	ResourceID = schema.StringAttribute{
		Description: "Used to identify data source at nexus",
		Computed:    true,
	}
)
