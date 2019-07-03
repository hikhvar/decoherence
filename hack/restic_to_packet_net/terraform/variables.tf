variable "auth_token" {
  type = string
}

variable "project_id" {
  type = string
}

variable "public_key" {
  type = string
}

variable "instance_type" {
  type    = string
  default = "s1.large.x86"
}

variable "facilities" {
  type    = list
  default = ["dfw2", "ams1", "ewr1", "nrt1"]
}

variable "spot_price" {
  type    = string
  default = 0.5
}