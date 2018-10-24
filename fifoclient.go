package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type FifoClient struct {
	ApiKey     string
	Endpoint   string
	Timeout    int
	MaxRetries int
	NetworkMap map[string]string
	PackageMap map[string]string
	DatasetMap map[string]string
}

type errorReply struct {
	Error *Error `json:"error"`
}

// ErrorItem is a detailed error code & message from the API frontend.
type ErrorItem struct {
	// Reason is the typed error code. For example: "some_example".
	Reason string `json:"reason"`
	// Message is the human-readable description of the error.
	Message string `json:"message"`
}

// Error contains an error response from the server.
type Error struct {
	// Code is the HTTP response status code and will always be populated.
	Code int `json:"code"`
	// Message is the server response message and is only populated when
	// explicitly referenced by the JSON server response.
	Message string `json:"message"`
	// Body is the raw response returned by the server.
	// It is often but not always JSON, depending on how the request fails.
	Body string
	// Header contains the response header fields from the server.
	Header http.Header

	Errors []ErrorItem
}

func (e *Error) Error() string {
	if len(e.Errors) == 0 && e.Message == "" {
		return fmt.Sprintf("Error: HTTP response code %d", e.Code)
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Error %d: ", e.Code)
	if e.Message != "" {
		fmt.Fprintf(&buf, "%s", e.Message)
	}
	if len(e.Errors) == 0 {
		return strings.TrimSpace(buf.String())
	}
	if len(e.Errors) == 1 && e.Errors[0].Message == e.Message {
		fmt.Fprintf(&buf, ", %s", e.Errors[0].Reason)
		return buf.String()
	}
	fmt.Fprintln(&buf, "\nMore details:")
	for _, v := range e.Errors {
		fmt.Fprintf(&buf, "Reason: %s, Message: %s\n", v.Reason, v.Message)
	}
	return buf.String()
}

func CheckResponse(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	slurp, err := ioutil.ReadAll(res.Body)
	if err == nil {
		jerr := new(errorReply)
		err = json.Unmarshal(slurp, jerr)
		if err == nil && jerr.Error != nil {
			if jerr.Error.Code == 0 {
				jerr.Error.Code = res.StatusCode
			}
			jerr.Error.Body = string(slurp)
			return jerr.Error
		}
	}
	return &Error{
		Code:   res.StatusCode,
		Body:   string(slurp),
		Header: res.Header,
	}
}

func (c *FifoClient) SendRequest(method string, api string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, c.Endpoint+api, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Authorization", "Bearer "+c.ApiKey)
	request.Header.Set("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("The HTTP request failed with error %s.\n", err)
	}

	if err := CheckResponse(response); err != nil {
		return nil, err
	}

	slurp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return slurp, nil
}

func (c *FifoClient) CachePackageList() error {
	response, err := c.SendRequest("GET", "/api/3/packages", nil)
	if err != nil {
		return err
	}

	var packages []string
	if err := json.Unmarshal(response, &packages); err != nil {
		return err
	}

	for _, packageUUID := range packages {
		pkg := Package{}

		pkg, err = c.GetPackage(packageUUID)
		if err != nil {
			return err
		}

		c.PackageMap[pkg.Name] = packageUUID
	}

	return nil
}

func (c *FifoClient) FindPackage(name string) (Package, bool, error) {
	if len(c.PackageMap) == 0 {
		c.PackageMap = make(map[string]string)
		err := c.CachePackageList()
		if err != nil {
			return Package{}, false, err
		}
	}

	packageUUID, found := c.PackageMap[name]
	if !found {
		return Package{}, false, nil
	}

	foundPackage, err := c.GetPackage(packageUUID)

	found = err == nil

	return foundPackage, found, err
}

func (c *FifoClient) GetPackage(uuid string) (Package, error) {
	pkg := Package{}

	response, err := c.SendRequest("GET", "/api/3/packages/"+uuid, nil)
	if err != nil {
		return pkg, err
	}

	if err := json.Unmarshal(response, &pkg); err != nil {
		return pkg, err
	}

	return pkg, nil
}

func (c *FifoClient) CacheNetworkList() error {
	response, err := c.SendRequest("GET", "/api/3/networks", nil)
	if err != nil {
		return err
	}

	var networks []string
	if err := json.Unmarshal(response, &networks); err != nil {
		return err
	}

	for _, uuid := range networks {
		nw := Network{}

		nw, err = c.GetNetwork(uuid)
		if err != nil {
			return err
		}

		key := nw.Name
		c.NetworkMap[key] = uuid
	}

	return nil
}

func (c *FifoClient) FindNetwork(name string) (Network, bool, error) {
	if len(c.NetworkMap) == 0 {
		c.NetworkMap = make(map[string]string)
		err := c.CacheNetworkList()
		if err != nil {
			return Network{}, false, err
		}
	}

	key := name
	uuid, found := c.NetworkMap[key]
	if !found {
		return Network{}, false, nil
	}

	foundNetwork, err := c.GetNetwork(uuid)

	found = err == nil

	return foundNetwork, found, err
}

func (c *FifoClient) GetNetwork(uuid string) (Network, error) {
	response, err := c.SendRequest("GET", "/api/3/networks/"+uuid, nil)
	if err != nil {
		return Network{}, err
	}

	nw := Network{}
	if err := json.Unmarshal(response, &nw); err != nil {
		return Network{}, err
	}

	return nw, nil

}

func (c *FifoClient) CacheDatasetList() error {
	response, err := c.SendRequest("GET", "/api/3/datasets", nil)
	if err != nil {
		return err
	}

	var datasets []string
	if err := json.Unmarshal(response, &datasets); err != nil {
		return err
	}

	for _, uuid := range datasets {
		ds := Dataset{}

		ds, err = c.GetDataset(uuid)
		if err != nil {
			return err
		}

		key := ds.Name + ":" + ds.Version
		c.DatasetMap[key] = uuid
	}

	return nil
}

func (c *FifoClient) FindDataset(name string, version string) (Dataset, bool, error) {
	if len(c.DatasetMap) == 0 {
		c.DatasetMap = make(map[string]string)
		err := c.CacheDatasetList()
		if err != nil {
			return Dataset{}, false, err
		}
	}

	key := name + ":" + version
	uuid, found := c.DatasetMap[key]
	if !found {
		return Dataset{}, false, nil
	}

	foundDataset, err := c.GetDataset(uuid)

	found = err == nil

	return foundDataset, found, err
}

func (c *FifoClient) GetDataset(uuid string) (Dataset, error) {
	response, err := c.SendRequest("GET", "/api/3/datasets/"+uuid, nil)
	if err != nil {
		return Dataset{}, err
	}

	ds := Dataset{}
	if err := json.Unmarshal(response, &ds); err != nil {
		return Dataset{}, err
	}

	return ds, nil

}
func (c *FifoClient) CreateVm(m *VMCreate) (string, error) {
	jsonDocument, _ := json.Marshal(m)

	response, err := c.SendRequest("POST", "/api/3/vms", bytes.NewBuffer(jsonDocument))
	if err != nil {
		return "", err
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(response, &result); err != nil {
		return "", err
	}

	var uuid = result["uuid"]

	return uuid.(string), nil
}

func (c *FifoClient) GetVm(uuid string) (VM, error) {
	vm := VM{}

	response, err := c.SendRequest("GET", "/api/3/vms/"+uuid, nil)
	if err != nil {
		return vm, err
	}

	if err := json.Unmarshal(response, &vm); err != nil {
		return vm, err
	}

	return vm, nil
}

func (c *FifoClient) VmExists(uuid string) bool {
	_, err := c.SendRequest("GET", "/api/3/vms/"+uuid, nil)

	if err != nil {
		return false
	}

	return true
}

func (c *FifoClient) DeleteVm(uuid string) error {
	_, err := c.SendRequest("DELETE", "/api/3/vms/"+uuid, nil)

	return err
}

func (c *FifoClient) CreateIpRange(m *IPRange) (string, error) {
	jsonDocument, _ := json.Marshal(m)

	response, err := c.SendRequest("POST", "/api/3/ipranges/", bytes.NewBuffer(jsonDocument))
	if err != nil {
		return "", err
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(response, &result); err != nil {
		return "", err
	}

	var uuid = result["uuid"]

	return uuid.(string), nil
}

func (c *FifoClient) GetIpRange(uuid string) (*IPRange, error) {
	response, err := c.SendRequest("GET", "/api/3/ipranges/"+uuid, nil)
	if err != nil {
		return nil, err
	}

	iprange := IPRange{}
	if err := json.Unmarshal(response, &iprange); err != nil {
		return nil, err
	}

	return &iprange, nil
}

func (c *FifoClient) UpdateIpRange(uuid string, m *IPRange) error {
	jsonDocument, _ := json.Marshal(m)

	_, err := c.SendRequest("PUT", "/api/3/ipranges/"+uuid, bytes.NewBuffer(jsonDocument))
	if err != nil {
		return err
	}

	return nil
}

func (c *FifoClient) DeleteIpRange(uuid string) error {
	_, err := c.SendRequest("DELETE", "/api/3/ipranges/"+uuid, nil)

	return err
}
