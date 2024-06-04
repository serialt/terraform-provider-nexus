package repository

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/serialt/terraform-provider-nexus/internal/schema/common"
	"github.com/serialt/terraform-provider-nexus/internal/schema/repository"
)

func DataSourceRepositoryGoGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing go group repository.",

		Read: dataSourceRepositoryGoGroupRead,
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

func dataSourceRepositoryGoGroupRead(resourceData *schema.ResourceData, m interface{}) error {
	resourceData.SetId(resourceData.Get("name").(string))

	return resourceGoGroupRepositoryRead(resourceData, m)
}
