provider "rackcorp" {
	api_uuid = "the-uuid-from-tf"
	api_secret = "the-secret-from-tf"
}

resource "rackcorp_server" "example" {
}