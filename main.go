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

// Source https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/provider.go#L20-L43
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema:         providerSchema(),
		ResourcesMap:   providerResources(),
		DataSourcesMap: providerDataSources(),
		ConfigureFunc:  providerConfigure,
	}
}

// List of supported configuration fields for the provider.
// More info in https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/schema.go#L29-L142
func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_key": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "API Key used to authenticate with the Project Fifo API",
			DefaultFunc: schema.MultiEnvDefaultFunc([]string{
				"TF_FIFO_APIKEY",
			}, nil),
		},
		"endpoint": &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The URL to the Project Fifo API",
			DefaultFunc: schema.MultiEnvDefaultFunc([]string{
				"TF_FIFO_ENDPOINT",
			}, nil),
		},
	}
}

// This is the function used to fetch the configuration params given
// to our provider which we will use to initialise a client that
// interacts with the Project Fifo API.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := FifoClient{
		ApiKey:   d.Get("api_key").(string),
		Endpoint: d.Get("endpoint").(string),
	}

	// You could have some field validations here, like checking that
	// the API Key is has not expired or that the username/password
	// combination is valid, etc.

	return &client, nil
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"projectfifo_vm": resourceVm(),
	}
}

func providerDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"projectfifo_iprange": datasourceIpRange(),
		"projectfifo_package": datasourcePackage(),
		"projectfifo_network": datasourceNetwork(),
		"projectfifo_dataset": datasourceDataset(),
	}
}
