package md_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/mnishiguchi/command-line-go/mdg/internal/md"
)

const (
	inputFile  = "../../testdata/test1.md"
	goldenFile = "../../testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	result, err := md.ParseContent(input, "")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}
}
