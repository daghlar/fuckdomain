# Professional Subdomain Finder

A comprehensive, modular subdomain enumeration and security analysis tool written in Go with advanced features and real-time web interface.

## 🚀 Features

### Core Functionality
- **Modular Architecture**: Clean separation of concerns with dedicated modules
- **DNS Resolution**: Multiple DNS record types support (A, CNAME, MX, TXT, NS, SOA)
- **HTTP/HTTPS Checking**: Advanced response analysis with status codes, headers, and content
- **Wordlist Support**: Built-in comprehensive wordlist with custom wordlist support
- **Concurrent Processing**: Multi-threaded subdomain enumeration with configurable threads
- **Multiple Output Formats**: Plain text, JSON, XML, and HTML report support
- **Real-time Web Interface**: Modern web UI for live monitoring and scanning

### Advanced Security Analysis
- **Port Scanning**: Comprehensive port scanning with service detection
- **SSL/TLS Analysis**: Certificate validation, expiration checks, and security grading
- **Technology Detection**: Automatic detection of web technologies and frameworks
- **Vulnerability Scanning**: Common web vulnerability detection and assessment
- **Screenshot Capture**: Automatic screenshot capture for visual analysis
- **Directory Brute-forcing**: Directory and file enumeration capabilities
- **Risk Assessment**: Automated risk level calculation and confidence scoring

### Professional Features
- **Progress Tracking**: Real-time progress bars and statistics
- **Rate Limiting**: Configurable rate limiting to avoid overwhelming targets
- **Retry Logic**: Intelligent retry mechanisms for failed requests
- **Error Handling**: Comprehensive error handling and logging
- **Configuration Management**: YAML-based configuration with CLI overrides
- **Docker Support**: Containerized deployment with Docker and Docker Compose
- **CI/CD Pipeline**: Automated testing and deployment with GitHub Actions

## 📦 Installation

### Prerequisites

- Go 1.21 or higher
- Internet connection for DNS resolution
- Chrome/Chromium for screenshot functionality (optional)

### Quick Start

```bash
git clone https://github.com/daghlar/fuckdomain.git
cd fuckdomain
go mod tidy
go build -o subdomain-finder
```

### Docker Installation

```bash
git clone https://github.com/daghlar/fuckdomain.git
cd fuckdomain
docker-compose up --build
```

### Binary Download

