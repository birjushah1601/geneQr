package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// Environment types
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"

	// Default values
	DefaultPort             = "8081"
	DefaultShutdownTimeout  = 30
	DefaultEnvironment      = EnvDevelopment
	DefaultLogLevel         = "info"
	DefaultTracingSampling  = 0.1
	DefaultDbMaxConnections = 10
	DefaultDbMaxIdleTime    = 5 * time.Minute
	DefaultDbTimeout        = 10 * time.Second
	DefaultKafkaBroker      = "kafka:9092"
	DefaultRedisAddr        = "redis:6379"
	DefaultKeycloakURL      = "http://keycloak:8080"
)

// Config holds all application configuration
type Config struct {
	// Core application settings
	Environment        string
	Version            string
	Port               string
	ShutdownTimeoutSec int
	EnabledModulesString string

	// CORS configuration
	CORS struct {
		AllowedOrigins []string
	}

	// Database configuration
	Database struct {
		Host         string
		Port         string
		User         string
		Password     string
		Name         string
		MaxConns     int
		MaxIdleTime  time.Duration
		QueryTimeout time.Duration
		SSLMode      string
	}

	// Redis configuration
	Redis struct {
		Addr     string
		Password string
		DB       int
	}

	// Kafka configuration
	Kafka struct {
		Brokers []string
		Topic   string
		GroupID string
	}

	// Keycloak configuration
	Keycloak struct {
		URL          string
		Realm        string
		ClientID     string
		ClientSecret string
	}

	// Observability configuration
	Observability struct {
		LogLevel       string
		TracingEnabled bool
		TracingURL     string
		SamplingRate   float64
		MetricsEnabled bool
	}

	// AI Services configuration
	AI struct {
		// Provider settings
		Provider         string
		FallbackProvider string

		// API Keys
		OpenAIAPIKey    string
		AnthropicAPIKey string

		// Model configuration
		OpenAIModel    string
		AnthropicModel string

		// Behavior
		MaxRetries     int
		TimeoutSeconds int
		Temperature    float64
		MaxTokens      int

		// Features
		CostTrackingEnabled bool

		// Feedback Learning
		FeedbackPatternThreshold  int
		FeedbackTestPeriodDays    int
		FeedbackDeployThreshold   int
		FeedbackRollbackThreshold int
	}
}

