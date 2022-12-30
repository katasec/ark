package cloudspace

type CloudspaceRequest struct {
	Id          int
	ProjectName string
	Name        string
	Tags        map[string]string
}
