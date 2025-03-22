package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/mnishiguchi/command-line-go/wcg/internal/filestats"
	"github.com/urfave/cli/v2"
)

type CountOptions struct {
	Lines bool
	Words bool
	Bytes bool
	Chars bool
}

// Execute defines the full CLI API and then runs it
func Execute(version string) {
	// Create a logger that writes to stderr without any prefix or timestamp.
	// This makes error messages cleaner and more consistent in CLI output.
	logger := log.New(os.Stderr, "", 0)

	// Disable global logger flags to prevent timestamps or prefixes from showing up.
	log.SetFlags(0)

	app := cli.App{
		Name:      "wcg",
		Version:   version,
		Usage:     "Print line, word, and byte counts for each file",
		UsageText: `wcg [flags] [files...]`,
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "lines", Aliases: []string{"l"}, Usage: "Count lines"},
			&cli.BoolFlag{Name: "words", Aliases: []string{"w"}, Usage: "Count words"},
			&cli.BoolFlag{Name: "bytes", Aliases: []string{"c"}, Usage: "Count bytes"},
			&cli.BoolFlag{Name: "chars", Aliases: []string{"m"}, Usage: "Count characters"},
		},
		Action: runWordCount,
	}

	if err := app.Run(os.Args); err != nil {
		// Print the error using the custom logger and exit with status code 1.
		logger.Printf("Error: %v", err)
		cli.OsExiter(1)
	}
}

func runWordCount(c *cli.Context) error {
	// Get input files from command-line arguments. If none are provided, use standard input.
	files := c.Args().Slice()
	if len(files) == 0 {
		files = append(files, "-") // "-" indicates standard input
	}

	// Collect user-specified flags into a structured CountOptions instance.
	opts := CountOptions{
		Lines: c.Bool("lines"),
		Words: c.Bool("words"),
		Bytes: c.Bool("bytes"),
		Chars: c.Bool("chars"),
	}

	// If no flags are specified, default to counting lines, words, and bytes (standard wc behavior).
	if !opts.Lines && !opts.Words && !opts.Bytes && !opts.Chars {
		opts.Lines, opts.Words, opts.Bytes = true, true, true
	}

	// Create a FileStats struct to keep track of accumulated statistics across multiple files.
	var tally filestats.FileStats

	// Process each file specified in the command arguments
	for _, file := range files {
		// Attempt to retrieve file statistics (line, word, byte, char counts).
		input, err := openInputSource(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wcg: %s: %v\n", file, err)
			continue
		}

		fileStats, err := func() (*filestats.FileStats, error) {
			// Use an inline function to ensure the file is closed immediately after processing,
			// avoiding deferred resource cleanup inside a loop which can lead to too many open files.
			defer input.Close()
			return filestats.GetFileStats(input)
		}()

		if err != nil {
			fmt.Fprintf(os.Stderr, "wcg: %s: %v\n", file, err)
			continue
		}

		// Print statistics for the current file.
		printCounts(fileStats, opts, file)

		// Aggregate the statistics into the `tally` struct to compute totals across files.
		tally.Lines += fileStats.Lines
		tally.Words += fileStats.Words
		tally.Bytes += fileStats.Bytes
		tally.Chars += fileStats.Chars
	}

	// If multiple files were processed, print the total statistics across all files.
	if len(files) > 1 {
		printCounts(&tally, opts, "total") // The label "total" is used instead of a filename.
	}

	return nil
}

// openInputSource returns an io.ReadCloser for the given filename.
// If the filename is "-", it returns stdin wrapped in a no-op closer.
func openInputSource(filename string) (io.ReadCloser, error) {
	if filename == "-" {
		return io.NopCloser(os.Stdin), nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func printCounts(stats *filestats.FileStats, opts CountOptions, label string) {
	var counts []string
	if opts.Lines {
		counts = append(counts, fmt.Sprintf("%3d", stats.Lines))
	}
	if opts.Words {
		counts = append(counts, fmt.Sprintf("%3d", stats.Words))
	}
	if opts.Bytes {
		counts = append(counts, fmt.Sprintf("%3d", stats.Bytes))
	}
	if opts.Chars {
		counts = append(counts, fmt.Sprintf("%3d", stats.Chars))
	}

	output := strings.Join(counts, " ")
	if label != "-" && label != "" {
		output += " " + label
	}

	fmt.Println(output)
}
