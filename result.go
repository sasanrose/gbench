package main

import (
    "fmt"
)

func showResult() {
    fmt.Println("\n\nBenchmark Result:");
    fmt.Printf("Total Transactions: %v\n", totalTransactions);
    fmt.Printf("Failed Transactions: %v\n", failedTransactions);
    fmt.Printf("Availability: %f %%\n", 100 - ((float64(failedTransactions) * 100) / float64(totalTransactions)));
    fmt.Printf("Elapsed Time: %f secs\n", totalResponseTime.Seconds());
    fmt.Printf("Transaction Rate: %f\n", transactionRate);
    fmt.Printf("Average Response Time: %f secs\n", averageResponseTime);
    fmt.Printf("Longest Response Time: %f secs\n", longestResponseTime.Seconds());
    fmt.Printf("Shortest Response Time: %f secs\n", shortestResponseTime.Seconds());
    fmt.Printf("Transferred Data: %f MB\n", transferredData);

    fmt.Println("\n\nLongest Response Times for each URL:");
    for url, time := range urlsResponseTimes {
        fmt.Printf("%v: %f secs\n", url, time.Seconds());
    }

    fmt.Println("\n\nResponse Status Codes Stats:");
    for code, count := range responseStats {
        fmt.Printf("%v: %v\n", code, count);
    }

    if (len(urlFailedStats) > 0) {
        fmt.Println("\n\nFailed Url Stats:");
        for url, count := range urlFailedStats {
            fmt.Printf("%v: %v\n", url, count);
        }
    }
}
