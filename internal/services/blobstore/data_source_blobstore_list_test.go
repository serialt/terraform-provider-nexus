package blobstore_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/serialt/terraform-provider-nexus/internal/acceptance"
)

var testAccDataSourceBlobStoreListConfig = `data "nexus_blobstore_list" "acceptance" {}`

func TestAccDataSourceBlobstoreList(t *testing.T) {
	dataSourceName := "data.nexus_blobstore_list.acceptance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlobStoreListConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet(dataSourceName, "items.0.name"),
						resource.TestCheckResourceAttrSet(dataSourceName, "items.0.type"),
					),
				),
			},
		},
	})
}
