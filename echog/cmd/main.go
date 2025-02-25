package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "echog",
		Usage:   "Go version of `echo`",
		Version: "0.1.0",
		Action: func(c *cli.Context) error {
			text := c.Args().Slice()
			omitNewline := c.Bool("n")

			output := strings.Join(text, " ")
			if omitNewline {
				fmt.Print(output)
			} else {
				fmt.Println(output)
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "n",
				Usage:   "do not print the trailing newline",
				Aliases: []string{"omit-newline"},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
