# Configure the Packet Provider.
provider "packet" {
  auth_token = "${var.auth_token}"
  version    = "~>2.2.1"
}
