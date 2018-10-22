package main

import "strings"
import "time"
import "github.com/hashicorp/terraform/helper/schema"

func resourceVm() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        vmCreateFunc,
		Read:          vmReadFunc,
		Update:        vmUpdateFunc,
		Delete:        vmDeleteFunc,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"dataset": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"package": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"config": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"autoboot": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"hostname": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"networks": &schema.Schema{
							Type:     schema.TypeMap,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func getVMNetworkConfig(networks map[string]interface{}) VMNetworkConfigCreate {
	networkConfig := VMNetworkConfigCreate{
		Net0: networks["net0"].(string),
	}

	return networkConfig
}

func getVMConfig(cfg map[string]interface{}) VMConfigCreate {
	config := VMConfigCreate{
		Alias:    cfg["alias"].(string),
		Autoboot: cfg["autoboot"].(bool),
		Hostname: cfg["hostname"].(string),
		Networks: getVMNetworkConfig(cfg["networks"].(map[string]interface{})),
	}

	return config
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

func vmCreateFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*FifoClient)
	vm := VMCreate{
		Dataset: d.Get("dataset").(string),
		Package: d.Get("package").(string),
		Config:  getVMConfig(d.Get("config").(*schema.Set).List()[0].(map[string]interface{})),
	}

	id, err := client.CreateVm(&vm)
	if err != nil {
		return err
	}

	for {
		vm := VM{}
		vm, err = client.GetVm(id)
		if err != nil {
			return err
		}

		state := strings.ToLower(vm.State)
		if state == "running" || state == "failed" {
			break
		}

		time.Sleep(1 * time.Second)
	}

	d.SetId(id)

	return nil
}

func vmReadFunc(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*FifoClient)

	vm, err := client.GetVm(d.Id())
	if err != nil {
		return err
	}

	d.Set("package", vm.Package)
	d.Set("dataset", vm.Dataset)
	d.Set("state", vm.State)

	/*
		d.Set("network", iprange.Network)
		d.Set("gateway", iprange.Gateway)
		d.Set("netmask", iprange.Netmask)
		d.Set("vlan", iprange.Vlan)
		d.Set("first", iprange.First)
		d.Set("last", iprange.Last)
		d.Set("uuid", iprange.UUID)
	*/
	return nil
}

func vmUpdateFunc(d *schema.ResourceData, meta interface{}) error {
	/*
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
	*/
	return nil
}

func vmDeleteFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*FifoClient)

	return client.DeleteVm(d.Id())
}
