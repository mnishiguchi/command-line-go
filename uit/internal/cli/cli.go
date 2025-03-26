package cli

import (
	"github.com/mnishiguchi/command-line-go/uit/internal/formatter"
	"github.com/urfave/cli/v2"
)

func NewApp(version string) *cli.App {
	return &cli.App{
		Name:    "uit",
		Usage:   "Replicate Uithub formatting locally",
		Version: version,
		Action: func(c *cli.Context) error {
			path := "."

			if c.Args().Len() > 0 {
				path = c.Args().First()
			}

			return formatter.Render(path)
		},
	}
}
