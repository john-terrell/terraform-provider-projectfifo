package main

import "github.com/hashicorp/terraform/helper/schema"

func datasourceIpRange() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Read:          iprangeDatasourceReadFunc,
		Schema: map[string]*schema.Schema{ // List of supported configuration fields for your resource
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"network": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"netmask": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"first": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func iprangeDatasourceReadFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*FifoClient)

	uuid := d.Get("uuid").(string)

	iprange, err := client.GetIpRange(uuid)
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
	d.Set("uuid", iprange.UUID)
	d.SetId(iprange.UUID)

	return nil
}
