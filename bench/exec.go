package bench

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Executes a benchmark. The context is used to cancel the benchmark at any
// given time.
func (b *Bench) Exec(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	waitChannel := make(chan struct{})

	client := b.getClient()
	remainingRequests := b.Requests
	t := time.Now()

	for remainingRequests > 0 {
		remainingConcurrent := b.Concurrency
		for remainingConcurrent > 0 && remainingRequests > 0 {
			for _, url := range b.Urls {
				req := b.buildRequest(url)
				req = req.WithContext(ctx)
				wg.Add(1)
				go b.runBench(wg, client, req)
			}
			remainingRequests--
			remainingConcurrent--
		}
	}

	go func() {
		defer close(waitChannel)
		wg.Wait()
	}()

	defer func() {
		b.Renderer.SetTotalDuration(time.Since(t))
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-waitChannel:
		return nil
	}

	return nil
}

func (b *Bench) runBench(wg *sync.WaitGroup, client *http.Client, req *http.Request) {
	defer wg.Done()

	tr := time.Now()
	resp, err := client.Do(req)
	responseTime := time.Since(tr)
	reqUrl := req.URL.String()

	if err != nil {
		if err, ok := err.(*url.Error); ok && err.Timeout() {
			b.Renderer.AddTimedoutResponse(reqUrl)
			return
		}

		b.printVerbosityMessage(fmt.Sprintf("Error for %s: %v\n", reqUrl, err))
		b.Renderer.AddFailedResponse(reqUrl)
		return
	}

	defer resp.Body.Close()

	b.Renderer.AddResponseTime(reqUrl, responseTime)
	b.Renderer.AddReceivedDataLength(reqUrl, resp.ContentLength)
	b.Renderer.AddResponseStatusCode(reqUrl, resp.StatusCode, b.isFailed(resp.StatusCode))
	b.printVerbosityMessage(fmt.Sprintf("Sent requests for %s in %v: %s\n", reqUrl, responseTime, http.StatusText(resp.StatusCode)))
}

func (b *Bench) printVerbosityMessage(msg string) {
	if b.VerbosityWriter != nil {
		b.VerbosityWriterLock.Lock()
		fmt.Fprint(b.VerbosityWriter, msg)
		b.VerbosityWriterLock.Unlock()
	}
}

func (b *Bench) isFailed(statusCode int) bool {
	for _, code := range b.SuccessStatusCodes {
		if statusCode == code {
			return false
		}
	}

	return true
}

func (b *Bench) getClient() *http.Client {
	tr := &http.Transport{
		Dial: b.getDial(),
	}

	if b.ResponseTimeout > 0 {
		tr.ResponseHeaderTimeout = b.ResponseTimeout
	}

	if b.Proxy != "" {
		p, _ := url.Parse(b.Proxy)
		tr.Proxy = http.ProxyURL(p)
	}

	return &http.Client{Transport: tr}
}

func (b *Bench) getDial() func(network, addr string) (net.Conn, error) {
	if b.ConnectionTimeout > 0 {
		return func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, b.ConnectionTimeout)
		}
	}

	return net.Dial
}

func (b *Bench) buildRequest(u *Url) *http.Request {
	req, err := newRequest(u)

	if err != nil {
		log.Fatalf("Could not create request for %s: %v", u, err)
	}

	if b.Auth != nil {
		req.SetBasicAuth(b.Auth.Username, b.Auth.Password)
	}

	if _, ok := b.Headers["User-Agent"]; !ok {
		req.Header.Add("User-Agent", "Gbench")
	}

	for key, value := range b.Headers {
		req.Header.Add(key, value)
	}

	if b.RawCookie != "" {
		req.Header.Add("Set-Cookie", b.RawCookie)
	}

	return req
}

func newRequest(u *Url) (*http.Request, error) {
	if u.Method == http.MethodGet || u.Method == http.MethodHead {
		return http.NewRequest(u.Method, u.Addr, nil)
	}

	values := url.Values{}

	for key, value := range u.Data {
		values.Add(key, value)
	}

	req, err := http.NewRequest(u.Method, u.Addr, bytes.NewBufferString(values.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}
