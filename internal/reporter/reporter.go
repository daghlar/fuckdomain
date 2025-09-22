package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"subdomain-finder/internal/types"
	"time"
)

type Reporter struct {
	outputDir string
}

func NewReporter(outputDir string) *Reporter {
	return &Reporter{
		outputDir: outputDir,
	}
}

func (r *Reporter) GenerateSummaryReport(results []types.Result) *types.ScanSummary {
	summary := &types.ScanSummary{
		TotalSubdomains:  len(results),
		FoundSubdomains:  0,
		OpenPorts:        0,
		Vulnerabilities:  0,
		HighRiskItems:    0,
		Technologies:     make([]types.Technology, 0),
		TopPorts:         make([]types.PortInfo, 0),
		RiskDistribution: make(map[string]int),
		TechnologyStats:  make(map[string]int),
		StartTime:        time.Now(),
		EndTime:          time.Now(),
		Metadata:         make(map[string]interface{}),
	}

	techMap := make(map[string]int)
	portMap := make(map[int]int)

	for _, result := range results {
		if result.IP != "" {
			summary.FoundSubdomains++
		}

		// Count open ports
		summary.OpenPorts += len(result.Ports)
		for _, port := range result.Ports {
			portMap[port.Port]++
		}

		// Count vulnerabilities
		summary.Vulnerabilities += len(result.Vulnerabilities)
		for _, vuln := range result.Vulnerabilities {
			if vuln.Severity == "Critical" || vuln.Severity == "High" {
				summary.HighRiskItems++
			}
		}

		// Count technologies
		for _, tech := range result.Technologies {
			techMap[tech.Name]++
		}

		// Count risk levels
		summary.RiskDistribution[result.RiskLevel]++
	}

	// Convert technology stats
	for tech, count := range techMap {
		summary.TechnologyStats[tech] = count
	}

	// Get top ports
	for port, count := range portMap {
		if count > 1 {
			summary.TopPorts = append(summary.TopPorts, types.PortInfo{
				Port:     port,
				Protocol: "tcp",
				State:    "open",
				Service:  r.getServiceName(port),
			})
		}
	}

	summary.ScanDuration = summary.EndTime.Sub(summary.StartTime)
	return summary
}

func (r *Reporter) getServiceName(port int) string {
	services := map[int]string{
		21:    "FTP",
		22:    "SSH",
		23:    "Telnet",
		25:    "SMTP",
		53:    "DNS",
		80:    "HTTP",
		110:   "POP3",
		143:   "IMAP",
		443:   "HTTPS",
		993:   "IMAPS",
		995:   "POP3S",
		3389:  "RDP",
		5432:  "PostgreSQL",
		3306:  "MySQL",
		6379:  "Redis",
		27017: "MongoDB",
	}

	if service, exists := services[port]; exists {
		return service
	}
	return "Unknown"
}

