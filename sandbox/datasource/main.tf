provider "infoblox" {
  user     = "${var.user}"
  password = "${var.pass}"
  server   = "${var.host}"
  allow_unverified_ssl = true  # default is false
}

data "infoblox_a_records" "free" {
  free_name_prefix = "swarm"
  free_name_domain = "${var.domain}"
  free_name_pad = "3"
}

resource "infoblox_a_record" "myrec" {
  domain = "${var.domain}"
  // name_prefix = "swarm"
  // name_index_pad = "3"
  // name_index_start = 100
  ipv4 = "10.0.0.220"
  name = "${data.infoblox_a_records.free.names[0]}"
  ttl = 600
}