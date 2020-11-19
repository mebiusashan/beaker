package common

import (
	"fmt"
	"os"
)

func Err(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(0)
}

func Assert(e error) {
	if e != nil {
		Err(e)
	}
}
