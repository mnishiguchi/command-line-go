package cli_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/mnishiguchi/command-line-go/mdg/internal/cli"
)

const (
	inputFile  = "../../testdata/test1.md"
	goldenFile = "../../testdata/test1.md.html"
)

func TestCLI_PreviewToTempFile(t *testing.T) {
	// Create an in-memory stdout buffer
	var stdout bytes.Buffer

	app := cli.NewApp("v-test")

	args := []string{"cmd", "preview", "--file", inputFile, "--skip-preview"}
	app.Writer = &stdout

	err := app.Run(args)
	if err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(stdout.String())
	t.Logf("Temp file generated: %s", resultFile)

	defer os.Remove(resultFile)

	result, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Errorf("Output does not match golden file")
	}
}
