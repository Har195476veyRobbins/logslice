// Package tail implements head/tail line-slicing over an io.Reader.
//
// It provides two functions:
//
//   - Lines: returns the last N lines from the input, behaving like
//     the Unix `tail -n` command.
//
//   - Head: returns the first N lines from the input, behaving like
//     the Unix `head -n` command.
//
// Both functions operate on any io.Reader, making them suitable for
// use with files, network streams, or in-memory buffers. They buffer
// only as many lines as requested, keeping memory usage bounded.
//
// Example usage:
//
//	lines, err := tail.Lines(os.Stdin, 100)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, l := range lines {
//	    fmt.Println(l)
//	}
package tail
