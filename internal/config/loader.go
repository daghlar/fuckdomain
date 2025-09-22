package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Loader struct {
	validator *validator.Validate
}

func NewLoader() *Loader {
	return &Loader{
		validator: validator.New(),
	}
}

func (l *Loader) LoadFromFile(filename string) (*AppConfig, error) {
	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := l.validator.Struct(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func (l *Loader) LoadFromViper() (*AppConfig, error) {
	config := DefaultConfig()

	if viper.IsSet("dns.servers") {
		config.DNS.Servers = viper.GetStringSlice("dns.servers")
	}
	if viper.IsSet("dns.timeout") {
		config.DNS.Timeout = viper.GetDuration("dns.timeout")
	}
	if viper.IsSet("dns.retries") {
		config.DNS.Retries = viper.GetInt("dns.retries")
	}
	if viper.IsSet("dns.rate_limit") {
		config.DNS.RateLimit = viper.GetInt("dns.rate_limit")
	}

	if viper.IsSet("http.timeout") {
		config.HTTP.Timeout = viper.GetDuration("http.timeout")
	}
	if viper.IsSet("http.user_agent") {
		config.HTTP.UserAgent = viper.GetString("http.user_agent")
	}
	if viper.IsSet("http.retries") {
		config.HTTP.Retries = viper.GetInt("http.retries")
	}
	if viper.IsSet("http.rate_limit") {
		config.HTTP.RateLimit = viper.GetInt("http.rate_limit")
	}
	if viper.IsSet("http.follow_redirects") {
		config.HTTP.FollowRedirects = viper.GetBool("http.follow_redirects")
	}

	if viper.IsSet("output.format") {
		config.Output.Format = viper.GetString("output.format")
	}
	if viper.IsSet("output.directory") {
		config.Output.Directory = viper.GetString("output.directory")
	}
	if viper.IsSet("output.filename") {
		config.Output.Filename = viper.GetString("output.filename")
	}
	if viper.IsSet("output.json") {
		config.Output.JSON = viper.GetBool("output.json")
	}
	if viper.IsSet("output.xml") {
		config.Output.XML = viper.GetBool("output.xml")
	}
	if viper.IsSet("output.csv") {
		config.Output.CSV = viper.GetBool("output.csv")
	}
	if viper.IsSet("output.color") {
		config.Output.Color = viper.GetBool("output.color")
	}
	if viper.IsSet("output.verbose") {
		config.Output.Verbose = viper.GetBool("output.verbose")
	}

	if viper.IsSet("log.level") {
		config.Log.Level = viper.GetString("log.level")
	}
	if viper.IsSet("log.format") {
		config.Log.Format = viper.GetString("log.format")
	}
	if viper.IsSet("log.file") {
		config.Log.File = viper.GetString("log.file")
	}

	if err := l.validator.Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func (l *Loader) SaveToFile(config *AppConfig, filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (l *Loader) CreateDefaultConfig(filename string) error {
	config := DefaultConfig()
	return l.SaveToFile(config, filename)
}

func (l *Loader) ValidateConfig(config *Config) error {
	return l.validator.Struct(config)
}
