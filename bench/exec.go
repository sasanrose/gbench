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
	client := b.getClient()
	remainingRequests := b.Requests
	t := time.Now()

	for remainingRequests > 0 {
		waitChannel := make(chan struct{})
		go b.runConcurrentJobs(ctx, waitChannel, client, &remainingRequests)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-waitChannel:
			continue
		}
	}

	defer func() {
		b.Renderer.SetTotalDuration(time.Since(t))
	}()

	return nil
}

func (b *Bench) runConcurrentJobs(ctx context.Context, waitChannel chan struct{}, client *http.Client, remainingRequests *int) {
	wg := &sync.WaitGroup{}
	remainingConcurrent := b.Concurrency
	for remainingConcurrent > 0 && *remainingRequests > 0 {
		for _, url := range b.Urls {
			req := b.buildRequest(url)
			req = req.WithContext(ctx)
			wg.Add(1)
			go b.runBench(wg, client, req)
		}
		(*remainingRequests)--
		remainingConcurrent--
	}
	wg.Wait()
	close(waitChannel)
}

func (b *Bench) runBench(wg *sync.WaitGroup, client *http.Client, req *http.Request) {
	defer wg.Done()

	tr := time.Now()
	resp, err := client.Do(req)
	responseTime := time.Since(tr)
	reqUrl := req.URL.String()

	if err != nil {
		if err, ok := err.(*url.Error); ok && err.Timeout() {
			go b.Renderer.AddTimedoutResponse(reqUrl)
			return
		}

		go b.printVerbosityMessage(fmt.Sprintf("Error for %s: %v\n", reqUrl, err))
		go b.Renderer.AddFailedResponse(reqUrl)
		return
	}

	defer resp.Body.Close()

	go b.updateResult(reqUrl, resp, responseTime)
}

func (b *Bench) updateResult(reqUrl string, resp *http.Response, responseTime time.Duration) {
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
