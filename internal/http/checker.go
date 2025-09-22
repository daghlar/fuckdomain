package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Checker struct {
	timeout time.Duration
	client  *http.Client
}

type HTTPResponse struct {
	StatusCode int
	Headers    map[string][]string
	Body       string
	Title      string
	Server     string
	Length     int
}

func NewChecker(timeoutSeconds int) *Checker {
	timeout := time.Duration(timeoutSeconds) * time.Second
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Checker{
		timeout: timeout,
		client:  client,
	}
}

func (c *Checker) Check(domain string) (string, string) {
	urls := []string{
		fmt.Sprintf("http://%s", domain),
		fmt.Sprintf("https://%s", domain),
	}

	for _, url := range urls {
		response := c.makeRequest(url)
		if response != nil {
			status := fmt.Sprintf("%d", response.StatusCode)
			info := fmt.Sprintf("Status: %d, Server: %s, Title: %s, Length: %d",
				response.StatusCode, response.Server, response.Title, response.Length)
			return status, info
		}
	}

	return "N/A", "No HTTP response"
}

func (c *Checker) makeRequest(url string) *HTTPResponse {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	response := &HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Server:     resp.Header.Get("Server"),
		Length:     int(resp.ContentLength),
	}

	if resp.ContentLength > 0 && resp.ContentLength < 1024*1024 {
		buffer := make([]byte, resp.ContentLength)
		resp.Body.Read(buffer)
		response.Body = string(buffer)
		response.Title = c.extractTitle(response.Body)
	}

	return response
}

func (c *Checker) extractTitle(body string) string {
	start := strings.Index(strings.ToLower(body), "<title>")
	if start == -1 {
		return ""
	}
	start += 7

	end := strings.Index(strings.ToLower(body[start:]), "</title>")
	if end == -1 {
		return ""
	}

	title := body[start : start+end]
	title = strings.TrimSpace(title)
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\r", " ")
	title = strings.ReplaceAll(title, "\t", " ")

	for strings.Contains(title, "  ") {
		title = strings.ReplaceAll(title, "  ", " ")
	}

	if len(title) > 100 {
		title = title[:100] + "..."
	}

	return title
}

func (c *Checker) CheckMultiple(domains []string) map[string]*HTTPResponse {
	results := make(map[string]*HTTPResponse)
	
	for _, domain := range domains {
		urls := []string{
			fmt.Sprintf("http://%s", domain),
			fmt.Sprintf("https://%s", domain),
		}

		for _, url := range urls {
			response := c.makeRequest(url)
			if response != nil {
				results[domain] = response
				break
			}
		}
	}

	return results
}

func (c *Checker) CheckWithCustomHeaders(domain string, headers map[string]string) *HTTPResponse {
	urls := []string{
		fmt.Sprintf("http://%s", domain),
		fmt.Sprintf("https://%s", domain),
	}

	for _, url := range urls {
		ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
		
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		response := &HTTPResponse{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Server:     resp.Header.Get("Server"),
			Length:     int(resp.ContentLength),
		}

		if resp.ContentLength > 0 && resp.ContentLength < 1024*1024 {
			buffer := make([]byte, resp.ContentLength)
			resp.Body.Read(buffer)
			response.Body = string(buffer)
			response.Title = c.extractTitle(response.Body)
		}

		return response
	}

	return nil
}
