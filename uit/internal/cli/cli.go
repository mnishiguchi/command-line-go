package cli

import (
	"fmt"
	"os"

	"github.com/mnishiguchi/command-line-go/uit/internal/formatter"
	"github.com/urfave/cli/v2"
)

// NewApp returns a CLI app instance for uit.
func NewApp(version string) *cli.App {
	return &cli.App{
		Name:    "uit",
		Usage:   "Replicate Uithub formatting locally",
		Version: version,
		Action: func(c *cli.Context) error {
			path := "."

			// Use argument as path if provided
			if c.Args().Len() > 0 {
				path = c.Args().First()
			}

			// Print Git-aware tree structure rooted at given path
			if err := formatter.RenderGitTree(path, os.Stdout); err == nil {
				fmt.Println() // spacer if tree was printed
			}

			// Check if the input path is a file or directory
			info, err := os.Stat(path)
			if err != nil {
				return fmt.Errorf("invalid path: %w", err)
			}

			if info.IsDir() {
				// Render all Git-tracked files under the directory
				files, err := formatter.ListGitFilesUnder(path)
				if err != nil {
					return fmt.Errorf("failed to list files: %w", err)
				}

				for _, f := range files {
					if err := formatter.RenderFileContent(f, os.Stdout); err != nil {
						return fmt.Errorf("failed to render file %s: %w", f, err)
					}
				}
			} else {
				// Render a single file
				if err := formatter.RenderFileContent(path, os.Stdout); err != nil {
					return fmt.Errorf("failed to render file %s: %w", path, err)
				}
			}

			return nil
		},
	}
}
