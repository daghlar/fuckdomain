package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
	"subdomain-finder/internal/finder"
	"subdomain-finder/internal/types"
)

type WebServer struct {
	port     int
	results  []types.Result
	summary  *types.ScanSummary
}

func NewWebServer(port int) *WebServer {
	return &WebServer{
		port:    port,
		results: make([]types.Result, 0),
		summary: &types.ScanSummary{},
	}
}

func (ws *WebServer) Start() error {
	http.HandleFunc("/", ws.handleIndex)
	http.HandleFunc("/api/results", ws.handleResults)
	http.HandleFunc("/api/summary", ws.handleSummary)
	http.HandleFunc("/api/scan", ws.handleScan)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	
	fmt.Printf("Web interface starting on http://localhost:%d\n", ws.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", ws.port), nil)
}

func (ws *WebServer) UpdateResults(results []types.Result, summary *types.ScanSummary) {
	ws.results = results
	ws.summary = summary
}

func (ws *WebServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Subdomain Finder - Web Interface</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            color: #333;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        
        .header {
            background: white;
            padding: 30px;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            margin-bottom: 30px;
            text-align: center;
        }
        
        .header h1 {
            color: #667eea;
            font-size: 2.5em;
            margin-bottom: 10px;
        }
        
        .header p {
            color: #666;
            font-size: 1.1em;
        }
        
        .scan-form {
            background: white;
            padding: 30px;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            margin-bottom: 30px;
        }
        
        .form-group {
            margin-bottom: 20px;
        }
        
        .form-group label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
            color: #333;
        }
        
        .form-group input, .form-group select {
            width: 100%;
            padding: 12px;
            border: 2px solid #ddd;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.3s ease;
        }
        
        .form-group input:focus, .form-group select:focus {
            outline: none;
            border-color: #667eea;
        }
        
        .btn {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 12px 30px;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            cursor: pointer;
            transition: transform 0.3s ease;
        }
        
        .btn:hover {
            transform: translateY(-2px);
        }
        
        .btn:disabled {
            opacity: 0.6;
            cursor: not-allowed;
        }
        
        .summary-cards {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        
        .card {
            background: white;
            padding: 25px;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            text-align: center;
            transition: transform 0.3s ease;
        }
        
        .card:hover {
            transform: translateY(-5px);
        }
        
        .card h3 {
            color: #667eea;
            margin-bottom: 10px;
            font-size: 1.2em;
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
        
        .results-section {
            background: white;
            border-radius: 15px;
            padding: 30px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
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
            border-radius: 10px;
            margin-bottom: 15px;
            overflow: hidden;
            transition: all 0.3s ease;
        }
        
        .subdomain-item:hover {
            box-shadow: 0 5px 20px rgba(0,0,0,0.1);
        }
        
        .subdomain-header {
            background: #f8f9fa;
            padding: 15px 20px;
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
        
        .loading {
            text-align: center;
            padding: 40px;
            color: #666;
        }
        
        .error {
            background: #f8d7da;
            color: #721c24;
            padding: 15px;
            border-radius: 8px;
            margin: 20px 0;
        }
        
        .success {
            background: #d4edda;
            color: #155724;
            padding: 15px;
            border-radius: 8px;
            margin: 20px 0;
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
            <h1>🔍 Subdomain Finder</h1>
            <p>Professional subdomain enumeration and security analysis</p>
        </div>
        
        <div class="scan-form">
            <h2>Start New Scan</h2>
            <form id="scanForm">
                <div class="form-group">
                    <label for="domain">Domain:</label>
                    <input type="text" id="domain" name="domain" placeholder="example.com" required>
                </div>
                <div class="form-group">
                    <label for="threads">Threads:</label>
                    <select id="threads" name="threads">
                        <option value="5">5</option>
                        <option value="10" selected>10</option>
                        <option value="20">20</option>
                        <option value="50">50</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="timeout">Timeout (seconds):</label>
                    <select id="timeout" name="timeout">
                        <option value="5">5</option>
                        <option value="10" selected>10</option>
                        <option value="30">30</option>
                    </select>
                </div>
                <button type="submit" class="btn" id="scanBtn">Start Scan</button>
            </form>
        </div>
        
        <div id="summary" class="summary-cards" style="display: none;">
            <div class="card">
                <h3>Total Subdomains</h3>
                <div class="number" id="totalSubdomains">0</div>
                <div class="label">Scanned</div>
            </div>
            <div class="card">
                <h3>Found Subdomains</h3>
                <div class="number" id="foundSubdomains">0</div>
                <div class="label">Active</div>
            </div>
            <div class="card">
                <h3>Open Ports</h3>
                <div class="number" id="openPorts">0</div>
                <div class="label">Discovered</div>
            </div>
            <div class="card">
                <h3>Vulnerabilities</h3>
                <div class="number" id="vulnerabilities">0</div>
                <div class="label">Found</div>
            </div>
        </div>
        
        <div class="results-section">
            <h2>Scan Results</h2>
            <div id="results">
                <div class="loading">No scan results yet. Start a scan to see results here.</div>
            </div>
        </div>
    </div>
    
    <script>
        let isScanning = false;
        
        document.getElementById('scanForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            if (isScanning) return;
            
            const domain = document.getElementById('domain').value;
            const threads = document.getElementById('threads').value;
            const timeout = document.getElementById('timeout').value;
            
            isScanning = true;
            document.getElementById('scanBtn').disabled = true;
            document.getElementById('scanBtn').textContent = 'Scanning...';
            
            try {
                const response = await fetch('/api/scan', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        domain: domain,
                        threads: parseInt(threads),
                        timeout: parseInt(timeout)
                    })
                });
                
                if (!response.ok) {
                    throw new Error('Scan failed');
                }
                
                const data = await response.json();
                updateResults(data.results, data.summary);
                
            } catch (error) {
                document.getElementById('results').innerHTML = 
                    '<div class="error">Scan failed: ' + error.message + '</div>';
            } finally {
                isScanning = false;
                document.getElementById('scanBtn').disabled = false;
                document.getElementById('scanBtn').textContent = 'Start Scan';
            }
        });
        
        function updateResults(results, summary) {
            // Update summary cards
            document.getElementById('totalSubdomains').textContent = summary.total_subdomains;
            document.getElementById('foundSubdomains').textContent = summary.found_subdomains;
            document.getElementById('openPorts').textContent = summary.open_ports;
            document.getElementById('vulnerabilities').textContent = summary.vulnerabilities;
            document.getElementById('summary').style.display = 'grid';
            
            // Update results
            let resultsHtml = '';
            if (results.length === 0) {
                resultsHtml = '<div class="loading">No subdomains found.</div>';
            } else {
                results.forEach(function(result) {
                    resultsHtml += '<div class="subdomain-item">' +
                        '<div class="subdomain-header" onclick="toggleDetails(this)">' +
                        '<div class="subdomain-name">' + result.subdomain + '</div>' +
                        '<div class="subdomain-status status-' + result.status + '">' + result.status + '</div>' +
                        '<span class="toggle-icon">▼</span>' +
                        '</div>' +
                        '<div class="subdomain-details">' +
                        '<div class="detail-grid">' +
                        '<div class="detail-item">' +
                        '<div class="detail-label">IP Address</div>' +
                        '<div class="detail-value">' + result.ip + '</div>' +
                        '</div>' +
                        '<div class="detail-item">' +
                        '<div class="detail-label">Server</div>' +
                        '<div class="detail-value">' + (result.server || 'Unknown') + '</div>' +
                        '</div>' +
                        '<div class="detail-item">' +
                        '<div class="detail-label">Title</div>' +
                        '<div class="detail-value">' + (result.title || 'N/A') + '</div>' +
                        '</div>' +
                        '<div class="detail-item">' +
                        '<div class="detail-label">Risk Level</div>' +
                        '<div class="detail-value">' + result.risk_level + '</div>' +
                        '</div>' +
                        '<div class="detail-item">' +
                        '<div class="detail-label">Confidence</div>' +
                        '<div class="detail-value">' + result.confidence + '%</div>' +
                        '</div>' +
                        '<div class="detail-item">' +
                        '<div class="detail-label">Response Time</div>' +
                        '<div class="detail-value">' + result.response_time + '</div>' +
                        '</div>' +
                        '</div>' +
                        '</div>' +
                        '</div>';
                });
            }
            
            document.getElementById('results').innerHTML = resultsHtml;
        }
        
        function toggleDetails(element) {
            const details = element.nextElementSibling;
            const icon = element.querySelector('.toggle-icon');
            
            if (details.classList.contains('active')) {
                details.classList.remove('active');
                icon.textContent = '▼';
            } else {
                details.classList.add('active');
                icon.textContent = '▲';
            }
        }
        
        // Load existing results on page load
        window.addEventListener('load', async function() {
            try {
                const response = await fetch('/api/results');
                if (response.ok) {
                    const results = await response.json();
                    const summaryResponse = await fetch('/api/summary');
                    if (summaryResponse.ok) {
                        const summary = await summaryResponse.json();
                        updateResults(results, summary);
                    }
                }
            } catch (error) {
                console.log('No existing results');
            }
        });
    </script>
</body>
</html>
`
	
	tmplParsed := template.Must(template.New("index").Parse(tmpl))
	tmplParsed.Execute(w, nil)
}

func (ws *WebServer) handleResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ws.results)
}

func (ws *WebServer) handleSummary(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ws.summary)
}

func (ws *WebServer) handleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var scanRequest struct {
		Domain  string `json:"domain"`
		Threads int    `json:"threads"`
		Timeout int    `json:"timeout"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&scanRequest); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Gerçek tarama yap
	results, summary := ws.performRealScan(scanRequest.Domain, scanRequest.Threads, scanRequest.Timeout)
	
	ws.UpdateResults(results, summary)
	
	response := map[string]interface{}{
		"results": results,
		"summary": summary,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (ws *WebServer) performRealScan(domain string, threads, timeout int) ([]types.Result, *types.ScanSummary) {
	// Gerçek tarama yapmak için finder modülünü kullan
	// Önce finder modülünü import edelim
	results, summary := ws.runActualScan(domain, threads, timeout)
	return results, summary
}

func (ws *WebServer) runActualScan(domain string, threads, timeout int) ([]types.Result, *types.ScanSummary) {
	// Gerçek tarama yap
	startTime := time.Now()
	
	// Önce mevcut sonuçları kontrol et
	jsonFile := fmt.Sprintf("results/%s.json", domain)
	if data, err := os.ReadFile(jsonFile); err == nil {
		var results []types.Result
		if err := json.Unmarshal(data, &results); err == nil {
			// Summary oluştur
			summary := &types.ScanSummary{
				TotalSubdomains: len(results),
				FoundSubdomains: 0,
				OpenPorts:       0,
				Vulnerabilities: 0,
				HighRiskItems:   0,
				ScanDuration:    time.Since(startTime),
				StartTime:       startTime,
				EndTime:         time.Now(),
			}
			
			for _, result := range results {
				if result.IP != "" {
					summary.FoundSubdomains++
				}
				summary.OpenPorts += len(result.Ports)
				summary.Vulnerabilities += len(result.Vulnerabilities)
				
				for _, vuln := range result.Vulnerabilities {
					if vuln.Severity == "Critical" || vuln.Severity == "High" {
						summary.HighRiskItems++
					}
				}
			}
			
			return results, summary
		}
	}
	
	// Eğer dosya yoksa, gerçek tarama yap
	config := finder.Config{
		Domain:     domain,
		Wordlist:   "wordlists/common.txt",
		Threads:    threads,
		Timeout:    timeout,
		RateLimit:  10,
		OutputFile: fmt.Sprintf("results/%s.txt", domain),
		Verbose:    false,
		JSON:       true,
		XML:        false,
		Progress:   false,
		Stats:      false,
		NoColor:    true,
		UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		Headers:    []string{},
		Retries:    3,
		Delay:      100,
	}
	
	// Finder oluştur ve gerçek tarama yap
	finderInstance := finder.NewFinder(config)
	results := finderInstance.Find()
	
	// Summary oluştur
	summary := &types.ScanSummary{
		TotalSubdomains: len(results),
		FoundSubdomains: 0,
		OpenPorts:       0,
		Vulnerabilities: 0,
		HighRiskItems:   0,
		ScanDuration:    time.Since(startTime),
		StartTime:       startTime,
		EndTime:         time.Now(),
	}
	
	for _, result := range results {
		if result.IP != "" {
			summary.FoundSubdomains++
		}
		summary.OpenPorts += len(result.Ports)
		summary.Vulnerabilities += len(result.Vulnerabilities)
		
		for _, vuln := range result.Vulnerabilities {
			if vuln.Severity == "Critical" || vuln.Severity == "High" {
				summary.HighRiskItems++
			}
		}
	}
	
	// Sonuçları JSON dosyasına kaydet
	if data, err := json.MarshalIndent(results, "", "  "); err == nil {
		os.WriteFile(jsonFile, data, 0644)
	}
	
	return results, summary
}
