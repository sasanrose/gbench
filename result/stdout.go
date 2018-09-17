package result

import (
	"fmt"
	"os"
)

type stdout struct {
	Result
}

func NewStdoutRenderer() Renderer {
	r := &stdout{}
	r.output = os.Stdout

	return r
}

func (r *stdout) Render() error {
	tableGen := &tableGenerator{&r.Result}
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
