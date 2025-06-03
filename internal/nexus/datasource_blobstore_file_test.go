package nexus

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestDatasourceBlobStoreFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestNexusProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDatasourceBlobStoreFileConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.scaffolding_example.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("example-id"),
					),
				},
			},
		},
	})
}

const testDatasourceBlobStoreFileConfig = `
data "nexus_blobstore_file" "test" {
  name = "test"
}
`
