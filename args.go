package main

import (
    "os"
    "fmt"
    "bufio"
    "regexp"
    "strings"
    "time"
    "net/url"
    "flag"
    "errors"
)

type singleUrlsList []Url;
var tempUrls singleUrlsList;

var Usage = func() {
    fmt.Println("Yata");
}

func (u *singleUrlsList) String() string {
    return fmt.Sprint(*u);
}

func (u *singleUrlsList) Set(singleUrl string) error {
    if (!urlExists(singleUrl)) {
        parsedUrl, err := parseUrl(singleUrl);

        if (err != nil) {
            return err;
        } else {
            *u = append(*u, parsedUrl);
        }
    } else {
        return errors.New("Duplicated Url: "+singleUrl);
    }

    return nil;
}

func (u *urlsList) String() string {
    return fmt.Sprint(*u);
}

func (u *urlsList) Set(file string) error {
    urlFile, err = os.Open(file);

    if (err != nil) {
        return errors.New("Error opening file: " + err.Error());
    }

    defer urlFile.Close();

    scanner := bufio.NewScanner(urlFile)

    for scanner.Scan() {
        url := scanner.Text();
        if (!urlExists(url)) {
            parsedUrl, err := parseUrl(url);

            if (err != nil) {
                return err;
            } else {
                *u = append(*u, parsedUrl);
            }
        }
    }

    return nil;
}

func (c *cookiesList) String() string {
    return fmt.Sprint(*c);
}

func (c *cookiesList) Set(cookie string) error {
    cookie = strings.Trim(cookie, "\"' ");

    *c = append(*c, cookie);

    return nil;
}

func (h *headersList) String() string {
    return fmt.Sprint(*h);
}

func (h *headersList) Set(header string) error {
    parts := strings.Split(header, ":");

    if (len(parts) != 2) {
        return errors.New("Wrong header format. Example: 'accept-encoding: gzip, deflate'");
    }

    key := strings.Trim(parts[0], "\"' ");
    value := strings.Trim(parts[1], "\"' ");

    (*h)[key] = value;

    return nil;
}

func (p *proxyURL) String() string {
    return fmt.Sprint(*p);
}

func (p *proxyURL) Set(enteredUrl string) error {
    enteredUrl = strings.Trim(enteredUrl, "\" ");

    var validArg = regexp.MustCompile("^localhost");

    if (enteredUrl == "localhost" || validArg.MatchString(enteredUrl)) {
        enteredUrl = "http://" + enteredUrl;
    }

    tempProxy, err := url.Parse(enteredUrl);

    if (err != nil) {
        return errors.New("Wrong proxy format. Example: http://localhost:8441");
    }

    if (tempProxy.Scheme == "") {
        tempProxy.Scheme = "http";
    }

    *p = proxyURL(tempProxy.String());

    return nil;
}

