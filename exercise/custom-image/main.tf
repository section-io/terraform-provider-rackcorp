provider "rackcorp" {}

variable "image_access_key" {}
variable "image_access_secret" {}

resource "rackcorp_server" "exercise" {
  name         = "custom-image"
  country      = "AU"
  location     = "GLOBALSWITCH-SYD1"
  server_class = "PERFORMANCE"
  cpu_count    = 1
  memory_gb    = 4
  user_data    = "${file("${path.module}/user_data.yaml")}"

  deploy_media_image_bucket        = "rackcorp-bucket-name"
  deploy_media_image_access_key    = "${var.image_access_key}"
  deploy_media_image_access_secret = "${var.image_access_secret}"
  deploy_media_image_path          = "image-b0093e2.qcow2"

  host_group_id = 4 # Not all hosts support deploy_media_image_* yet

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
