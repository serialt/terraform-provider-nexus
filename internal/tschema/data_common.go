package tschema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	DataSourceID = schema.StringAttribute{
		Description:         "Used to identify data source at nexus",
		MarkdownDescription: "Used to identify data source at nexus",
		Computed:            true,
	}

	DataSourceName = ResourceName

	DataSourceOnline = schema.BoolAttribute{
		Description:         "Whether this repository accepts incoming requests",
		MarkdownDescription: "Whether this repository accepts incoming requests",
		Computed:            true,
	}
	DataSourceRoutingRule = schema.StringAttribute{
		Description:         "The name of the routing rule assigned to this repository",
		MarkdownDescription: "The name of the routing rule assigned to this repository",
		Computed:            true,
	}
)
