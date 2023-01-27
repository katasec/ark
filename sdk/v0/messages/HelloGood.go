package messages

type Base struct {
	Id string
}

type HelloSuccess struct {
	Base
	Message string
}
