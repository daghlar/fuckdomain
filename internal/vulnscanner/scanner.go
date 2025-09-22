package vulnscanner

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type VulnScanner struct {
	client  *http.Client
	timeout time.Duration
}

type VulnCheck struct {
	Name        string
	Description string
	Severity    string
	CheckFunc   func(string, *http.Response) *Vulnerability
}

type Vulnerability struct {
	Name        string   `json:"name"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
	CVSS        string   `json:"cvss"`
	CVE         string   `json:"cve"`
	Solution    string   `json:"solution"`
	References  []string `json:"references"`
	Evidence    string   `json:"evidence"`
	Confidence  int      `json:"confidence"`
}

func NewVulnScanner(timeout time.Duration) *VulnScanner {
	return &VulnScanner{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

func (vs *VulnScanner) ScanURL(url string) ([]Vulnerability, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := vs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var vulnerabilities []Vulnerability

	// HTTP Security Headers Check
	vulns := vs.checkSecurityHeaders(resp)
	vulnerabilities = append(vulnerabilities, vulns...)

	// Server Information Disclosure
	vulns = vs.checkServerInfo(resp)
	vulnerabilities = append(vulnerabilities, vulns...)

	// Directory Traversal
	vulns = vs.checkDirectoryTraversal(url, resp)
	vulnerabilities = append(vulnerabilities, vulns...)

	// SQL Injection
	vulns = vs.checkSQLInjection(url, resp)
	vulnerabilities = append(vulnerabilities, vulns...)

	// XSS
	vulns = vs.checkXSS(url, resp)
	vulnerabilities = append(vulnerabilities, vulns...)

	// Information Disclosure
	vulns = vs.checkInformationDisclosure(string(body), resp)
	vulnerabilities = append(vulnerabilities, vulns...)

	// SSL/TLS Issues
	vulns = vs.checkSSLIssues(url, resp)
	vulnerabilities = append(vulnerabilities, vulns...)

	return vulnerabilities, nil
}

func (vs *VulnScanner) checkSecurityHeaders(resp *http.Response) []Vulnerability {
	var vulns []Vulnerability

	// Missing Security Headers
	securityHeaders := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":       "1; mode=block",
		"Strict-Transport-Security": "max-age=31536000",
		"Content-Security-Policy": "default-src 'self'",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
	}

	for header, expected := range securityHeaders {
		if value := resp.Header.Get(header); value == "" {
			vulns = append(vulns, Vulnerability{
				Name:        "Missing Security Header: " + header,
				Severity:    "Medium",
				Description: fmt.Sprintf("Missing security header: %s", header),
				Solution:    fmt.Sprintf("Add %s header with value: %s", header, expected),
				Confidence:  90,
			})
		}
	}

	// Weak HSTS
	if hsts := resp.Header.Get("Strict-Transport-Security"); hsts != "" {
		if !strings.Contains(hsts, "includeSubDomains") {
			vulns = append(vulns, Vulnerability{
				Name:        "Weak HSTS Configuration",
				Severity:    "Low",
				Description: "HSTS header missing includeSubDomains directive",
				Solution:    "Add includeSubDomains directive to HSTS header",
				Confidence:  80,
			})
		}
	}

	return vulns
}

func (vs *VulnScanner) checkServerInfo(resp *http.Response) []Vulnerability {
	var vulns []Vulnerability

	server := resp.Header.Get("Server")
	if server != "" {
		// Server version disclosure
		if strings.Contains(server, "/") {
			vulns = append(vulns, Vulnerability{
				Name:        "Server Version Disclosure",
				Severity:    "Low",
				Description: fmt.Sprintf("Server version disclosed: %s", server),
				Solution:    "Remove or obfuscate server version information",
				Confidence:  95,
			})
		}

		// Outdated server versions
		if strings.Contains(server, "Apache/2.2") || strings.Contains(server, "Apache/2.0") {
			vulns = append(vulns, Vulnerability{
				Name:        "Outdated Apache Version",
				Severity:    "High",
				Description: fmt.Sprintf("Outdated Apache version: %s", server),
				Solution:    "Update Apache to latest version",
				Confidence:  90,
			})
		}
	}

	// X-Powered-By disclosure
	if poweredBy := resp.Header.Get("X-Powered-By"); poweredBy != "" {
		vulns = append(vulns, Vulnerability{
			Name:        "Technology Disclosure",
			Severity:    "Low",
			Description: fmt.Sprintf("Technology disclosed: %s", poweredBy),
			Solution:    "Remove X-Powered-By header",
			Confidence:  95,
		})
	}

	return vulns
}

func (vs *VulnScanner) checkDirectoryTraversal(url string, resp *http.Response) []Vulnerability {
	var vulns []Vulnerability

	// Check for directory traversal patterns
	patterns := []string{
		"../",
		"..\\",
		"....//",
		"....\\\\",
		"%2e%2e%2f",
		"%2e%2e%5c",
	}

	for _, pattern := range patterns {
		testURL := url + "/" + pattern + "etc/passwd"
		req, _ := http.NewRequest("GET", testURL, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
		
		testResp, err := vs.client.Do(req)
		if err == nil {
			defer testResp.Body.Close()
			body, _ := io.ReadAll(testResp.Body)
			
			if strings.Contains(string(body), "root:") || strings.Contains(string(body), "bin:") {
				vulns = append(vulns, Vulnerability{
					Name:        "Directory Traversal",
					Severity:    "High",
					Description: "Directory traversal vulnerability detected",
					Solution:    "Implement proper input validation and path sanitization",
					Confidence:  85,
				})
				break
			}
		}
	}

	return vulns
}

func (vs *VulnScanner) checkSQLInjection(url string, resp *http.Response) []Vulnerability {
	var vulns []Vulnerability

	// SQL injection test patterns
	patterns := []string{
		"' OR '1'='1",
		"' UNION SELECT NULL--",
		"'; DROP TABLE users--",
		"' OR 1=1--",
		"admin'--",
		"admin'/*",
	}

	for _, pattern := range patterns {
		testURL := url + "?id=" + pattern
		req, _ := http.NewRequest("GET", testURL, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
		
		testResp, err := vs.client.Do(req)
		if err == nil {
			defer testResp.Body.Close()
			body, _ := io.ReadAll(testResp.Body)
			
			errorPatterns := []string{
				"mysql_fetch_array",
				"mysql_num_rows",
				"ORA-01756",
				"Microsoft OLE DB Provider",
				"ODBC SQL Server Driver",
				"SQLServer JDBC Driver",
				"PostgreSQL query failed",
				"Warning: mysql_",
				"valid MySQL result",
				"MySqlClient.",
			}

			for _, errorPattern := range errorPatterns {
				if strings.Contains(strings.ToLower(string(body)), strings.ToLower(errorPattern)) {
					vulns = append(vulns, Vulnerability{
						Name:        "SQL Injection",
						Severity:    "Critical",
						Description: "SQL injection vulnerability detected",
						Solution:    "Use parameterized queries and input validation",
						Confidence:  80,
					})
					break
				}
			}
		}
	}

	return vulns
}

func (vs *VulnScanner) checkXSS(url string, resp *http.Response) []Vulnerability {
	var vulns []Vulnerability

	// XSS test patterns
	patterns := []string{
		"<script>alert('XSS')</script>",
		"<img src=x onerror=alert('XSS')>",
		"javascript:alert('XSS')",
		"<svg onload=alert('XSS')>",
		"<iframe src=javascript:alert('XSS')>",
	}

	for _, pattern := range patterns {
		testURL := url + "?q=" + pattern
		req, _ := http.NewRequest("GET", testURL, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
		
		testResp, err := vs.client.Do(req)
		if err == nil {
			defer testResp.Body.Close()
			body, _ := io.ReadAll(testResp.Body)
			
			if strings.Contains(string(body), pattern) {
				vulns = append(vulns, Vulnerability{
					Name:        "Cross-Site Scripting (XSS)",
					Severity:    "High",
					Description: "XSS vulnerability detected",
					Solution:    "Implement proper output encoding and input validation",
					Confidence:  75,
				})
				break
			}
		}
	}

	return vulns
}

func (vs *VulnScanner) checkInformationDisclosure(body string, resp *http.Response) []Vulnerability {
	var vulns []Vulnerability

	// Check for sensitive information in response
	sensitivePatterns := map[string]string{
		"password":     "Password found in response",
		"api_key":      "API key found in response",
		"secret":       "Secret found in response",
		"token":        "Token found in response",
		"database":     "Database information found",
		"config":       "Configuration information found",
		"error":        "Error information disclosed",
		"stack trace":  "Stack trace disclosed",
		"exception":    "Exception information disclosed",
	}

	bodyLower := strings.ToLower(body)
	for pattern, description := range sensitivePatterns {
		if strings.Contains(bodyLower, pattern) {
			vulns = append(vulns, Vulnerability{
				Name:        "Information Disclosure",
				Severity:    "Medium",
				Description: description,
				Solution:    "Remove sensitive information from responses",
				Confidence:  70,
			})
		}
	}

	// Check for debug information
	if strings.Contains(bodyLower, "debug") || strings.Contains(bodyLower, "development") {
		vulns = append(vulns, Vulnerability{
			Name:        "Debug Information Disclosure",
			Severity:    "Low",
			Description: "Debug information found in response",
			Solution:    "Disable debug mode in production",
			Confidence:  80,
		})
	}

	return vulns
}

func (vs *VulnScanner) checkSSLIssues(url string, resp *http.Response) []Vulnerability {
	var vulns []Vulnerability

	// Check if HTTPS is used
	if !strings.HasPrefix(url, "https://") {
		vulns = append(vulns, Vulnerability{
			Name:        "HTTP Instead of HTTPS",
			Severity:    "High",
			Description: "Site is not using HTTPS",
			Solution:    "Implement HTTPS and redirect HTTP to HTTPS",
			Confidence:  100,
		})
	}

	// Check for mixed content
	if strings.HasPrefix(url, "https://") {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
		
		testResp, err := vs.client.Do(req)
		if err == nil {
			defer testResp.Body.Close()
			body, _ := io.ReadAll(testResp.Body)
			
			if strings.Contains(string(body), "http://") {
				vulns = append(vulns, Vulnerability{
					Name:        "Mixed Content",
					Severity:    "Medium",
					Description: "Mixed content detected (HTTP resources on HTTPS page)",
					Solution:    "Use HTTPS for all resources",
					Confidence:  85,
				})
			}
		}
	}

	return vulns
}

func (vs *VulnScanner) ScanMultiple(urls []string) map[string][]Vulnerability {
	results := make(map[string][]Vulnerability)
	
	for _, url := range urls {
		if vulns, err := vs.ScanURL(url); err == nil {
			results[url] = vulns
		}
	}
	
	return results
}
