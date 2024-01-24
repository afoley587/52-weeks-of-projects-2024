# nginx-server.pkr.hcl

locals {
  build_date = formatdate("YYYYMMDDHHmm", timestamp())
}

data "amazon-ami" "debian" {
  filters = {
    name                = "debian-11-amd64*"
    root-device-type    = "ebs"
    virtualization-type = "hvm"
  }
  most_recent = true
  owners      = ["136693071363"]
}

source "amazon-ebs" "debian" {
  ami_name              = "${var.ami_name}-${local.build_date}"
  instance_type         = "${var.instance_type}"
  communicator          = "ssh"
  encrypt_boot          = false
  force_delete_snapshot = true
  force_deregister      = true
  source_ami            = "${data.amazon-ami.debian.id}"
  ssh_username          = "${var.ssh_username}"

  tags = {
    Name      = "${var.ami_name}"
    buildDate = "${local.build_date}"
  }
}

build {
  sources = ["source.amazon-ebs.debian"]

  provisioner "ansible" {
    ansible_env_vars = ["ANSIBLE_DIFF_ALWAYS=1", "ANSIBLE_FORCE_COLOR=1", "ANSIBLE_HOST_KEY_CHECKING=False", "ANSIBLE_SSH_ARGS='-o ForwardAgent=yes -o ControlMaster=auto -o ControlPersist=60s -o ServerAliveInterval=30'"]
    extra_arguments  = ["--extra-vars", "ansible_python_interpreter=/usr/bin/python3"]
    playbook_file    = "../ansible/main.yml"
    user             = "${var.ssh_username}"
  }
}