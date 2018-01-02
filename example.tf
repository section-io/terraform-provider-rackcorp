provider "rackcorp" {
  api_uuid    = "the-uuid-from-tf"
  api_secret  = "the-secret-from-tf"
  customer_id = "001122"
}

resource "rackcorp_server" "example" {
  country          = "the-country-from-tf"
  server_class     = "the-server_class-from-tf"
  operating_system = "the-operating-system-from-tf"
  cpu_count        = 1
  memory_gb        = 4
  root_password    = "a-secret-password"

  // post_install_script = "${file("a-script.sh")}"

  /*
  storage = [
    {
      size_mb = 50
      type    = "SSD" // or MAGNETIC
    },
  ]
  */
}
