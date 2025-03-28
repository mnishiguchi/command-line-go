package cli

import (
	"fmt"
	"os"

	"github.com/mnishiguchi/command-line-go/uit/internal/formatter"
	"github.com/urfave/cli/v2"
)

// Config holds CLI options.
type Config struct {
	Path       string
	ShowBinary bool
	HeadLines  int
}

// NewApp returns a CLI app instance for uit.
func NewApp(version string) *cli.App {
	return &cli.App{
		Name:    "uit",
		Usage:   "Replicate Uithub formatting locally",
		Version: version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "show-binary",
				Usage: "Show binary file contents",
				Value: false,
			},
			&cli.IntFlag{
				Name:  "head",
				Usage: "Limit the number of lines printed per file",
			},
		},
		Action: func(c *cli.Context) error {
			cfg := Config{
				Path:       ".",
				ShowBinary: c.Bool("show-binary"),
				HeadLines:  c.Int("head"),
			}

			// Use argument as path if provided
			if c.Args().Len() > 0 {
				cfg.Path = c.Args().First()
			}

			return Run(cfg)
		},
	}
}

// Run executes the main logic using the given config.
func Run(cfg Config) error {
	// Print Git-aware tree structure rooted at given path
	if err := formatter.RenderGitTree(cfg.Path, os.Stdout); err == nil {
		fmt.Println() // spacer if tree was printed
	}

	// Check if the input path is a file or directory
	info, err := os.Stat(cfg.Path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	if info.IsDir() {
		// Render all Git-tracked files under the directory
		files, err := formatter.ListGitFilesUnder(cfg.Path)
		if err != nil {
			return fmt.Errorf("failed to list files: %w", err)
		}

		for _, f := range files {
			if err := formatter.RenderFileContent(f, os.Stdout, cfg.ShowBinary, cfg.HeadLines); err != nil {
				return fmt.Errorf("failed to render file %s: %w", f, err)
			}
		}
	} else {
		// Render a single file
		if err := formatter.RenderFileContent(cfg.Path, os.Stdout, cfg.ShowBinary, cfg.HeadLines); err != nil {
			return fmt.Errorf("failed to render file %s: %w", cfg.Path, err)
		}
	}

	return nil
}
