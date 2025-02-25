package tests

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to run the echog command
func runEchog(args ...string) (string, error) {
	cmd := exec.Command("../build/echog", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestEchogBasic(t *testing.T) {
	output, err := runEchog("Hello there")
	assert.NoError(t, err)
	assert.Equal(t, "Hello there\n", output)
}

func TestEchogMultipleArgs(t *testing.T) {
	output, err := runEchog("Hello", "there")
	assert.NoError(t, err)
	assert.Equal(t, "Hello there\n", output)
}

func TestEchogNoNewline(t *testing.T) {
	output, err := runEchog("-n", "Hello", "there")
	assert.NoError(t, err)
	assert.Equal(t, "Hello there", output) // No newline
}
