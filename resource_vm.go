package main

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVm() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        vmCreateFunc,
		Read:          vmReadFunc,
		Delete:        vmDeleteFunc,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dataset": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"package": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"autoboot": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
							ForceNew: true,
						},
						"hostname": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"networks": &schema.Schema{
							Type:     schema.TypeMap,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							ForceNew: true,
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
			if len(vm.Config.Networks) > 0 && vm.Config.Networks[0].IP != "" {
				d.Set("package", vm.Package)
				d.Set("dataset", vm.Dataset)
				d.Set("state", vm.State)
				d.Set("ip", vm.Config.Networks[0].IP)

				break
			}
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
	d.Set("ip", vm.Config.Networks[0].IP)

	return nil
}

func vmDeleteFunc(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*FifoClient)

	err := client.DeleteVm(d.Id())
	if err != nil {
		return err
	}

	for {
		if !client.VmExists(d.Id()) {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
