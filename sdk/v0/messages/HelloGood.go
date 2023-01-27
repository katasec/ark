package messages

type Base struct {
	Id string
}

type HelloGood struct {
	Base
	Message string
}
