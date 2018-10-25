package bench

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Exec executes a benchmark. The context is used to cancel the benchmark at any
// given time.
func (b *Bench) Exec(ctx context.Context) error {
	client := b.getClient()
	remainingRequests := b.Requests
	t := time.Now()

	b.Report.SetStartTime(t)

	defer func() {
		te := time.Now()
		b.Report.SetTotalDuration(te.Sub(t))
		b.Report.SetEndTime(te)
	}()

	for remainingRequests > 0 {
		waitChannel := make(chan struct{})
		doneReqs := b.Requests - remainingRequests
		b.printOutputMessage(fmt.Sprintf("%d of %d (%.1f%%)\n", doneReqs, b.Requests, float64(doneReqs*b.Requests)/100))
		go b.runConcurrentJobs(ctx, waitChannel, client, &remainingRequests)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-waitChannel:
			continue
		}
	}

	return nil
}

func (b *Bench) runConcurrentJobs(ctx context.Context, waitChannel chan struct{}, client *http.Client, remainingRequests *int) {
	wg := &sync.WaitGroup{}
	remainingConcurrent := b.Concurrency
	for remainingConcurrent > 0 && *remainingRequests > 0 {
		for _, url := range b.URLs {
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
	reqURL := req.URL.String()

	if err != nil {
		if err, ok := err.(*url.Error); ok && err.Timeout() {
			b.printOutputMessage(fmt.Sprintf("Timed out request for %s: %v\n", reqURL, err))
			b.Report.AddTimedoutResponse(reqURL)
			return
		}

		b.printOutputMessage(fmt.Sprintf("Error for %s: %v\n", reqURL, err))
		b.Report.AddFailedResponse(reqURL)
		return
	}

	defer resp.Body.Close()

	contentLength := 0
	body, err := ioutil.ReadAll(resp.Body)

	if err == nil {
		contentLength = len(body)
	}

	b.Report.AddResponseTime(reqURL, responseTime)
	b.Report.AddReceivedDataLength(reqURL, int64(contentLength))
	b.Report.AddResponseStatusCode(reqURL, resp.StatusCode, b.isFailed(resp.StatusCode))
	b.printOutputMessage(fmt.Sprintf("Received response for sent requests to %s in %v. Status: %s\n", reqURL, responseTime, http.StatusText(resp.StatusCode)))
}

func (b *Bench) printOutputMessage(msg string) {
	if b.OutputWriter != nil {
		b.OutputWriterLock.Lock()
		defer b.OutputWriterLock.Unlock()
		fmt.Fprint(b.OutputWriter, msg)
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

func (b *Bench) buildRequest(u *URL) *http.Request {
	req, err := b.newRequest(u)

	if err != nil {
		log.Fatalf("Could not create request for %s: %v", u, err)
	}

	auth := b.getAuth(u)

	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}

	if _, ok := b.Headers["User-Agent"]; !ok {
		req.Header.Add("User-Agent", "Gbench")
	}

	headers := b.getHeaders(u)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	rawCookie := b.getRawCookie(u)

	if rawCookie != "" {
		req.Header.Add("Set-Cookie", rawCookie)
	}

	return req
}

func (b *Bench) newRequest(u *URL) (*http.Request, error) {
	if len(u.Data) == 0 ||
		(u.Method != http.MethodPost && u.Method != http.MethodPatch && u.Method != http.MethodPut) {
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

func (b *Bench) getAuth(u *URL) *Auth {
	if u.Auth != nil {
		return u.Auth
	}

	return b.Auth
}

func (b *Bench) getHeaders(u *URL) map[string]string {
	headers := make(map[string]string)
	for key, value := range b.Headers {
		headers[key] = value
	}

	for key, value := range u.Headers {
		headers[key] = value
	}

	return headers
}

func (b *Bench) getRawCookie(u *URL) string {
	if u.RawCookie != "" {
		return u.RawCookie
	}

	return b.RawCookie
}
