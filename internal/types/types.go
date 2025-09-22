package types

import "time"

type Result struct {
	Subdomain       string                 `json:"subdomain"`
	IP              string                 `json:"ip"`
	Status          string                 `json:"status"`
	Response        string                 `json:"response"`
	Title           string                 `json:"title"`
	Server          string                 `json:"server"`
	ContentLength   int64                  `json:"content_length"`
	ResponseTime    time.Duration          `json:"response_time"`
	Technologies    []Technology           `json:"technologies"`
	Ports           []PortInfo             `json:"ports"`
	SSL             *SSLInfo               `json:"ssl"`
	Vulnerabilities []Vulnerability        `json:"vulnerabilities"`
	Headers         map[string]string      `json:"headers"`
	Cookies         []Cookie               `json:"cookies"`
	Redirects       []Redirect             `json:"redirects"`
	DNS             *DNSInfo               `json:"dns"`
	GeoLocation     *GeoLocation           `json:"geo_location"`
	RiskLevel       string                 `json:"risk_level"`
	Confidence      int                    `json:"confidence"`
	Timestamp       time.Time              `json:"timestamp"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type Technology struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Category    string `json:"category"`
	Confidence  int    `json:"confidence"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Icon        string `json:"icon"`
}

type PortInfo struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	State    string `json:"state"`
	Service  string `json:"service"`
	Banner   string `json:"banner"`
	Version  string `json:"version"`
}

type SSLInfo struct {
	Valid              bool      `json:"valid"`
	Expired            bool      `json:"expired"`
	ExpiresSoon        bool      `json:"expires_soon"`
	DaysUntilExpiry    int       `json:"days_until_expiry"`
	Issuer             string    `json:"issuer"`
	Subject            string    `json:"subject"`
	SerialNumber       string    `json:"serial_number"`
	SignatureAlgorithm string    `json:"signature_algorithm"`
	PublicKeyAlgorithm string    `json:"public_key_algorithm"`
	KeySize            int       `json:"key_size"`
	Grade              string    `json:"grade"`
	Vulnerabilities    []string  `json:"vulnerabilities"`
	NotBefore          time.Time `json:"not_before"`
	NotAfter           time.Time `json:"not_after"`
}

type Vulnerability struct {
	Name        string   `json:"name"`
	Severity    string   `json:"severity"`
	Description string   `json:"description"`
	CVSS        string   `json:"cvss"`
	CVE         string   `json:"cve"`
	Solution    string   `json:"solution"`
	References  []string `json:"references"`
}

type Cookie struct {
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	Domain   string    `json:"domain"`
	Path     string    `json:"path"`
	Expires  time.Time `json:"expires"`
	Secure   bool      `json:"secure"`
	HttpOnly bool      `json:"http_only"`
	SameSite string    `json:"same_site"`
}

type Redirect struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Location   string `json:"location"`
}

type DNSInfo struct {
	ARecords     []string `json:"a_records"`
	AAAARecords  []string `json:"aaaa_records"`
	CNAMERecords []string `json:"cname_records"`
	MXRecords    []string `json:"mx_records"`
	TXTRecords   []string `json:"txt_records"`
	NSRecords    []string `json:"ns_records"`
	SOARecord    string   `json:"soa_record"`
}

type GeoLocation struct {
	Country      string  `json:"country"`
	CountryCode  string  `json:"country_code"`
	Region       string  `json:"region"`
	City         string  `json:"city"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Timezone     string  `json:"timezone"`
	ISP          string  `json:"isp"`
	ASN          string  `json:"asn"`
	Organization string  `json:"organization"`
}

type ScanSummary struct {
	TotalSubdomains  int                    `json:"total_subdomains"`
	FoundSubdomains  int                    `json:"found_subdomains"`
	OpenPorts        int                    `json:"open_ports"`
	Vulnerabilities  int                    `json:"vulnerabilities"`
	HighRiskItems    int                    `json:"high_risk_items"`
	Technologies     []Technology           `json:"technologies"`
	TopPorts         []PortInfo             `json:"top_ports"`
	RiskDistribution map[string]int         `json:"risk_distribution"`
	TechnologyStats  map[string]int         `json:"technology_stats"`
	ScanDuration     time.Duration          `json:"scan_duration"`
	StartTime        time.Time              `json:"start_time"`
	EndTime          time.Time              `json:"end_time"`
	Metadata         map[string]interface{} `json:"metadata"`
}
