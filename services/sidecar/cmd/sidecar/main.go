package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "sidecar",
	Usage: "the specular sidecar service",
	Action: func(*cli.Context) error {
		fmt.Println("coming soon...")
		return nil
	},
}

func main() {

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