func init() {
    // Make HTTP Custom headers map
    headers = make(headersList);

    // Number of concurrent requests
    defaultUsage := "The number of concurrent requests (Default is 1)";
    flag.IntVar(&concurrent, "concurrent", 1, defaultUsage);
    flag.IntVar(&concurrent, "c", 1, defaultUsage + " - Shorthand");

    // Number of total requests to be sent
    defaultUsage = "The number of total requests to send (Default is 1)";
    flag.IntVar(&requests, "requests", 1, defaultUsage);
    flag.IntVar(&requests, "r", 1, defaultUsage + " - Shorthand");

    // To be verbose or not to be
    defaultUsage = "Verbosity (Default is False)";
    flag.BoolVar(&verbose, "verbose", false, defaultUsage);
    flag.BoolVar(&verbose, "v", false, defaultUsage + " - Shorthand");

    // To be verbose or not to be
    defaultUsage = "Prevents re-use of TCP connections (Default is False)";
    flag.BoolVar(&disableKeepAlive, "disable-keep-alive", false, defaultUsage);
    flag.BoolVar(&disableKeepAlive, "d", false, defaultUsage + " - Shorthand");

    // Random interval
    defaultUsage = "A random interval in seconds between 0 and the defined number (Default is 0)";
    flag.IntVar(&delay, "delay", 0, defaultUsage);
    flag.IntVar(&delay, "de", 0, defaultUsage + " - Shorthand");

    // Username for basic HTTP Authentication
    defaultUsage = "Username for basic HTTP Authentication";
    flag.StringVar(&basicAuthUsername, "username", "", defaultUsage);
    flag.StringVar(&basicAuthUsername, "user", "", defaultUsage + " - Shorthand");

    // Password for basic HTTP Authentication
    defaultUsage = "Password for basic HTTP Authentication";
    flag.StringVar(&basicAuthPassword, "password", "", defaultUsage);
    flag.StringVar(&basicAuthPassword, "pass", "", defaultUsage + " - Shorthand");

    // Response header timeout
    var tempResponseTimeout int;
    defaultUsage = "Specifies the amount of time to wait for a server's response headers (Default is 0 which means no timeout)";
    flag.IntVar(&tempResponseTimeout, "response-timeout", 0, defaultUsage);
    flag.IntVar(&tempResponseTimeout, "t", 0, defaultUsage + " - Shorthand");
    responseTimeout = time.Duration(tempResponseTimeout) * time.Second;

    // Connection timeout
    var tempConnectionTimeout int;
    defaultUsage = "Specifies the amount of time to wait for establishment of connection (Default is 0 which means no timeout)";
    flag.IntVar(&tempConnectionTimeout, "connection-timeout", 0, defaultUsage);
    flag.IntVar(&tempConnectionTimeout, "ct", 0, defaultUsage + " - Shorthand");
    connectionTimeout = time.Duration(tempConnectionTimeout) * time.Second;

    // Benchmarking url file
    defaultUsage = "Address of file containing urls to benchmark. This field is repeatable.";
    flag.Var(&urls, "file", defaultUsage);
    flag.Var(&urls, "f", defaultUsage + " - Shorthand");

    // Single URL to benchmark
    defaultUsage = "Single URL to benchmark. This field is repeatable.";
    flag.Var(&tempUrls, "url", defaultUsage);
    flag.Var(&tempUrls, "u", defaultUsage + " - Shorthand");

    // Custome cookies
    defaultUsage = "Custom cookies. This field is repeatable.";
    flag.Var(&cookies, "cookie", defaultUsage);
    flag.Var(&cookies, "co", defaultUsage + " - Shorthand");

    // Custom Headers
    defaultUsage = "Additional header information to send. This field is repeatable.";
    flag.Var(&headers, "header", defaultUsage);
    flag.Var(&headers, "he", defaultUsage + " - Shorthand");

    // HTTP proxy
    defaultUsage = "Proxy server to use.";
    flag.Var(&proxyUrl, "proxy", defaultUsage);
    flag.Var(&proxyUrl, "p", defaultUsage + " - Shorthand");

    flag.Usage = func() {
        flagDefaults();
    }

    flag.Parse();

    for _, singleUrl := range tempUrls {
        urls = append(urls, singleUrl);
    }

    if (requests == 0 || concurrent == 0 || requests < concurrent || len(urls) == 0) {
        flagDefaults();
        os.Exit(1);
    }
}

func flagDefaults () {
    fmt.Fprintf(os.Stderr, "\nUsage of gbnech:\n");
    flag.VisitAll(func(f *flag.Flag) {
        fmt.Printf("-%v: %v\n", f.Name, f.Usage);
    });

    fmt.Println("\nURL format should be as follow:");
    fmt.Println("METHOD|URL|POSTDATA or METHOD|URL");
    fmt.Println("Sample URLs: GET|www.google.com?search=test or POST|www.google.com|search=test or HEAD|www.google.com");
    fmt.Println("\nExamples:");
    fmt.Println("gbench -file ~/benchmarkurl.txt -r 100 -c 10 -v");
    fmt.Println("gbench -url www.google.com -url www.yahoo.com -r 100 -c 10 -v");
}

func parseUrl(enteredUrl string) (details Url, parseErr error) {

    enteredUrl = strings.Trim(enteredUrl, "\" ");
    parts := strings.Split(enteredUrl, "|");

    if ((len(parts) != 2 && len(parts) != 3) || (strings.ToUpper(parts[0]) != "GET" && strings.ToUpper(parts[0]) != "POST" && strings.ToUpper(parts[0]) != "HEAD")) {
        parseErr = errors.New("Wrong URL format. Example: GET|www.google.com?search=test or POST|www.google.com|search=test or HEAD|www.google.com");
    } else {

        details.method = strings.ToUpper(parts[0]);

        urlStruct, err := url.Parse(parts[1]);

        if (err != nil) {
            parseErr = errors.New("Wrong URL format: " + err.Error());
        } else {
            if (urlStruct.Scheme == "") {
                urlStruct.Scheme = "http";
            }

            details.url = urlStruct.String();

            if (details.method == "POST") {
                if (len(parts) == 3) {
                    details.data = make(map[string]string);
                    dataParts := strings.Split(parts[2], "&");

                    for i := 0; i < len(dataParts); i++ {
                        keyValue := strings.Split(dataParts[i], "=");
                        details.data[keyValue[0]] = keyValue[1];
                    }
                } else {
                    parseErr = errors.New("Wrong URL format. You need to provide post data. Example: POST|www.google.com|name=sasan&lastname=rose");
                }
            }
        }
    }

    return details, parseErr;
}

func urlExists(url string) bool {
    for _, u := range urls {
        if (url == u.url) {
            return true;
        }
    }

    return false;
}
