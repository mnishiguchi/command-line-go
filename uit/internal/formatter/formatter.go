package formatter

import (
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// TreeNode represents a node in the directory tree.
type TreeNode struct {
	Name     string
	IsFile   bool
	Children map[string]*TreeNode
}

// RenderGitTree builds and prints a Git-tracked file tree to the provided writer.
func RenderGitTree(root string, w io.Writer) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	tree, err := buildGitTree(absRoot)
	if err != nil {
		return fmt.Errorf("failed to build tree from Git files: %w", err)
	}

	printTree(tree, w)
	return nil
}

// buildGitTree constructs a directory tree from Git-tracked files.
func buildGitTree(root string) (*TreeNode, error) {
	cmd := exec.Command("git", "-C", root, "ls-files")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	rootNode := &TreeNode{
		Name:     filepath.Base(root),
		IsFile:   false,
		Children: make(map[string]*TreeNode),
	}

	for _, path := range lines {
		addPath(rootNode, strings.Split(path, "/"))
	}

	return rootNode, nil
}

// addPath inserts a file path (split into parts) into the tree recursively.
func addPath(node *TreeNode, parts []string) {
	if len(parts) == 0 {
		return
	}

	name := parts[0]
	child, exists := node.Children[name]
	if !exists {
		child = &TreeNode{
			Name:     name,
			IsFile:   len(parts) == 1,
			Children: make(map[string]*TreeNode),
		}
		node.Children[name] = child
	}

	addPath(child, parts[1:])
}

// printTree prints the tree starting from the root node.
func printTree(node *TreeNode, w io.Writer) {
	fmt.Fprintf(w, "%s\n", node.Name)
	printChildren(node, "", true, w)
}

// printChildren prints child nodes of the given tree node recursively.
func printChildren(node *TreeNode, prefix string, isLast bool, w io.Writer) {
	_ = isLast // Reserved for future enhancements (e.g., formatting tweaks)

	var keys []string
	for k := range node.Children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, key := range keys {
		child := node.Children[key]

		connector := "├──"
		nextPrefix := prefix + "│   "
		if i == len(keys)-1 {
			connector = "└──"
			nextPrefix = prefix + "    "
		}

		fmt.Fprintf(w, "%s%s %s\n", prefix, connector, child.Name)

		if !child.IsFile {
			printChildren(child, nextPrefix, i == len(keys)-1, w)
		}
	}
}
