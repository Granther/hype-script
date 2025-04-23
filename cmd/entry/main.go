package main

import (
	"fmt"
	"os"
	"hype-script/internal/mainhype"
)

func main() {
	hype := mainhype.NewHype()
	err := hype.Start()
	if err != nil {
		fmt.Println("Unable to get Hype with it: ", err)
		os.Exit(1)
	}
}
