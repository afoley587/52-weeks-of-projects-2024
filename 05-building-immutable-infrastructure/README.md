# Building and Deploying Immutable Infrastructure On AWS
## A tale of two tools: Terraform and Packer

In the world of infrastructure, two main types have reigned king: 
mutable and immutable. We can think of mutable infrastructure as a piece
of infrastructure that we apply upgrades to without destroying. For example, 
let's say you have Nginx running on an EC2 instance. In the mutable world, 
you would use a configuration management tool, such as ansible, to perform
the following steps:

1. Stopping the service
2. Upgrading the binaries with a package manager (think `apt-get` or `dnf`)
3. Making any configuration changes
4. Restarting the service

In the immutable world, on the other hand, you would build a new image with
some image builder tool like packer. Packer would do steps 1-4 above, but it 
would do it on a base AMI which you would then deploy with terraform. Terraform
would likely create the new instance, swap your load balancer or DNS names to 
the new instance, and then remove the old instance.

One side leans on the side convenience (mutable) where the other side leans on
the side of ephemerality and statelessness. Now, where have I __really__ seen this
come in handy? In the case of big system upgrades. Have you ever tried to upgrade
an old version of debian, say from 10 to 11? Well, no, it's not the end of the world.
But, you do have to:

* backup the server
* update the packages
* update the sources.list
* do a full upgrade
* verify the upgrade

Wouldn't it be easier if you could just build a new image, deploy that, and be done
with it? Well, we're going to go over immutable infrastructure today. We will build
a base image with packer and ansible and then we will deploy it with terraform.

I also want to go on the record for saying that mutable infrastructure does have
its place in a lot of places! My blog is not to say __**DON'T USE MUTABLE INFRA**__.
I just simply want to explain immutable and show you what it is... hopefully you 
find it useful!

## Prerequisites
We will be using three tools to deploy a piece of immutable infrastructure:

* [packer]()
* [ansible]()
* [terraform]()

I will also be deploying an EC2 instance into AWS. If you want to follow along,
and I hope you do, please make sure you have these three tools installed and an 
AWS account setup.

## Directory Structure / Project Setup
The directory layout is below. I have tried to label all directories according
to the tools that they pertain to:

```shell
.
├── README.md
├── ansible
│   ├── main.yml
│   └── tasks
│       ├── nginx.yml
│       └── ufw.yml
├── packer
│   ├── nginx-server.pkr.hcl
│   ├── plugins.pkr.hcl
│   └── variables.pkr.hcl
└── terraform
    ├── ec2.tf
    ├── iam.tf
    └── versions.tf

5 directories, 10 files
```

So, the `packer` directory is all of the `packer` HCL which will define our 
image build. The `anisble` directory will define our playbooks which will
configure that image. And the `terraform` directory will be used to hold
our `terraform` HCL to deploy the image.

# Building The Image
## Packer: Defining The Build
First, we will be building the base image. As previously noted, we will be 
using terraform and ansible to accomplish this task. Packer will be used
to clone an existing AMI and more-or-less make it our own by adding tags to
it, configuring initial ssh connections to it, giving the AMI a name, etc.
Packer will then call ansible to do more of the nitty-gritty configuration
such as installing packages, configuring os-level firewall rules, etc.

Our packer directory has three files:

1. `plugins.pkr.hcl` - Defines the plugins our packer project will use. For
    example, an aws and ansbile plugin.
2. `variables.pkr.hcl` - Defines the variables that we will be using in our
    packer build. For example, ssh usernames and ami names.
3. `nginx-server.pkr.hcl` - Defines the bulk of the build. For example,
    finding the base AMI on AWS, definining our new AMI, and calling
    ansible.

Our `plugins.pkr.hcl` just tells packer which plugins to install
locally, which versions to use, etc. We can think of this similar
to a python TOML file:

```hcl
packer {
  required_plugins {
    amazon = {
      version = ">= 1.2.8"
      source  = "github.com/hashicorp/amazon"
    }
    ansible = {
      version = "~> 1"
      source  = "github.com/hashicorp/ansible"
    }
  }
}
```

In the above, we are telling packer to install at least version 1.2.8
of the AWS plugin and any ansible version that is from major version 1.

Our `variables.pkr.hcl` file is another super simple file. It just tells
packer about variables that we define. We're only going to define three, 
but you could define just about as many as you'd like:

```hcl
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
```

As you can see, we have three variables:

1. `ami_name` which is a string that defaults to `immutable-infra`
2. `instance_type` which is a string that defaults to `t2.micro`
3. `ssh_username` which is a string that defaults to `admin`

The least simple file in the packer configuration is `nginx-server.pkr.hcl`.
Like I said before, it defines the bulk of the logic for this build. Let's 
take it step by step so that we understand what we're looking at.

```hcl
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
```

```hcl
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
```

```hcl
build {
  sources = ["source.amazon-ebs.debian"]

  provisioner "ansible" {
    ansible_env_vars = [
      "ANSIBLE_DIFF_ALWAYS=1",
      "ANSIBLE_FORCE_COLOR=1",
      "ANSIBLE_HOST_KEY_CHECKING=False"
    ]
    extra_arguments = [
      "--extra-vars",
      "ansible_python_interpreter=/usr/bin/python3",
      "--scp-extra-args",
      "'-O'"
    ]
    playbook_file = "../ansible/main.yml"
    user          = "${var.ssh_username}"
  }
}
```