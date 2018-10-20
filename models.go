package main

type IPRange struct {
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Network string `json:"network"`
	Gateway string `json:"gateway"`
	Netmask string `json:"netmask"`
	Vlan    int    `json:"vlan"`
	First   string `json:"first"`
	Last    string `json:"last"`
}
