package config

import (
	"time"
)

type Config struct {
	Domain     string   `validate:"required"`
	Wordlist   string
	Threads    int      `validate:"min=1,max=1000"`
	Timeout    int      `validate:"min=1,max=60"`
	RateLimit  int      `validate:"min=0"`
	OutputFile string
	Verbose    bool
	JSON       bool
	XML        bool
	Progress   bool
	Stats      bool
	NoColor    bool
	UserAgent  string
	Headers    []string
	Retries    int      `validate:"min=0,max=10"`
	Delay      int      `validate:"min=0"`
}

type DNSConfig struct {
	Servers     []string      `yaml:"servers"`
	Timeout     time.Duration `yaml:"timeout"`
	Retries     int           `yaml:"retries"`
	RateLimit   int           `yaml:"rate_limit"`
}

type HTTPConfig struct {
	Timeout     time.Duration `yaml:"timeout"`
	UserAgent   string        `yaml:"user_agent"`
	Headers     map[string]string `yaml:"headers"`
	Retries     int           `yaml:"retries"`
	RateLimit   int           `yaml:"rate_limit"`
	FollowRedirects bool      `yaml:"follow_redirects"`
}

type OutputConfig struct {
	Format      string `yaml:"format"`
	Directory   string `yaml:"directory"`
	Filename    string `yaml:"filename"`
	JSON        bool   `yaml:"json"`
	XML         bool   `yaml:"xml"`
	CSV         bool   `yaml:"csv"`
	Color       bool   `yaml:"color"`
	Verbose     bool   `yaml:"verbose"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	File   string `yaml:"file"`
}

type AppConfig struct {
	DNS    DNSConfig    `yaml:"dns"`
	HTTP   HTTPConfig   `yaml:"http"`
	Output OutputConfig `yaml:"output"`
	Log    LogConfig    `yaml:"log"`
}

func DefaultConfig() *AppConfig {
	return &AppConfig{
		DNS: DNSConfig{
			Servers:   []string{"8.8.8.8:53", "1.1.1.1:53", "8.8.4.4:53"},
			Timeout:   5 * time.Second,
			Retries:   3,
			RateLimit: 0,
		},
		HTTP: HTTPConfig{
			Timeout:         5 * time.Second,
			UserAgent:       "SubdomainFinder/1.0.0",
			Headers:         make(map[string]string),
			Retries:         3,
			RateLimit:       0,
			FollowRedirects: false,
		},
		Output: OutputConfig{
			Format:    "text",
			Directory: "./results",
			Filename:  "",
			JSON:      false,
			XML:       false,
			CSV:       false,
			Color:     true,
			Verbose:   false,
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
			File:   "",
		},
	}
}
