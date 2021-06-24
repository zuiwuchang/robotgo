package main

import (
	"io/ioutil"
	"os"

	"github.com/zuiwuchang/robotgo/clipboard"
)

func main() {
	out, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	if err := clipboard.WriteAll(string(out)); err != nil {
		panic(err)
	}
}
