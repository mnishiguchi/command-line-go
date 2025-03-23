package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mnishiguchi/command-line-go/todog/internal/todo"
	"github.com/urfave/cli/v2"
)

func Execute(version string) {
	log.SetFlags(0)
	logger := log.New(os.Stderr, "", 0)

	app := &cli.App{
		Name:    "todog",
		Version: version,
		Usage:   "Manage your todo list from the command line",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List all tasks",
				UsageText: "todog list",
				Action: func(c *cli.Context) error {
					list, _, err := loadTodoList()
					if err != nil {
						return err
					}

					if len(*list) == 0 {
						fmt.Println("No tasks found.")
						return nil
					}

					printTasks(list)

					return nil
				},
			},
			{
				Name:      "task",
				Usage:     "Add a new task",
				UsageText: "todog task <task description>",
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return fmt.Errorf("please provide a task description")
					}

					list, file, err := loadTodoList()
					if err != nil {
						return err
					}

					task := strings.Join(c.Args().Slice(), " ")
					list.Add(task)

					if err := list.Save(file); err != nil {
						return fmt.Errorf("failed to save task: %w", err)
					}

					fmt.Printf("Added task: %q\n", task)
					return nil
				},
			},
			{
				Name:      "complete",
				Usage:     "Mark a task as complete",
				UsageText: "todog complete <task number>",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return fmt.Errorf("please provide a task number to complete")
					}

					num, err := strconv.Atoi(c.Args().First())
					if err != nil || num <= 0 {
						return fmt.Errorf("invalid task number: %s", c.Args().First())
					}

					list, file, err := loadTodoList()
					if err != nil {
						return err
					}

					if err := list.Complete(num); err != nil {
						return fmt.Errorf("failed to complete task: %w", err)
					}

					if err := list.Save(file); err != nil {
						return fmt.Errorf("failed to save list: %w", err)
					}

					fmt.Printf("Marked task #%d as completed.\n", num)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Printf("Error: %v", err)
		cli.OsExiter(1)
	}
}

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

func loadTodoList() (*todo.List, string, error) {
	file := getTodoFileName()
	list := &todo.List{}
	if err := list.Get(file); err != nil {
		return nil, "", fmt.Errorf("failed to load tasks: %w", err)
	}
	return list, file, nil
}

func getTodoFileName() string {
	if path := os.Getenv("TODOG_FILE"); path != "" {
		return path
	}

	switch os.Getenv("TODOG_ENV") {
	case "development":
		tmpPath := filepath.Join(".", "tmp")
		_ = os.MkdirAll(tmpPath, 0755)
		return filepath.Join(tmpPath, "todo.json")
	case "test":
		fmt.Fprintln(os.Stderr, "ERROR: TODOG_FILE must be set in test environment.")
		os.Exit(1)
	}

	if home, err := os.UserHomeDir(); err == nil {
		configPath := filepath.Join(home, ".todog")
		_ = os.MkdirAll(configPath, 0755)
		return filepath.Join(configPath, "todo.json")
	}

	return ".todo.json"
}
