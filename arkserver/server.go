package arkserver

import (
	"fmt"
	"net/http"

	"github.com/katasec/ark/arkrouter"
)

func Start() {

	r := arkrouter.NewChiRouter()

	r.GET("/", home)

	r.SERVE(":8080")
}

func home(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "Hello World!")
}
