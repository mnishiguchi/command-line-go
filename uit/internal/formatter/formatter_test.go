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