// Load reads configuration from environment variables with sensible defaults
func Load() (*Config, error) {
	cfg := &Config{}

	// Core application settings
	cfg.Environment = getEnv("APP_ENV", DefaultEnvironment)
	cfg.Version = getEnv("APP_VERSION", "dev")
	cfg.Port = getEnv("PORT", DefaultPort)
	cfg.ShutdownTimeoutSec = getEnvAsInt("SHUTDOWN_TIMEOUT_SEC", DefaultShutdownTimeout)
	cfg.EnabledModulesString = getEnv("ENABLED_MODULES", "*")

	// CORS configuration
	cfg.CORS.AllowedOrigins = getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"})

	// Database configuration
	cfg.Database.Host = getEnv("DB_HOST", "postgres")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "postgres")
	cfg.Database.Name = getEnv("DB_NAME", "medplatform")
	cfg.Database.MaxConns = getEnvAsInt("DB_MAX_CONNS", DefaultDbMaxConnections)
	cfg.Database.MaxIdleTime = getEnvAsDuration("DB_MAX_IDLE_TIME", DefaultDbMaxIdleTime)
	cfg.Database.QueryTimeout = getEnvAsDuration("DB_QUERY_TIMEOUT", DefaultDbTimeout)
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// Redis configuration
	cfg.Redis.Addr = getEnv("REDIS_ADDR", DefaultRedisAddr)
	cfg.Redis.Password = getEnv("REDIS_PASSWORD", "")
	cfg.Redis.DB = getEnvAsInt("REDIS_DB", 0)

	// Kafka configuration
	cfg.Kafka.Brokers = getEnvAsSlice("KAFKA_BROKERS", []string{DefaultKafkaBroker})
	cfg.Kafka.Topic = getEnv("KAFKA_TOPIC", "medical-platform")
	cfg.Kafka.GroupID = getEnv("KAFKA_GROUP_ID", "medical-platform-group")

	// Keycloak configuration
	cfg.Keycloak.URL = getEnv("KEYCLOAK_URL", DefaultKeycloakURL)
	cfg.Keycloak.Realm = getEnv("KEYCLOAK_REALM", "master")
	cfg.Keycloak.ClientID = getEnv("KEYCLOAK_CLIENT_ID", "api-gateway")
	cfg.Keycloak.ClientSecret = getEnv("KEYCLOAK_CLIENT_SECRET", "")

	// Observability configuration
	cfg.Observability.LogLevel = getEnv("LOG_LEVEL", DefaultLogLevel)
	cfg.Observability.TracingEnabled = getEnvAsBool("TRACING_ENABLED", true)
	cfg.Observability.TracingURL = getEnv("TRACING_URL", "http://otel-collector:4317")
	cfg.Observability.SamplingRate = getEnvAsFloat("TRACING_SAMPLING_RATE", DefaultTracingSampling)
	cfg.Observability.MetricsEnabled = getEnvAsBool("METRICS_ENABLED", true)

	// AI Services configuration
	cfg.AI.Provider = getEnv("AI_PROVIDER", "openai")
	cfg.AI.FallbackProvider = getEnv("AI_FALLBACK_PROVIDER", "anthropic")
	cfg.AI.OpenAIAPIKey = getEnv("AI_OPENAI_API_KEY", "")
	cfg.AI.AnthropicAPIKey = getEnv("AI_ANTHROPIC_API_KEY", "")
	cfg.AI.OpenAIModel = getEnv("AI_OPENAI_MODEL", "gpt-4")
	cfg.AI.AnthropicModel = getEnv("AI_ANTHROPIC_MODEL", "claude-3-opus-20240229")
	cfg.AI.MaxRetries = getEnvAsInt("AI_MAX_RETRIES", 3)
	cfg.AI.TimeoutSeconds = getEnvAsInt("AI_TIMEOUT_SECONDS", 30)
	cfg.AI.Temperature = getEnvAsFloat("AI_TEMPERATURE", 0.7)
	cfg.AI.MaxTokens = getEnvAsInt("AI_MAX_TOKENS", 2000)
	cfg.AI.CostTrackingEnabled = getEnvAsBool("AI_COST_TRACKING_ENABLED", true)
	cfg.AI.FeedbackPatternThreshold = getEnvAsInt("AI_FEEDBACK_PATTERN_THRESHOLD", 3)
	cfg.AI.FeedbackTestPeriodDays = getEnvAsInt("AI_FEEDBACK_TEST_PERIOD_DAYS", 7)
	cfg.AI.FeedbackDeployThreshold = getEnvAsInt("AI_FEEDBACK_DEPLOY_THRESHOLD", 5)
	cfg.AI.FeedbackRollbackThreshold = getEnvAsInt("AI_FEEDBACK_ROLLBACK_THRESHOLD", -5)

	// Validate critical configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate environment
	if c.Environment != EnvDevelopment && c.Environment != EnvStaging && c.Environment != EnvProduction {
		return fmt.Errorf("invalid environment: %s", c.Environment)
	}

	// Validate port
	if _, err := strconv.Atoi(c.Port); err != nil {
		return fmt.Errorf("invalid port: %s", c.Port)
	}

	// In production, ensure Keycloak client secret is set
	if c.Environment == EnvProduction && c.Keycloak.ClientSecret == "" {
		return fmt.Errorf("keycloak client secret must be set in production")
	}

	// Validate database configuration in production
	if c.Environment == EnvProduction {
		if c.Database.Password == "postgres" {
			return fmt.Errorf("default database password cannot be used in production")
		}
		if c.Database.SSLMode == "disable" {
			return fmt.Errorf("SSL must be enabled for database in production")
		}
	}

	return nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == EnvDevelopment
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == EnvProduction
}

// Helper functions to get environment variables with defaults

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
