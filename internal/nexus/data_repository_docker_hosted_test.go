package nexus

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	// "github.com/hashicorp/terraform-plugin-testing/knownvalue"
	// "github.com/hashicorp/terraform-plugin-testing/statecheck"
	// "github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestDatasource_RepositoryDockerHosted(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: TestNexusProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDatasourceDockerHostConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.nexus_repository_docker_hosted.docker-local", "name", "docker-local"),
				),
			},
		},
	})
}

const testDatasourceDockerHostConfig = `

data "nexus_repository_docker_hosted" "docker-local" {
  name = "docker-local"
}
`
