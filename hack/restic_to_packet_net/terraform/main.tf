# Create a spot market request
resource "packet_spot_market_request" "req" {
  project_id       = "${var.project_id}"
  max_bid_price    = "${var.spot_price}"
  facilities       = var.facilities
  wait_for_devices = true
  devices_min      = 1
  devices_max      = 1

  instance_parameters {
    hostname         = "restic-restore"
    billing_cycle    = "hourly"
    operating_system = "ubuntu_18_04"
    plan             = "${var.instance_type}"
    project_ssh_keys = ["${packet_project_ssh_key.restic_restore.id}"]
  }
}

resource "packet_project_ssh_key" "restic_restore" {
  name       = "restic_restore"
  project_id = "${var.project_id}"
  public_key = "${var.public_key}"
}