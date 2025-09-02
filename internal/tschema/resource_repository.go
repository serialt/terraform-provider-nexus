package tschema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	ResourceDockerHostedStorage = schema.SingleNestedBlock{
		Description:         "The storage configuration of the repository",
		MarkdownDescription: "The storage configuration of the repository",
		Attributes: map[string]schema.Attribute{
			"blob_store_name": schema.StringAttribute{
				Description:         "Blob store used to store repository contents",
				MarkdownDescription: "Blob store used to store repository contents",
				Required:            true,
			},
			"strict_content_type_validation": schema.BoolAttribute{
				Description:         "Whether to validate uploaded content's MIME type appropriate for the repository format",
				MarkdownDescription: "Whether to validate uploaded content's MIME type appropriate for the repository format",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"write_policy": schema.StringAttribute{
				Description:         "Controls if deployments of and updates to assets are allowed",
				MarkdownDescription: "Controls if deployments of and updates to assets are allowed",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("ALLOW"),
			},
			"latest_policy": schema.BoolAttribute{
				Description:         "Whether to allow redeploying the 'latest' tag but defer to the Deployment Policy for all other tags. Only usable with write_policy \"ALLOW_ONCE\"",
				MarkdownDescription: "Whether to allow redeploying the 'latest' tag but defer to the Deployment Policy for all other tags. Only usable with write_policy \"ALLOW_ONCE\"",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
	ResourceDocker = schema.SingleNestedBlock{
		Description:         "docker contains the configuration of the docker repository",
		MarkdownDescription: "docker contains the configuration of the docker repository",
		Attributes: map[string]schema.Attribute{
			"force_basic_auth": schema.BoolAttribute{
				Description:         "Whether to force authentication (Docker Bearer Token Realm required if false)",
				MarkdownDescription: "Whether to force authentication (Docker Bearer Token Realm required if false)",
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"http_port": schema.Int64Attribute{
				Description:         "Create an HTTP connector at specified port",
				MarkdownDescription: "Create an HTTP connector at specified port",
				Computed:            true,
			},
			"https_port": schema.Int64Attribute{
				Description:         "Create an HTTPS connector at specified port",
				MarkdownDescription: "Create an HTTPS connector at specified port",
				Computed:            true,
			},
			"v1_enabled": schema.BoolAttribute{
				Description:         "Whether to allow clients to use the V1 API to interact with this repository",
				MarkdownDescription: "Whether to allow clients to use the V1 API to interact with this repository",
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"subdomain": schema.StringAttribute{
				Description:         "Pro-only: Whether to allow clients to use subdomain routing connector",
				MarkdownDescription: "Pro-only: Whether to allow clients to use subdomain routing connector",
				Computed:            true,
			},
		},
	}
	ResourceCleanUp = schema.SingleNestedBlock{
		MarkdownDescription: "Cleanup policies",
		Description:         "Cleanup policies",
		Attributes: map[string]schema.Attribute{
			"policy_names": schema.ListAttribute{
				Description:         "List of policy names",
				MarkdownDescription: "List of policy names",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
	ResourceComponent = schema.SingleNestedBlock{
		Description:         "Component configuration for the hosted repository",
		MarkdownDescription: "Component configuration for the hosted repository",
		Attributes: map[string]schema.Attribute{
			"proprietary_components": schema.BoolAttribute{
				Description:         "Components in this repository count as proprietary for namespace conflict attacks (requires Sonatype Nexus Firewall)",
				MarkdownDescription: "Components in this repository count as proprietary for namespace conflict attacks (requires Sonatype Nexus Firewall)",
				Computed:            true,
			},
		},
	}
)
