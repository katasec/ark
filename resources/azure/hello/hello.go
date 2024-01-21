package hello

import (
	"encoding/json"
	"fmt"
)

type Hello struct {
	ID   int
	Name string
}

func (res *Hello) ToJson() string {
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(b)
}
