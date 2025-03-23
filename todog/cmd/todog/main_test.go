package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	binName = "todog"
	binPath string
)

// TestMain builds the CLI binary before tests run, and cleans it up after.
func TestMain(m *testing.M) {
	fmt.Println("Building binary...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	cwd, err := os.Getwd()
	require.NoError(nil, err, "failed to get working directory")

	binPath = filepath.Join(cwd, binName)

	cmd := exec.Command("go", "build", "-o", binPath, ".")
	require.NoError(nil, cmd.Run(), "failed to build CLI")

	code := m.Run()

	fmt.Println("Cleaning up...")
	_ = os.Remove(binPath)

	os.Exit(code)
}

func TestTodoCLI(t *testing.T) {
	task := "test task number 1"

	// Create an isolated file for testing
	tmpFile, err := os.CreateTemp("", "todo-test-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	t.Run("AddNewTask", func(t *testing.T) {
		_, err := runCommand(tmpFile.Name(), "task", task)
		require.NoError(t, err, "should add a new task without error")
	})

	t.Run("ListTasks", func(t *testing.T) {
		output, err := runCommand(tmpFile.Name(), "list")
		require.NoError(t, err, "should list tasks without error")

		expected := fmt.Sprintf("1. [ ] %s\n", task)
		assert.Contains(t, output, expected, "expected task to appear in output")
	})

	t.Run("CompleteTask", func(t *testing.T) {
		_, err := runCommand(tmpFile.Name(), "complete", "1")
		require.NoError(t, err, "should complete task without error")
	})

	t.Run("ListTasksAfterComplete", func(t *testing.T) {
		output, err := runCommand(tmpFile.Name(), "list")
		require.NoError(t, err, "should list tasks after completion without error")

		expected := fmt.Sprintf("1. [x] %s\n", task)
		assert.Contains(t, output, expected, "expected completed task in output")
	})
}

func runCommand(todoFile string, args ...string) (string, error) {
	cmd := exec.Command(binPath, args...)
	cmd.Env = append(os.Environ(), "TODOG_FILE="+todoFile)

	out, err := cmd.CombinedOutput()
	return string(out), err
}
