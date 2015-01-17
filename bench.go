package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, connectionTimeout)
}

func request(inputUrl chan Url, wait chan bool) {
	var response *http.Response
	var req *http.Request
	var err error

	go func() {
		request := <-inputUrl
		tr := &http.Transport{
			ResponseHeaderTimeout: responseTimeout,
			Dial:              dialTimeout,
			DisableKeepAlives: disableKeepAlive,
		}

		if proxyUrl != "" {
			p, _ := url.Parse(string(proxyUrl))
			tr.Proxy = http.ProxyURL(p)
		}

		client := &http.Client{Transport: tr}
		t1 := time.Now()

		if request.method == "GET" || request.method == "HEAD" {
			req, err = http.NewRequest(request.method, request.url, nil)
		} else {
			values := url.Values{}

			for key, value := range request.data {
				values.Add(key, value)
			}

			req, err = http.NewRequest("POST", request.url, bytes.NewBufferString(values.Encode()))
		}

		if basicAuthPassword != "" && basicAuthUsername != "" {
			req.SetBasicAuth(basicAuthUsername, basicAuthPassword)
		}

		if err == nil {
			req.Header.Add("User-Agent", "GBench")

			for key, value := range headers {
				req.Header.Add(key, value)
			}

			for _, cookie := range cookies {
				req.Header.Add("Set-Cookie", cookie)
			}

			response, err = client.Do(req)
			totalTransactions++
		}

		t2 := time.Now()

		responseTime := t2.Sub(t1)

		if err != nil {
			failedTransactions++

			if verbose {
				fmt.Printf("Error in calling url: %v\n", err)
			}
		} else {
			if longestResponseTime <= responseTime {
				longestResponseTime = responseTime
			}

			if shortestResponseTime >= responseTime || shortestResponseTime == 0 {
				shortestResponseTime = responseTime
			}

			if urlsResponseTimes[request.url] <= responseTime {
				urlsResponseTimes[request.url] = responseTime
			}

			if response.StatusCode != 200 {
				failedTransactions++
				urlFailedStats[request.url]++
			}

			responseStats[response.Status]++

			totalLength += response.ContentLength

			if verbose {
				fmt.Println(request.url)
				fmt.Printf("Response Status: %v ", response.Status)
				fmt.Printf("Response Time: %f secs\n", responseTime.Seconds())
			}

			defer response.Body.Close()
		}

		wait <- true
	}()
}

func startBench() {
	fmt.Println("Start benchmarking...")

	mod := int(math.Mod(float64(requests), float64(concurrent)))
	floor := int(math.Floor(float64(requests / concurrent)))

	var count int = 1
	var wait = make(chan bool)
	var urlInput = make(chan Url)
	responseStats = make(map[string]int)
	urlFailedStats = make(map[string]int)
	urlsResponseTimes = make(map[string]time.Duration)

	startTime = time.Now()
	for i := 0; i <= floor; i++ {

		var max int

		if (i*concurrent)+mod >= requests {
			max = mod
		} else {
			max = concurrent
		}

		for j := 0; j < max; j++ {
			count++
			request(urlInput, wait)
		}

		for j := 0; j < max; j++ {
			urlInput <- urls[rand.Intn(len(urls))]
		}

		if delay > 0 {
			time.Sleep(time.Duration(rand.Intn(delay)) * time.Second)
		}
	}

	for i := 0; i < requests; i++ {
		<-wait
	}
}
