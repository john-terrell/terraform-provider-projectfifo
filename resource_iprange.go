package main

import "github.com/hashicorp/terraform/helper/schema"

func resourceIpRange() *schema.Resource {
	return &schema.Resource{
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
	}
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

func iprangeCreateFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*FifoClient)
	ipRange := IPRange{
		Name:    d.Get("name").(string),
		Tag:     d.Get("tag").(string),
		Network: d.Get("network").(string),
		Gateway: d.Get("gateway").(string),
		Netmask: d.Get("netmask").(string),
		Vlan:    d.Get("vlan").(int),
		First:   d.Get("first").(string),
		Last:    d.Get("last").(string),
	}

	id, err := client.CreateIpRange(&ipRange)
	if err != nil {
		return err
	}

	d.SetId(id)

	return nil
}

func iprangeReadFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*FifoClient)

	iprange, err := client.GetIpRange(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", iprange.Name)
	d.Set("tag", iprange.Tag)
	d.Set("network", iprange.Network)
	d.Set("gateway", iprange.Gateway)
	d.Set("netmask", iprange.Netmask)
	d.Set("vlan", iprange.Vlan)
	d.Set("first", iprange.First)
	d.Set("last", iprange.Last)

	return nil
}

func iprangeUpdateFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*FifoClient)
	ipRange := IPRange{
		Name:    d.Get("name").(string),
		Tag:     d.Get("tag").(string),
		Network: d.Get("network").(string),
		Gateway: d.Get("gateway").(string),
		Netmask: d.Get("netmask").(string),
		Vlan:    d.Get("vlan").(int),
		First:   d.Get("first").(string),
		Last:    d.Get("last").(string),
	}

	err := client.UpdateIpRange(d.Id(), &ipRange)
	if err != nil {
		return err
	}

	return nil
}

func iprangeDeleteFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*FifoClient)

	return client.DeleteIpRange(d.Id())
}
