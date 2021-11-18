package main

import (
	"os"
)

func main() {
	app := newApp()

	if err := app.Run(os.Args); err != nil {
		println(err.Error())
	}
}
