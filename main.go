package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	opts := plugin.ServeOpts{
		ProviderFunc: Provider,
	}
	plugin.Serve(&opts)
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{ // Source https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/provider.go#L20-L43
		Schema:        providerSchema(),
		ResourcesMap:  providerResources(),
		ConfigureFunc: providerConfigure,
	}
}

// List of supported configuration fields for your provider.
// Here we define a linked list of all the fields that we want to
// support in our provider (api_key, endpoint, timeout & max_retries).
// More info in https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/schema.go#L29-L142
func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_key": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "API Key used to authenticate with the Project Fifo Cluster",
		},
		"endpoint": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The URL to the Fifo Cluster API",
		},
		"timeout": &schema.Schema{
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Max. wait time we should wait for a successful connection to the API",
		},
		"max_retries": &schema.Schema{
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The max. amount of times we will retry to connect to the API",
		},
	}
}

// This is the function used to fetch the configuration params given
// to our provider which we will use to initialise a dummy client that
// interacts with the API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := FifoClient{
		ApiKey:     d.Get("api_key").(string),
		Endpoint:   d.Get("endpoint").(string),
		Timeout:    d.Get("timeout").(int),
		MaxRetries: d.Get("max_retries").(int),
	}

	// You could have some field validations here, like checking that
	// the API Key is has not expired or that the username/password
	// combination is valid, etc.

	return &client, nil
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"projectfifo_iprange": &schema.Resource{
			SchemaVersion: 1,
			Create:        iprangeCreateFunc,
			Read:          iprangeReadFunc,
			Update:        iprangeUpdateFunc,
			Delete:        iprangeDeleteFunc,
			Schema: map[string]*schema.Schema{ // List of supported configuration fields for your resource
				"name": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"tag": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"network": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"gateway": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"netmask": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"vlan": &schema.Schema{
					Type:     schema.TypeInt,
					Required: true,
				},
				"first": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"last": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}