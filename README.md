# Subdomain Finder

A powerful and modular subdomain enumeration tool written in Go.

## Features

- **Modular Architecture**: Clean separation of concerns with dedicated modules
- **DNS Resolution**: Multiple DNS record types support (A, CNAME, MX, TXT)
- **HTTP Checking**: HTTP/HTTPS response analysis with status codes and headers
- **Wordlist Support**: Built-in wordlist with option to use custom wordlists
- **Concurrent Processing**: Multi-threaded subdomain enumeration
- **Multiple Output Formats**: Plain text, JSON, and XML output support
- **Colored Output**: Beautiful terminal output with color coding
- **Verbose Mode**: Detailed information about found subdomains

## Installation

### Prerequisites

- Go 1.21 or higher
- Internet connection for DNS resolution

### Build from Source

```bash
git clone <repository-url>
cd subdomain-finder
go mod tidy
go build -o subdomain-finder
```

## Usage

### Basic Usage

```bash
./subdomain-finder -domain example.com
```

### Advanced Usage

```bash
./subdomain-finder -domain example.com -wordlist custom-wordlist.txt -threads 20 -timeout 10 -output results.txt -verbose -json
```

### Command Line Options

- `-domain`: Target domain to find subdomains for (required)
- `-wordlist`: Path to custom wordlist file (optional)
- `-threads`: Number of concurrent threads (default: 10)
- `-timeout`: Timeout in seconds for DNS/HTTP requests (default: 5)
- `-output`: Output file to save results (optional)
- `-verbose`: Enable verbose output (default: false)
- `-json`: Save results as JSON format (default: false)
- `-xml`: Save results as XML format (default: false)

## Examples

### Basic Subdomain Enumeration

```bash
./subdomain-finder -domain google.com
```

### Using Custom Wordlist

```bash
./subdomain-finder -domain example.com -wordlist /path/to/wordlist.txt
```

### High Performance Scan

```bash
./subdomain-finder -domain example.com -threads 50 -timeout 3
```

### Save Results in Multiple Formats

```bash
./subdomain-finder -domain example.com -output results.txt -json -xml
```

## Project Structure

```
subdomain-finder/
├── main.go                 # Main entry point
├── go.mod                  # Go module file
├── internal/
│   ├── finder/            # Main finder logic
│   │   └── finder.go
│   ├── dns/               # DNS resolution module
│   │   └── resolver.go
│   ├── http/              # HTTP checking module
│   │   └── checker.go
│   ├── wordlist/          # Wordlist management
│   │   └── wordlist.go
│   └── output/            # Output formatting
│       └── output.go
└── README.md
```

## Modules

### Finder Module
- Main orchestration logic
- Coordinates between DNS, HTTP, and wordlist modules
- Manages concurrent processing

### DNS Module
- A record resolution
- CNAME record resolution
- MX record resolution
- TXT record resolution
- Multiple DNS server support

### HTTP Module
- HTTP/HTTPS response checking
- Status code analysis
- Header extraction
- Title extraction
- Response length analysis

### Wordlist Module
- Built-in default wordlist
- Custom wordlist support
- Word management functions

### Output Module
- Colored terminal output
- Multiple file format support (TXT, JSON, XML)
- Progress tracking
- Summary reporting

## Dependencies

- `github.com/miekg/dns`: DNS client library
- `github.com/fatih/color`: Terminal color output

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Disclaimer

This tool is for educational and authorized testing purposes only. Always ensure you have permission to test the target domain before using this tool.
