package bruteforce

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type BruteforceConfig struct {
	Threads    int
	Timeout    time.Duration
	UserAgent  string
	Headers    map[string]string
	Extensions []string
	StatusCodes []int
}

type BruteforceResult struct {
	URL          string
	StatusCode   int
	ContentLength int64
	Title        string
	Server       string
	ResponseTime time.Duration
	Found        bool
}

type DirectoryBruteforcer struct {
	config BruteforceConfig
	client *http.Client
}

func NewDirectoryBruteforcer(config BruteforceConfig) *DirectoryBruteforcer {
	return &DirectoryBruteforcer{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (db *DirectoryBruteforcer) Bruteforce(baseURL string, wordlist []string) map[string]*BruteforceResult {
	results := make(map[string]*BruteforceResult)
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, db.config.Threads)

	for _, word := range wordlist {
		wg.Add(1)
		go func(w string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			urls := db.generateURLs(baseURL, w)
			for _, url := range urls {
				result := db.checkURL(url)
				if result.Found {
					mu.Lock()
					results[url] = result
					mu.Unlock()
				}
			}
		}(word)
	}

	wg.Wait()
	return results
}

func (db *DirectoryBruteforcer) generateURLs(baseURL, word string) []string {
	var urls []string
	
	// Clean base URL
	baseURL = strings.TrimSuffix(baseURL, "/")
	
	// Directory bruteforce
	urls = append(urls, fmt.Sprintf("%s/%s/", baseURL, word))
	urls = append(urls, fmt.Sprintf("%s/%s", baseURL, word))
	
	// File bruteforce with extensions
	for _, ext := range db.config.Extensions {
		urls = append(urls, fmt.Sprintf("%s/%s%s", baseURL, word, ext))
	}
	
	return urls
}

func (db *DirectoryBruteforcer) checkURL(url string) *BruteforceResult {
	start := time.Now()
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &BruteforceResult{URL: url, Found: false}
	}

	req.Header.Set("User-Agent", db.config.UserAgent)
	for key, value := range db.config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := db.client.Do(req)
	if err != nil {
		return &BruteforceResult{URL: url, Found: false}
	}
	defer resp.Body.Close()

	responseTime := time.Since(start)
	
	// Check if status code is in allowed list
	allowed := false
	for _, code := range db.config.StatusCodes {
		if resp.StatusCode == code {
			allowed = true
			break
		}
	}
	
	if !allowed {
		return &BruteforceResult{URL: url, Found: false}
	}

	// Read limited content for title extraction
	body := make([]byte, 1024)
	n, _ := io.ReadFull(resp.Body, body)
	content := string(body[:n])
	
	title := db.extractTitle(content)
	server := resp.Header.Get("Server")

	return &BruteforceResult{
		URL:          url,
		StatusCode:   resp.StatusCode,
		ContentLength: resp.ContentLength,
		Title:        title,
		Server:       server,
		ResponseTime: responseTime,
		Found:        true,
	}
}

func (db *DirectoryBruteforcer) extractTitle(content string) string {
	start := strings.Index(strings.ToLower(content), "<title>")
	if start == -1 {
		return ""
	}
	start += 7

	end := strings.Index(strings.ToLower(content[start:]), "</title>")
	if end == -1 {
		return ""
	}

	title := content[start : start+end]
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

func (db *DirectoryBruteforcer) BruteforceWithContext(ctx context.Context, baseURL string, wordlist []string) map[string]*BruteforceResult {
	results := make(map[string]*BruteforceResult)
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, db.config.Threads)

	for _, word := range wordlist {
		select {
		case <-ctx.Done():
			return results
		default:
		}

		wg.Add(1)
		go func(w string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			urls := db.generateURLs(baseURL, w)
			for _, url := range urls {
				select {
				case <-ctx.Done():
					return
				default:
				}

				result := db.checkURL(url)
				if result.Found {
					mu.Lock()
					results[url] = result
					mu.Unlock()
				}
			}
		}(word)
	}

	wg.Wait()
	return results
}

func (db *DirectoryBruteforcer) BruteforceCommon(baseURL string) map[string]*BruteforceResult {
	commonPaths := []string{
		"admin", "administrator", "login", "wp-admin", "wp-login", "dashboard",
		"panel", "control", "manage", "manager", "admin.php", "login.php",
		"index.php", "home.php", "about.php", "contact.php", "services.php",
		"products.php", "blog.php", "news.php", "support.php", "help.php",
		"api", "v1", "v2", "api/v1", "api/v2", "rest", "graphql",
		"config", "configuration", "settings", "setup", "install",
		"backup", "backups", "files", "uploads", "images", "css", "js",
		"assets", "static", "public", "private", "secure", "test", "dev",
		"staging", "beta", "alpha", "demo", "sandbox", "playground",
		"docs", "documentation", "wiki", "help", "faq", "support",
		"status", "health", "ping", "monitor", "metrics", "stats",
		"logs", "log", "debug", "trace", "error", "errors",
		"robots.txt", "sitemap.xml", "crossdomain.xml", "security.txt",
		".env", ".git", ".svn", ".hg", ".bzr", ".cvs",
		"phpinfo.php", "info.php", "test.php", "debug.php",
		"readme.txt", "readme.md", "changelog.txt", "license.txt",
		"version.txt", "version.json", "package.json", "composer.json",
		"yarn.lock", "package-lock.json", "requirements.txt",
		"docker-compose.yml", "dockerfile", "Dockerfile",
		"k8s.yaml", "kubernetes.yaml", "helm.yaml",
		"terraform.tf", "ansible.yml", "puppet.pp",
		"vagrantfile", "Vagrantfile", "Makefile",
		"gruntfile.js", "gulpfile.js", "webpack.config.js",
		"tsconfig.json", "babel.config.js", "eslint.config.js",
		"prettier.config.js", "jest.config.js", "karma.conf.js",
		"protractor.conf.js", "cypress.json", "playwright.config.js",
		"vitest.config.js", "vite.config.js", "rollup.config.js",
		"parcel.config.js", "snowpack.config.js", "esbuild.config.js",
		"swc.config.js", "turbo.json", "nx.json", "lerna.json",
		"rush.json", "pnpm-workspace.yaml", "yarn.lock",
		"package-lock.json", "npm-shrinkwrap.json", "yarn-error.log",
		"npm-debug.log", "lerna-debug.log", "rush-debug.log",
		"pnpm-debug.log", "yarn-debug.log", "npm-debug.log",
		"lerna-debug.log", "rush-debug.log", "pnpm-debug.log",
	}

	return db.Bruteforce(baseURL, commonPaths)
}

func (db *DirectoryBruteforcer) BruteforceWithExtensions(baseURL string, wordlist []string, extensions []string) map[string]*BruteforceResult {
	originalExtensions := db.config.Extensions
	db.config.Extensions = extensions
	defer func() { db.config.Extensions = originalExtensions }()

	return db.Bruteforce(baseURL, wordlist)
}

func (db *DirectoryBruteforcer) BruteforceWithStatusCodes(baseURL string, wordlist []string, statusCodes []int) map[string]*BruteforceResult {
	originalStatusCodes := db.config.StatusCodes
	db.config.StatusCodes = statusCodes
	defer func() { db.config.StatusCodes = originalStatusCodes }()

	return db.Bruteforce(baseURL, wordlist)
}
