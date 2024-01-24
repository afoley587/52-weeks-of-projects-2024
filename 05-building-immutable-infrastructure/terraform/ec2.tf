// ec2.tf
data "aws_ami" "nginx" {
  most_recent = true

  filter {
    name   = "name"
    values = ["immutable-infra-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["self"]
}

resource "tls_private_key" "nginx_ec2" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "nginx_ec2" {
  key_name   = "nginx-ssh-key"
  public_key = tls_private_key.nginx_ec2.public_key_openssh
}


resource "aws_instance" "nginx_ec2" {
  ami                  = data.aws_ami.nginx.id
  instance_type        = "t2.nano"
  iam_instance_profile = aws_iam_instance_profile.nginx_iam.id
  key_name             = aws_key_pair.nginx_ec2.key_name
}