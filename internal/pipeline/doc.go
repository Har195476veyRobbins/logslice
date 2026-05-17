// Package pipeline provides a high-level orchestration layer for logslice.
//
// It connects the individual processing stages — parsing, filtering,
// aggregation and output formatting — into a single, easy-to-use Run
// function that consumes an io.Reader and writes results to an io.Writer.
//
// Basic usage:
//
//	res, err := pipeline.Run(os.Stdin, os.Stdout, pipeline.Config{
//		Format: "text",
//		Filter: filter.Filter{
//			Level:    "error",
//			Keywords: []string{"timeout"},
//		},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("matched %d lines\n", res.LinesWritten)
//
// When AggregateOnly is set to true no log lines are written; only the
// Stats field of the returned Result is populated.
package pipeline
