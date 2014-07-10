package ARHelper

import (
	"fmt"
	"os"
	"strings"
)

func CheckGoEnv() {

	var (
		goPath = os.Getenv("GOPATH")
		goRoot = os.Getenv("GOROOT")
	)
	var error = false
	if len(goPath) == 0 {
		fmt.Println("GOPATH environment variable is not set. ",
			"Please refer to http://golang.org/doc/code.html to configure your Go environment.")
		error = true
	}

	if strings.Contains(goRoot, goPath) {
		fmt.Println("GOPATH (%s) must not include your GOROOT (%s). "+
			"Please refer to http://golang.org/doc/code.html to configure your Go environment.",
			goRoot, goPath)
		error = true
	}
	if error == false {
		fmt.Println("golang ENV check is ok")
	}
}

func checkErr(e error) {
	if e != nil {
		fmt.Println("error is ", e)
	}
}
