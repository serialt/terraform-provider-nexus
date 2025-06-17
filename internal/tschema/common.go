package tschema

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

var (
	DataSourceID = schema.StringAttribute{
		Description:         "Used to identify data source at nexus",
		MarkdownDescription: "Used to identify data source at nexus",
		Computed:            true,
	}
	ResourceID = schema.StringAttribute{
		Description: "Used to identify data source at nexus",
		Computed:    true,
	}

	ResourceName = schema.StringAttribute{
		Description:         "A unique identifier for this repository",
		MarkdownDescription: "A unique identifier for this repository",
		Required:            true,
	}
	DataSourceName = ResourceName

	ResourceOnline = schema.BoolAttribute{
		Description:         "Whether this repository accepts incoming requests",
		MarkdownDescription: "Whether this repository accepts incoming requests",
		Computed:            true,
	}
	DataSourceOnline = schema.BoolAttribute{
		Description:         "Whether this repository accepts incoming requests",
		MarkdownDescription: "Whether this repository accepts incoming requests",
		Computed:            true,
	}
)
