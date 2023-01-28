package describe

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/katasec/ark/config"
	"github.com/katasec/ark/utils"
	"github.com/olekukonko/tablewriter"
)

var (
	arkConfig = config.ReadConfig()
)

func Start(cloudspace string) {

	client := retryablehttp.NewClient()
	client.RetryMax = 10
	client.Logger = nil

	org := "katasec"
	stack := "dev"
	resource := "azurecloudspace"
	url := fmt.Sprintf("http://localhost:%s/status/%s/%s/%s", arkConfig.ApiServer.Port, org, resource, stack)

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	r, _ := utils.JsonUnmarshall[StatusResponse](string(body))

	table1 := tablewriter.NewWriter(os.Stdout)
	fmt.Println("Result:")
	table1.SetHeader([]string{"Result", "StartTime", "EndTime"})

	data1 := []string{r.Result, r.StartTime.String(), r.EndTime.String()}

	table1.Append(data1)
	table1.Render()

	fmt.Println("Resource Details:")
	table2 := tablewriter.NewWriter(os.Stdout)
	table2.SetHeader([]string{"DeleteCount", "CreateCount", "SameCount", "UpdateCount", "UpdateUrl"})

	data2 := []string{
		strconv.Itoa(r.DeleteCount),
		strconv.Itoa(r.CreateCount),
		strconv.Itoa(r.SameCount),
		strconv.Itoa(r.UpdateCount),
		r.UpdateUrl,
	}

	table2.Append(data2)
	table2.Render()
}

type StatusResponse struct {
	UpdateID    string
	StartTime   time.Time
	EndTime     time.Time
	Result      string
	DeleteCount int
	CreateCount int
	SameCount   int
	UpdateCount int
	UpdateUrl   string
}
