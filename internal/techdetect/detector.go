package techdetect

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Technology struct {
	Name        string
	Version     string
	Category    string
	Confidence  int
	Description string
	Website     string
}

type TechResult struct {
	URL         string
	Technologies []Technology
	Server      string
	Framework   string
	Database    string
	CDN         string
	Analytics   string
	Widgets     []string
	Languages   []string
	OS          string
}

type TechDetector struct {
	client  *http.Client
	timeout time.Duration
}

func NewTechDetector(timeout time.Duration) *TechDetector {
	return &TechDetector{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

func (td *TechDetector) Detect(url string) (*TechResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), td.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := td.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &TechResult{
		URL:          url,
		Technologies: make([]Technology, 0),
		Server:       resp.Header.Get("Server"),
		Widgets:      make([]string, 0),
		Languages:    make([]string, 0),
	}

	td.detectFromHeaders(resp.Header, result)
	td.detectFromBody(string(body), result)
	td.detectFromURL(url, result)

	return result, nil
}

func (td *TechDetector) detectFromHeaders(headers http.Header, result *TechResult) {
	server := headers.Get("Server")
	if server != "" {
		result.Technologies = append(result.Technologies, Technology{
			Name:        "Web Server",
			Version:     td.extractVersion(server),
			Category:    "Web Server",
			Confidence:  100,
			Description: "Web server technology",
		})
	}

	xPoweredBy := headers.Get("X-Powered-By")
	if xPoweredBy != "" {
		result.Technologies = append(result.Technologies, Technology{
			Name:        "X-Powered-By",
			Version:     td.extractVersion(xPoweredBy),
			Category:    "Framework",
			Confidence:  90,
			Description: "Application framework",
		})
	}

	xAspNetVersion := headers.Get("X-AspNet-Version")
	if xAspNetVersion != "" {
		result.Technologies = append(result.Technologies, Technology{
			Name:        "ASP.NET",
			Version:     xAspNetVersion,
			Category:    "Framework",
			Confidence:  100,
			Description: "Microsoft ASP.NET framework",
		})
	}
}

func (td *TechDetector) detectFromBody(body string, result *TechResult) {
	patterns := map[string]Technology{
		`<meta name="generator" content="([^"]+)"`: {
			Name:        "Generator",
			Category:    "CMS",
			Confidence:  95,
			Description: "Content management system",
		},
		`<script[^>]*src="[^"]*jquery[^"]*\.js[^"]*"`: {
			Name:        "jQuery",
			Category:    "JavaScript Library",
			Confidence:  90,
			Description: "JavaScript library",
		},
		`<script[^>]*src="[^"]*bootstrap[^"]*\.js[^"]*"`: {
			Name:        "Bootstrap",
			Category:    "CSS Framework",
			Confidence:  90,
			Description: "CSS framework",
		},
		`<script[^>]*src="[^"]*react[^"]*\.js[^"]*"`: {
			Name:        "React",
			Category:    "JavaScript Framework",
			Confidence:  90,
			Description: "JavaScript framework",
		},
		`<script[^>]*src="[^"]*angular[^"]*\.js[^"]*"`: {
			Name:        "Angular",
			Category:    "JavaScript Framework",
			Confidence:  90,
			Description: "JavaScript framework",
		},
		`<script[^>]*src="[^"]*vue[^"]*\.js[^"]*"`: {
			Name:        "Vue.js",
			Category:    "JavaScript Framework",
			Confidence:  90,
			Description: "JavaScript framework",
		},
		`<script[^>]*src="[^"]*wordpress[^"]*\.js[^"]*"`: {
			Name:        "WordPress",
			Category:    "CMS",
			Confidence:  95,
			Description: "Content management system",
		},
		`<script[^>]*src="[^"]*drupal[^"]*\.js[^"]*"`: {
			Name:        "Drupal",
			Category:    "CMS",
			Confidence:  95,
			Description: "Content management system",
		},
		`<script[^>]*src="[^"]*joomla[^"]*\.js[^"]*"`: {
			Name:        "Joomla",
			Category:    "CMS",
			Confidence:  95,
			Description: "Content management system",
		},
		`<script[^>]*src="[^"]*google-analytics[^"]*\.js[^"]*"`: {
			Name:        "Google Analytics",
			Category:    "Analytics",
			Confidence:  100,
			Description: "Web analytics service",
		},
		`<script[^>]*src="[^"]*gtag[^"]*\.js[^"]*"`: {
			Name:        "Google Tag Manager",
			Category:    "Analytics",
			Confidence:  100,
			Description: "Tag management system",
		},
		`<script[^>]*src="[^"]*facebook[^"]*\.js[^"]*"`: {
			Name:        "Facebook SDK",
			Category:    "Social Media",
			Confidence:  90,
			Description: "Facebook integration",
		},
		`<script[^>]*src="[^"]*twitter[^"]*\.js[^"]*"`: {
			Name:        "Twitter Widget",
			Category:    "Social Media",
			Confidence:  90,
			Description: "Twitter integration",
		},
		`<script[^>]*src="[^"]*cloudflare[^"]*\.js[^"]*"`: {
			Name:        "Cloudflare",
			Category:    "CDN",
			Confidence:  95,
			Description: "Content delivery network",
		},
		`<script[^>]*src="[^"]*amazonaws[^"]*\.js[^"]*"`: {
			Name:        "Amazon Web Services",
			Category:    "Cloud",
			Confidence:  90,
			Description: "Cloud computing platform",
		},
		`<script[^>]*src="[^"]*stripe[^"]*\.js[^"]*"`: {
			Name:        "Stripe",
			Category:    "Payment",
			Confidence:  95,
			Description: "Payment processing",
		},
		`<script[^>]*src="[^"]*paypal[^"]*\.js[^"]*"`: {
			Name:        "PayPal",
			Category:    "Payment",
			Confidence:  95,
			Description: "Payment processing",
		},
	}

	for pattern, tech := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(body)
		if len(matches) > 0 {
			tech.Version = td.extractVersion(matches[0])
			result.Technologies = append(result.Technologies, tech)
		}
	}

	td.detectLanguages(body, result)
	td.detectDatabases(body, result)
	td.detectFrameworks(body, result)
}

