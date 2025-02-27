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
		Action:  runEcho,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "n",
				Usage:   "do not print the trailing newline",
				Aliases: []string{"omit-newline"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runEcho(ctx *cli.Context) error {
	message := strings.Join(ctx.Args().Slice(), " ")
	omitNewline := ctx.Bool("n")

	if omitNewline {
		fmt.Print(message)
	} else {
		fmt.Println(message)
	}

	return nil
}
