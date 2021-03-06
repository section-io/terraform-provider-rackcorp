provider "rackcorp" {}

variable "image_access_key" {}
variable "image_access_secret" {}
variable "image_bucket" {}
variable "image_path" {}

resource "rackcorp_server" "exercise" {
  name         = "custom-image"
  country      = "AU"
  location     = "GLOBALSWITCH-SYD1"
  server_class = "PERFORMANCE"
  cpu_count    = 1
  memory_gb    = 4
  user_data    = "${file("${path.module}/user_data.yaml")}"

  operating_system                 = "SELFINSTALLEDFROMISO"
  deploy_media_image_access_key    = "${var.image_access_key}"
  deploy_media_image_access_secret = "${var.image_access_secret}"
  deploy_media_image_bucket        = "${var.image_bucket}"
  deploy_media_image_path          = "${var.image_path}"

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
