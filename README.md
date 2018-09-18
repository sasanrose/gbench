GBench
======

HTTP Load Testing And Benchmarking Tool inspired by Apache Benchmark and Siege.

## Requirements

You need [GO](https://golang.org) installed and ready on your system.

## Installation

```bash
git clone git@github.com/sasanrose/gbench
cd gbench
go build cmd/* -o gbench
```

## Usage of Gbench:                                                                                                            
Command line parameters:
      --auth-password string          Password for basic HTTP authentication.                                                                                                                
      --auth-username string          Username for basic HTTP authentication.
  -c, --concurrent int                Number of concurrent requests. (default 1)
  -C, --connection-timeout duration   Connection timeout (0 means no timeout).
  -f, --file string                   Path to the file containing list of urls to benchmark.
  -H, --header strings                HTTP header in format of key=value. This can be used multiple times.
      --proxy-url string              Proxy server url.
      --raw-cookie string             A string to be sent as raw cookie (In the format of Set-Cookie HTTP header).
  -r, --requests int                  Number of requests to send. (default 1)
  -R, --response-timeout duration     Response timeout (0 means no timeout).
  -s, --status-codes ints             Define what should be considered as a successful status code. (default [200,202,201])
  -u, --url strings                   Url to benchmark. This can be used multiple times.
  -v, --verbose                       Turn on verbosity mode.

URL format should be as follow:

METHOD|URL|POSTDATA or METHOD|URL

Sample URLs: GET|www.google.com?search=test or POST|www.google.com|search=test or HEAD|www.google.com

Examples:

gbench -file ~/benchmarkurl.txt -r 100 -c 10 -v
gbench -url 'GET|www.google.com' -url 'GET|www.google.com/path2' -r 100 -c 10 -v
