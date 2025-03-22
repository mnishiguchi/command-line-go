package filestats

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileStats(t *testing.T) {
	input := "Hello, world!\nThis\tis a test.\r\n123 456 789\n"
	reader := strings.NewReader(input)

	fmt.Println("Byte Analysis:")
	for i, b := range []byte(input) {
		fmt.Printf("  Byte[%d]: %d (%q)\n", i, b, string(b))
	}

	fmt.Println("Rune Analysis:")
	for i, r := range []rune(input) {
		fmt.Printf("  Rune[%d]: %q (Unicode: %U)\n", i, r, r)
	}

	stats, err := GetFileStats(reader)
  assert.NoError(t, err)

	fmt.Printf("Expected bytes: %d, Actual bytes: %d\n", len(input), stats.Bytes)
	fmt.Printf("Expected chars: %d, Actual chars: %d\n", len([]rune(input)), stats.Chars)

	assert.Equal(t, 3, stats.Lines)
	assert.Equal(t, 9, stats.Words)
	assert.Equal(t, len(input), stats.Bytes)
	assert.Equal(t, len([]rune(input)), stats.Chars)
}
