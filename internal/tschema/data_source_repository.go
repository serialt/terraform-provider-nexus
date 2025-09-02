package tschema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	DSHttpClient = schema.SingleNestedBlock{
		Description:         "HTTP Client configuration for proxy repositories",
		MarkdownDescription: "HTTP Client configuration for proxy repositories",
		Attributes: map[string]schema.Attribute{
			"authentication": schema.SingleNestedAttribute{
				Description:         "Authentication configuration of the HTTP client",
				MarkdownDescription: "Authentication configuration of the HTTP client",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description:         "Authentication type. Possible values: `ntlm` or `username`",
						MarkdownDescription: "Authentication type. Possible values: `ntlm` or `username`",
						Computed:            true,
					},
					"username": schema.StringAttribute{
						Description:         "The username used by the proxy repository",
						MarkdownDescription: "The username used by the proxy repository",
						Computed:            true,
					},
					"password": schema.StringAttribute{
						Description:         "The password used by the proxy repository",
						MarkdownDescription: "The password used by the proxy repository",
						Computed:            true,
						// Sensitive:           true,
					},
					"ntlm_domain": schema.StringAttribute{
						Description:         "The ntlm domain to connect",
						MarkdownDescription: "The ntlm domain to connect",
						Computed:            true,
					},
					"ntlm_host": schema.StringAttribute{
						Description:         "The ntlm host to connect",
						MarkdownDescription: "The ntlm host to connect",
						Computed:            true,
					},
				},
			},
			"connection": schema.SingleNestedAttribute{
				Description:         "Connection configuration of the HTTP client",
				MarkdownDescription: "Connection configuration of the HTTP client",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"enable_circular_redirects": schema.BoolAttribute{
						Description:         "Whether to enable redirects to the same location (may be required by some servers)",
						MarkdownDescription: "Whether to enable redirects to the same location (may be required by some servers)",
						Computed:            true,
					},
					"enable_cookies": schema.BoolAttribute{
						Description:         "Whether to allow cookies to be stored and used",
						MarkdownDescription: "Whether to allow cookies to be stored and used",
						Computed:            true,
					},
					"retries": schema.Int64Attribute{
						Description:         "Total retries if the initial connection attempt suffers a timeout",
						MarkdownDescription: "Total retries if the initial connection attempt suffers a timeout",
						Computed:            true,
					},
					"timeout": schema.Int64Attribute{
						Description:         "Seconds to wait for activity before stopping and retrying the connection",
						MarkdownDescription: "Seconds to wait for activity before stopping and retrying the connection",
						Computed:            true,
					},
					"user_agent_suffix": schema.StringAttribute{
						Description:         "Custom fragment to append to User-Agent header in HTTP requests",
						MarkdownDescription: "Custom fragment to append to User-Agent header in HTTP requests",
						Computed:            true,
					},
					"use_trust_store": schema.BoolAttribute{
						Description:         "Use certificates stored in the Nexus Repository Manager truststore to connect to external systems",
						MarkdownDescription: "Use certificates stored in the Nexus Repository Manager truststore to connect to external systems",
						Computed:            true,
					},
				},
			},
			"auto_block": schema.BoolAttribute{
				Description:         "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
				MarkdownDescription: "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
				Computed:            true,
			},
			"blocked": schema.BoolAttribute{
				Description:         "Whether to block outbound connections on the repository",
				MarkdownDescription: "Whether to block outbound connections on the repository",
				Computed:            true,
			},
		},
	}
	DSProxy = schema.SingleNestedBlock{
		Description:         "Configuration for the proxy repository",
		MarkdownDescription: "Configuration for the proxy repository",
		Attributes: map[string]schema.Attribute{
			"content_max_age": schema.Int64Attribute{
				Description:         "How long (in minutes) to cache artifacts before rechecking the remote repository",
				MarkdownDescription: "How long (in minutes) to cache artifacts before rechecking the remote repository",
				Computed:            true,
			},
			"metadata_max_age": schema.Int64Attribute{
				Description:         "How long (in minutes) to cache metadata before rechecking the remote repository.",
				MarkdownDescription: "How long (in minutes) to cache metadata before rechecking the remote repository.",
				Computed:            true,
			},
			"remote_url": schema.StringAttribute{
				Description:         "Location of the remote repository being proxied",
				MarkdownDescription: "Location of the remote repository being proxied",
				Computed:            true,
			},
		},
	}
	DSCleanUp = schema.SingleNestedBlock{
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
	DSComponent = schema.SingleNestedBlock{
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

	DSSigning = schema.SingleNestedBlock{
		Description:         "Contains signing data of repositores",
		MarkdownDescription: "Contains signing data of repositores",
		Attributes: map[string]schema.Attribute{
			"keypair": schema.StringAttribute{
				Description: `PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor)
							If passphrase is unset, the keypair cannot be read from the nexus api.
							When reading the resource, the keypair will be read from the previous state,
							so external changes won't be detected in this case.`,
				MarkdownDescription: `PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor)
							If passphrase is unset, the keypair cannot be read from the nexus api.
							When reading the resource, the keypair will be read from the previous state,
							so external changes won't be detected in this case.`,
				Computed:  true,
				Sensitive: true,
			},
			"passphrase": schema.StringAttribute{
				Description: `Passphrase to access PGP signing key.
							This value cannot be read from the nexus api.
							When reading the resource, the value will be read from the previous state,
							so external changes won't be detected.`,
				MarkdownDescription: `Passphrase to access PGP signing key.
							This value cannot be read from the nexus api.
							When reading the resource, the value will be read from the previous state,
							so external changes won't be detected.`,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
	DSDockerHostedStorage = schema.SingleNestedBlock{
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
			"write_policy": schema.StringAttribute{
				Description:         "Controls if deployments of and updates to assets are allowed",
				MarkdownDescription: "Controls if deployments of and updates to assets are allowed",
				Computed:            true,
			},
			"latest_policy": schema.BoolAttribute{
				Description:         "Whether to allow redeploying the 'latest' tag but defer to the Deployment Policy for all other tags. Only usable with write_policy \"ALLOW_ONCE\"",
				MarkdownDescription: "Whether to allow redeploying the 'latest' tag but defer to the Deployment Policy for all other tags. Only usable with write_policy \"ALLOW_ONCE\"",
				Computed:            true,
			},
		},
	}
	DSGroup = schema.SingleNestedBlock{
		Description:         "The group configuration of the repository",
		MarkdownDescription: "The group configuration of the repository",
		Attributes: map[string]schema.Attribute{
			"member_names": schema.ListAttribute{
				Description:         "Member repositories names",
				MarkdownDescription: "Member repositories names",
				ElementType:         types.StringType,
				Required:            true,
				Validators:          []validator.List{listvalidator.SizeAtLeast(1)},
			},
		},
	}
	DSDocker = schema.SingleNestedBlock{
		Description:         "docker contains the configuration of the docker repository",
		MarkdownDescription: "docker contains the configuration of the docker repository",
		Attributes: map[string]schema.Attribute{
			"force_basic_auth": schema.BoolAttribute{
				Description:         "Whether to force authentication (Docker Bearer Token Realm required if false)",
				MarkdownDescription: "Whether to force authentication (Docker Bearer Token Realm required if false)",
				Computed:            true,
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
			},
			"subdomain": schema.StringAttribute{
				Description:         "Pro-only: Whether to allow clients to use subdomain routing connector",
				MarkdownDescription: "Pro-only: Whether to allow clients to use subdomain routing connector",
				Computed:            true,
			},
		},
	}
	DSHostedStorage = schema.SingleNestedBlock{
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
			"write_policy": schema.StringAttribute{
				Description:         "Controls if deployments of and updates to assets are allowed",
				MarkdownDescription: "Controls if deployments of and updates to assets are allowed",
				Computed:            true,
			},
		},
	}
)
