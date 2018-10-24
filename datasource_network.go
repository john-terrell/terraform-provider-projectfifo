package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func datasourceNetwork() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Read:          networkDatasourceReadFunc,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"uuid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func networkDatasourceReadFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*FifoClient)

	name := d.Get("name").(string)
	nw, found, err := client.FindNetwork(name)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("Network %s was not found", name)
	}

	d.Set("uuid", nw.UUID)
	d.SetId(nw.UUID)

	return nil
}
