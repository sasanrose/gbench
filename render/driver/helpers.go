package driver

import (
	"fmt"
	"math"
	"time"

	"github.com/apcera/termtables"
	"github.com/sasanrose/gbench/report"
	"github.com/ttacon/chalk"
)

type tableGenerator struct {
	r *report.Result
}

func (g *tableGenerator) getBenchResultTable() *termtables.Table {
	table := termtables.CreateTable()

	successRate := (float64(g.r.SuccessfulRequests*100) / float64(g.r.TotalRequests))
	failureRate := (float64(g.r.FailedRequests*100) / float64(g.r.TotalRequests))
	timedoutRate := (float64(g.r.TimedOutRequests*100) / float64(g.r.TotalRequests))
	averageResponseTime := time.Duration(g.r.TotalResponseTime.Nanoseconds() / int64(g.r.ResponseTimesTotalCount))
	transferredData := float64(g.r.TotalReceivedDataLength) / math.Pow(2, 20)

	table.AddTitle(g.getColoredString("Final benchmark result", chalk.Blue))
	g.addColoredRow(table, chalk.Cyan, "Start time", g.r.StartTime.Format(time.RFC1123))
	g.addColoredRow(table, chalk.Cyan, "End time", g.r.EndTime.Format(time.RFC1123))
	g.addColoredRow(table, chalk.Cyan, "Total requests sent", g.r.TotalRequests)
	g.addColoredRow(table, chalk.Cyan, "Total data received", fmt.Sprintf("%.5f MB", transferredData))
	g.addColoredRow(table, chalk.Green, "Total successful requests", g.r.SuccessfulRequests)
	g.addColoredRow(table, chalk.Red, "Total failed requests", g.r.FailedRequests)
	g.addColoredRow(table, chalk.Yellow, "Total timedout requests", g.r.TimedOutRequests)
	g.addColoredRow(table, chalk.Green, "Success rate", fmt.Sprintf("%%%.2f", successRate))
	g.addColoredRow(table, chalk.Red, "Failure rate", fmt.Sprintf("%%%.2f", failureRate))
	g.addColoredRow(table, chalk.Yellow, "Timedout rate", fmt.Sprintf("%%%.2f", timedoutRate))
	g.addColoredRow(table, chalk.Cyan, "Total benchmark time", g.r.TotalTime)
	g.addColoredRow(table, chalk.Cyan, "Sum of all response times", g.r.TotalResponseTime)
	g.addColoredRow(table, chalk.Cyan, "Shortest response time", g.r.ShortestResponseTime)
	g.addColoredRow(table, chalk.Cyan, "Longest response time", g.r.LongestResponseTime)
	g.addColoredRow(table, chalk.Cyan, "Average response time", averageResponseTime)

	return table
}

func (g *tableGenerator) getUrlTables() []*termtables.Table {
	urlTables := make([]*termtables.Table, 0)

	for url := range g.r.Urls {
		urlTable := termtables.CreateTable()
		urlTable.AddTitle(g.getColoredString(fmt.Sprintf("Final result for %s", url), chalk.Blue))

		if length, ok := g.r.ReceivedDataLength[url]; ok {
			transferredData := float64(length) / math.Pow(2, 20)
			g.addColoredRow(urlTable, chalk.Cyan, "Total data received", fmt.Sprintf("%.5f MB", transferredData))
		}

		if _, ok := g.r.ResponseStatusCode[url]; ok {
			for statusCode, count := range g.r.ResponseStatusCode[url] {
				g.addColoredRow(urlTable, chalk.Green, fmt.Sprintf("Response with status code %d", statusCode), count)
			}
		}

		if _, ok := g.r.FailedResponseStatusCode[url]; ok {
			for statusCode, count := range g.r.FailedResponseStatusCode[url] {
				g.addColoredRow(urlTable, chalk.Red, fmt.Sprintf("Response with status code %d", statusCode), count)
			}
		}

		averageResponseTime := time.Duration(g.r.ResponseTime[url].Nanoseconds() / int64(g.r.ResponseTimesCount[url]))

		g.addColoredRow(urlTable, chalk.Red, "Failed requests", g.r.FailedResponse[url])
		g.addColoredRow(urlTable, chalk.Yellow, "Timedout requests", g.r.TimedoutResponse[url])
		g.addColoredRow(urlTable, chalk.Cyan, "Sum response times", g.r.ResponseTime[url])
		g.addColoredRow(urlTable, chalk.Cyan, "Shortest response time", g.r.ShortestResponseTimes[url])
		g.addColoredRow(urlTable, chalk.Cyan, "Longest response time", g.r.LongestResponseTimes[url])
		g.addColoredRow(urlTable, chalk.Cyan, "Average response time", averageResponseTime)

		urlTables = append(urlTables, urlTable)
	}

	return urlTables
}

func (g *tableGenerator) getConcurrencyTables() map[int]*termtables.Table {
	concurrencyTables := make(map[int]*termtables.Table)

	for url, concurrencyResults := range g.r.ConcurrencyResult {
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
				g.getColoredString(concurrencyResult.TotalRequests, chalk.Cyan),
				g.getColoredString(concurrencyResult.SuccessfulRequests, chalk.Green),
				g.getColoredString(concurrencyResult.FailedRequests, chalk.Red),
				g.getColoredString(concurrencyResult.TimedOutRequests, chalk.Yellow))
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
