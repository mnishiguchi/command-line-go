package cli_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/mnishiguchi/command-line-go/mdg/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	inputFile  = "../../testdata/test1.md"
	goldenFile = "../../testdata/test1.md.html"
)

func TestCLI_PreviewToTempFile(t *testing.T) {
	var stdout bytes.Buffer
	app := cli.NewApp("v-test")
	app.Writer = &stdout

	args := []string{"cmd", "preview", "--file", inputFile, "--skip-preview"}
	err := app.Run(args)
	require.NoError(t, err, "CLI command should run without error")

	resultFile := strings.TrimSpace(stdout.String())
	t.Logf("Temp file generated: %s", resultFile)

	// Cleanup file after test
	t.Cleanup(func() {
		_ = os.Remove(resultFile)
	})

	require.FileExists(t, resultFile, "temp output file should exist")

	result, err := os.ReadFile(resultFile)
	require.NoError(t, err, "should read result HTML file")

	expected, err := os.ReadFile(goldenFile)
	require.NoError(t, err, "should read golden HTML file")

	assert.Equal(t, string(expected), string(result), "CLI output should match golden file")
}

