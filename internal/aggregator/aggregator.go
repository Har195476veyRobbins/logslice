// Package aggregator provides log entry counting and grouping by level or field.
package aggregator

import (
	"sort"

	"github.com/logslice/logslice/internal/parser"
)

// LevelCount holds the count of log entries for a given level.
type LevelCount struct {
	Level string
	Count int
}

// Summary holds aggregated statistics over a set of log entries.
type Summary struct {
	Total      int
	ByLevel    []LevelCount
	ByField    map[string]map[string]int
}

// Aggregate computes a Summary from a slice of parsed log entries.
func Aggregate(entries []parser.Entry) Summary {
	levelMap := make(map[string]int)
	fieldMap := make(map[string]map[string]int)

	for _, e := range entries {
		level := e.Level
		if level == "" {
			level = "UNKNOWN"
		}
		levelMap[level]++

		for k, v := range e.Fields {
			if _, ok := fieldMap[k]; !ok {
				fieldMap[k] = make(map[string]int)
			}
			fieldMap[k][v]++
		}
	}

	byLevel := make([]LevelCount, 0, len(levelMap))
	for lvl, cnt := range levelMap {
		byLevel = append(byLevel, LevelCount{Level: lvl, Count: cnt})
	}
	sort.Slice(byLevel, func(i, j int) bool {
		if byLevel[i].Count != byLevel[j].Count {
			return byLevel[i].Count > byLevel[j].Count
		}
		return byLevel[i].Level < byLevel[j].Level
	})

	return Summary{
		Total:   len(entries),
		ByLevel: byLevel,
		ByField: fieldMap,
	}
}
