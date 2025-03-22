package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/mnishiguchi/command-line-go/wcg/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestPrintCounts(t *testing.T) {
  tests := []struct {
    name     string
    stats    *internal.FileStats
    opts     CountOptions
    label    string
    expected string
  }{
    {
      name:     "all_fields_on",
      stats:    &internal.FileStats{Lines: 2, Words: 7, Bytes: 48, Chars: 48},
      opts:     CountOptions{Lines: true, Words: true, Bytes: true, Chars: true},
      label:    "test.txt",
      expected: "  2   7  48  48 test.txt\n",
    },
    {
      name:     "chars_only",
      stats:    &internal.FileStats{Chars: 99},
      opts:     CountOptions{Chars: true},
      label:    "data.txt",
      expected: " 99 data.txt\n",
    },
    {
      name:     "stdin_label_omitted",
      stats:    &internal.FileStats{Lines: 1, Words: 1, Bytes: 5},
      opts:     CountOptions{Lines: true, Words: true, Bytes: true},
      label:    "-",
      expected: "  1   1   5\n",
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      // Capture stdout
      r, w, _ := os.Pipe()
      oldStdout := os.Stdout
      os.Stdout = w

      printCounts(tt.stats, tt.opts, tt.label)

      w.Close()
      os.Stdout = oldStdout

      var buf bytes.Buffer
      _, _ = io.Copy(&buf, r)

      assert.Equal(t, tt.expected, buf.String())
    })
  }
}

func TestRunWordCount(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := "hello world\nthis is a test\n"
	_, err = tmpFile.WriteString(content)
	require.NoError(t, err)
	tmpFile.Close()

	// Capture stdout
	var buf bytes.Buffer
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	app := &cli.App{
		Name:   "wcg",
		Usage:  "test",
		Action: runWordCount,
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{Name: "lines", Aliases: []string{"l"}},
		&cli.BoolFlag{Name: "words", Aliases: []string{"w"}},
		&cli.BoolFlag{Name: "bytes", Aliases: []string{"c"}},
		&cli.BoolFlag{Name: "chars", Aliases: []string{"m"}},
	}

	err = app.Run([]string{"wcg", tmpFile.Name()})
	require.NoError(t, err)

	w.Close()
	os.Stdout = stdout
	_, _ = io.Copy(&buf, r)

	expected := fmt.Sprintf("  2   6  27 %s\n", tmpFile.Name())
	assert.Equal(t, expected, buf.String(), "CLI output should match expected format for file statistics")
}
