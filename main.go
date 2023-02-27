package main

import (
	"log"

	"github.com/thecodedproject/servicegen/internal"
)

func main() {

	err := internal.Generate()
	if err != nil {
		log.Fatal(err.Error())
	}
}
