package cmd

import (
	"fmt"
	"os"

	"subdomain-finder/internal/config"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
	Long:  "Manage configuration files and settings for subdomain-finder",
}

var initConfigCmd = &cobra.Command{
	Use:   "init [filename]",
	Short: "Initialize a new configuration file",
	Long:  "Create a new configuration file with default settings",
	Args:  cobra.MaximumNArgs(1),
	Run:   runInitConfig,
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current configuration settings",
	Run:   runShowConfig,
}

var validateConfigCmd = &cobra.Command{
	Use:   "validate [filename]",
	Short: "Validate a configuration file",
	Long:  "Check if a configuration file is valid",
	Args:  cobra.ExactArgs(1),
	Run:   runValidateConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(initConfigCmd)
	configCmd.AddCommand(showConfigCmd)
	configCmd.AddCommand(validateConfigCmd)
}

func runInitConfig(cmd *cobra.Command, args []string) {
	filename := ".subdomain-finder.yaml"
	if len(args) > 0 {
		filename = args[0]
	}

	loader := config.NewLoader()
	if err := loader.CreateDefaultConfig(filename); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration file created: %s\n", filename)
}

func runShowConfig(cmd *cobra.Command, args []string) {
	loader := config.NewLoader()
	cfg, err := loader.LoadFromViper()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Current Configuration:")
	fmt.Println("====================")
	fmt.Printf("DNS Servers: %v\n", cfg.DNS.Servers)
	fmt.Printf("DNS Timeout: %v\n", cfg.DNS.Timeout)
	fmt.Printf("DNS Retries: %d\n", cfg.DNS.Retries)
	fmt.Printf("DNS Rate Limit: %d\n", cfg.DNS.RateLimit)
	fmt.Println()
	fmt.Printf("HTTP Timeout: %v\n", cfg.HTTP.Timeout)
	fmt.Printf("HTTP User Agent: %s\n", cfg.HTTP.UserAgent)
	fmt.Printf("HTTP Retries: %d\n", cfg.HTTP.Retries)
	fmt.Printf("HTTP Rate Limit: %d\n", cfg.HTTP.RateLimit)
	fmt.Printf("HTTP Follow Redirects: %t\n", cfg.HTTP.FollowRedirects)
	fmt.Println()
	fmt.Printf("Output Format: %s\n", cfg.Output.Format)
	fmt.Printf("Output Directory: %s\n", cfg.Output.Directory)
	fmt.Printf("Output Filename: %s\n", cfg.Output.Filename)
	fmt.Printf("Output JSON: %t\n", cfg.Output.JSON)
	fmt.Printf("Output XML: %t\n", cfg.Output.XML)
	fmt.Printf("Output CSV: %t\n", cfg.Output.CSV)
	fmt.Printf("Output Color: %t\n", cfg.Output.Color)
	fmt.Printf("Output Verbose: %t\n", cfg.Output.Verbose)
	fmt.Println()
	fmt.Printf("Log Level: %s\n", cfg.Log.Level)
	fmt.Printf("Log Format: %s\n", cfg.Log.Format)
	fmt.Printf("Log File: %s\n", cfg.Log.File)
}

func runValidateConfig(cmd *cobra.Command, args []string) {
	filename := args[0]

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Config file does not exist: %s\n", filename)
		os.Exit(1)
	}

	loader := config.NewLoader()
	cfg, err := loader.LoadFromFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config validation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration file is valid: %s\n", filename)
	fmt.Printf("DNS Servers: %v\n", cfg.DNS.Servers)
	fmt.Printf("Output Directory: %s\n", cfg.Output.Directory)
}
