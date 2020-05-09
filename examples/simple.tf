# Example Project Fifo VM specification in Terraform.
# This example creates three VMs with differing OS configurations:
#   1) Smart OS Base 64 Zone
#   2) Ubuntu 18.04 KVM Instance
#   3) FreeBSD 11 KVM Instance
#

# Provider declaration for Project Fifo
provider "projectfifo" {
#    endpoint - The URL of the REST endpoint.   If not specified, the TF_FIFO_ENDPOINT environment variable will be used. 
#    example: 
#       endpoint = "http://<your endpoint IP>"

#    api_key - The API key used to authenticate with the REST endpoint.  If not specified, the TF_FIFO_APIKEY enviornment variable
#              will be used
#    example:
#       api_key = "1234567890ABCD"
}

# Data source declaration for existing Project Fifo data.
# (These 'data' declarations bring existing read-only data sources and makes them available to Terraform)

data "projectfifo_package" "example_package" {
    # Expose the "t2.micro" package in Project Fifo to Terraform.
    name = "t2.micro"
}

data "projectfifo_dataset" "smartos" {
    # Expose the a dataset with the given name and version to Terraform.
    name = "base-multiarch"
    version = "18.3.0"
}

data "projectfifo_dataset" "ubuntu" {
    # Expose the a dataset with the given name and version to Terraform.
    name = "ubuntu-certified-18.04"
    version = "20180808"
}

data "projectfifo_dataset" "freebsd" {
    # Expose the a dataset with the given name and version to Terraform.
    name = "freebsd-11"
    version = "20180213"
}

data "projectfifo_network" "default" {
    # Expose the a network with the given name and version to Terraform.
    name = "ClusterNetwork"
}

resource "projectfifo_vm" "example_vm1" {
    count = 1
    name = "Example SmartOS Zone"
    dataset = "${data.projectfifo_dataset.smartos.uuid}"
    package = "${data.projectfifo_package.example_package.uuid}"
    config = {
        alias = "vm1"
        networks = {
            net0 = "${data.projectfifo_network.default.uuid}"
        }
        hostname = "vm1"
    }
}

resource "projectfifo_vm" "example_vm2" {
    count = 1
    name = "Example Ubuntu VM"
    dataset = "${data.projectfifo_dataset.ubuntu.uuid}"
    package = "${data.projectfifo_package.example_package.uuid}"
    config = {
        alias = "vm2"
        networks = {
            net0 = "${data.projectfifo_network.default.uuid}"
        }
        hostname = "vm2"
    }
}

resource "projectfifo_vm" "example_vm3" {
    count = 1
    name = "Example FreeBSD VM"
    dataset = "${data.projectfifo_dataset.freebsd.uuid}"
    package = "${data.projectfifo_package.example_package.uuid}"
    config = {
        alias = "vm3"
        networks = {
            net0 = "${data.projectfifo_network.default.uuid}"
        }
        hostname = "vm3"
    }
}

/*
output "vm_ip1" {
    value = "${projectfifo_vm.example_vm1.ip}"
}

output "vm_ip2" {
    value = "${projectfifo_vm.example_vm2.ip}"
}

output "vm_ip3" {
    value = "${projectfifo_vm.example_vm3.ip}"
}
*/