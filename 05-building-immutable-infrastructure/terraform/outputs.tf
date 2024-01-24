# outputs.tf
output "private_key" {
  value     = tls_private_key.nginx_ec2.private_key_pem
  sensitive = true
}

# terraform output -raw private_key > /tmp/nginx-ssh
