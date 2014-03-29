package main

import (
    "os"
    "fmt"
    "bufio"
    "strconv"
    "regexp"
    "strings"
    "time"
    "net/url"
)

func parsArgs() {
    args := os.Args;
    headers = make(map[string]string);

    if (len(args) == 1) {
        fmt.Println("Invalid arguments");
        os.Exit(1);
    }

    for i := 0; i<len(args); i++ {
        switch args[i] {
            case "--file", "-f":
                checkArgs(args, i);
                getFile(args[i+1]);
            case "--url", "-u":
                checkArgs(args, i);
                if (!urlExists(args[i+1])) {
                    urls = append(urls, parseUrl(args[i+1]));
                }
            case "--concurrent", "-c":
                checkArgs(args, i);
                concurrent,_ = strconv.Atoi(args[i+1]);
            case "--requests", "-r":
                checkArgs(args, i);
                requests,_ = strconv.Atoi(args[i+1]);
            case "--response-timeout", "-t":
                checkArgs(args, i);
                timeout, _ := strconv.Atoi(args[i+1]);
                responseTimeout = time.Duration(timeout) * time.Second;
            case "--connection-timeout", "-ct":
                checkArgs(args, i);
                timeout, _ := strconv.Atoi(args[i+1]);
                connectionTimeout = time.Duration(timeout) * time.Second;
            case "--verbose", "-v":
                verbose = true;
            case "--cookie", "-co":
                checkArgs(args, i);
                addCookie(args[i+1]);
            case "--username", "-user":
                checkArgs(args, i);
                basicAuthUsername = args[i+1];
            case "--password", "-pass":
                checkArgs(args, i);
                basicAuthPassword = args[i+1];
            case "--proxy", "-p":
                checkArgs(args, i);
                setProxy(args[i+1]);
            case "--custom-header", "-h":
                checkArgs(args, i);
                addHeader(args[i+1]);
            case "--disable-keep-alive", "-d":
                disableKeepAlive = true;
            case "--delay", "-de":
                checkArgs(args, i);
                delay, _ = strconv.Atoi(args[i+1]);
            case "--":
                break;
        }
    }

    if (requests == 0 || concurrent == 0 || requests < concurrent || len(urls) == 0) {
        fmt.Println("Invalid arguments");
        os.Exit(1);
    }
}

func addCookie(cookie string) {
    cookie = strings.Trim(cookie, "\"' ");

    cookies = append(cookies, cookie);
}

func setProxy(enteredUrl string) {
    enteredUrl = strings.Trim(enteredUrl, "\" ");

    var validArg = regexp.MustCompile("^localhost");

    if (enteredUrl == "localhost" || validArg.MatchString(enteredUrl)) {
        enteredUrl = "http://" + enteredUrl;
    }

    proxyUrl, err = url.Parse(enteredUrl);

    if (err != nil) {
        fmt.Printf("%v\n", enteredUrl);
        fmt.Println("Wrong proxy format. Example: http://localhost:8441");
        os.Exit(1);
    }

    if (proxyUrl.Scheme == "") {
        proxyUrl.Scheme = "http";
    }
}

func addHeader(header string) {
    parts := strings.Split(header, ":");

    if (len(parts) != 2) {
        fmt.Printf("%v\n", header);
        fmt.Println("Wrong header format. Example: 'accept-encoding: gzip, deflate'");
        os.Exit(1);
    }

    key := strings.Trim(parts[0], "\"' ");
    value := strings.Trim(parts[1], "\"' ");

    headers[key] = value;
}

func parseUrl(enteredUrl string) (details Url) {

    enteredUrl = strings.Trim(enteredUrl, "\" ");
    parts := strings.Split(enteredUrl, "|");

    if (len(parts) != 2 || (strings.ToUpper(parts[0]) != "GET" && strings.ToUpper(parts[0]) != "POST" && strings.ToUpper(parts[0]) != "HEAD")) {
        fmt.Printf("%v\n", enteredUrl);
        fmt.Println("Wrong URL format. Example: GET|www.google.com or POST|www.google.com?search=test or HEAD|www.google.com");
        os.Exit(1);
    }

    details.method = strings.ToUpper(parts[0]);

    urlParts := strings.Split(parts[1], "?");
    urlStruct, err := url.Parse(urlParts[0]);

    if (err != nil) {
        fmt.Printf("%v\n", enteredUrl);
        fmt.Println("Wrong URL format: %v", err);
        os.Exit(1);
    }

    if (urlStruct.Scheme == "") {
        urlStruct.Scheme = "http";
    }

    details.url = urlStruct.String();

    if (details.method == "POST") {
        if (len(urlParts) == 2) {
            details.data = make(map[string]string);
            dataParts := strings.Split(urlParts[1], "&");

            for i := 0; i < len(dataParts); i++ {
                keyValue := strings.Split(dataParts[i], "=");
                details.data[keyValue[0]] = keyValue[1];
            }
        } else {
            fmt.Printf("%v\n", enteredUrl);
            fmt.Println("Wrong URL format. You need to provide post data. Example: POST|www.google.com?search=test");
            os.Exit(1);
        }
    }

    return details;
}

func urlExists(url string) bool {
    for _, u := range urls {
        if (url == u.url) {
            return true;
        }
    }

    return false;
}

func checkArgs(args []string, index int) {
    var validArg = regexp.MustCompile("^--");

    if (len(args) <= index+1 || validArg.MatchString(args[index+1])) {
        fmt.Println("Invalid arguments");
        os.Exit(1);
    }
}

func getFile(file string) {
    urlFile, err = os.Open(file);

    if (err != nil) {
        fmt.Printf("Error opening file: %v\n", err);
        os.Exit(1);
    }

    defer urlFile.Close();

    scanner := bufio.NewScanner(urlFile)

    for scanner.Scan() {
        url := scanner.Text();
        if (!urlExists(url)) {
            urls = append(urls, parseUrl(url));
        }
    }
}
