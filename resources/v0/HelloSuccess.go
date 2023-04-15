package resources

type Base struct {
	UpdateID string
	Id       string
}

type HelloSuccess struct {
	Base
	Message string
}
