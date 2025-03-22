package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mnishiguchi/command-line-go/todog/internal/todo"
	"github.com/urfave/cli/v2"
)

// Execute sets up the CLI app and runs it.
func Execute(version string) {
	log.SetFlags(0)
	logger := log.New(os.Stderr, "", 0)

	app := &cli.App{
		Name:    "todog",
		Version: version,
		Usage:   "Manage your todo list from the command line",
		Description: `todog is a simple CLI tool for adding and listing tasks.

Examples:
  todog                   Show all tasks
  todog buy groceries     Add a new task

All data is saved to a todo.json file based on environment.`,
		Action: runCli,
	}

	if err := app.Run(os.Args); err != nil {
		logger.Printf("Error: %v", err)
		cli.OsExiter(1)
	}
}

func runCli(c *cli.Context) error {
	args := c.Args().Slice()
	list := &todo.List{}
	file := getTodoFileName()

	// Load tasks from file
	if err := list.Get(file); err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	if len(args) == 0 {
		return printTasks(list)
	}

	return addTask(list, args, file)
}

// getTodoFileName determines where to save/load the todo list.
func getTodoFileName() string {
	// 1. Allow override via environment
	if path := os.Getenv("TODOG_FILE"); path != "" {
		return path
	}

	switch os.Getenv("TODOG_ENV") {
	case "development":
		// Save to local tmp directory
		tmpPath := filepath.Join(".", "tmp")
		_ = os.MkdirAll(tmpPath, 0755)
		return filepath.Join(tmpPath, "todo.json")

	case "test":
		// Enforce explicit file for test mode
		fmt.Fprintln(os.Stderr, "ERROR: TODOG_FILE must be set in test environment.")
		os.Exit(1)
	}

	// Production: save to user's config directory
	if home, err := os.UserHomeDir(); err == nil {
		configPath := filepath.Join(home, ".todog")
		_ = os.MkdirAll(configPath, 0755)
		return filepath.Join(configPath, "todo.json")
	}

	// Fallback (should rarely happen)
	return ".todo.json"
}

// printTasks outputs all tasks in the list.
func printTasks(list *todo.List) error {
	if len(*list) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	for i, item := range *list {
		status := "[ ]"
		if item.Done {
			status = "[x]"
		}
		fmt.Printf("%d. %s %s\n", i+1, status, item.Task)
	}

	return nil
}

// addTask appends a new task and saves the updated list.
func addTask(list *todo.List, args []string, file string) error {
	task := strings.Join(args, " ")
	list.Add(task)

	if err := list.Save(file); err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	fmt.Printf("Added task: %q\n", task)
	return nil
}

