GBench
======

HTTP Load Testing And Benchmarking Tool inspired by Apache Benchmark and Siege.

## Requirements

You need [GO](https://golang.org) installed and ready on your system.

## Installation

```bash
git clone git@github.com/sasanrose/gbench
cd gbench
go build
```

## Usage of GBench

-c: The number of concurrent requests (Default is 1) - Shorthand

-co: Custom cookies. This field is repeatable. - Shorthand

-concurrent: The number of concurrent requests (Default is 1)

-connection-timeout: Specifies the amount of time to wait for establishment of connection (Default Zero)

-cookie: Custom cookies. This field is repeatable.

-ct: Specifies the amount of time to wait for establishment of connection (Default Zero) - Shorthand

-d: Prevents re-use of TCP connections (Default is False) - Shorthand

-de: A random interval in seconds between 0 and the defined number (Default is 0) - Shorthand

-delay: A random interval in seconds between 0 and the defined number (Default is 0)

-disable-keep-alive: Prevents re-use of TCP connections (Default is False)

-f: File containing urls to benchmark. This field is repeatable. - Shorthand

-file: File containing urls to benchmark. This field is repeatable.

-he: Additional header information to send. This field is repeatable. - Shorthand

-header: Additional header information to send. This field is repeatable.

-p: Proxy server to use. - Shorthand

-pass: Password for basic HTTP Authentication - Shorthand

-password: Password for basic HTTP Authentication

-proxy: Proxy server to use.

-r: The number of total requests (Default is 0) - Shorthand

-requests: The number of total requests (Default is 0)

-response-timeout: Specifies the amount of time to wait for a server's response headers (Default Zero)

-t: Specifies the amount of time to wait for a server's response headers (Default Zero) - Shorthand

-u: Single URL to benchmark. This field is repeatable. - Shorthand

-url: Single URL to benchmark. This field is repeatable.

-user: Username for basic HTTP Authentication - Shorthand

-username: Username for basic HTTP Authentication

-v: Verbosity (Default is False) - Shorthand

-verbose: Verbosity (Default is False)

### URL format should be as follow

METHOD|URL|POSTDATA or METHOD|URL

Sample URLs: GET|www.google.com?search=test or POST|www.google.com|search=test or HEAD|www.google.com

### Examples

gbench -file ~/benchmarkurl.txt -r 100 -c 10 -v

gbench -url 'GET|www.google.com' -url 'GET|www.yahoo.com' -r 100 -c 10 -v

## Sample Output

```bash
Benchmark Result:
Total Transactions: 694
Failed Transactions: 36
Availability: 94.812680 %
Elapsed Time: 23.899729 secs
Transaction Rate: 29.037987
Average Response Time: 0.034438 secs
Longest Response Time: 23.809487 secs
Shortest Response Time: 0.138656 secs
Transferred Data: 0.055639 MB


Longest Response Times for each URL:
http://sample.com/api/v11/profile/me/l1p05x282EF6C6CA35652A1FFC860D1DE6E0A34Cjnraveaf/username/l3yhjbwj74D51411558C8C0AECC643FB9D7FAA0551ts0x9i: 6.205469 secs
http://sample.com/api/v11/timeline/hash/54a841d7d1ee7b6424d750efc26659e6946d28fd: 23.809487 secs
http://sample.com/api/v11/run/game/com.leagem.chesslive: 23.606921 secs
http://sample.com/api/v11/searchfriend/hash/54a841d7d1ee7b6424d750efc26659e6946d28fd?username=mbdp3vd9665EB4DF966EEB700D27BDA2B7F6BBD5erzbpe8u: 21.108662 secs
http://sample.com/api/v11/notificationscount/username/mbdp3vd9665EB4DF966EEB700D27BDA2B7F6BBD5erzbpe8u/hash/54a841d7d1ee7b6424d750efc26659e6946d28fd: 23.700143 secs


Response Status Codes Stats:
200 OK: 658
500 Internal Server Error: 4


Failed Url Stats:
http://sample.com/api/v11/searchfriend/hash/54a841d7d1ee7b6424d750efc26659e6946d28fd?username=mbdp3vd9665EB4DF966EEB700D27BDA2B7F6BBD5erzbpe8u: 3
http://sample.com/api/v11/timeline/hash/54a841d7d1ee7b6424d750efc26659e6946d28fd: 1
Got Signal: interrupt
```
