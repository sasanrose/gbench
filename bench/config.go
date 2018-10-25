package bench

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/sasanrose/gbench/report"
)

// WithConcurrency creates a config to set number of concurrent
// requests per URL.
func WithConcurrency(n int) func(*Bench) {
	return func(b *Bench) {
		b.Concurrency = n
	}
}

// WithRequests creates a config to set total number of requests
// to send per URL.
func WithRequests(n int) func(*Bench) {
	return func(b *Bench) {
		b.Requests = n
	}
}

// WithURL adds an endpoint to benchmark.
func WithURL(u *URL) func(*Bench) {
	return func(b *Bench) {
		b.URLs = append(b.URLs, u)
	}
}

// WithSuccessStatusCode defines what should be considered as a success
// status code.
// Default values are: 200, 201, 202
func WithSuccessStatusCode(code int) func(*Bench) {
	return func(b *Bench) {
		b.SuccessStatusCodes = append(b.SuccessStatusCodes, code)
	}
}

// WithAuthUserPass adds basic HTTP authentication based on a string with
// username:password format.
func WithAuthUserPass(userPass string) (func(*Bench), error) {
	user, pass, err := parseUserPath(userPass)

	if err != nil {
		return nil, err
	}

	return WithAuth(user, pass), nil
}

// WithAuth adds basic HTTP authentication.
// Note: This will be used for all the provided urls.
func WithAuth(username, password string) func(*Bench) {
	return func(b *Bench) {
		b.Auth = &Auth{username, password}
	}
}

// WithOutput defines the output writer. All the messages during benchmarking
// will be written to this writer.
func WithOutput(w io.Writer) func(*Bench) {
	return func(b *Bench) {
		b.OutputWriter = w
	}
}

// WithProxy defines a proxy server address to use.
// Note: This will be used for all the provided urls.
func WithProxy(addr string) func(*Bench) {
	return func(b *Bench) {
		b.Proxy = addr
	}
}

// WithURLSettings sets a benchmarking endpoint using sepcific URL settings.
func WithURLSettings(requestedUrl,
	method string,
	data []string,
	headers []string,
	rawCookie string,
	userPass string,
) (func(*Bench), error) {
	parsedURL, err := url.Parse(requestedUrl)
	if err != nil {
		return nil, fmt.Errorf("Invalid URL provided: %v", err)
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, errors.New("Only http and https schemes are supported")
	}

	method = strings.ToUpper(method)

	if len(data) > 0 && method != http.MethodPost && method != http.MethodPut && method != http.MethodPatch {
		return nil, errors.New("Request data is only allowed with POST, PUT and PATCH request methods")
	}

	parsedData, err := parseData(data)

	if err != nil {
		return nil, err
	}

	endpoint := &URL{
		Addr:      parsedURL.String(),
		Method:    method,
		Data:      parsedData,
		RawCookie: rawCookie,
		Headers:   make(map[string]string),
	}

	if userPass != "" {
		user, pass, err := parseUserPath(userPass)

		if err != nil {
			return nil, err
		}

		endpoint.Auth = &Auth{user, pass}
	}

	if len(headers) > 0 {
		for _, header := range headers {
			key, value, err := parseHeaderString(header)

			if err != nil {
				return nil, err
			}

			endpoint.Headers[key] = value
		}
	}

	return WithURL(endpoint), nil
}

// WithConnectionTimeout sets connection timeout.
func WithConnectionTimeout(t time.Duration) func(*Bench) {
	return func(b *Bench) {
		b.ConnectionTimeout = t
	}
}

// WithResponseTimeout sets response timeout.
func WithResponseTimeout(t time.Duration) func(*Bench) {
	return func(b *Bench) {
		b.ResponseTimeout = t
	}
}

// WithRawCookie sets a raw cookie string.
// Note: This will be used for all the provided urls.
func WithRawCookie(cookie string) func(*Bench) {
	return func(b *Bench) {
		b.RawCookie = cookie
	}
}

// WithHeader sets http headers.
// Note: This will be used for all the provided urls.
func WithHeader(key, value string) func(*Bench) {
	return func(b *Bench) {
		b.Headers[key] = value
	}
}

// WithHeaderString sets http headers using a string in key=value format.
// Note: This will be used for all the provided urls.
func WithHeaderString(header string) (func(*Bench), error) {
	key, value, err := parseHeaderString(header)

	if err != nil {
		return nil, err
	}

	return WithHeader(key, value), nil
}

// WithFile sets endpoints using a file. The file is expected to have a
// string as defined in WithURLString in each line.
/*func WithFile(path string) (func(*Bench), error) {
	file, err := fs.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	s := bufio.NewScanner(file)
	s.Split(bufio.ScanLines)

	urls := make([]*URL, 0)

	for s.Scan() {
		parsedURL, err := parseURL(s.Text())

		if err != nil {
			return nil, err
		}

		urls = append(urls, parsedURL)
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("Did not find any url in the %s", path)
	}

	return func(b *Bench) {
		b.URLs = append(b.URLs, urls...)
	}, nil
}*/

// WithReport sets a result report
func WithReport(report report.Report) func(*Bench) {
	return func(b *Bench) {
		b.Report = report
	}
}

func parseData(formData []string) (map[string]string, error) {
	data := make(map[string]string)
	for _, formValue := range formData {
		dataParts := strings.Split(formValue, "&")
		for _, dataPart := range dataParts {
			keyValue := strings.Split(dataPart, "=")
			if len(keyValue) != 2 {
				return data, fmt.Errorf("Wrong key value format for request data: %s", dataPart)
			}

			data[keyValue[0]] = keyValue[1]
		}
	}

	return data, nil
}

func parseUserPath(userPass string) (user, pass string, err error) {
	m := regexp.MustCompile(`^([^:]+):(.+)$`)
	if !m.MatchString(userPass) {
		err = fmt.Errorf("Wrong auth credentials format: %s", userPass)
		return user, pass, err
	}

	matches := m.FindStringSubmatch(userPass)

	user, pass = matches[1], matches[2]

	return user, pass, err
}

func parseHeaderString(header string) (key, value string, err error) {
	keyValue := strings.Split(header, ":")

	if len(keyValue) != 2 && len(keyValue) != 1 {
		err = fmt.Errorf("%s is not a correct 'key: value;' format", header)
		return key, value, err
	}

	m := regexp.MustCompile(`^[a-zA-Z0-9-_]+;$`)

	if len(keyValue) == 1 && !m.MatchString(keyValue[0]) {
		err = fmt.Errorf("%s is not a correct 'key;' format", header)
		return key, value, err
	}

	key = strings.Trim(keyValue[0], " ")

	if len(keyValue) > 1 {
		value = strings.Trim(strings.Trim(keyValue[1], ";"), " ")
	}

	return key, value, err
}
