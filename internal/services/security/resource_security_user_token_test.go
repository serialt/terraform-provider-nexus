package security_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/serialt/terraform-provider-nexus/internal/acceptance"
	"github.com/serialt/terraform-provider-nexus/internal/tools"
)

func TestAccResourceSecurityUserToken(t *testing.T) {
	if tools.GetEnv("SKIP_PRO_TESTS", "false") == "true" {
		t.Skip("Skipping Nexus Pro tests")
	}

	resName := "nexus_security_user_token.acceptance"

	token := security.UserTokenConfiguration{
		Enabled:        true,
		ProtectContent: false,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityUserTokenConfig(token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "enabled", strconv.FormatBool(token.Enabled)),
					resource.TestCheckResourceAttr(resName, "protect_content", strconv.FormatBool(token.ProtectContent)),
				),
			},
		},
	})
}

func testAccResourceSecurityUserTokenConfig(token security.UserTokenConfiguration) string {
	return fmt.Sprintf(`
resource "nexus_security_user_token" "acceptance" {
	enabled         = %t
	protect_content = %t
}
`, token.Enabled, token.ProtectContent)
}
