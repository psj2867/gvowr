package main

import (
	"gvowr/api"
)

func main() {
	if err := api.Run(); err != nil {
		panic(err)
	}
}
