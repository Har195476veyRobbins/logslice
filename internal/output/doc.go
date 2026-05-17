// Package output provides formatting utilities for log entries.
//
// It supports two output formats:
//
//   - Text: human-readable, optionally coloured output with aligned fields.
//   - JSON: structured JSON output suitable for machine consumption or
//     piping into other tools.
//
// # Usage
//
//	formatter := output.New(output.Options{
//		Format: output.FormatText,
//		Color:  true,
//	})
//
//	for _, entry := range entries {
//		line, err := formatter.Format(entry)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Println(line)
//	}
//
// Formatters are safe for concurrent use.
package output
