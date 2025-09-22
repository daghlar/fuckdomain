package finder

import (
	"sync"
	"time"

	"subdomain-finder/internal/dns"
	"subdomain-finder/internal/http"
	"subdomain-finder/internal/portscanner"
	"subdomain-finder/internal/ssl"
	"subdomain-finder/internal/techdetect"
	"subdomain-finder/internal/types"
	"subdomain-finder/internal/vulnscanner"
	"subdomain-finder/internal/wordlist"
)

type Config struct {
	Domain     string
	Wordlist   string
	Threads    int
	Timeout    int
	RateLimit  int
	OutputFile string
	Verbose    bool
	JSON       bool
	XML        bool
	Progress   bool
	Stats      bool
	NoColor    bool
	UserAgent  string
	Headers    []string
	Retries    int
	Delay      int
}

type Finder struct {
	config       Config
	dns          *dns.Resolver
	http         *http.Checker
	portScanner  *portscanner.PortScanner
	sslAnalyzer  *ssl.SSLAnalyzer
	techDetector *techdetect.TechDetector
	vulnScanner  *vulnscanner.VulnScanner
	wordlist     *wordlist.Wordlist
}

func NewFinder(config Config) *Finder {
	dnsResolver := dns.NewResolver(config.Timeout)
	httpChecker := http.NewChecker(config.Timeout)
	portScanner := portscanner.NewPortScanner(time.Duration(config.Timeout)*time.Second, config.Threads)
	sslAnalyzer := ssl.NewSSLAnalyzer(time.Duration(config.Timeout) * time.Second)
	techDetector := techdetect.NewTechDetector(time.Duration(config.Timeout) * time.Second)
	vulnScanner := vulnscanner.NewVulnScanner(time.Duration(config.Timeout) * time.Second)
	wordlistManager := wordlist.NewWordlist(config.Wordlist)

	return &Finder{
		config:       config,
		dns:          dnsResolver,
		http:         httpChecker,
		portScanner:  portScanner,
		sslAnalyzer:  sslAnalyzer,
		techDetector: techDetector,
		vulnScanner:  vulnScanner,
		wordlist:     wordlistManager,
	}
}

func (f *Finder) Find() []types.Result {
	words := f.wordlist.GetWords()
	results := make([]types.Result, 0)
	resultsChan := make(chan types.Result, len(words))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, f.config.Threads)

	for _, word := range words {
		wg.Add(1)
		go func(w string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			subdomain := w + "." + f.config.Domain
			result := f.checkSubdomain(subdomain)

			if result.Subdomain != "" {
				resultsChan <- result
			}
		}(word)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

func (f *Finder) checkSubdomain(subdomain string) types.Result {
	startTime := time.Now()
	result := types.Result{
		Subdomain: subdomain,
		Timestamp: startTime,
		Metadata:  make(map[string]interface{}),
	}

	// DNS Resolution
	ip, err := f.dns.Resolve(subdomain)
	if err != nil {
		return types.Result{}
	}
	result.IP = ip

	// HTTP Check
	status, response := f.http.Check(subdomain)
	result.Status = status
	result.Response = response

	// Port Scanning
	portResult := f.portScanner.QuickScan(ip)
	if portResult != nil {
		result.Ports = make([]types.PortInfo, 0)
		for _, port := range portResult.Ports {
			if port.State == "open" {
				result.Ports = append(result.Ports, types.PortInfo{
					Port:     port.Port,
					Protocol: port.Protocol,
					State:    port.State,
					Service:  port.Service,
					Banner:   port.Banner,
					Version:  "",
				})
			}
		}
	}

	// SSL Analysis
	if sslResult, err := f.sslAnalyzer.Analyze(subdomain, 443); err == nil {
		result.SSL = &types.SSLInfo{
			Valid:              sslResult.IsSecure,
			Expired:            sslResult.Certificate.IsExpired,
			ExpiresSoon:        sslResult.Certificate.IsExpiringSoon,
			DaysUntilExpiry:    sslResult.Certificate.DaysUntilExpiry,
			Issuer:             sslResult.Certificate.Issuer,
			Subject:            sslResult.Certificate.Subject,
			SerialNumber:       sslResult.Certificate.SerialNumber,
			SignatureAlgorithm: sslResult.Certificate.SignatureAlgorithm,
			PublicKeyAlgorithm: sslResult.Certificate.PublicKeyAlgorithm,
			Grade:              sslResult.Grade,
			Vulnerabilities:    sslResult.Certificate.Vulnerabilities,
			NotBefore:          sslResult.Certificate.NotBefore,
			NotAfter:           sslResult.Certificate.NotAfter,
		}
	}

	// Technology Detection
	if techResult, err := f.techDetector.Detect("https://" + subdomain); err == nil {
		result.Technologies = make([]types.Technology, 0)
		for _, tech := range techResult.Technologies {
			result.Technologies = append(result.Technologies, types.Technology{
				Name:        tech.Name,
				Version:     tech.Version,
				Category:    tech.Category,
				Confidence:  tech.Confidence,
				Description: tech.Description,
				Website:     tech.Website,
			})
		}
		result.Server = techResult.Server
	}

	// Vulnerability Scanning
	if vulns, err := f.vulnScanner.ScanURL("https://" + subdomain); err == nil {
		result.Vulnerabilities = make([]types.Vulnerability, 0)
		for _, vuln := range vulns {
			result.Vulnerabilities = append(result.Vulnerabilities, types.Vulnerability{
				Name:        vuln.Name,
				Severity:    vuln.Severity,
				Description: vuln.Description,
				CVSS:        vuln.CVSS,
				CVE:         vuln.CVE,
				Solution:    vuln.Solution,
				References:  vuln.References,
			})
		}
	}

	// Risk Assessment
	result.RiskLevel = f.assessRisk(result)
	result.Confidence = f.calculateConfidence(result)
	result.ResponseTime = time.Since(startTime)

	return result
}

func (f *Finder) assessRisk(result types.Result) string {
	riskScore := 0

	// Check vulnerabilities
	for _, vuln := range result.Vulnerabilities {
		switch vuln.Severity {
		case "Critical":
			riskScore += 10
		case "High":
			riskScore += 7
		case "Medium":
			riskScore += 4
		case "Low":
			riskScore += 1
		}
	}

	// Check SSL issues
	if result.SSL != nil {
		if result.SSL.Expired {
			riskScore += 8
		}
		if result.SSL.ExpiresSoon {
			riskScore += 3
		}
		if len(result.SSL.Vulnerabilities) > 0 {
			riskScore += 5
		}
	}

	// Check open ports
	if len(result.Ports) > 10 {
		riskScore += 3
	}

	// Check status codes
	switch result.Status {
	case "403":
		riskScore += 2
	case "500":
		riskScore += 5
	}

	if riskScore >= 15 {
		return "high"
	} else if riskScore >= 8 {
		return "medium"
	} else if riskScore >= 3 {
		return "low"
	}
	return "info"
}

func (f *Finder) calculateConfidence(result types.Result) int {
	confidence := 50

	// Base confidence from DNS resolution
	if result.IP != "" {
		confidence += 20
	}

	// HTTP response confidence
	if result.Status != "" {
		confidence += 15
	}

	// Technology detection confidence
	if len(result.Technologies) > 0 {
		confidence += 10
	}

	// SSL analysis confidence
	if result.SSL != nil {
		confidence += 5
	}

	if confidence > 100 {
		confidence = 100
	}

	return confidence
}
