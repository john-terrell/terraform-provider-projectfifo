package main

import "github.com/hashicorp/terraform/helper/schema"

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
	return nil
}

func iprangeUpdateFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func iprangeDeleteFunc(d *schema.ResourceData, meta interface{}) error {
	return nil
}
