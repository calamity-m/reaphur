package main

import (
	"fmt"

	"github.com/calamity-m/reaphur/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
