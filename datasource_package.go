package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func datasourcePackage() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Read:          packageDatasourceReadFunc,
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

func packageDatasourceReadFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*FifoClient)

	name := d.Get("name").(string)
	pkg, found, err := client.FindPackage(name)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("Package %s was not found", name)
	}

	d.Set("uuid", pkg.UUID)
	d.SetId(pkg.UUID)

	return nil
}
