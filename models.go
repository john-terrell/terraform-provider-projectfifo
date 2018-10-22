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
	UUID    string `json:"uuid"`
}

type VMNetworkConfigCreate struct {
	Net0 string `json:"net0"`
}

type VMConfigCreate struct {
	Alias    string                `json:"alias"`
	Autoboot bool                  `json:"autoboot"`
	Hostname string                `json:"hostname"`
	Networks VMNetworkConfigCreate `json:"networks"`
}

type VMCreate struct {
	Dataset string         `json:"dataset"`
	Package string         `json:"package"`
	Config  VMConfigCreate `json:"config"`
}

type VMConfig struct {
	Alias    string `json:"alias"`
	Autoboot bool   `json:"autoboot"`
	Hostname string `json:"hostname"`
}

type VM struct {
	Dataset string   `json:"dataset"`
	Package string   `json:"package"`
	Config  VMConfig `json:"config"`
	UUID    string   `json:"uuid"`
	State   string   `json:"state"`
}