func (td *TechDetector) detectFromURL(url string, result *TechResult) {
	if strings.Contains(url, "wordpress") {
		result.Technologies = append(result.Technologies, Technology{
			Name:        "WordPress",
			Category:    "CMS",
			Confidence:  80,
			Description: "Content management system",
		})
	}
}

func (td *TechDetector) detectLanguages(body string, result *TechResult) {
	if strings.Contains(body, "<?php") {
		result.Languages = append(result.Languages, "PHP")
		result.Technologies = append(result.Technologies, Technology{
			Name:        "PHP",
			Category:    "Programming Language",
			Confidence:  95,
			Description: "Server-side programming language",
		})
	}
	if strings.Contains(body, "asp.net") || strings.Contains(body, "ASP.NET") {
		result.Languages = append(result.Languages, "ASP.NET")
	}
	if strings.Contains(body, "jsp") || strings.Contains(body, "JSP") {
		result.Languages = append(result.Languages, "JSP")
	}
	if strings.Contains(body, "python") || strings.Contains(body, "django") {
		result.Languages = append(result.Languages, "Python")
	}
	if strings.Contains(body, "ruby") || strings.Contains(body, "rails") {
		result.Languages = append(result.Languages, "Ruby")
	}
}

func (td *TechDetector) detectDatabases(body string, result *TechResult) {
	if strings.Contains(body, "mysql") || strings.Contains(body, "MySQL") {
		result.Database = "MySQL"
		result.Technologies = append(result.Technologies, Technology{
			Name:        "MySQL",
			Category:    "Database",
			Confidence:  80,
			Description: "Relational database management system",
		})
	}
	if strings.Contains(body, "postgresql") || strings.Contains(body, "PostgreSQL") {
		result.Database = "PostgreSQL"
	}
	if strings.Contains(body, "mongodb") || strings.Contains(body, "MongoDB") {
		result.Database = "MongoDB"
	}
}

func (td *TechDetector) detectFrameworks(body string, result *TechResult) {
	if strings.Contains(body, "laravel") || strings.Contains(body, "Laravel") {
		result.Framework = "Laravel"
		result.Technologies = append(result.Technologies, Technology{
			Name:        "Laravel",
			Category:    "Framework",
			Confidence:  90,
			Description: "PHP web framework",
		})
	}
	if strings.Contains(body, "symfony") || strings.Contains(body, "Symfony") {
		result.Framework = "Symfony"
	}
	if strings.Contains(body, "codeigniter") || strings.Contains(body, "CodeIgniter") {
		result.Framework = "CodeIgniter"
	}
}

func (td *TechDetector) extractVersion(text string) string {
	re := regexp.MustCompile(`(\d+\.\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func (td *TechDetector) DetectMultiple(urls []string) map[string]*TechResult {
	results := make(map[string]*TechResult)
	
	for _, url := range urls {
		if result, err := td.Detect(url); err == nil {
			results[url] = result
		}
	}
	
	return results
}
