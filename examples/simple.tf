provider "projectfifo" {
#    endpoint - The URL of the REST endpoint.   If not specified, the TF_FIFO_ENDPOINT environment variable will be used. 
#    example: 
#       endpoint = "http://<your endpoint IP>"

#    api_key - The API key used to authenticate with the REST endpoint.  If not specified, the TF_FIFO_APIKEY enviornment variable
#              will be used
#    example:
#       api_key = "1234567890ABCD"
}

data "projectfifo_package" "example_package" {
    name = "t2.micro"
}

data "projectfifo_dataset" "ubuntu" {
    name = "ubuntu-certified-18.04"
    version = "20180808"
}

data "projectfifo_dataset" "centos" {
    name = "centos-7"
    version = "20181003"
}

data "projectfifo_dataset" "smartos" {
    name = "base-multiarch"
    version = "18.3.0"
}

data "projectfifo_dataset" "freebsd" {
    name = "freebsd-11"
    version = "20180213"
}

data "projectfifo_network" "default" {
    name = "ClusterNetwork"
}

resource "projectfifo_vm" "example_vm" {
    count = 1
    name = "example VM"
    dataset = "${data.projectfifo_dataset.smartos.uuid}"
    package = "${data.projectfifo_package.example_package.uuid}"
    config = {
        alias = "vm2"
        networks = {
            net0 = "${data.projectfifo_network.default.uuid}"
        }
        hostname = "vm2"

        autoboot = true
    }
}

output "vm_ip" {
    value = "${projectfifo_vm.example_vm.ip}"
}