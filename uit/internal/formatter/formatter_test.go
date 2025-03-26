package formatter_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/mnishiguchi/command-line-go/uit/internal/formatter"
	"github.com/stretchr/testify/assert"
)

func TestRenderGitTree(t *testing.T) {
	var buf bytes.Buffer

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	err = formatter.RenderGitTree(cwd, &buf)
	assert.NoError(t, err)

	output := buf.String()

	t.Run("prints root directory name", func(t *testing.T) {
		expectedRoot := filepath.Base(cwd)
		assert.Contains(t, output, expectedRoot)
	})

	t.Run("includes known file", func(t *testing.T) {
		assert.Contains(t, output, "formatter.go")
	})
}

func TestRenderFileContent(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "binary.dat")

	// Write null bytes to simulate binary
	err := os.WriteFile(tmpFile, []byte{0x00, 0x01, 0x02, 0xFF}, 0644)
	assert.NoError(t, err)

	t.Run("skips binary by default", func(t *testing.T) {
		var buf bytes.Buffer
		err := formatter.RenderFileContent(tmpFile, &buf, false)
		assert.NoError(t, err)
		assert.Empty(t, buf.String(), "binary output should be empty by default")
	})

	t.Run("shows binary if flag is true", func(t *testing.T) {
		var buf bytes.Buffer
		err := formatter.RenderFileContent(tmpFile, &buf, true)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "/binary.dat")
	})
}
