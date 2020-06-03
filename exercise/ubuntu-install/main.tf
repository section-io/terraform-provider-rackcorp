provider "rackcorp" {}

resource "rackcorp_server" "exercise" {
  name             = "ubuntu-install"
  country          = "AU"
  server_class     = "PERFORMANCE"
  operating_system = "UBUNTU18.04_64"
  cpu_count        = 1
  memory_gb        = 4
  root_password    = "a-secret-password"

  storage = [
    {
      size_gb = 20
      type    = "SSD"
    },
  ]
}
