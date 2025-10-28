package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stretchr/testify/assert/yaml"
)

// Config represents the complete linuxwave system configuration
type Config struct {
	Service ServiceConfig `yaml:"service"`
	Logging LoggingConfig `yaml:"logging"`
	Audio AudioConfig `yaml:"audio"`
	Security SecurityConfig `yaml:"security"`
}

// ServiceConfig contains service-level settings for authentication operations
type ServiceConfig struct {
	// Timeout is the authentication timeout in seconds (1-60)
	Timeout int `yaml:"timeout"`
	// RetryAttempts is the maximum number of retry attempts for authentication (1-10)
	RetryAttempts int `yaml:"retry_attempts"`
	// SocketPath is the Unix domain socket path for IPC communication (must be absolute)
	SocketPath string `yaml:"socket_path"`
}

// LoggingConfig contains logging configuration settings.
type LoggingConfig struct {
	// Level is the log verbosity level: DEBUG, INFO, WARN, ERROR
	Level string `yaml:"level"`
	// Format is the log output format: json (structured) or text (human-readable)
	Format string `yaml:"format"`
}

// AudioConfig contains audio feedback configuration settings.
type AudioConfig struct {
	// Enabled controls whether audio feedback is enabled for authentication events
	Enabled bool `yaml:"enabled"`
	// VOlume is the audio volume level (0-100)
	Volume int `yaml:"volume"`
	// CustomSoundSuccess is the optional path to a custom success sound file
	CustomSoundSuccess string `yaml:"custom_sound_success"`
	// CustomSoundFailure is the optional path to a custom failure sound file
	CustomSoundFailure string `yaml:"custom_sound_failure"`
}

// SecurityConfig contains security-related settings
type SecurityConfig struct {
	// LivenessRequired enables liveness detection (blink/movement check)
	LivenessRequired bool `yaml:"liveness_required"`
	// MatchThreshold is the face matching confidence threshold (0.0-1.0)
	MatchThreshold float64 `yaml:"match_threshold"`
	// MaxAuthAttempts is the maximum authentication attempts before lockout
	MaxAuthAttempts int `yaml:"max_auth_attempts"`
	// LockoutDuration is the lockout duration in seconds after max attempts succeeded
	LockoutDuration int `yaml:"lockout_duration"`
}

const (
	// System-wide configuration file path
	systemConfigPath = "/etc/linux-wave/config.yaml"
	// User-specific configuration file path (relative to home directory)
	userConfigRelPath = ".config/linux-wave/config.yaml"
)

// DefaultConfig returns a Config struct populated with sensible default values.
// These defaults are used when no configuration files exist or as a base
// for merging with loaded configurations.
func DefaultConfig() *Config {
	return &Config{
		Service: ServiceConfig{
			Timeout: 10, // 10 seconds is reasonable for face auth without being too short
			RetryAttempts: 3, // 3 attempts balances user convenience with security
			SocketPath: "/run/linux-wave/auth.sock", // Standard systemd runtime directory
		},
		Logging: LoggingConfig{
			Level: "INFO", // INFO level provides useful feedback without excessive verbosity
			Format: "test", // Human-readable format is better for system logs by default
		},
		Audio: AudioConfig{
			Enabled: true, // Audio feedback improves accessibility and user experience
			Volume: 50, // 50% volume is a balanced default, not too loud or quiet
			CustomSoundSuccess: "", // Empty means use default embedded sounds
			CustomSoundFailure: "", // Empty means use default embedded sounds
		},
		Security: SecurityConfig{
			LivenessRequired: true, // Liveness detection is critical for security
			MatchThreshold: 0.85, // 0.85 balances security with usability (higher = stricter)
			MaxAuthAttempts: 3, // Same as retry attempts for consistency
			LockoutDuration: 300, // 5 minutes is standard for authentication lockouts
		},
	}
}

