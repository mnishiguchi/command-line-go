package formatter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Render(targetPath string) error {
	info, err := os.Stat(targetPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return renderDirectory(targetPath)
	} else {
		return renderFile(targetPath)
	}
}

func renderDirectory(root string) error {
	fmt.Println("└──", filepath.Base(root))
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == root {
			return nil
		}

		rel, _ := filepath.Rel(root, path)
		parts := strings.Split(rel, string(os.PathSeparator))
		prefix := strings.Repeat("    ", len(parts)-1)

		connector := "├──"
		if info.IsDir() {
			fmt.Printf("%s%s %s\n", prefix, connector, info.Name())
		} else {
			fmt.Printf("%s%s %s\n", prefix, connector, info.Name())
		}
		return nil
	})
}

func renderFile(path string) error {
	fmt.Printf("\n%s:\n", path)
	fmt.Println(strings.Repeat("-", 80))

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		fmt.Printf("%2d | %s\n", i+1, line)
	}

	fmt.Println("\n" + strings.Repeat("-", 80))
	return nil
}
