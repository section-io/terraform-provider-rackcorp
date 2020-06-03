provider "rackcorp" {}

resource "rackcorp_server" "exercise" {
  name         = "ubuntu-image"
  country      = "AU"
  location     = "GLOBALSWITCH-SYD1"
  server_class = "PERFORMANCE"
  cpu_count    = 1
  memory_gb    = 4
  user_data    = "${file("${path.module}/user_data.yaml")}"

  meta_data    = "instance-id: ubuntu-image-123\nlocal-hostname: my-ubuntu-image"

  deploy_media_image_id = "233" # Ubuntu 18.04 cloud image VMDK

  storage = [
    {
      size_gb = 20
      type    = "SSD"
    },
  ]
}

output "primary_ip" {
  value = "${rackcorp_server.exercise.primary_ip}"
}
