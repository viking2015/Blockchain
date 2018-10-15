package main

import (
	"log"

	cli "github.com/NlaakStudios/Blockchain/api/cli"
)

func main() {
	// Create a new MVC Application object
	cli, err := cli.NewClient("./data", "blockchain")
	if err != nil {
		log.Fatal(err)
	} else {
		cli.Run()
	}

}
