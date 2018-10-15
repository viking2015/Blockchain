package main

import (
	"log"

	cli "github.com/NlaakStudios/Blockchain/api/cli"
)

func main() {
	//cli := cli.Client{}
	//cli.Run()

	// Create a new MVC Application object
	cli, err := cli.NewClient("./data", "blockchain")
	if err != nil {
		log.Fatal(err)
	} else {
		cli.Run()
	}

}
