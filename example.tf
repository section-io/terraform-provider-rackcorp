provider "rackcorp" {
  api_uuid    = "the-uuid-from-tf"
  api_secret  = "the-secret-from-tf"
  customer_id = "001122"
}

resource "rackcorp_server" "example" {
  country          = "the-country-from-tf"
  server_class     = "PERFORMANCE"
  operating_system = "the-operating-system-from-tf"
  cpu_count        = 1
  memory_gb        = 4
  root_password    = "a-secret-password"

  /*
  nics = [
    {
      name      = "public"
      vlan      = 1
      speed     = 1000
      ipv4      = 1
      pool_ipv4 = 1
      ipv6      = 1
      pool_ipv6 = 1
    },
  ]
  */

  /*
  firewall_policies = [
    {
      direction       = "INBOUND"
      policy          = "ALLOW"
      protocol        = "TCP"
      port_to         = "80"
      port_from       = "80"
      ip_address_from = "203.0.113.5"
      ip_address_to   = "203.0.113.6"
      comment         = "HTTP"
      order           = 0
    },
  ]
  */

  // data_center_id = 19

  // name = "the-hostname-from-tf"

  // post_install_script = "${file("a-script.sh")}"

  // traffic_gb = 10

  /*
  storage = [
    {
      size_gb = 50
      type    = "SSD" // or MAGNETIC
    },
  ]
  */
}
