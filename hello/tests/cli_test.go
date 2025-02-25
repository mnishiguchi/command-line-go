package tests

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloCLI(t *testing.T) {
	cmd := exec.Command("../build/hello")
	output, err := cmd.CombinedOutput()

	assert.Nil(t, err)
	assert.Contains(t, string(output), "元氣が一番")
}
