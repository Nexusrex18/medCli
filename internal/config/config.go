package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Cache   CacheConfig   `mapstructure:"cache"`
	Display DisplayConfig `mapstructure:"display"`
	CSV     CSVConfig     `mapstructure:"csv"`
}

type CSVConfig struct {
	FilePath string `mapstructure:"file_path"`
}

type CacheConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	TTL      string `mapstructure:"ttl"`
	MaxItems int    `mapstructure:"max_items"`
}

type DisplayConfig struct {
	Theme       string `mapstructure:"theme"`
	Animations  bool   `mapstructure:"animations"`
	PageSize    int    `mapstructure:"page_size"`
	AutoRefresh bool   `mapstructure:"auto_refresh"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Search in multiple locations
	v.AddConfigPath(".")                      // Current directory
	v.AddConfigPath("$HOME/.medCli")          // User config
	v.AddConfigPath("/etc/medCli/")           // System config
	v.AddConfigPath("/usr/local/etc/medCli/") // Local config

	// Set defaults
	v.SetDefault("display.theme", "dark")
	v.SetDefault("display.animations", true)
	v.SetDefault("display.page_size", 10)
	v.SetDefault("display.auto_refresh", true)
	v.SetDefault("cache.enabled", true)
	v.SetDefault("cache.ttl", "1h")
	v.SetDefault("cache.max_items", 1000)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found, using defaults")
		} else {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Always find the CSV file dynamically
	foundCSVPath := findCSVFile()
	log.Printf("Using CSV file at: %s", foundCSVPath)
	cfg.CSV.FilePath = foundCSVPath

	return &cfg, nil
}

// findCSVFile tries to locate the CSV file in various locations
func findCSVFile() string {
	// Get executable path to find data relative to binary
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Could not get executable path: %v", err)
		exePath = ""
	}

	possiblePaths := []string{
		// System installation paths (highest priority for installed binaries)
		"/usr/local/share/medCli/medicine_data.csv", // System-wide installation
		"/usr/share/medCli/medicine_data.csv",       // System share
		"/etc/medCli/medicine_data.csv",             // System config

		// User paths
		"$HOME/.medCli/medicine_data.csv", // User config

		// Development paths (lowest priority)
		"data/medicine_data.csv",   // Local development
		"./data/medicine_data.csv", // Current directory
		"medicine_data.csv",        // Current directory
	}

	// Add paths relative to executable (medium priority)
	if exePath != "" {
		exeDir := filepath.Dir(exePath)

		// If installed in system directories, prioritize system data paths
		if isSystemBinary(exeDir) {
			// Prepend system paths for installed binaries
			systemPaths := []string{
				"/usr/local/share/medCli/medicine_data.csv",
				"/usr/share/medCli/medicine_data.csv",
				"/etc/medCli/medicine_data.csv",
			}
			possiblePaths = append(systemPaths, possiblePaths...)
		}

		// Always check relative to executable
		relativePaths := []string{
			filepath.Join(exeDir, "medicine_data.csv"),
			filepath.Join(exeDir, "data/medicine_data.csv"),
			filepath.Join(filepath.Dir(exeDir), "share/medCli/medicine_data.csv"),
		}
		possiblePaths = append(relativePaths, possiblePaths...)
	}

	for _, path := range possiblePaths {
		expandedPath := os.ExpandEnv(path)
		if fileExists(expandedPath) {
			log.Printf("Found CSV file at: %s", expandedPath)
			return expandedPath
		}
	}

	log.Printf("CSV file not found in any standard location")
	return "medicine_data.csv" // Fallback
}

// isSystemBinary checks if the binary is installed in a system directory
func isSystemBinary(exeDir string) bool {
	systemDirs := []string{
		"/usr/local/bin",
		"/usr/bin",
		"/bin",
		"/sbin",
		"/usr/sbin",
	}

	for _, dir := range systemDirs {
		if exeDir == dir {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return !info.IsDir()
}

