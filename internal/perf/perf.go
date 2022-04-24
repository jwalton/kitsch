// Package perf contains tools for recording the performance of an operation or operations.
package perf

import (
	"fmt"
	"strings"
	"time"

	"github.com/jwalton/gchalk"
)

// Record is used to store stats about execution times.
type Record struct {
	// Description is a description of what is being measured.
	Description string
	// Duration is the time taken to execute this item.
	Duration time.Duration
	// Children is the time taken to execute each of this item's children.
	Children []Record
}

// Performance is used to record the performance of a set of operations.
type Performance struct {
	// Records is a list of PerfRecords.
	Records []Record
	// startTime is a map of start times.
	startTime map[string]time.Time
	lastStart time.Time
}

// New creates a new Performance object, initialized to hold size records.
func New(size int) *Performance {
	return &Performance{
		Records:   make([]Record, 0, size),
		startTime: make(map[string]time.Time, size),
		lastStart: time.Now(),
	}
}

// Start starts recording performance for the given item.
func (p *Performance) Start(description string) {
	p.lastStart = time.Now()
	p.startTime[description] = p.lastStart
}

// End finishes recording performance for the given item.
// If this item was never `Start()`ed, this will use the time since the
// last call to `Start()` or `End()`.`
func (p *Performance) End(description string) {
	p.EndWithChildren(description, nil)
}

// EndWithChildren finishes recording performance for the given item, and adds the given children.
func (p *Performance) EndWithChildren(description string, childPerfs *Performance) {
	startTime, ok := p.startTime[description]
	if !ok {
		startTime = p.lastStart
	}

	duration := time.Since(startTime)
	p.Add(description, duration, childPerfs)
	p.lastStart = time.Now()
}

// Add adds execution time for an item to this Performance object.
// `description` is a unique name for this item, `duration` is the time the item
// took to execute, and `childPerfs` are any child items that were executed.
func (p *Performance) Add(
	description string,
	duration time.Duration,
	childPerfs *Performance,
) {
	var children []Record
	if childPerfs != nil {
		children = childPerfs.Records
	}

	p.Records = append(p.Records, Record{
		Description: description,
		Duration:    duration,
		Children:    children,
	})
}

// Print will print the performance data to stdout.
func (p *Performance) Print() {
	renderPerf(p.Records, 0)
}

func renderPerf(durations []Record, indent int) {
	for _, duration := range durations {
		printDuration := duration.Duration.String()
		if duration.Duration > 250000000 {
			printDuration = gchalk.Red(printDuration)
		} else if duration.Duration > 1000000 {
			printDuration = gchalk.Yellow(printDuration)
		} else {
			printDuration = gchalk.Green(printDuration)
		}

		fmt.Printf("%s%s - %s\n",
			strings.Repeat(" ", indent),
			duration.Description,
			printDuration,
		)
		if len(duration.Children) > 0 {
			renderPerf(duration.Children, indent+2)
		}
	}
}
