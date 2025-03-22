package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	binName = "todog" // Name of the binary to build
	binPath string    // Full path to the compiled binary (set in TestMain)
)

// TestMain builds the CLI binary before tests run, and cleans it up after.
func TestMain(m *testing.M) {
	fmt.Println("Building binary...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	// Get absolute path to the binary (e.g., ./todog)
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get working dir: %v\n", err)
		os.Exit(1)
	}
	binPath = filepath.Join(cwd, binName)

	// Build the binary from current directory
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build CLI: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Clean up the binary
	fmt.Println("Cleaning up...")
	_ = os.Remove(binPath)

	os.Exit(code)
}

// TestTodoCLI runs basic end-to-end tests against the compiled CLI binary.
func TestTodoCLI(t *testing.T) {
	task := "test task number 1"

	// Create a temporary file for isolated task storage
	tmpFile, err := os.CreateTemp("", "todo-test-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	t.Run("AddNewTask", func(t *testing.T) {
		// Simulate running: todog test task number 1
		cmd := exec.Command(binPath, strings.Split(task, " ")...)
		cmd.Env = append(os.Environ(), "TODOG_FILE="+tmpFile.Name())

		err := cmd.Run()
		require.NoError(t, err)
	})

	t.Run("ListTasks", func(t *testing.T) {
		// Simulate running: todog
		cmd := exec.Command(binPath)
		cmd.Env = append(os.Environ(), "TODOG_FILE="+tmpFile.Name())

		out, err := cmd.CombinedOutput()
		require.NoError(t, err)

		output := string(out)
		expected := fmt.Sprintf("1. [ ] %s\n", task)

		assert.Contains(t, output, expected, "Expected task to appear in output")
	})
}
