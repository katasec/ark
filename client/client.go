package client

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/katasec/ark/config"
	"github.com/katasec/ark/sdk/v0/messages"
)

var (
	arkConfig = config.ReadConfig()
	dtLayout  = "2006-01-02 15:04:05"
)

func init() {
	//fmt.Println("This is init()")
}

func Start() {
	fmt.Println("client Start()")

}

type ArkClient struct {
	client   *retryablehttp.Client
	org      string
	stack    string
	resource string
}

func NewArkClient() *ArkClient {

	// Create retryable client with 10 retries and suppress debug log
	c := retryablehttp.NewClient()
	c.RetryMax = 10
	c.Logger = nil

	return &ArkClient{
		client: c,
		org:    "katasec",
		stack:  "dev",
	}
}

func (c *ArkClient) AddCloudSpaces(cs messages.AzureCloudspace) error {

	// Construct Url
	url := fmt.Sprintf("http://localhost:%s/azure/cloudspaces/%s", arkConfig.ApiServer.Port, cs.Name)

	fmt.Println("The url is:" + url)

	// Convert to json
	postBody, err := json.Marshal(cs)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Post to ArkServer
	resp, err := c.client.Post(url, "application/json", postBody)
	if err != nil {
		fmt.Println("Error posting data to crete cloudspace", err)
		return err
	}

	// Read response
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err.Error())
		return err
	}

	return nil
}
