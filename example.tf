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

  // data_center_id = 19

  // name = "the-hostname-from-tf"

  // post_install_script = "${file("a-script.sh")}"

  // traffic_gb = 10

  /*
  storage = [
    {
      size_mb = 50
      type    = "SSD" // or MAGNETIC
    },
  ]
  */
}
