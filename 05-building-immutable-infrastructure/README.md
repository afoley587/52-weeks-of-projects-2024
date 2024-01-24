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

I will also be deploying an EC2 instance into AWS. We will be staying __well__
within the free-tier. If you want to follow along, and I hope you do, please
make sure you have these three tools installed and an AWS account setup.