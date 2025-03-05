package internal

import (
	"bufio"
	"io"
	"strings"
)

// FileStats represents the line, word, byte, and character counts for a file.
type FileStats struct {
	Lines int
	Words int
	Bytes int
	Chars int
}

// GetFileStats reads from an `io.Reader` and calculates statistics.
// This makes it testable without needing actual files.
func GetFileStats(r io.Reader) (*FileStats, error) {
	stats := &FileStats{}
	reader := bufio.NewReader(r)

	// Read the input line by line, preserving newlines
	for {
		line, err := reader.ReadString('\n') // Read up to and including `\n`

		// Handle errors properly
		if err != nil {
			// If EOF and no remaining data, exit loop
			if err == io.EOF && line == "" {
				break
			}
			return nil, err
		}

		// Increment line count
		stats.Lines++

		// Count words using space-based tokenization
		stats.Words += len(strings.Fields(line))

		// Count characters (runes) correctly, including multibyte UTF-8
		runes := []rune(line)
		stats.Chars += len(runes)

		// Count bytes exactly, including newlines
		stats.Bytes += len(line)
	}

	return stats, nil
}
