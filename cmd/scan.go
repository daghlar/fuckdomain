package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"subdomain-finder/internal/finder"
	"subdomain-finder/internal/logger"
	"subdomain-finder/internal/output"
)

var scanCmd = &cobra.Command{
	Use:   "scan [flags] <domain>",
	Short: "Scan for subdomains of the target domain",
	Long: `Scan for subdomains using various enumeration techniques.
This command will perform DNS resolution and HTTP checking on discovered subdomains.

Examples:
  subdomain-finder scan example.com
  subdomain-finder scan example.com --wordlist custom.txt --threads 50
  subdomain-finder scan example.com --output results.txt --json --xml
  subdomain-finder scan example.com --timeout 10 --rate-limit 100`,
	Args: cobra.ExactArgs(1),
	Run:  runScan,
}

var (
	wordlist     string
	threads      int
	timeout      int
	rateLimit    int
	outputFile   string
	jsonOutput   bool
	xmlOutput    bool
	progress     bool
	stats        bool
	noColor      bool
	userAgent    string
	headers      []string
	retries      int
	delay        int
)

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().StringVarP(&wordlist, "wordlist", "w", "", "Path to custom wordlist file")
	scanCmd.Flags().IntVarP(&threads, "threads", "t", 10, "Number of concurrent threads")
	scanCmd.Flags().IntVar(&timeout, "timeout", 5, "Timeout in seconds for DNS/HTTP requests")
	scanCmd.Flags().IntVarP(&rateLimit, "rate-limit", "r", 0, "Rate limit (requests per second, 0 = no limit)")
	scanCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file to save results")
	scanCmd.Flags().BoolVar(&jsonOutput, "json", false, "Save results as JSON")
	scanCmd.Flags().BoolVar(&xmlOutput, "xml", false, "Save results as XML")
	scanCmd.Flags().BoolVar(&progress, "progress", true, "Show progress bar")
	scanCmd.Flags().BoolVar(&stats, "stats", true, "Show statistics")
	scanCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	scanCmd.Flags().StringVar(&userAgent, "user-agent", "SubdomainFinder/1.0.0", "Custom User-Agent string")
	scanCmd.Flags().StringArrayVar(&headers, "header", []string{}, "Custom headers (format: key:value)")
	scanCmd.Flags().IntVar(&retries, "retries", 3, "Number of retries for failed requests")
	scanCmd.Flags().IntVar(&delay, "delay", 0, "Delay between requests in milliseconds")

	viper.BindPFlag("scan.wordlist", scanCmd.Flags().Lookup("wordlist"))
	viper.BindPFlag("scan.threads", scanCmd.Flags().Lookup("threads"))
	viper.BindPFlag("scan.timeout", scanCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("scan.rate_limit", scanCmd.Flags().Lookup("rate-limit"))
	viper.BindPFlag("scan.output", scanCmd.Flags().Lookup("output"))
	viper.BindPFlag("scan.json", scanCmd.Flags().Lookup("json"))
	viper.BindPFlag("scan.xml", scanCmd.Flags().Lookup("xml"))
	viper.BindPFlag("scan.progress", scanCmd.Flags().Lookup("progress"))
	viper.BindPFlag("scan.stats", scanCmd.Flags().Lookup("stats"))
	viper.BindPFlag("scan.no_color", scanCmd.Flags().Lookup("no-color"))
	viper.BindPFlag("scan.user_agent", scanCmd.Flags().Lookup("user-agent"))
	viper.BindPFlag("scan.headers", scanCmd.Flags().Lookup("header"))
	viper.BindPFlag("scan.retries", scanCmd.Flags().Lookup("retries"))
	viper.BindPFlag("scan.delay", scanCmd.Flags().Lookup("delay"))
}

func runScan(cmd *cobra.Command, args []string) {
	domain := args[0]

	cfg := finder.Config{
		Domain:     domain,
		Wordlist:   wordlist,
		Threads:    threads,
		Timeout:    timeout,
		RateLimit:  rateLimit,
		OutputFile: outputFile,
		Verbose:    viper.GetBool("verbose"),
		JSON:       jsonOutput,
		XML:        xmlOutput,
		Progress:   progress,
		Stats:      stats,
		NoColor:    noColor,
		UserAgent:  userAgent,
		Headers:    headers,
		Retries:    retries,
		Delay:      delay,
	}

	log := logger.NewLogger(viper.GetString("log.level"), viper.GetString("log.format"))
	
	outputter := output.NewOutputter(cfg, log)
	finder := finder.NewFinder(cfg)

	log.Info("Starting subdomain enumeration", "domain", domain)
	
	startTime := time.Now()
	results := finder.Find()
	duration := time.Since(startTime)

	log.Info("Subdomain enumeration completed", 
		"domain", domain, 
		"found", len(results), 
		"duration", duration.String())

	outputter.PrintSummary(len(results), duration)

	if outputFile != "" {
		outputDir := viper.GetString("output.dir")
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Error("Failed to create output directory", "error", err)
			return
		}

		fullPath := filepath.Join(outputDir, outputFile)
		outputter.SaveToFile(results, fullPath)
	}

	if jsonOutput {
		outputDir := viper.GetString("output.dir")
		jsonFile := filepath.Join(outputDir, fmt.Sprintf("%s.json", domain))
		outputter.SaveAsJSON(results, jsonFile)
	}

	if xmlOutput {
		outputDir := viper.GetString("output.dir")
		xmlFile := filepath.Join(outputDir, fmt.Sprintf("%s.xml", domain))
		outputter.SaveAsXML(results, xmlFile)
	}
}
