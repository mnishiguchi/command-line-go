package cli

import (
	"fmt"
	"os"

	"github.com/mnishiguchi/command-line-go/mdg/internal/md"
	"github.com/urfave/cli/v2"
)

func NewApp(version string) *cli.App {
	return &cli.App{
		Name:    "mdg",
		Usage:   "Preview a markdown file in your browser",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:  "preview",
				Usage: "Convert markdown to HTML and preview in browser",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Aliases: []string{"f"}, Usage: "Markdown file", Required: true},
					&cli.StringFlag{Name: "template", Aliases: []string{"t"}, Usage: "Optional HTML template"},
					&cli.BoolFlag{Name: "skip-preview", Aliases: []string{"s"}, Usage: "Skip opening browser"},
				},
				Action: func(c *cli.Context) error {
					file := c.String("file")
					template := c.String("template")
					skip := c.Bool("skip-preview")

					input, err := os.ReadFile(file)
					if err != nil {
						return err
					}

					htmlData, err := md.ParseContent(input, template)
					if err != nil {
						return err
					}

					temp, err := os.CreateTemp("", "mdg*.html")
					if err != nil {
						return err
					}
					outName := temp.Name()

					if err := temp.Close(); err != nil {
						return err
					}
					fmt.Fprintln(c.App.Writer, outName)

					if err := os.WriteFile(outName, htmlData, 0644); err != nil {
						return err
					}

					if skip {
						return nil
					}

					defer os.Remove(outName)

					return md.Preview(outName)
				},
			},
		},
	}
}
