package aws

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"regexp"
	"testing"
)

func TestAccDataSourceAwsWafV2IPSet_Basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-name-test")
	scope := acctest.RandomWithPrefix("tf-scope-test")
	resourceName := "aws_wafv2_ip_set.ipset"
	datasourceName := "data.aws_wafv2_ip_set.ipset"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAwsWafV2IPSet_NonExistent,
				ExpectError: regexp.MustCompile(`WAFv2 IP Set not found`),
			},
			{
				Config: testAccDataSourceAwsWafV2IPSet_Name(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
				),
			},
			{
				Config: testAccDataSourceAwsWafV2IPSet_NameAndScope(name, scope),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "scope", resourceName, "scope"),
				),
			},
		},
	})
}

func testAccDataSourceAwsWafV2IPSet_Name(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_ip_set" "ipset" {
  name = %[1]q
}
data "aws_wafv2_ip_set" "ipset" {
  name = "aws_wafv2_ip_set.ipset.name"
}
`, name)
}

func testAccDataSourceAwsWafV2IPSet_NameAndScope(name string, scope string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_ip_set" "ipset" {
  name  = %[1]q
  scope = %[2]q
}
data "aws_wafv2_ip_set" "ipset" {
	scope = "aws_wafv2_ip_set.ipset.name"
}
`, name, scope)
}

const testAccDataSourceAwsWafV2IPSet_NonExistent = `
data "aws_wafv2_ip_set" "ipset" {
  name = "tf-name-test-does-not-exist"
}
`
