// Package driver contains different drivers that can be used for rending the
// report of a benchmark.
package driver

import (
	"fmt"
	"io"
	"os"

	"github.com/sasanrose/gbench/render"
	"github.com/sasanrose/gbench/report"
)

type cli struct {
	output io.Writer
}

// NewCli creates a new cli renderer for benchmark report.
func NewCli() render.Renderer {
	return &cli{os.Stdout}
}

// Render will output the result of the report to cli.
func (r *cli) Render(result *report.Result) error {
	tableGen := &tableGenerator{result}
	table := tableGen.getBenchResultTable()
	urlTables := tableGen.getURLTables()
	concurrencyTables := tableGen.getConcurrencyTables()

	fmt.Fprint(r.output, table.Render())

	for _, urlTable := range urlTables {
		fmt.Fprint(r.output, urlTable.Render())
	}

	for i := 0; i < len(concurrencyTables); i++ {
		fmt.Fprint(r.output, concurrencyTables[i].Render())
	}

	return nil
}
