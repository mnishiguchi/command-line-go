package md_test

import (
	"os"
	"testing"

	"github.com/mnishiguchi/command-line-go/mdg/internal/md"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	inputFile  = "../../testdata/test1.md"
	goldenFile = "../../testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	require.NoError(t, err, "should read input markdown file")

	result, err := md.ParseContent(input, "")
	require.NoError(t, err, "should parse markdown content")

	expected, err := os.ReadFile(goldenFile)
	require.NoError(t, err, "should read golden HTML file")

	assert.Equal(t, string(expected), string(result), "parsed HTML should match golden output")
}

