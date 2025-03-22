package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"
)

// Item represents a single to-do task.
type Item struct {
	Task        string    `json:"task"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
}

// List is a collection of to-do items.
type List []Item

// Add creates a new task and appends it to the list.
func (l *List) Add(task string) Item {
	item := Item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, item)
	return item
}

// Complete marks the i-th task as done.
func (l *List) Complete(i int) error {
	if i <= 0 || i > len(*l) {
		return fmt.Errorf("item %d does not exist", i)
	}

	(*l)[i-1].Done = true
	(*l)[i-1].CompletedAt = time.Now()
	return nil
}

// Delete removes the i-th task from the list.
func (l *List) Delete(i int) error {
	if i <= 0 || i > len(*l) {
		return fmt.Errorf("item %d does not exist", i)
	}

	*l = slices.Delete(*l, i-1, i)
	return nil
}

// Save writes the list to a file in JSON format.
func (l *List) Save(filename string) error {
	data, err := json.MarshalIndent(l, "", "  ") // prettier formatting
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// Get reads the list from a JSON file, if it exists.
func (l *List) Get(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// No file yet; treat as empty list
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, l)
}
