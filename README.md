Terraform Provider for Project Fifo Cloud Platform
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10+
-	[Go](https://golang.org/doc/install) 1.11.0 or higher

Building the provider
---------------------

Clone repository to: `$GOPATH/src/github.com/johnterrell/terraform-provider-projectfifo`

```sh
$ mkdir -p $GOPATH/src/github.com/johnterrell/terraform-provider-projectfifo; cd $GOPATH/src/github.com/johnterrell/terraform-provider-projectfifo
$ git clone git@github.com:john-terrell/terraform-provider-projectfifo.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/johnterrell/terraform-provider-projectfifo
$ make build
```

Developing the provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-projectfifo
...
```
