package scaleway

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	vpcgw "github.com/scaleway/scaleway-sdk-go/api/vpcgw/v1"
)

func TestAccScalewayVPCPublicGatewayDHCP_Basic(t *testing.T) {
	tt := NewTestTools(t)
	defer tt.Cleanup()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: tt.ProviderFactories,
		CheckDestroy:      testAccCheckScalewayVPCPublicGatewayDHCPDestroy(tt),
		Steps: []resource.TestStep{
			{
				Config: `
					resource scaleway_vpc_public_gateway_dhcp main {
						subnet = "192.168.1.0/24"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCPublicGatewayDHCPExists(tt, "scaleway_vpc_public_gateway_dhcp.main"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "subnet", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "enable_dynamic", "true"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "valid_lifetime", "3600"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "renew_timer", "3000"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "rebind_timer", "3060"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "push_default_route", "true"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "push_dns_server", "true"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "dns_server_override.#", "0"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "dns_search.#", "0"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "dns_local_name"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "pool_low"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "pool_high"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "updated_at"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "zone"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "organization_id"),
				),
			},
			{
				Config: `
					resource scaleway_vpc_public_gateway_dhcp main {
						subnet = "192.168.1.0/24"
						valid_lifetime = 3000
						renew_timer = 2000
						rebind_timer = 2060
						push_default_route = false
						push_dns_server = false
						enable_dynamic = false
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCPublicGatewayDHCPExists(tt, "scaleway_vpc_public_gateway_dhcp.main"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "subnet", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "push_default_route", "false"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "push_dns_server", "false"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "enable_dynamic", "false"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "valid_lifetime", "3000"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "renew_timer", "2000"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "rebind_timer", "2060"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "dns_server_override.#", "0"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "dns_search.#", "0"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "dns_local_name"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "pool_low"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "pool_high"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "updated_at"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "zone"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "organization_id"),
				),
			},
			{
				Config: `
					resource "scaleway_vpc_public_gateway_dhcp" main {
					  subnet = "192.168.1.0/24"
					  push_default_route = true
					  push_dns_server = true
					  enable_dynamic = true
					  dns_servers_override = ["192.168.1.2"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCPublicGatewayDHCPExists(tt, "scaleway_vpc_public_gateway_dhcp.main"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "subnet", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "enable_dynamic", "true"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "valid_lifetime", "3000"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "renew_timer", "2000"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "rebind_timer", "2060"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "push_default_route", "true"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "push_dns_server", "true"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "dns_servers_override.#", "1"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "dns_servers_override.0", "192.168.1.2"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "dns_local_name"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "pool_low"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "pool_high"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "created_at"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "updated_at"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "zone"),
					resource.TestCheckResourceAttrSet("scaleway_vpc_public_gateway_dhcp.main", "organization_id"),
				),
			},
			{
				Config: `
					resource "scaleway_vpc_public_gateway_dhcp" main {
					  subnet = "192.168.1.0/24"
					  push_default_route = true
					  push_dns_server = true
					  enable_dynamic = true
					  dns_servers_override = ["192.168.1.3"]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScalewayVPCPublicGatewayDHCPExists(tt, "scaleway_vpc_public_gateway_dhcp.main"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "subnet", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "dns_servers_override.#", "1"),
					resource.TestCheckResourceAttr("scaleway_vpc_public_gateway_dhcp.main", "dns_servers_override.0", "192.168.1.3"),
				),
			},
		},
	})
}

func testAccCheckScalewayVPCPublicGatewayDHCPExists(tt *TestTools, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource not found: %s", n)
		}

		vpcgwAPI, zone, ID, err := vpcgwAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = vpcgwAPI.GetDHCP(&vpcgw.GetDHCPRequest{
			DHCPID: ID,
			Zone:   zone,
		})

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckScalewayVPCPublicGatewayDHCPDestroy(tt *TestTools) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != "scaleway_vpc_public_gateway_dhcp" {
				continue
			}

			vpcgwAPI, zone, ID, err := vpcgwAPIWithZoneAndID(tt.Meta, rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = vpcgwAPI.GetDHCP(&vpcgw.GetDHCPRequest{
				DHCPID: ID,
				Zone:   zone,
			})

			if err == nil {
				return fmt.Errorf(
					"VPC public gateway DHCP config %s still exists",
					rs.Primary.ID,
				)
			}

			// Unexpected api error we return it
			if !is404Error(err) {
				return err
			}
		}

		return nil
	}
}
