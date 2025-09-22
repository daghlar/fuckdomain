package reporter

import (
	"html/template"
	"os"
	"path/filepath"
	"subdomain-finder/internal/types"
	"time"
)

type HTMLReporter struct {
	templateDir string
	outputDir   string
}

func NewHTMLReporter(templateDir, outputDir string) *HTMLReporter {
	return &HTMLReporter{
		templateDir: templateDir,
		outputDir:   outputDir,
	}
}

func (hr *HTMLReporter) GenerateReport(summary *types.ScanSummary, results []types.Result, filename string) error {
	if err := os.MkdirAll(hr.outputDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(hr.outputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl := hr.getReportTemplate()
	if err := tmpl.Execute(file, map[string]interface{}{
		"Summary":     summary,
		"Results":     results,
		"GeneratedAt": time.Now(),
	}); err != nil {
		return err
	}

	return nil
}

func (hr *HTMLReporter) getReportTemplate() *template.Template {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Subdomain Security Report</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f5f5f5;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px 0;
            text-align: center;
            border-radius: 10px;
            margin-bottom: 30px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
        }
        
        .header p {
            font-size: 1.2em;
            opacity: 0.9;
        }
        
        .summary-cards {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        
        .card {
            background: white;
            padding: 25px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            text-align: center;
            transition: transform 0.3s ease;
        }
        
        .card:hover {
            transform: translateY(-5px);
        }
        
        .card h3 {
            color: #667eea;
            margin-bottom: 10px;
            font-size: 1.5em;
        }
        
        .card .number {
            font-size: 2.5em;
            font-weight: bold;
            color: #333;
        }
        
        .card .label {
            color: #666;
            margin-top: 5px;
        }
        
        .risk-high { color: #e74c3c; }
        .risk-medium { color: #f39c12; }
        .risk-low { color: #27ae60; }
        .risk-info { color: #3498db; }
        
        .results-section {
            background: white;
            border-radius: 10px;
            padding: 30px;
            margin-bottom: 30px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        
        .results-section h2 {
            color: #333;
            margin-bottom: 20px;
            font-size: 1.8em;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
        }
        
        .subdomain-item {
            border: 1px solid #ddd;
            border-radius: 8px;
            margin-bottom: 15px;
            overflow: hidden;
            transition: all 0.3s ease;
        }
        
        .subdomain-item:hover {
            box-shadow: 0 4px 15px rgba(0,0,0,0.1);
        }
        
        .subdomain-header {
            background: #f8f9fa;
            padding: 15px 20px;
            border-bottom: 1px solid #ddd;
            cursor: pointer;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .subdomain-name {
            font-weight: bold;
            color: #333;
            font-size: 1.1em;
        }
        
        .subdomain-status {
            padding: 5px 15px;
            border-radius: 20px;
            font-size: 0.9em;
            font-weight: bold;
        }
        
        .status-200 { background: #d4edda; color: #155724; }
        .status-301 { background: #fff3cd; color: #856404; }
        .status-302 { background: #fff3cd; color: #856404; }
        .status-403 { background: #f8d7da; color: #721c24; }
        .status-404 { background: #d1ecf1; color: #0c5460; }
        .status-500 { background: #f8d7da; color: #721c24; }
        
        .subdomain-details {
            padding: 20px;
            display: none;
            background: white;
        }
        
        .subdomain-details.active {
            display: block;
        }
        
        .detail-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-bottom: 15px;
        }
        
        .detail-item {
            background: #f8f9fa;
            padding: 10px;
            border-radius: 5px;
        }
        
        .detail-label {
            font-weight: bold;
            color: #666;
            font-size: 0.9em;
        }
        
        .detail-value {
            color: #333;
            margin-top: 5px;
        }
        
        .technologies {
            margin-top: 15px;
        }
        
        .tech-tag {
            display: inline-block;
            background: #667eea;
            color: white;
            padding: 5px 10px;
            border-radius: 15px;
            font-size: 0.8em;
            margin: 2px;
        }
        
        .vulnerabilities {
            margin-top: 15px;
        }
        
        .vuln-item {
            background: #fff5f5;
            border-left: 4px solid #e74c3c;
            padding: 10px;
            margin: 5px 0;
            border-radius: 0 5px 5px 0;
        }
        
        .vuln-severity {
            font-weight: bold;
            color: #e74c3c;
        }
        
        .footer {
            text-align: center;
            color: #666;
            margin-top: 40px;
            padding: 20px;
            border-top: 1px solid #ddd;
        }
        
        .toggle-icon {
            transition: transform 0.3s ease;
        }
        
        .toggle-icon.rotated {
            transform: rotate(180deg);
        }
        
        @media (max-width: 768px) {
            .container {
                padding: 10px;
            }
            
            .header h1 {
                font-size: 2em;
            }
            
            .summary-cards {
                grid-template-columns: 1fr;
            }
            
            .detail-grid {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 Subdomain Security Report</h1>
            <p>Generated on {{.GeneratedAt.Format "January 2, 2006 at 15:04:05 MST"}}</p>
        </div>
        
        <div class="summary-cards">
            <div class="card">
                <h3>Total Subdomains</h3>
                <div class="number">{{.Summary.TotalSubdomains}}</div>
                <div class="label">Scanned</div>
            </div>
            <div class="card">
                <h3>Found Subdomains</h3>
                <div class="number">{{.Summary.FoundSubdomains}}</div>
                <div class="label">Active</div>
            </div>
            <div class="card">
                <h3>Open Ports</h3>
                <div class="number">{{.Summary.OpenPorts}}</div>
                <div class="label">Discovered</div>
            </div>
            <div class="card">
                <h3>Vulnerabilities</h3>
                <div class="number risk-high">{{.Summary.Vulnerabilities}}</div>
                <div class="label">Found</div>
            </div>
            <div class="card">
                <h3>High Risk Items</h3>
                <div class="number risk-high">{{.Summary.HighRiskItems}}</div>
                <div class="label">Critical</div>
            </div>
            <div class="card">
                <h3>Scan Duration</h3>
                <div class="number">{{.Summary.ScanDuration}}</div>
                <div class="label">Time</div>
            </div>
        </div>
        
        <div class="results-section">
            <h2>📊 Detailed Results</h2>
            {{range .Results}}
            <div class="subdomain-item">
                <div class="subdomain-header" onclick="toggleDetails(this)">
                    <div class="subdomain-name">{{.Subdomain}}</div>
                    <div class="subdomain-status status-{{.Status}}">{{.Status}}</div>
                    <span class="toggle-icon">▼</span>
                </div>
                <div class="subdomain-details">
                    <div class="detail-grid">
                        <div class="detail-item">
                            <div class="detail-label">IP Address</div>
                            <div class="detail-value">{{.IP}}</div>
                        </div>
                        <div class="detail-item">
                            <div class="detail-label">Server</div>
                            <div class="detail-value">{{.Server}}</div>
                        </div>
                        <div class="detail-item">
                            <div class="detail-label">Title</div>
                            <div class="detail-value">{{.Title}}</div>
                        </div>
                        <div class="detail-item">
                            <div class="detail-label">Content Length</div>
                            <div class="detail-value">{{.ContentLength}}</div>
                        </div>
                        <div class="detail-item">
                            <div class="detail-label">Response Time</div>
                            <div class="detail-value">{{.ResponseTime}}</div>
                        </div>
                        <div class="detail-item">
                            <div class="detail-label">Risk Level</div>
                            <div class="detail-value risk-{{.RiskLevel}}">{{.RiskLevel}}</div>
                        </div>
                    </div>
                    
                    {{if .Technologies}}
                    <div class="technologies">
                        <strong>Technologies:</strong><br>
                        {{range .Technologies}}
                        <span class="tech-tag">{{.Name}} {{.Version}}</span>
                        {{end}}
                    </div>
                    {{end}}
                    
                    {{if .Vulnerabilities}}
                    <div class="vulnerabilities">
                        <strong>Vulnerabilities:</strong>
                        {{range .Vulnerabilities}}
                        <div class="vuln-item">
                            <span class="vuln-severity">{{.Severity}}</span> - {{.Name}}
                            <br><small>{{.Description}}</small>
                        </div>
                        {{end}}
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}
        </div>
        
        <div class="footer">
            <p>Report generated by Subdomain Finder v1.0.0</p>
            <p>For security purposes, this report should be kept confidential</p>
        </div>
    </div>
    
    <script>
        function toggleDetails(element) {
            const details = element.nextElementSibling;
            const icon = element.querySelector('.toggle-icon');
            
            if (details.classList.contains('active')) {
                details.classList.remove('active');
                icon.classList.remove('rotated');
            } else {
                details.classList.add('active');
                icon.classList.add('rotated');
            }
        }
    </script>
</body>
</html>
`

	return template.Must(template.New("report").Parse(tmpl))
}
