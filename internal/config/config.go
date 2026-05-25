package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL    string `mapstructure:"DATABASE_URL"`
	SendGridAPIKey string `mapstructure:"SENDGRID_API_KEY"`
	SendGridFrom   string `mapstructure:"SENDGRID_FROM"`
	ApproverEmail  string `mapstructure:"APPROVER_EMAIL"`
	GeminiAPIKey   string `mapstructure:"GEMINI_API_KEY"`
	Port           string `mapstructure:"PORT"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error reading config file: %w", err)
	}

	// Set environment variables
	v.AutomaticEnv()

	// Set default values
	v.SetDefault("PORT", "8080")

	// Validate required fields
	requiredFields := []string{"DATABASE_URL", "SENDGRID_API_KEY", "SENDGRID_FROM", "APPROVER_EMAIL", "GEMINI_API_KEY"}
	for _, field := range requiredFields {
		if !v.IsSet(field) {
			return nil, fmt.Errorf("missing required environment variable: %s", field)
		}
	}

	// Print loaded configuration for debugging
	fmt.Printf("Loaded configuration\n")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg,nil
}
