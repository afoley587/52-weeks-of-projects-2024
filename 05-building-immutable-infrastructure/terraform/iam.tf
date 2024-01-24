# iam.tf
# Nginx For EC2
resource "aws_iam_instance_profile" "nginx_iam" {
  name = "nginxIam"
  role = aws_iam_role.nginx_iam.name
}

resource "aws_iam_role" "nginx_iam" {
  name               = "nginxIam"
  assume_role_policy = data.aws_iam_policy_document.nginx_iam.json
}

data "aws_iam_policy_document" "nginx_iam" {
  statement {
    effect = "Allow"

    principals {
      type = "Service"

      identifiers = [
        "ec2.amazonaws.com",
        "ssm.amazonaws.com",
      ]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role_policy_attachment" "nginx_iam" {
  role       = aws_iam_role.nginx_iam.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}