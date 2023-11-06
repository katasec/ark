package describe

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/katasec/ark/config"
	"github.com/katasec/ark/utils"
	//"github.com/olekukonko/tablewriter"
)

var (
	arkConfig *config.Config
	dtLayout  string
)

func init() {
	//arkConfig = config.ReadConfig()
	dtLayout = "2006-01-02 15:04:05"
}
func Start(cloudspace string) {

	cmd := NewDescribeCmd()
	resp := cmd.DescribeAzureCloudSpace()

	resp.printTable()

}

type DescribeCmd struct {
	response *AzureCsStatusResponse
	client   *retryablehttp.Client
}

func NewDescribeCmd() *DescribeCmd {

	// Create retryable client with 10 retries and suppress debug log
	c := retryablehttp.NewClient()
	c.RetryMax = 10
	c.Logger = nil

	return &DescribeCmd{
		client:   c,
		response: &AzureCsStatusResponse{},
	}
}

func (d DescribeCmd) DescribeAzureCloudSpace() AzureCsStatusResponse {

	org := "katasec"
	stack := "dev"
	resource := "azurecloudspace"

	// Construct Url
	url := fmt.Sprintf("http://localhost:%s/status/%s/%s/%s", arkConfig.ApiServer.Port, org, resource, stack)

	// Create make GET request
	resp, err := d.client.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Save response
	status, err := utils.JsonUnmarshall[AzureCsStatusResponse](string(body))
	if err != nil {
		fmt.Println("error, could not get response")
	}
	*d.response = status

	return *d.response
}

func fmtTime(theTime time.Time) string {
	return theTime.Local().Format(dtLayout)
}
