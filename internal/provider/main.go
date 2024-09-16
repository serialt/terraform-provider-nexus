package provider

import (
	nexus "github.com/datadrivers/go-nexus-client/nexus3"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/serialt/terraform-provider-nexus/internal/services/blobstore"
	"github.com/serialt/terraform-provider-nexus/internal/services/other"
	"github.com/serialt/terraform-provider-nexus/internal/services/repository"
	"github.com/serialt/terraform-provider-nexus/internal/services/security"
)

// Provider returns a terraform.Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"nexus_blobstore_azure":                       blobstore.DataSourceBlobstoreAzure(),
			"nexus_blobstore_file":                        blobstore.DataSourceBlobstoreFile(),
			"nexus_blobstore_group":                       blobstore.DataSourceBlobstoreGroup(),
			"nexus_blobstore_s3":                          blobstore.DataSourceBlobstoreS3(),
			"nexus_blobstore_list":                        blobstore.DataSourceBlobstoreList(),
			"nexus_repository_apt_hosted":                 repository.DataSourceRepositoryAptHosted(),
			"nexus_repository_apt_proxy":                  repository.DataSourceRepositoryAptProxy(),
			"nexus_repository_bower_group":                repository.DataSourceRepositoryBowerGroup(),
			"nexus_repository_bower_hosted":               repository.DataSourceRepositoryBowerHosted(),
			"nexus_repository_bower_proxy":                repository.DataSourceRepositoryBowerProxy(),
			"nexus_repository_cocoapods_proxy":            repository.DataSourceRepositoryCocoapodsProxy(),
			"nexus_repository_conan_proxy":                repository.DataSourceRepositoryConanProxy(),
			"nexus_repository_conda_proxy":                repository.DataSourceRepositoryCondaProxy(),
			"nexus_repository_docker_group":               repository.DataSourceRepositoryDockerGroup(),
			"nexus_repository_docker_hosted":              repository.DataSourceRepositoryDockerHosted(),
			"nexus_repository_docker_proxy":               repository.DataSourceRepositoryDockerProxy(),
			"nexus_repository_gitlfs_hosted":              repository.DataSourceRepositoryGitlfsHosted(),
			"nexus_repository_go_group":                   repository.DataSourceRepositoryGoGroup(),
			"nexus_repository_go_proxy":                   repository.DataSourceRepositoryGoProxy(),
			"nexus_repository_helm_hosted":                repository.DataSourceRepositoryHelmHosted(),
			"nexus_repository_helm_proxy":                 repository.DataSourceRepositoryHelmProxy(),
			"nexus_repository_list":                       repository.DataSourceRepositoryList(),
			"nexus_repository_maven_group":                repository.DataSourceRepositoryMavenGroup(),
			"nexus_repository_maven_hosted":               repository.DataSourceRepositoryMavenHosted(),
			"nexus_repository_maven_proxy":                repository.DataSourceRepositoryMavenProxy(),
			"nexus_repository_npm_group":                  repository.DataSourceRepositoryNpmGroup(),
			"nexus_repository_npm_hosted":                 repository.DataSourceRepositoryNpmHosted(),
			"nexus_repository_npm_proxy":                  repository.DataSourceRepositoryNpmProxy(),
			"nexus_repository_nuget_group":                repository.DataSourceRepositoryNugetGroup(),
			"nexus_repository_nuget_hosted":               repository.DataSourceRepositoryNugetHosted(),
			"nexus_repository_nuget_proxy":                repository.DataSourceRepositoryNugetProxy(),
			"nexus_repository_p2_proxy":                   repository.DataSourceRepositoryP2Proxy(),
			"nexus_repository_pypi_group":                 repository.DataSourceRepositoryPypiGroup(),
			"nexus_repository_pypi_hosted":                repository.DataSourceRepositoryPypiHosted(),
			"nexus_repository_pypi_proxy":                 repository.DataSourceRepositoryPypiProxy(),
			"nexus_repository_r_group":                    repository.DataSourceRepositoryRGroup(),
			"nexus_repository_r_hosted":                   repository.DataSourceRepositoryRHosted(),
			"nexus_repository_r_proxy":                    repository.DataSourceRepositoryRProxy(),
			"nexus_repository_raw_group":                  repository.DataSourceRepositoryRawGroup(),
			"nexus_repository_raw_hosted":                 repository.DataSourceRepositoryRawHosted(),
			"nexus_repository_raw_proxy":                  repository.DataSourceRepositoryRawProxy(),
			"nexus_repository_rubygems_group":             repository.DataSourceRepositoryRubygemsGroup(),
			"nexus_repository_rubygems_hosted":            repository.DataSourceRepositoryRubygemsHosted(),
			"nexus_repository_rubygems_proxy":             repository.DataSourceRepositoryRubygemsProxy(),
			"nexus_repository_yum_group":                  repository.DataSourceRepositoryYumGroup(),
			"nexus_repository_yum_hosted":                 repository.DataSourceRepositoryYumHosted(),
			"nexus_repository_yum_proxy":                  repository.DataSourceRepositoryYumProxy(),
			"nexus_routing_rule":                          other.DataSourceRoutingRule(),
			"nexus_security_anonymous":                    security.DataSourceSecurityAnonymous(),
			"nexus_security_content_selector":             security.DataSourceSecurityContentSelector(),
			"nexus_security_ldap":                         security.DataSourceSecurityLDAP(),
			"nexus_security_realms":                       security.DataSourceSecurityRealms(),
			"nexus_security_role":                         security.DataSourceSecurityRole(),
			"nexus_security_saml":                         security.DataSourceSecuritySAML(),
			"nexus_security_ssl":                          security.DataSourceSecuritySSL(),
			"nexus_security_ssl_truststore":               security.DataSourceSecuritySSLTrustStore(),
			"nexus_security_user":                         security.DataSourceSecurityUser(),
			"nexus_security_user_token":                   security.DataSourceSecurityUserToken(),
			"nexus_mail_config":                           other.DataSourceMailConfig(),
			"nexus_privilege_script":                      security.DataSourceSecurityPrivilegeScript(),
			"nexus_privilege_wildcard":                    security.DataSourceSecurityPrivilegeWildcard(),
			"nexus_privilege_application":                 security.DataSourceSecurityPrivilegeApplication(),
			"nexus_privilege_repository_view":             security.DataSourceSecurityPrivilegeRepositoryView(),
			"nexus_privilege_repository_admin":            security.DataSourceSecurityPrivilegeRepositoryAdmin(),
			"nexus_privilege_repository_content_selector": security.DataSourceSecurityPrivilegeRepositoryContentSelector(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"nexus_blobstore_azure":                       blobstore.ResourceBlobstoreAzure(),
			"nexus_blobstore_file":                        blobstore.ResourceBlobstoreFile(),
			"nexus_blobstore_group":                       blobstore.ResourceBlobstoreGroup(),
			"nexus_blobstore_s3":                          blobstore.ResourceBlobstoreS3(),
			"nexus_repository_apt_hosted":                 repository.ResourceRepositoryAptHosted(),
			"nexus_repository_apt_proxy":                  repository.ResourceRepositoryAptProxy(),
			"nexus_repository_bower_group":                repository.ResourceRepositoryBowerGroup(),
			"nexus_repository_bower_hosted":               repository.ResourceRepositoryBowerHosted(),
			"nexus_repository_bower_proxy":                repository.ResourceRepositoryBowerProxy(),
			"nexus_repository_cocoapods_proxy":            repository.ResourceRepositoryCocoapodsProxy(),
			"nexus_repository_conan_proxy":                repository.ResourceRepositoryConanProxy(),
			"nexus_repository_conda_proxy":                repository.ResourceRepositoryCondaProxy(),
			"nexus_repository_docker_group":               repository.ResourceRepositoryDockerGroup(),
			"nexus_repository_docker_hosted":              repository.ResourceRepositoryDockerHosted(),
			"nexus_repository_docker_proxy":               repository.ResourceRepositoryDockerProxy(),
			"nexus_repository_gitlfs_hosted":              repository.ResourceRepositoryGitlfsHosted(),
			"nexus_repository_go_group":                   repository.ResourceRepositoryGoGroup(),
			"nexus_repository_go_proxy":                   repository.ResourceRepositoryGoProxy(),
			"nexus_repository_helm_hosted":                repository.ResourceRepositoryHelmHosted(),
			"nexus_repository_helm_proxy":                 repository.ResourceRepositoryHelmProxy(),
			"nexus_repository_maven_group":                repository.ResourceRepositoryMavenGroup(),
			"nexus_repository_maven_hosted":               repository.ResourceRepositoryMavenHosted(),
			"nexus_repository_maven_proxy":                repository.ResourceRepositoryMavenProxy(),
			"nexus_repository_npm_group":                  repository.ResourceRepositoryNpmGroup(),
			"nexus_repository_npm_hosted":                 repository.ResourceRepositoryNpmHosted(),
			"nexus_repository_npm_proxy":                  repository.ResourceRepositoryNpmProxy(),
			"nexus_repository_nuget_group":                repository.ResourceRepositoryNugetGroup(),
			"nexus_repository_nuget_hosted":               repository.ResourceRepositoryNugetHosted(),
			"nexus_repository_nuget_proxy":                repository.ResourceRepositoryNugetProxy(),
			"nexus_repository_p2_proxy":                   repository.ResourceRepositoryP2Proxy(),
			"nexus_repository_pypi_group":                 repository.ResourceRepositoryPypiGroup(),
			"nexus_repository_pypi_hosted":                repository.ResourceRepositoryPypiHosted(),
			"nexus_repository_pypi_proxy":                 repository.ResourceRepositoryPypiProxy(),
			"nexus_repository_r_group":                    repository.ResourceRepositoryRGroup(),
			"nexus_repository_r_hosted":                   repository.ResourceRepositoryRHosted(),
			"nexus_repository_r_proxy":                    repository.ResourceRepositoryRProxy(),
			"nexus_repository_raw_group":                  repository.ResourceRepositoryRawGroup(),
			"nexus_repository_raw_hosted":                 repository.ResourceRepositoryRawHosted(),
			"nexus_repository_raw_proxy":                  repository.ResourceRepositoryRawProxy(),
			"nexus_repository_rubygems_group":             repository.ResourceRepositoryRubygemsGroup(),
			"nexus_repository_rubygems_hosted":            repository.ResourceRepositoryRubygemsHosted(),
			"nexus_repository_rubygems_proxy":             repository.ResourceRepositoryRubygemsProxy(),
			"nexus_repository_yum_group":                  repository.ResourceRepositoryYumGroup(),
			"nexus_repository_yum_hosted":                 repository.ResourceRepositoryYumHosted(),
			"nexus_repository_yum_proxy":                  repository.ResourceRepositoryYumProxy(),
			"nexus_routing_rule":                          other.ResourceRoutingRule(),
			"nexus_script":                                other.ResourceScript(),
			"nexus_cleanup_policy":                        other.ResourceCleanUpPolicy(),
			"nexus_security_anonymous":                    security.ResourceSecurityAnonymous(),
			"nexus_security_content_selector":             security.ResourceSecurityContentSelector(),
			"nexus_security_ldap":                         security.ResourceSecurityLDAP(),
			"nexus_security_ldap_order":                   security.ResourceSecurityLDAPOrder(),
			"nexus_security_realms":                       security.ResourceSecurityRealms(),
			"nexus_security_role":                         security.ResourceSecurityRole(),
			"nexus_security_ssl_truststore":               security.ResourceSecuritySSLTruststore(),
			"nexus_security_saml":                         security.ResourceSecuritySAML(),
			"nexus_security_user":                         security.ResourceSecurityUser(),
			"nexus_security_user_token":                   security.ResourceSecurityUserToken(),
			"nexus_mail_config":                           other.ResourceMailConfig(),
			"nexus_privilege_application":                 security.ResourceSecurityPrivilegeApplication(),
			"nexus_privilege_script":                      security.ResourceSecurityPrivilegeScript(),
			"nexus_privilege_repository_view":             security.ResourceSecurityPrivilegeRepositoryView(),
			"nexus_privilege_repository_admin":            security.ResourceSecurityPrivilegeRepositoryAdmin(),
			"nexus_privilege_repository_content_selector": security.ResourceSecurityPrivilegeRepositoryContentSelector(),
			"nexus_privilege_wildcard":                    security.ResourceSecurityPrivilegeWildcard(),
		},
		Schema: map[string]*schema.Schema{
			"insecure": {
				Description: "Boolean to specify wether insecure SSL connections are allowed or not. Reading environment variable NEXUS_INSECURE_SKIP_VERIFY. Default:`true`",
				Default:     false,
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_INSECURE_SKIP_VERIFY", "true"),
				Optional:    true,
				Type:        schema.TypeBool,
			},
			"password": {
				Description: "Password of user to connect to API. Reading environment variable NEXUS_PASSWORD. Default:`admin123`",
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_PASSWORD", "admin123"),
				Optional:    true,
				Type:        schema.TypeString,
			},
			"url": {
				Description: "URL of Nexus to reach API. Reading environment variable NEXUS_URL. Default:`http://127.0.0.1:8080`",
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_URL", "http://127.0.0.1:8080"),
				Optional:    true,
				Type:        schema.TypeString,
			},
			"username": {
				Description: "Username used to connect to API. Reading environment variable NEXUS_USERNAME. Default:`admin`",
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_USERNAME", "admin"),
				Optional:    true,
				Type:        schema.TypeString,
			},
			"timeout": {
				Description: "Timeout in seconds to connect to API. Reading environment variable NEXUS_TIMEOUT. Default:`30`",
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_TIMEOUT", 30),
				Optional:    true,
				Type:        schema.TypeInt,
			},
			"client_cert_path": {
				Description: "Path to a client PEM certificate to load for mTLS. Reading environment variable NEXUS_CLIENT_CERT_PATH. Default:``",
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_CLIENT_CERT_PATH", ""),
				Optional:    true,
				Type:        schema.TypeString,
			},
			"client_key_path": {
				Description: "Path to a client PEM key to load for mTLS. Reading environment variable NEXUS_CLIENT_KEY_PATH. Default:``",
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_CLIENT_KEY_PATH", ""),
				Optional:    true,
				Type:        schema.TypeString,
			},
			"root_ca_path": {
				Description: "Path to a root CA certificate to load for mTLS. Reading environment variable NEXUS_ROOT_CA_PATH. Default:``",
				DefaultFunc: schema.EnvDefaultFunc("NEXUS_ROOT_CA_PATH", ""),
				Optional:    true,
				Type:        schema.TypeString,
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	timeout := d.Get("timeout").(int)
	clientCertPath := d.Get("client_cert_path").(string)
	clientKeyPath := d.Get("client_key_path").(string)
	rootCaPath := d.Get("root_ca_path").(string)
	config := client.Config{
		Insecure:              d.Get("insecure").(bool),
		Password:              d.Get("password").(string),
		URL:                   d.Get("url").(string),
		Username:              d.Get("username").(string),
		Timeout:               &timeout,
		ClientCertificatePath: &clientCertPath,
		ClientKeyPath:         &clientKeyPath,
		RootCAPath:            &rootCaPath,
	}

	return nexus.NewClient(config), nil
}
