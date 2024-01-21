package requests

type CreateHelloRequest struct {
	Name string
}

// GetResourceType For e.g. "azurecloudspace"
func (r CreateHelloRequest) GetResourceType() string {
	return "hello"
}

// GetRequestType For e.g. "create"
func (r CreateHelloRequest) GetActionType() string {
	return "create"
}
