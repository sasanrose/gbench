package cmd

import (
	"net/http"
	"time"
)

var (
	defaultConcurreny int = 1
	defaultRequests   int = 1

	defaultStatusCodes []int  = []int{http.StatusOK, http.StatusAccepted, http.StatusCreated}
	defaultMethod      string = http.MethodGet

	headers                            []string
	authUserPass, proxyURL, rawCookie  string
	concurrency, requests              int
	successStatusCodes                 []int
	connectionTimeout, responseTimeout time.Duration
)

type JsonConfig struct {
	Host            string        `json:"host"`
	Concurrency     int           `json:"concurrency"`
	Requests        int           `json:"requests"`
	StatusCodes     []int         `json:"status-codes"`
	AuthUserPass    string        `json:"user"`
	Proxy           string        `json:"proxy"`
	ConnectTimeout  time.Duration `json:"connect-timeout"`
	ResponseTimeout time.Duration `json:"response-timeout"`
	Headers         []string      `json:"headers"`
	RawCookie       string        `json:"cookie"`
	Paths           []*PathConfig `json:"paths"`
}

type PathConfig struct {
	Path         string   `json:"path"`
	Method       string   `json:"method"`
	Headers      []string `json:"headers"`
	Data         []string `json:"data"`
	RawCookie    string   `json:"cookie"`
	AuthUserPass string   `json:"user"`
}