Download pre-compiled binaries from the [Releases](https://github.com/daghlar/fuckdomain/releases) page.

## 🎯 Usage

### CLI Usage

#### Basic Subdomain Enumeration
```bash
./subdomain-finder scan example.com
```

#### Advanced Security Scan
```bash
./subdomain-finder scan example.com --wordlist custom-wordlist.txt --threads 20 --timeout 10 --output results.txt --verbose --json
```

#### Web Interface
```bash
./subdomain-finder web --port 8080
```
Then open http://localhost:8080 in your browser

### Command Line Options

#### Scan Command
- `--domain`: Target domain to find subdomains for (required)
- `--wordlist`: Path to custom wordlist file (default: wordlists/common.txt)
- `--threads`: Number of concurrent threads (default: 10)
- `--timeout`: Timeout in seconds for DNS/HTTP requests (default: 5)
- `--output`: Output file to save results (optional)
- `--verbose`: Enable verbose output (default: false)
- `--json`: Save results as JSON format (default: false)
- `--xml`: Save results as XML format (default: false)
- `--progress`: Show progress bar (default: true)
- `--stats`: Show detailed statistics (default: false)
- `--no-color`: Disable colored output (default: false)
- `--user-agent`: Custom User-Agent string
- `--headers`: Custom HTTP headers (format: "Header:Value")
- `--retries`: Number of retries for failed requests (default: 3)
- `--delay`: Delay between requests in milliseconds (default: 100)
- `--rate-limit`: Maximum requests per second (default: 10)

#### Web Command
- `--port`: Web interface port (default: 8080)
- `--host`: Web interface host (default: localhost)

#### Config Command
- `--init`: Initialize configuration file
- `--show`: Show current configuration

## 📋 Examples

### Basic Subdomain Enumeration
```bash
./subdomain-finder scan google.com
```

### Using Custom Wordlist
```bash
./subdomain-finder scan example.com --wordlist /path/to/wordlist.txt
```

### High Performance Security Scan
```bash
./subdomain-finder scan example.com --threads 50 --timeout 3 --verbose --stats
```

### Save Results in Multiple Formats
```bash
./subdomain-finder scan example.com --output results.txt --json --xml
```

### Web Interface with Custom Port
```bash
./subdomain-finder web --port 9090
```

### Comprehensive Security Analysis
```bash
./subdomain-finder scan target.com --threads 20 --timeout 10 --verbose --json --stats --progress
```

### Docker Deployment
```bash
docker-compose up -d
# Web interface available at http://localhost:8080
```

## 🏗️ Project Structure

```
fuckdomain/
├── main.go                    # Main entry point
├── go.mod                     # Go module file
├── go.sum                     # Go module checksums
├── Dockerfile                 # Docker container definition
├── docker-compose.yml         # Multi-container setup
├── Makefile                   # Build automation
├── .github/workflows/         # CI/CD pipeline
├── cmd/                       # CLI commands
│   ├── root.go               # Root command
│   ├── scan.go               # Scan command
│   ├── web.go                # Web interface command
│   └── config.go             # Configuration command
├── internal/                  # Internal modules
│   ├── finder/               # Main orchestration
│   ├── dns/                  # DNS resolution
│   ├── http/                 # HTTP/HTTPS checking
│   ├── portscanner/          # Port scanning
│   ├── ssl/                  # SSL/TLS analysis
│   ├── techdetect/           # Technology detection
│   ├── vulnscanner/          # Vulnerability scanning
│   ├── screenshot/           # Screenshot capture
│   ├── bruteforce/           # Directory brute-forcing
│   ├── wordlist/             # Wordlist management
│   ├── output/               # Output formatting
│   ├── reporter/             # Report generation
│   ├── web/                  # Web interface
│   ├── types/                # Data structures
│   ├── config/               # Configuration management
│   ├── logger/               # Logging system
│   ├── limiter/              # Rate limiting
│   ├── progress/             # Progress tracking
│   └── errors/               # Error handling
├── wordlists/                 # Wordlist files
│   └── common.txt            # Default wordlist
└── README.md                 # This file
```

## 🔧 Modules

### Core Modules
- **Finder**: Main orchestration logic, coordinates all modules
- **DNS**: A, CNAME, MX, TXT, NS, SOA record resolution
- **HTTP**: HTTP/HTTPS response checking, status codes, headers
- **Wordlist**: Built-in and custom wordlist management

### Security Analysis Modules
- **Port Scanner**: Comprehensive port scanning with service detection
- **SSL Analyzer**: Certificate validation, expiration checks, security grading
- **Tech Detector**: Automatic technology and framework detection
- **Vuln Scanner**: Common web vulnerability detection and assessment

### Advanced Modules
- **Screenshot**: Automatic screenshot capture for visual analysis
- **Brute Force**: Directory and file enumeration capabilities
- **Reporter**: HTML, PDF, and other format report generation
- **Web Interface**: Real-time web UI for monitoring and scanning

### Utility Modules
- **Output**: Colored terminal output, multiple file formats
- **Config**: YAML-based configuration management
- **Logger**: Structured logging with multiple levels
- **Limiter**: Rate limiting and retry mechanisms
- **Progress**: Real-time progress bars and statistics
- **Types**: Comprehensive data structures and types

## 📦 Dependencies

### Core Dependencies
- `github.com/spf13/cobra`: CLI framework
- `github.com/spf13/viper`: Configuration management
- `github.com/sirupsen/logrus`: Structured logging
- `github.com/miekg/dns`: DNS client library
- `github.com/fatih/color`: Terminal color output

### Security & Analysis
- `github.com/chromedp/chromedp`: Screenshot capture
- `crypto/tls`: SSL/TLS analysis
- `crypto/x509`: Certificate parsing

### Utilities
- `github.com/cheggaaa/pb/v3`: Progress bars
- `github.com/go-playground/validator/v10`: Input validation
- `gopkg.in/yaml.v3`: YAML configuration
- `github.com/gobwas/ws`: WebSocket support

## 🚀 Quick Start

1. **Clone the repository**:
   ```bash
   git clone https://github.com/daghlar/fuckdomain.git
   cd fuckdomain
   ```

2. **Build the tool**:
   ```bash
   go mod tidy
   go build -o subdomain-finder
   ```

3. **Run a basic scan**:
   ```bash
   ./subdomain-finder scan example.com
   ```

4. **Start web interface**:
   ```bash
   ./subdomain-finder web
   # Open http://localhost:8080
   ```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This tool is for **educational and authorized testing purposes only**. Always ensure you have explicit permission to test the target domain before using this tool. Unauthorized scanning may violate laws and terms of service.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/daghlar/fuckdomain/issues)
- **Discussions**: [GitHub Discussions](https://github.com/daghlar/fuckdomain/discussions)
- **Documentation**: [Wiki](https://github.com/daghlar/fuckdomain/wiki)

## 🌟 Features Roadmap

- [ ] Database integration for result storage
- [ ] REST API endpoints
- [ ] Email/Slack notifications
- [ ] Performance optimizations
- [ ] Additional vulnerability checks
- [ ] Machine learning-based subdomain prediction
- [ ] Integration with popular security tools

---

**Made with ❤️ by [daghlar](https://github.com/daghlar)**
