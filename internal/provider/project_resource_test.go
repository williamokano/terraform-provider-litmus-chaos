package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "litmus-chaos_project" "main_project" {
  name = "Main Project"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item
					resource.TestCheckResourceAttr("litmus-chaos_project.main_project", "name", "Main Project"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("litmus-chaos_project.main_project", "id"),
					resource.TestCheckResourceAttrSet("litmus-chaos_project.main_project", "last_updated"),
				),
			},
		},
	})
}
