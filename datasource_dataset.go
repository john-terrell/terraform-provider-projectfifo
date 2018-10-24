package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func datasourceDataset() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Read:          datasetDatasourceReadFunc,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"version": &schema.Schema{
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

func datasetDatasourceReadFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*FifoClient)

	name := d.Get("name").(string)
	version := d.Get("version").(string)
	pkg, found, err := client.FindDataset(name, version)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("Dataset %s (version %s) was not found", name, version)
	}

	d.Set("uuid", pkg.UUID)
	d.SetId(pkg.UUID)

	return nil
}