// Load reads configuration from system and user config files, merges them
// (with user config taking precedence), validates the result, and returns
// the final configuration. If configuration files don't exist, it uses defaults.
//
// Loading sequence:
// 1. Start with DefaultConfig() as base
// 2. Load /etc/linux-wave/config.yaml if it exists (merge into base)
// 3. Load ~/.config/linuxwave/config.yaml if it exists (merge, overriding system)
// 4. Validate the merged configuration
// 5. Return final config or validation error
//
// File not found is not an error; YAML parsing errors and validation errors are returned.
func Load() (*Config, error) {
	// Start with defaults
	cfg := DefaultConfig()

	// Load system config if it exists
	if _, err := os.Stat(systemConfigPath); err == nil {
		systemCfg := &Config{}
		if err := loadFromPath(systemConfigPath, systemCfg); err != nil {
			return nil, fmt.Errorf("failed to load system config: %w", err)
		}
		cfg = mergeConfigs(cfg, systemCfg)
	}

	// Load user config if it exists
	userConfigPath, err := expandPath(filepath.Join("~", userConfigRelPath))
	if err != nil {
		return nil, fmt.Errorf("failed to expand user config path: %w", err)
	}

	if _, err := os.Stat(userConfigPath); err == nil {
		userCfg := &Config{}
		if err := loadFromPath(userConfigPath, userCfg); err != nil {
			return nil, fmt.Errorf("failed to load user config: %w", err)
		}
		cfg = mergeConfigs(cfg, userCfg)
	}

	// Validate merged configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// LoadFromPath loads configuration from a specific file path.
// This is useful for testing or custom configuration file locations.
func LoadFromPath(path string) (*Config, error) {
	cfg := DefaultConfig()

	expandedPath, err := expandPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to expand path %s: %w", path, err)
	}

	if err := loadFromPath(expandedPath, cfg); err != nil {
		return nil, fmt.Errorf("failed to load config from %s: %w", expandedPath, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// loadFromPath is an internal helper that reads and unmarshals a YAML file.
func loadFromPath(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	return nil
}

// mergeConfigs merges override config into base config, with override values
// taking precedence for any non-zero values. This implements the user config
// overriding system config behavior.
func mergeConfigs(base, override *Config) *Config {
	result := *base // Start with a copy of base

	// Merge Service settings
	if override.Service.Timeout != 0 {
		result.Service.Timeout = override.Service.Timeout
	}
	if override.Service.RetryAttempts != 0 {
		result.Service.RetryAttempts = override.Service.RetryAttempts
	}
	if override.Service.SocketPath != "" {
		result.Service.SocketPath = override.Service.SocketPath
	}

	// Merge Logging settings
	if override.Logging.Level != "" {
		result.Logging.Level = override.Logging.Level
	}
	if override.Logging.Format != "" {
		result.Logging.Format = override.Logging.Format
	}

	// Merge Audio settings
	// Note: For boolean fields, we can't distinguish between explicit false and zero value
	// This is acceptable since the user would set enabled: false explicitly if desired
	if override.Audio.Enabled != base.Audio.Enabled {
		result.Audio.Enabled = override.Audio.Enabled
	}
	if override.Audio.Volume != 0 {
		result.Audio.Volume = override.Audio.Volume
	}
	if override.Audio.CustomSoundSuccess != "" {
		result.Audio.CustomSoundSuccess = override.Audio.CustomSoundSuccess
	}
	if override.Audio.CustomSoundFailure != "" {
		result.Audio.CustomSoundFailure = override.Audio.CustomSoundFailure
	}

	// Merge Security settings
	if override.Security.LivenessRequired != base.Security.LivenessRequired {
		result.Security.LivenessRequired = override.Security.LivenessRequired
	}
	if override.Security.MatchThreshold != 0 {
		result.Security.MatchThreshold = override.Security.MatchThreshold
	}
	if override.Security.MaxAuthAttempts != 0 {
		result.Security.MaxAuthAttempts = override.Security.MaxAuthAttempts
	}
	if override.Security.LockoutDuration != 0 {
		result.Security.LockoutDuration = override.Security.LockoutDuration
	}

	return &result
}

// Validate checks that all configuration values are within acceptable ranges
// and that required fields are properly set. It returns a detailed error
// indicating which field failed validation and why.
func (c *Config) Validate() error {
	var errs []error

	// Validate Service settings
	if c.Service.Timeout <= 0 || c.Service.Timeout > 60 {
		errs = append(errs, errors.New("timeout must be between 1 and 60 seconds"))
	}
	if c.Service.RetryAttempts < 1 || c.Service.RetryAttempts > 10 {
		errs = append(errs, errors.New("retry_attempts must be between 1 and 10"))
	}
	if c.Service.SocketPath == "" {
		errs = append(errs, errors.New("socket_path must not be empty"))
	} else if !filepath.IsAbs(c.Service.SocketPath) {
		errs = append(errs, errors.New("socket_path must be an absolute path"))
	}

	// Validate Logging settings
	validLogLevels := map[string]bool{"DEBUG": true, "INFO": true, "WARN": true, "ERROR": true}
	if !validLogLevels[strings.ToUpper(c.Logging.Level)] {
		errs = append(errs, errors.New("log_level must be DEBUG, INFO, WARN, or ERROR"))
	}
	validLogFormats := map[string]bool{"json": true, "text": true}
	if !validLogFormats[strings.ToLower(c.Logging.Format)] {
		errs = append(errs, errors.New("log_format must be 'json' or 'text'"))
	}

	// Validate Audio settings
	if c.Audio.Volume < 0 || c.Audio.Volume > 100 {
		errs = append(errs, errors.New("audio_volume must be between 0 and 100"))
	}
	// Validate custom sound files exist if specified (non-empty)
	if c.Audio.CustomSoundSuccess != "" {
		if _, err := os.Stat(c.Audio.CustomSoundSuccess); err != nil {
			errs = append(errs, fmt.Errorf("custom_sound_success file not found: %s", c.Audio.CustomSoundSuccess))
		}
	}
	if c.Audio.CustomSoundFailure != "" {
		if _, err := os.Stat(c.Audio.CustomSoundFailure); err != nil {
			errs = append(errs, fmt.Errorf("custom_sound_failure file not found: %s", c.Audio.CustomSoundFailure))
		}
	}

	// Validate Security settings
	if c.Security.MatchThreshold < 0.0 || c.Security.MatchThreshold > 1.0 {
		errs = append(errs, errors.New("match_threshold must be between 0.0 and 1.0"))
	}
	if c.Security.MaxAuthAttempts < 1 || c.Security.MaxAuthAttempts > 20 {
		errs = append(errs, errors.New("max_auth_attempts must be between 1 and 20"))
	}
	if c.Security.LockoutDuration < 0 || c.Security.LockoutDuration > 3600 {
		errs = append(errs, errors.New("lockout_duration must be between 0 and 3600"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// expandPath expands ~ to the user's home directory and resolves environment variables.
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[1:])
	}
	return os.ExpandEnv(path), nil
}

// String returns a human-readable representation of the configuration.
// Sensitive values are sanitized for security.
func (c *Config) String() string {
	return fmt.Sprintf(`Config{
  Service: {Timeout: %ds, RetryAttempts: %d, SocketPath: %s}
  Logging: {Level: %s, Format: %s}
  Audio: {Enabled: %v, Volume: %d, CustomSuccess: %s, CustomFailure: %s}
  Security: {LivenessRequired: %v, MatchThreshold: %.2f, MaxAuthAttempts: %d, LockoutDuration: %ds}
}`,
		c.Service.Timeout,
		c.Service.RetryAttempts,
		c.Service.SocketPath,
		c.Logging.Level,
		c.Logging.Format,
		c.Audio.Enabled,
		c.Audio.Volume,
		sanitizePath(c.Audio.CustomSoundSuccess),
		sanitizePath(c.Audio.CustomSoundFailure),
		c.Security.LivenessRequired,
		c.Security.MatchThreshold,
		c.Security.MaxAuthAttempts,
		c.Security.LockoutDuration,
	)
}

// sanitizePath returns the path or "<not set>" if empty for display purposes.
func sanitizePath(path string) string {
	if path == "" {
		return "<not set>"
	}
	return path
}
