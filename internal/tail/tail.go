// Package tail provides functionality for reading the last N lines
// from a log stream or file, similar to the Unix tail command.
package tail

import (
	"bufio"
	"fmt"
	"io"
)

// ErrInvalidCount is returned when the requested line count is invalid.
var ErrInvalidCount = fmt.Errorf("tail: line count must be greater than zero")

// Lines reads all lines from r and returns the last n lines.
// If n is greater than the total number of lines, all lines are returned.
// Returns ErrInvalidCount if n < 1.
func Lines(r io.Reader, n int) ([]string, error) {
	if n < 1 {
		return nil, ErrInvalidCount
	}

	scanner := bufio.NewScanner(r)
	buf := make([]string, 0, n)

	for scanner.Scan() {
		line := scanner.Text()
		if len(buf) >= n {
			// Shift left: drop oldest entry
			copy(buf, buf[1:])
			buf[n-1] = line
		} else {
			buf = append(buf, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("tail: scanner error: %w", err)
	}

	return buf, nil
}

// Head reads all lines from r and returns the first n lines.
// If n is greater than the total number of lines, all lines are returned.
// Returns ErrInvalidCount if n < 1.
func Head(r io.Reader, n int) ([]string, error) {
	if n < 1 {
		return nil, ErrInvalidCount
	}

	scanner := bufio.NewScanner(r)
	buf := make([]string, 0, n)

	for scanner.Scan() {
		if len(buf) >= n {
			break
		}
		buf = append(buf, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("tail: scanner error: %w", err)
	}

	return buf, nil
}
