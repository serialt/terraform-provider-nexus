package repository_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/serialt/terraform-provider-nexus/internal/acceptance"
)

var testAccDataSourceRepositoryListConfig = `data "nexus_repository_list" "acceptance" {}`

func TestAccDataSourceRepositoryList(t *testing.T) {
	dataSourceName := "data.nexus_repository_list.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRepositoryListConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet(dataSourceName, "items.0.name"),
						resource.TestCheckResourceAttrSet(dataSourceName, "items.0.format"),
						resource.TestCheckResourceAttrSet(dataSourceName, "items.0.type"),
						resource.TestCheckResourceAttrSet(dataSourceName, "items.0.url"),
					),
				),
			},
		},
	})
}
