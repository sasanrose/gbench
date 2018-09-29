package render

import (
	"fmt"
	"io"
	"os"

	"github.com/sasanrose/gbench/report"
)

type stdout struct {
	output io.Writer
}

func NewStdoutRenderer() Renderer {
	return &stdout{os.Stdout}
}

func (r *stdout) Render(result *report.Result) error {
	tableGen := &tableGenerator{result}
	table := tableGen.getBenchResultTable()
	urlTables := tableGen.getUrlTables()
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
