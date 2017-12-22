provider "rackcorp" {
	api_uuid = "the-uuid-from-tf"
	api_secret = "the-secret-from-tf"
	customer_id = "001122"
}

resource "rackcorp_server" "example" {
}