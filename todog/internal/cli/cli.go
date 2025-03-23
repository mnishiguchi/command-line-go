package cli

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
		Flags:   []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:      "list",
				Usage:     "List all tasks",
				UsageText: "todog list",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "verbose",
						Usage:   "Enable verbose output",
						Aliases: []string{"v"},
					},
					&cli.BoolFlag{
						Name:  "hide-completed",
						Usage: "Hide tasks marked as completed",
					},
				},
				Action: func(c *cli.Context) error {
					list, _, err := loadTodoList()
					if err != nil {
						return err
					}

					if len(*list) == 0 {
						fmt.Println("No tasks found.")
						return nil
					}

					verbose := c.Bool("verbose")
					hideCompleted := c.Bool("hide-completed")

					taskCount := 0

					for i, item := range *list {
						if hideCompleted && item.Done {
							continue
						}

						taskCount++

						status := "[ ]"
						if item.Done {
							status = "[x]"
						}
						fmt.Printf("%d. %s %s\n", i+1, status, item.Task)

						if verbose {
							fmt.Printf("    Created:\t%s\n", item.CreatedAt.Format(time.RFC3339))
							if item.Done {
								fmt.Printf("    Completed:\t%s\n", item.CompletedAt.Format(time.RFC3339))
							}
						}
					}

					if taskCount == 0 {
						fmt.Println("No tasks to display.")
					}

					return nil
				},
			},
			{
				Name:      "add",
				Usage:     "Add a new task (from args or stdin)",
				UsageText: "todog add [task description] [--multiline]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "multiline",
						Usage: "Enable multiline STDIN input (one task per line)",
					},
				},
				Action: func(c *cli.Context) error {
					list, file, err := loadTodoList()
					if err != nil {
						return err
					}

					var tasks []string

					if c.Bool("multiline") {
						tasks, err = getTasksMultiline(os.Stdin)
					} else {
						tasks, err = getTask(os.Stdin, c.Args().Slice()...)
					}
					if err != nil {
						return err
					}

					for _, task := range tasks {
						list.Add(task)
						fmt.Printf("Added task: %q\n", task)
					}

					if err := list.Save(file); err != nil {
						return fmt.Errorf("failed to save tasks: %w", err)
					}

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
			{
				Name:      "delete",
				Usage:     "Delete a task by its number",
				UsageText: "todog delete <task number>",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return fmt.Errorf("please provide a task number to delete")
					}

					num, err := strconv.Atoi(c.Args().First())
					if err != nil || num <= 0 {
						return fmt.Errorf("invalid task number: %s", c.Args().First())
					}

					list, file, err := loadTodoList()
					if err != nil {
						return err
					}

					if err := list.Delete(num); err != nil {
						return fmt.Errorf("failed to delete task: %w", err)
					}

					if err := list.Save(file); err != nil {
						return fmt.Errorf("failed to save list: %w", err)
					}

					fmt.Printf("Deleted task #%d.\n", num)
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

func getTask(r io.Reader, args ...string) ([]string, error) {
	if len(args) > 0 {
		return []string{strings.Join(args, " ")}, nil
	}

	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("task input is required")
	}

	line := strings.TrimSpace(scanner.Text())
	if line == "" {
		return nil, fmt.Errorf("task cannot be blank")
	}

	return []string{line}, nil
}

func getTasksMultiline(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var tasks []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			tasks = append(tasks, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("no task input provided")
	}

	return tasks, nil
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
