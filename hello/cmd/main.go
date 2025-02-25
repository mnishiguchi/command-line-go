package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "hello",
		Usage:   "A simple greeting program",
		Version: "0.1.0",
		Action: func(c *cli.Context) error {
			fmt.Println("元氣が一番、元氣があれば何でもできる！")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

