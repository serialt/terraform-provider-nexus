package security_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/serialt/terraform-provider-nexus/internal/acceptance"
)

func TestAccDataSourceSecurityRealms(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecurityRealmsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.nexus_security_realms.acceptance", "active.#"),
					resource.TestCheckResourceAttrSet("data.nexus_security_realms.acceptance", "available.#"),
				),
			},
		},
	})

}

func testAccDataSourceSecurityRealmsConfig() string {
	return `
data "nexus_security_realms" "acceptance" {}
`
}
