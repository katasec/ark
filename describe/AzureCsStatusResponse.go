package describe

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/andanhm/go-prettytime"
)

type AzureCsStatusResponse struct {
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

func (resp AzureCsStatusResponse) printTable() {

	//endTime := resp.EndTime

	//Create tab writer to format table
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)

	inProgress := "in-progress"

	if resp.Result == inProgress {
		// Include normal date/time and pretty time
		fmt.Fprint(w, "Result\tStart time\tEndTime\tUpdateUrl\n")
		fmt.Fprint(w, "------\t----------\t-------\t---------\n")
		fmt.Fprintf(w, "%s\t%s (%s)\t(%s)\t%s\n", resp.Result, fmtTime(resp.StartTime), prettytime.Format(resp.StartTime), inProgress, resp.UpdateUrl)

	} else {
		// Include normal date/time and pretty time
		fmt.Fprint(w, "Result\tStart time\tEndTime\tUpdateUrl\n")
		fmt.Fprint(w, "------\t----------\t-------\t---------\n")
		fmt.Fprintf(w, "%s\t%s (%s)\t%s (%s)\t%s\n", resp.Result, fmtTime(resp.StartTime), prettytime.Format(resp.StartTime), fmtTime(resp.EndTime), prettytime.Format(resp.EndTime), resp.UpdateUrl)

	}

	// Flush output to stdout
	w.Flush()
}
