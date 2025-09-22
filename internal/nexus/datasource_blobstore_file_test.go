package nexus

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDatasource_BlobStoreFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: TestNexusProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDatasourceBlobStoreFileDefaultConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nexus_blobstore_file.default", "name", "default"),
				),
			},
			{
				Config: testDatasourceBlobStoreFileLocalConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nexus_blobstore_file.local", "name", "local"),
					resource.TestCheckResourceAttr("data.nexus_blobstore_file.local", "soft_quota.type", "spaceUsedQuota"),
					resource.TestCheckResourceAttr("data.nexus_blobstore_file.local", "soft_quota.limit", "9"),
				),
			},
		},
	})
}

const testDatasourceBlobStoreFileDefaultConfig = `
data "nexus_blobstore_file" "default" {
  name = "default"
}
`

const testDatasourceBlobStoreFileLocalConfig = `
data "nexus_blobstore_file" "local" {
  name = "local"
}
`
