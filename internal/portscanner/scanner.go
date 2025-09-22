package portscanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PortScanner struct {
	timeout     time.Duration
	threads     int
	commonPorts []int
}

type PortResult struct {
	Port     int
	Protocol string
	State    string
	Service  string
	Banner   string
}

type ScanResult struct {
	Host       string
	Ports      []PortResult
	OpenPorts  int
	TotalPorts int
}

func NewPortScanner(timeout time.Duration, threads int) *PortScanner {
	return &PortScanner{
		timeout: timeout,
		threads: threads,
		commonPorts: []int{
			21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 443, 993, 995, 1723, 3306, 3389, 5432, 5900, 8080, 8443, 8888, 9000, 9090, 9200, 9300, 11211, 27017, 6379, 5984, 9200, 9300, 11211, 27017, 6379, 5984,
		},
	}
}

func (ps *PortScanner) ScanHost(host string, ports []int) *ScanResult {
	if len(ports) == 0 {
		ports = ps.commonPorts
	}

	result := &ScanResult{
		Host:       host,
		Ports:      make([]PortResult, 0),
		TotalPorts: len(ports),
	}

	semaphore := make(chan struct{}, ps.threads)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			portResult := ps.scanPort(host, p)

			mu.Lock()
			if portResult.State == "open" {
				result.OpenPorts++
			}
			result.Ports = append(result.Ports, portResult)
			mu.Unlock()
		}(port)
	}

	wg.Wait()
	return result
}

func (ps *PortScanner) scanPort(host string, port int) PortResult {
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", address, ps.timeout)
	if err != nil {
		return PortResult{
			Port:     port,
			Protocol: "tcp",
			State:    "closed",
			Service:  ps.getServiceName(port),
		}
	}
	defer conn.Close()

	banner := ps.getBanner(conn, port)
	service := ps.getServiceName(port)

	return PortResult{
		Port:     port,
		Protocol: "tcp",
		State:    "open",
		Service:  service,
		Banner:   banner,
	}
}

func (ps *PortScanner) getBanner(conn net.Conn, port int) string {
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return ""
	}

	banner := string(buffer[:n])
	banner = strings.TrimSpace(banner)
	banner = strings.ReplaceAll(banner, "\r\n", " ")
	banner = strings.ReplaceAll(banner, "\n", " ")

	if len(banner) > 200 {
		banner = banner[:200] + "..."
	}

	return banner
}

func (ps *PortScanner) getServiceName(port int) string {
	services := map[int]string{
		21:    "ftp",
		22:    "ssh",
		23:    "telnet",
		25:    "smtp",
		53:    "dns",
		80:    "http",
		110:   "pop3",
		111:   "rpcbind",
		135:   "msrpc",
		139:   "netbios-ssn",
		143:   "imap",
		443:   "https",
		993:   "imaps",
		995:   "pop3s",
		1723:  "pptp",
		3306:  "mysql",
		3389:  "rdp",
		5432:  "postgresql",
		5900:  "vnc",
		8080:  "http-proxy",
		8443:  "https-alt",
		8888:  "http-alt",
		9000:  "http-alt",
		9090:  "http-alt",
		9200:  "elasticsearch",
		9300:  "elasticsearch",
		11211: "memcached",
		27017: "mongodb",
		6379:  "redis",
		5984:  "couchdb",
	}

	if service, exists := services[port]; exists {
		return service
	}
	return "unknown"
}

func (ps *PortScanner) ScanMultipleHosts(hosts []string, ports []int) map[string]*ScanResult {
	results := make(map[string]*ScanResult)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, host := range hosts {
		wg.Add(1)
		go func(h string) {
			defer wg.Done()
			result := ps.ScanHost(h, ports)

			mu.Lock()
			results[h] = result
			mu.Unlock()
		}(host)
	}

	wg.Wait()
	return results
}

func (ps *PortScanner) QuickScan(host string) *ScanResult {
	quickPorts := []int{21, 22, 23, 25, 53, 80, 110, 135, 139, 143, 443, 993, 995, 1723, 3306, 3389, 5432, 5900, 8080, 8443, 8888, 9000, 9090}
	return ps.ScanHost(host, quickPorts)
}

func (ps *PortScanner) FullScan(host string) *ScanResult {
	fullPorts := make([]int, 0, 65535)
	for i := 1; i <= 65535; i++ {
		fullPorts = append(fullPorts, i)
	}
	return ps.ScanHost(host, fullPorts)
}

func (ps *PortScanner) CustomScan(host string, portRange string) *ScanResult {
	ports := ps.parsePortRange(portRange)
	return ps.ScanHost(host, ports)
}

func (ps *PortScanner) parsePortRange(portRange string) []int {
	var ports []int

	if strings.Contains(portRange, "-") {
		parts := strings.Split(portRange, "-")
		if len(parts) == 2 {
			start, err1 := strconv.Atoi(parts[0])
			end, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil {
				for i := start; i <= end; i++ {
					ports = append(ports, i)
				}
			}
		}
	} else if strings.Contains(portRange, ",") {
		parts := strings.Split(portRange, ",")
		for _, part := range parts {
			if port, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
				ports = append(ports, port)
			}
		}
	} else {
		if port, err := strconv.Atoi(portRange); err == nil {
			ports = append(ports, port)
		}
	}

	return ports
}
