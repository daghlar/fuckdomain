package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"subdomain-finder/internal/finder"
	"subdomain-finder/internal/logger"
	"subdomain-finder/internal/types"

	"github.com/fatih/color"
)

type Outputter struct {
	config  finder.Config
	logger  *logger.Logger
	results []types.Result
}

func NewOutputter(cfg finder.Config, log *logger.Logger) *Outputter {
	return &Outputter{
		config:  cfg,
		logger:  log,
		results: make([]types.Result, 0),
	}
}

func (o *Outputter) PrintResult(result types.Result, verbose bool) {
	o.results = append(o.results, result)

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	fmt.Printf("[%s] %s -> %s",
		green("FOUND"),
		white(result.Subdomain),
		blue(result.IP))

	if result.Status != "N/A" {
		fmt.Printf(" [%s]", yellow(result.Status))
	}

	if verbose && result.Response != "" {
		fmt.Printf(" | %s", result.Response)
	}

	fmt.Println()
}

func (o *Outputter) PrintHeader(domain string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	fmt.Println()
	fmt.Printf("%s %s %s\n",
		cyan("="),
		bold("SUBDOMAIN FINDER"),
		cyan("="))
	fmt.Printf("Target: %s\n", bold(domain))
	fmt.Printf("Started: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()
}

func (o *Outputter) PrintSummary(totalFound int, duration time.Duration) {
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	fmt.Println()
	fmt.Printf("%s %s %s\n",
		cyan("="),
		bold("SUMMARY"),
		cyan("="))
	fmt.Printf("Total subdomains found: %s\n", green(totalFound))
	fmt.Printf("Duration: %s\n", duration.String())
	fmt.Println()
}

func (o *Outputter) SaveToFile(results []types.Result, filename string) {
	if filename == "" {
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	for _, result := range results {
		line := fmt.Sprintf("%s,%s,%s,%s\n",
			result.Subdomain,
			result.IP,
			result.Status,
			strings.ReplaceAll(result.Response, ",", ";"))
		_, _ = file.WriteString(line)
	}

	fmt.Printf("Results saved to: %s\n", filename)
}

func (o *Outputter) SaveAsJSON(results []types.Result, filename string) {
	if filename == "" {
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating JSON file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(results); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}

	fmt.Printf("JSON results saved to: %s\n", filename)
}

func (o *Outputter) SaveAsXML(results []types.Result, filename string) {
	if filename == "" {
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating XML file: %v\n", err)
		return
	}
	defer file.Close()

	_, _ = file.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	_, _ = file.WriteString("<subdomains>\n")

	for _, result := range results {
		file.WriteString("  <subdomain>\n")
		file.WriteString(fmt.Sprintf("    <name>%s</name>\n", result.Subdomain))
		file.WriteString(fmt.Sprintf("    <ip>%s</ip>\n", result.IP))
		file.WriteString(fmt.Sprintf("    <status>%s</status>\n", result.Status))
		file.WriteString(fmt.Sprintf("    <response>%s</response>\n", result.Response))
		file.WriteString("  </subdomain>\n")
	}

	file.WriteString("</subdomains>\n")
	fmt.Printf("XML results saved to: %s\n", filename)
}

func (o *Outputter) PrintProgress(current, total int) {
	percent := float64(current) / float64(total) * 100
	bar := strings.Repeat("=", int(percent/2))
	spaces := strings.Repeat(" ", 50-int(percent/2))

	fmt.Printf("\r[%s%s] %.1f%% (%d/%d)",
		bar, spaces, percent, current, total)
}

func (o *Outputter) PrintError(message string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("[%s] %s\n", red("ERROR"), message)
}

func (o *Outputter) PrintWarning(message string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("[%s] %s\n", yellow("WARNING"), message)
}

func (o *Outputter) PrintInfo(message string) {
	blue := color.New(color.FgBlue).SprintFunc()
	fmt.Printf("[%s] %s\n", blue("INFO"), message)
}
