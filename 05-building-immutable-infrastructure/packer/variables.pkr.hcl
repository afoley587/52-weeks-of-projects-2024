variable "ami_name" {
  type    = string
  default = "immutable-infra"
}

variable "instance_type" {
  type    = string
  default = "t2.micro"
}

variable "ssh_username" {
  type    = string
  default = "admin"
}