package result

import (
	"fmt"
	"math"
	"time"

	"github.com/apcera/termtables"
	"github.com/ttacon/chalk"
)

type tableGenerator struct {
	r *Result
}

func (g *tableGenerator) getBenchResultTable() *termtables.Table {
	table := termtables.CreateTable()

	successRate := (float64(g.r.successfulRequests*100) / float64(g.r.totalRequests))
	failureRate := (float64(g.r.failedRequests*100) / float64(g.r.totalRequests))
	timedoutRate := (float64(g.r.timedOutRequests*100) / float64(g.r.totalRequests))
	averageResponseTime := time.Duration(g.r.totalResponseTime.Nanoseconds() / int64(g.r.responseTimesTotalCount))
	transferredData := float64(g.r.totalReceivedDataLength) / math.Pow(2, 20)

	table.AddTitle(g.getColoredString("Final benchmark result", chalk.Blue))
	g.addColoredRow(table, chalk.Cyan, "Total requests sent", g.r.totalRequests)
	g.addColoredRow(table, chalk.Cyan, "Total data received", fmt.Sprintf("%.5f MB", transferredData))
	g.addColoredRow(table, chalk.Green, "Total successful requests", g.r.successfulRequests)
	g.addColoredRow(table, chalk.Red, "Total failed requests", g.r.failedRequests)
	g.addColoredRow(table, chalk.Yellow, "Total timedout requests", g.r.timedOutRequests)
	g.addColoredRow(table, chalk.Green, "Success rate", fmt.Sprintf("%%%.2f", successRate))
	g.addColoredRow(table, chalk.Red, "Failure rate", fmt.Sprintf("%%%.2f", failureRate))
	g.addColoredRow(table, chalk.Yellow, "Timedout rate", fmt.Sprintf("%%%.2f", timedoutRate))
	g.addColoredRow(table, chalk.Cyan, "Total benchmark time", g.r.totalTime)
	g.addColoredRow(table, chalk.Cyan, "Sum of all response times", g.r.totalResponseTime)
	g.addColoredRow(table, chalk.Cyan, "Shortest response time", g.r.shortestResponseTime)
	g.addColoredRow(table, chalk.Cyan, "Longest response time", g.r.longestResponseTime)
	g.addColoredRow(table, chalk.Cyan, "Average response time", averageResponseTime)

	return table
}

func (g *tableGenerator) getUrlTables() []*termtables.Table {
	urlTables := make([]*termtables.Table, 0)

	for url, _ := range g.r.urls {
		urlTable := termtables.CreateTable()
		urlTable.AddTitle(g.getColoredString(fmt.Sprintf("Final result for %s", url), chalk.Blue))

		if length, ok := g.r.receivedDataLength[url]; ok {
			transferredData := float64(length) / math.Pow(2, 20)
			g.addColoredRow(urlTable, chalk.Cyan, "Total data received", fmt.Sprintf("%.5f MB", transferredData))
		}

		if _, ok := g.r.responseStatusCode[url]; ok {
			for statusCode, count := range g.r.responseStatusCode[url] {
				g.addColoredRow(urlTable, chalk.Green, fmt.Sprintf("Response with status code %d", statusCode), count)
			}
		}

		if _, ok := g.r.failedResponseStatusCode[url]; ok {
			for statusCode, count := range g.r.failedResponseStatusCode[url] {
				g.addColoredRow(urlTable, chalk.Red, fmt.Sprintf("Response with status code %d", statusCode), count)
			}
		}

		averageResponseTime := time.Duration(g.r.responseTime[url].Nanoseconds() / int64(g.r.responseTimesCount[url]))

		g.addColoredRow(urlTable, chalk.Red, "Failed requests", g.r.failedResponse[url])
		g.addColoredRow(urlTable, chalk.Yellow, "Timedout requests", g.r.timedoutResponse[url])
		g.addColoredRow(urlTable, chalk.Cyan, "Sum response times", g.r.responseTime[url])
		g.addColoredRow(urlTable, chalk.Cyan, "Shortest response time", g.r.shortestResponseTimes[url])
		g.addColoredRow(urlTable, chalk.Cyan, "Longest response time", g.r.longestResponseTimes[url])
		g.addColoredRow(urlTable, chalk.Cyan, "Average response time", averageResponseTime)

		urlTables = append(urlTables, urlTable)
	}

	return urlTables
}

func (g *tableGenerator) getConcurrencyTables() map[int]*termtables.Table {
	concurrencyTables := make(map[int]*termtables.Table)

	for url, concurrencyResults := range g.r.concurrencyResult {
		for index, concurrencyResult := range concurrencyResults {
			if _, ok := concurrencyTables[index]; !ok {
				concurrencyTables[index] = termtables.CreateTable()
				concurrencyTables[index].AddTitle(g.getColoredString(fmt.Sprintf("Result for concurrent requests batch %d", index+1), chalk.Blue))
				concurrencyTables[index].AddHeaders(g.getColoredString("Url", chalk.Cyan))
				concurrencyTables[index].AddHeaders(g.getColoredString("Total", chalk.Cyan))
				concurrencyTables[index].AddHeaders(g.getColoredString("Success", chalk.Green))
				concurrencyTables[index].AddHeaders(g.getColoredString("Failed", chalk.Red))
				concurrencyTables[index].AddHeaders(g.getColoredString("Timedout", chalk.Yellow))
			}

			concurrencyTables[index].AddRow(g.getColoredString(url, chalk.Cyan),
				g.getColoredString(concurrencyResult.totalRequests, chalk.Cyan),
				g.getColoredString(concurrencyResult.successfulRequests, chalk.Green),
				g.getColoredString(concurrencyResult.failedRequests, chalk.Red),
				g.getColoredString(concurrencyResult.timedOutRequests, chalk.Yellow))
		}
	}

	return concurrencyTables
}

func (g *tableGenerator) addColoredRow(table *termtables.Table, color chalk.Color, values ...interface{}) {
	cells := []string{}

	for _, value := range values {
		cells = append(cells, g.getColoredString(value, color))
	}

	args := make([]interface{}, len(cells))

	for k, v := range cells {
		args[k] = v
	}

	table.AddRow(args...)
}

func (g *tableGenerator) getColoredString(value interface{}, color chalk.Color) string {
	return fmt.Sprintf("%s%v%s", color, value, chalk.Reset)
}
