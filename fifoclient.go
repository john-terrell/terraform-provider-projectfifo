package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type FifoClient struct {
	ApiKey     string
	Endpoint   string
	Timeout    int
	MaxRetries int
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

func (c *FifoClient) SendRequest(method string, api string, body io.Reader) (string, error) {
	request, err := http.NewRequest(method, c.Endpoint+api, body)
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Authorization", "Bearer "+c.ApiKey)
	request.Header.Set("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("The HTTP request failed with error %s\n", err)
	}

	if err := CheckResponse(response); err != nil {
		return "", err
	}

	slurp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(slurp), nil
}

func (c *FifoClient) CreateIpRange(m *IPRange) (string, error) {

	jsonValue, _ := json.Marshal(m)

	request, err := http.NewRequest("POST", c.Endpoint+"/api/3/ipranges", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Authorization", "Bearer "+c.ApiKey)
	request.Header.Set("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("The HTTP request failed with error %s\n", err)
	}

	if err := CheckResponse(response); err != nil {
		return "", err
	}

	slurp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	log.Printf("Unmarsalling json from: %s\n", slurp)

	result := make(map[string]interface{})
	if err := json.Unmarshal(slurp, &result); err != nil {
		return "", err
	}

	var uuid = result["uuid"]

	return uuid.(string), nil
}
