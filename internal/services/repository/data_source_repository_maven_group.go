package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/serialt/terraform-provider-nexus/internal/schema/common"
	"github.com/serialt/terraform-provider-nexus/internal/schema/repository"
)

func DataSourceRepositoryMavenGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing maven group repository.",

		Read: dataSourceRepositoryMavenGroupRead,
		Schema: map[string]*schema.Schema{
			// Common schemas
			"id":     common.DataSourceID,
			"name":   repository.DataSourceName,
			"online": repository.DataSourceOnline,
			// Group schemas
			"group":   repository.DataSourceGroup,
			"storage": repository.DataSourceStorage,
		},
	}
}

func dataSourceRepositoryMavenGroupRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceMavenGroupRepositoryRead(resourceData, m)
}