func (r *Reporter) SaveAsJSON(results []types.Result, filename string) error {
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(r.outputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func (r *Reporter) SaveAsXML(results []types.Result, filename string) error {
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(r.outputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, _ = file.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	_, _ = file.WriteString("<subdomain-scan-results>\n")
	_, _ = file.WriteString(fmt.Sprintf("  <scan-info>\n"))
	_, _ = file.WriteString(fmt.Sprintf("    <total-subdomains>%d</total-subdomains>\n", len(results)))
	_, _ = file.WriteString(fmt.Sprintf("    <scan-date>%s</scan-date>\n", time.Now().Format(time.RFC3339)))
	_, _ = file.WriteString(fmt.Sprintf("  </scan-info>\n"))

	for _, result := range results {
		_, _ = file.WriteString("  <subdomain>\n")
		_, _ = file.WriteString(fmt.Sprintf("    <name>%s</name>\n", result.Subdomain))
		_, _ = file.WriteString(fmt.Sprintf("    <ip>%s</ip>\n", result.IP))
		file.WriteString(fmt.Sprintf("    <status>%s</status>\n", result.Status))
		file.WriteString(fmt.Sprintf("    <server>%s</server>\n", result.Server))
		file.WriteString(fmt.Sprintf("    <title>%s</title>\n", result.Title))
		file.WriteString(fmt.Sprintf("    <risk-level>%s</risk-level>\n", result.RiskLevel))
		file.WriteString(fmt.Sprintf("    <confidence>%d</confidence>\n", result.Confidence))
		file.WriteString(fmt.Sprintf("    <response-time>%s</response-time>\n", result.ResponseTime))

		if len(result.Ports) > 0 {
			file.WriteString("    <ports>\n")
			for _, port := range result.Ports {
				file.WriteString("      <port>\n")
				file.WriteString(fmt.Sprintf("        <number>%d</number>\n", port.Port))
				file.WriteString(fmt.Sprintf("        <protocol>%s</protocol>\n", port.Protocol))
				file.WriteString(fmt.Sprintf("        <state>%s</state>\n", port.State))
				file.WriteString(fmt.Sprintf("        <service>%s</service>\n", port.Service))
				file.WriteString("      </port>\n")
			}
			file.WriteString("    </ports>\n")
		}

		if len(result.Technologies) > 0 {
			file.WriteString("    <technologies>\n")
			for _, tech := range result.Technologies {
				file.WriteString("      <technology>\n")
				file.WriteString(fmt.Sprintf("        <name>%s</name>\n", tech.Name))
				file.WriteString(fmt.Sprintf("        <version>%s</version>\n", tech.Version))
				file.WriteString(fmt.Sprintf("        <category>%s</category>\n", tech.Category))
				file.WriteString(fmt.Sprintf("        <confidence>%d</confidence>\n", tech.Confidence))
				file.WriteString("      </technology>\n")
			}
			file.WriteString("    </technologies>\n")
		}

		if len(result.Vulnerabilities) > 0 {
			file.WriteString("    <vulnerabilities>\n")
			for _, vuln := range result.Vulnerabilities {
				file.WriteString("      <vulnerability>\n")
				file.WriteString(fmt.Sprintf("        <name>%s</name>\n", vuln.Name))
				file.WriteString(fmt.Sprintf("        <severity>%s</severity>\n", vuln.Severity))
				file.WriteString(fmt.Sprintf("        <description>%s</description>\n", vuln.Description))
				file.WriteString(fmt.Sprintf("        <solution>%s</solution>\n", vuln.Solution))
				file.WriteString("      </vulnerability>\n")
			}
			file.WriteString("    </vulnerabilities>\n")
		}

		file.WriteString("  </subdomain>\n")
	}

	file.WriteString("</subdomain-scan-results>\n")
	return nil
}

func (r *Reporter) SaveAsCSV(results []types.Result, filename string) error {
	if err := os.MkdirAll(r.outputDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(r.outputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header
	file.WriteString("Subdomain,IP,Status,Server,Title,Risk Level,Confidence,Response Time,Open Ports,Technologies,Vulnerabilities\n")

	for _, result := range results {
		ports := ""
		for i, port := range result.Ports {
			if i > 0 {
				ports += ";"
			}
			ports += fmt.Sprintf("%d:%s", port.Port, port.Service)
		}

		technologies := ""
		for i, tech := range result.Technologies {
			if i > 0 {
				technologies += ";"
			}
			technologies += tech.Name
		}

		vulnerabilities := ""
		for i, vuln := range result.Vulnerabilities {
			if i > 0 {
				vulnerabilities += ";"
			}
			vulnerabilities += vuln.Name
		}

		line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%d,%s,%s,%s,%s\n",
			result.Subdomain,
			result.IP,
			result.Status,
			result.Server,
			result.Title,
			result.RiskLevel,
			result.Confidence,
			result.ResponseTime,
			ports,
			technologies,
			vulnerabilities,
		)
		file.WriteString(line)
	}

	return nil
}
