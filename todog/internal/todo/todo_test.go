package todo_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mnishiguchi/command-line-go/todog/internal/todo"
)

func TestAdd(t *testing.T) {
	var list todo.List

	taskName := "New Task"
	item := list.Add(taskName)

	assert.Equal(t, taskName, item.Task)
	assert.Len(t, list, 1)
	assert.False(t, list[0].Done)
	assert.True(t, list[0].CompletedAt.IsZero()) // Should be zero since not completed
}

func TestComplete(t *testing.T) {
	var list todo.List
	list.Add("Buy groceries")

	err := list.Complete(1)

	assert.NoError(t, err)
	assert.True(t, list[0].Done)
	assert.False(t, list[0].CompletedAt.IsZero())
}

func TestDelete(t *testing.T) {
	list := todo.List{}
	list.Add("Task 1")
	list.Add("Task 2")
	list.Add("Task 3")

	err := list.Delete(2)

	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "Task 1", list[0].Task)
	assert.Equal(t, "Task 3", list[1].Task)
}

func TestSaveGet(t *testing.T) {
	var list1 todo.List
	var list2 todo.List

	taskName := "Write tests"
	list1.Add(taskName)

	tmpFile, err := os.CreateTemp("", "todo-test-")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	err = list1.Save(tmpFile.Name())
	require.NoError(t, err)

	err = list2.Get(tmpFile.Name())
	require.NoError(t, err)

	require.Len(t, list2, 1)
	assert.Equal(t, list1[0].Task, list2[0].Task)
	assert.Equal(t, list1[0].Done, list2[0].Done)
	assert.WithinDuration(t, list1[0].CreatedAt, list2[0].CreatedAt, 2*time.Second)
}
