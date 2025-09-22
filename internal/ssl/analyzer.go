package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strings"
	"time"
)

type CertificateInfo struct {
	Subject            string
	Issuer             string
	SerialNumber       string
	NotBefore          time.Time
	NotAfter           time.Time
	SignatureAlgorithm string
	PublicKeyAlgorithm string
	KeyUsage           []string
	DNSNames           []string
	EmailAddresses     []string
	IPAddresses        []string
	IsValid            bool
	DaysUntilExpiry    int
	IsSelfSigned       bool
	IsWildcard         bool
	IsExpired          bool
	IsExpiringSoon     bool
	Strength           string
	Vulnerabilities    []string
}

type SSLResult struct {
	Host           string
	Port           int
	Protocol       string
	Certificate    *CertificateInfo
	SupportedCiphers []string
	SupportedProtocols []string
	IsSecure       bool
	Grade          string
	Recommendations []string
}

type SSLAnalyzer struct {
	timeout time.Duration
}

func NewSSLAnalyzer(timeout time.Duration) *SSLAnalyzer {
	return &SSLAnalyzer{
		timeout: timeout,
	}
}

func (sa *SSLAnalyzer) Analyze(host string, port int) (*SSLResult, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	
	conn, err := net.DialTimeout("tcp", address, sa.timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tlsConn := tls.Client(conn, &tls.Config{
		ServerName: host,
		InsecureSkipVerify: true,
	})

	if err := tlsConn.Handshake(); err != nil {
		return nil, err
	}

	state := tlsConn.ConnectionState()
	cert := state.PeerCertificates[0]

	certInfo := sa.analyzeCertificate(cert)
	supportedCiphers := sa.getSupportedCiphers(tlsConn)
	supportedProtocols := sa.getSupportedProtocols(tlsConn)
	
	isSecure := sa.isSecure(certInfo, supportedCiphers, supportedProtocols)
	grade := sa.calculateGrade(certInfo, supportedCiphers, supportedProtocols)
	recommendations := sa.getRecommendations(certInfo, supportedCiphers, supportedProtocols)

	return &SSLResult{
		Host:                host,
		Port:                port,
		Protocol:            "TLS",
		Certificate:         certInfo,
		SupportedCiphers:    supportedCiphers,
		SupportedProtocols:  supportedProtocols,
		IsSecure:            isSecure,
		Grade:               grade,
		Recommendations:     recommendations,
	}, nil
}

func (sa *SSLAnalyzer) analyzeCertificate(cert *x509.Certificate) *CertificateInfo {
	now := time.Now()
	daysUntilExpiry := int(cert.NotAfter.Sub(now).Hours() / 24)
	
	info := &CertificateInfo{
		Subject:            cert.Subject.String(),
		Issuer:             cert.Issuer.String(),
		SerialNumber:       cert.SerialNumber.String(),
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
		KeyUsage:           sa.getKeyUsage(cert.KeyUsage),
		DNSNames:           cert.DNSNames,
		EmailAddresses:     cert.EmailAddresses,
		IPAddresses:        sa.getIPAddresses(cert.IPAddresses),
		IsValid:            now.After(cert.NotBefore) && now.Before(cert.NotAfter),
		DaysUntilExpiry:    daysUntilExpiry,
		IsSelfSigned:       cert.Issuer.String() == cert.Subject.String(),
		IsWildcard:         sa.isWildcard(cert.DNSNames),
		IsExpired:          now.After(cert.NotAfter),
		IsExpiringSoon:     daysUntilExpiry < 30,
		Strength:           sa.getKeyStrength(cert.PublicKeyAlgorithm),
		Vulnerabilities:    sa.checkVulnerabilities(cert),
	}

	return info
}

func (sa *SSLAnalyzer) getKeyUsage(keyUsage x509.KeyUsage) []string {
	var usage []string
	
	if keyUsage&x509.KeyUsageDigitalSignature != 0 {
		usage = append(usage, "Digital Signature")
	}
	if keyUsage&x509.KeyUsageContentCommitment != 0 {
		usage = append(usage, "Content Commitment")
	}
	if keyUsage&x509.KeyUsageKeyEncipherment != 0 {
		usage = append(usage, "Key Encipherment")
	}
	if keyUsage&x509.KeyUsageDataEncipherment != 0 {
		usage = append(usage, "Data Encipherment")
	}
	if keyUsage&x509.KeyUsageKeyAgreement != 0 {
		usage = append(usage, "Key Agreement")
	}
	if keyUsage&x509.KeyUsageCertSign != 0 {
		usage = append(usage, "Certificate Sign")
	}
	if keyUsage&x509.KeyUsageCRLSign != 0 {
		usage = append(usage, "CRL Sign")
	}
	if keyUsage&x509.KeyUsageEncipherOnly != 0 {
		usage = append(usage, "Encipher Only")
	}
	if keyUsage&x509.KeyUsageDecipherOnly != 0 {
		usage = append(usage, "Decipher Only")
	}
	
	return usage
}

func (sa *SSLAnalyzer) getIPAddresses(ips []net.IP) []string {
	var ipStrings []string
	for _, ip := range ips {
		ipStrings = append(ipStrings, ip.String())
	}
	return ipStrings
}

func (sa *SSLAnalyzer) isWildcard(dnsNames []string) bool {
	for _, name := range dnsNames {
		if strings.HasPrefix(name, "*.") {
			return true
		}
	}
	return false
}

func (sa *SSLAnalyzer) getKeyStrength(algorithm x509.PublicKeyAlgorithm) string {
	switch algorithm {
	case x509.RSA:
		return "RSA (strength depends on key size)"
	case x509.DSA:
		return "DSA (deprecated)"
	case x509.ECDSA:
		return "ECDSA (elliptic curve)"
	case x509.Ed25519:
		return "Ed25519 (modern)"
	default:
		return "Unknown"
	}
}

func (sa *SSLAnalyzer) checkVulnerabilities(cert *x509.Certificate) []string {
	var vulnerabilities []string
	
	if cert.SignatureAlgorithm == x509.MD5WithRSA {
		vulnerabilities = append(vulnerabilities, "MD5 signature (weak)")
	}
	if cert.SignatureAlgorithm == x509.SHA1WithRSA {
		vulnerabilities = append(vulnerabilities, "SHA1 signature (weak)")
	}
	if cert.PublicKeyAlgorithm == x509.DSA {
		vulnerabilities = append(vulnerabilities, "DSA public key (deprecated)")
	}
	
	return vulnerabilities
}

func (sa *SSLAnalyzer) getSupportedCiphers(conn *tls.Conn) []string {
	state := conn.ConnectionState()
	var ciphers []string
	
	if state.CipherSuite != 0 {
		ciphers = append(ciphers, tls.CipherSuiteName(state.CipherSuite))
	}
	
	return ciphers
}

func (sa *SSLAnalyzer) getSupportedProtocols(conn *tls.Conn) []string {
	state := conn.ConnectionState()
	var protocols []string
	
	switch state.Version {
	case tls.VersionTLS10:
		protocols = append(protocols, "TLS 1.0")
	case tls.VersionTLS11:
		protocols = append(protocols, "TLS 1.1")
	case tls.VersionTLS12:
		protocols = append(protocols, "TLS 1.2")
	case tls.VersionTLS13:
		protocols = append(protocols, "TLS 1.3")
	}
	
	return protocols
}

func (sa *SSLAnalyzer) isSecure(certInfo *CertificateInfo, ciphers, protocols []string) bool {
	if !certInfo.IsValid || certInfo.IsExpired {
		return false
	}
	
	if certInfo.IsSelfSigned {
		return false
	}
	
	for _, vuln := range certInfo.Vulnerabilities {
		if strings.Contains(vuln, "weak") || strings.Contains(vuln, "deprecated") {
			return false
		}
	}
	
	return true
}

func (sa *SSLAnalyzer) calculateGrade(certInfo *CertificateInfo, ciphers, protocols []string) string {
	score := 100
	
	if certInfo.IsExpired {
		score -= 50
	}
	if certInfo.IsExpiringSoon {
		score -= 20
	}
	if certInfo.IsSelfSigned {
		score -= 30
	}
	if certInfo.IsWildcard {
		score -= 10
	}
	
	for _, vuln := range certInfo.Vulnerabilities {
		if strings.Contains(vuln, "weak") {
			score -= 20
		}
		if strings.Contains(vuln, "deprecated") {
			score -= 15
		}
	}
	
	if len(protocols) == 0 || !sa.hasModernProtocol(protocols) {
		score -= 25
	}
	
	if score >= 90 {
		return "A+"
	} else if score >= 80 {
		return "A"
	} else if score >= 70 {
		return "B"
	} else if score >= 60 {
		return "C"
	} else if score >= 50 {
		return "D"
	} else {
		return "F"
	}
}

func (sa *SSLAnalyzer) hasModernProtocol(protocols []string) bool {
	for _, protocol := range protocols {
		if protocol == "TLS 1.2" || protocol == "TLS 1.3" {
			return true
		}
	}
	return false
}

func (sa *SSLAnalyzer) getRecommendations(certInfo *CertificateInfo, ciphers, protocols []string) []string {
	var recommendations []string
	
	if certInfo.IsExpired {
		recommendations = append(recommendations, "Certificate is expired - renew immediately")
	}
	if certInfo.IsExpiringSoon {
		recommendations = append(recommendations, "Certificate expires soon - plan renewal")
	}
	if certInfo.IsSelfSigned {
		recommendations = append(recommendations, "Use a trusted CA certificate instead of self-signed")
	}
	if !sa.hasModernProtocol(protocols) {
		recommendations = append(recommendations, "Upgrade to TLS 1.2 or 1.3")
	}
	if len(certInfo.Vulnerabilities) > 0 {
		recommendations = append(recommendations, "Fix certificate vulnerabilities")
	}
	if certInfo.IsWildcard {
		recommendations = append(recommendations, "Consider using specific certificates for better security")
	}
	
	return recommendations
}

func (sa *SSLAnalyzer) AnalyzeMultiple(hosts []string, port int) map[string]*SSLResult {
	results := make(map[string]*SSLResult)
	
	for _, host := range hosts {
		if result, err := sa.Analyze(host, port); err == nil {
			results[host] = result
		}
	}
	
	return results
}
