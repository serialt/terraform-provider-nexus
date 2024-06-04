package other_test

import (
	"fmt"
	"testing"

	"github.com/serialt/terraform-provider-nexus/internal/services/other"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/serialt/terraform-provider-nexus/internal/acceptance"
)

func TestAccResourceCleanUpPolicy(t *testing.T) {
	resName := "nexus_cleanup_policy.acceptance"

	r := other.CleanUpPolicy{
		Name:   "acc-test",
		Format: "docker",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCleanUpPolicyConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", r.Name),
					resource.TestCheckResourceAttr(resName, "format", r.Format),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateId:     r.Name,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceCleanUpPolicyConfig(c other.CleanUpPolicy) string {
	return fmt.Sprintf(`
resource "nexus_cleanup_policy" "acceptance" {
	name    = "%s"
	format = "%s"
}
`, c.Name, c.Format)
}
