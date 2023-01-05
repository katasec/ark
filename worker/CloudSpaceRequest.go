package worker

import "time"

type CloudSpaceRequest struct {
	ProjectName string
	DtTimeStamp time.Time
}
