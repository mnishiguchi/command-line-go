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
	task1 := "test task number 1"
	task2 := "second task from stdin"

	// Create an isolated file for testing
	tmpFile, err := os.CreateTemp("", "todo-test-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	t.Run("AddNewTaskFromArguments", func(t *testing.T) {
		_, err := runCommand(tmpFile.Name(), "add", task1)
		require.NoError(t, err, "should add a new task from arguments")
	})

	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		_, err := runCommandWithStdin(tmpFile.Name(), task2, "add")
		require.NoError(t, err, "should add a new task from stdin")
	})

	t.Run("ListTasks", func(t *testing.T) {
		output, err := runCommand(tmpFile.Name(), "list")
		require.NoError(t, err, "should list tasks without error")

		expected1 := fmt.Sprintf("1. [ ] %s", task1)
		expected2 := fmt.Sprintf("2. [ ] %s", task2)

		assert.Contains(t, output, expected1, "expected first task in output")
		assert.Contains(t, output, expected2, "expected second task in output")
	})

	t.Run("CompleteTask", func(t *testing.T) {
		_, err := runCommand(tmpFile.Name(), "complete", "1")
		require.NoError(t, err, "should complete task without error")
	})

	t.Run("ListTasksAfterComplete", func(t *testing.T) {
		output, err := runCommand(tmpFile.Name(), "list")
		require.NoError(t, err, "should list tasks after completion without error")

		expected1 := fmt.Sprintf("1. [x] %s", task1)
		expected2 := fmt.Sprintf("2. [ ] %s", task2)

		assert.Contains(t, output, expected1, "expected completed task in output")
		assert.Contains(t, output, expected2, "expected second task to remain incomplete")
	})

	t.Run("DeleteSecondTask", func(t *testing.T) {
		_, err := runCommand(tmpFile.Name(), "delete", "2")
		require.NoError(t, err, "should delete second task")
	})

	t.Run("ListTasksAfterDelete", func(t *testing.T) {
		output, err := runCommand(tmpFile.Name(), "list")
		require.NoError(t, err, "should list tasks after deletion without error")

		assert.Contains(t, output, task1, "first task should remain")
		assert.NotContains(t, output, task2, "second task should be deleted")
	})
}

func runCommand(todoFile string, args ...string) (string, error) {
	cmd := exec.Command(binPath, args...)
	cmd.Env = append(os.Environ(), "TODOG_FILE="+todoFile)

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runCommandWithStdin(todoFile, stdin string, args ...string) (string, error) {
	cmd := exec.Command(binPath, args...)
	cmd.Env = append(os.Environ(), "TODOG_FILE="+todoFile)

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	go func() {
		defer stdinPipe.Close()
		_, _ = stdinPipe.Write([]byte(stdin))
	}()

	out, err := cmd.CombinedOutput()
	return string(out), err
}
