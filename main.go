package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

type FifoClient struct {
	ApiKey     string
	Endpoint   string
	Timeout    int
	MaxRetries int
}

type IPRange struct {
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Subnet  string `json:"subnet"`
	Gateway string `json:"gateway"`
	Netmask string `json:"netmask"`
	Vlan    int    `json:"vlan"`
	First   string `json:"first"`
	Last    string `json:"last"`
}

func (m *IPRange) Id() string {
	return "id-" + m.Name + "!"
}

func (c *FifoClient) CreateIpRange(m *IPRange) error {
	return nil
}

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

// List of supported resources and their configuration fields.
// Here we define da linked list of all the resources that we want to
// support in our provider. As an example, if you were to write an AWS provider
// which supported resources like ec2 instances, elastic balancers and things of that sort
// then this would be the place to declare them.
// More info here https://github.com/hashicorp/terraform/blob/v0.6.6/helper/schema/resource.go#L17-L81
func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"projectfifo_iprange": &schema.Resource{
			SchemaVersion: 1,
			Create:        createFunc,
			Read:          readFunc,
			Update:        updateFunc,
			Delete:        deleteFunc,
			Schema: map[string]*schema.Schema{ // List of supported configuration fields for your resource
				"name": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"tag": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"subnet": &schema.Schema{
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
					Type:     schema.TypeInt,
					Required: true,
				},
				"last": &schema.Schema{
					Type:     schema.TypeInt,
					Required: true,
				},
			},
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

// The methods defined below will get called for each resource that needs to
// get created (createFunc), read (readFunc), updated (updateFunc) and deleted (deleteFunc).
// For example, if 10 resources need to be created then `createFunc`
// will get called 10 times every time with the information for the proper
// resource that is being mapped.
//
// If at some point any of these functions returns an error, Terraform will
// imply that something went wrong with the modification of the resource and it
// will prevent the execution of further calls that depend on that resource
// that failed to be created/updated/deleted.

func createFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*FifoClient)
	ipRange := IPRange{
		Name:    d.Get("name").(string),
		Tag:     d.Get("tag").(string),
		Subnet:  d.Get("subnet").(string),
		Gateway: d.Get("gateway").(string),
		Netmask: d.Get("netmask").(string),
		Vlan:    d.Get("vlan").(int),
		First:   d.Get("first").(string),
		Last:    d.Get("last").(string),
	}

	err := client.CreateIpRange(&ipRange)
	if err != nil {
		return err
	}

	d.SetId(ipRange.Id())

	return nil
}

func readFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func updateFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func deleteFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}
