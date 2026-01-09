package envconfig

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	// DefaultEnvFile is the default path to the .env file
	DefaultEnvFile = ".env"
	// EnvFileKey is the environment variable key for custom .env file path
	EnvFileKey = "ENV_FILE"
)

// Load loads environment variables from a .env file.
// By default, it looks for a file named ".env" in the current directory.
// You can specify a custom file path by setting the ENV_FILE environment variable.
func Load() error {
	envFile := Get(EnvFileKey, DefaultEnvFile)
	return godotenv.Load(envFile)
}

// LoadStruct loads configuration from environment variables into a struct.
// It uses the "env" tag to specify the environment variable name
// and the "default" tag to specify a default value.
//
// Supported field types: string, bool, int, int64, []int, []int64, and arrays of int.
//
// Example:
//
//	type Config struct {
//	    Host string `env:"HOST" default:"localhost"`
//	    Port int    `env:"PORT" default:"8080"`
//	}
//
//	var cfg Config
//	if err := envconfig.LoadStruct(&cfg); err != nil {
//	    log.Fatal(err)
//	}
func LoadStruct(cfg any) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cfg must be pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		envName := fieldType.Tag.Get("env")
		if envName == "" {
			continue
		}

		envValue := getEnvValue(envName, fieldType.Tag.Get("default"))
		if err := setValue(field, envValue); err != nil {
			return fmt.Errorf("env %s: %w", envName, err)
		}
	}

	return nil
}

// getEnvValue retrieves the environment variable value or returns the default.
func getEnvValue(envName, defaultValue string) string {
	if value, exists := os.LookupEnv(envName); exists {
		return value
	}
	return defaultValue
}

// ToList splits a string into a list of strings by the specified separator.
func ToList(value string, separator string) ([]string, error) {
	return strings.Split(value, separator), nil
}

// Get retrieves a string value from environment variables with a default value.
// Returns the default value if the environment variable is not set or is empty.
func Get(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetBool retrieves a boolean value from environment variables with a default value.
// Returns the default value if the environment variable is not set, is empty, or cannot be parsed.
//
// Supported values: true, false, 1, 0, t, f, T, F, TRUE, FALSE, True, False
func GetBool(key string, defaultValue bool) bool {
	value := Get(key, "")
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

// GetInt retrieves an integer value from environment variables with a default value.
// Returns the default value if the environment variable is not set, is empty, or cannot be parsed.
func GetInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

// GetInt64 retrieves a 64-bit integer value from environment variables with a default value.
// Returns the default value if the environment variable is not set, is empty, or cannot be parsed.
func GetInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return parsed
}

// GetIntSlice retrieves a slice of integers from environment variables with a default value.
// Values should be comma-separated. Spaces around values are automatically trimmed.
// Returns the default value if the environment variable is not set, is empty, or contains invalid values.
//
// Example:
//
//	PORTS=8080,8081,8082
//	ports := GetIntSlice("PORTS", []int{3000, 3001})
func GetIntSlice(key string, defaultValue []int) []int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	slice, err := parseIntSlice(value)
	if err != nil {
		return defaultValue
	}

	return slice
}

// GetInt64Slice retrieves a slice of 64-bit integers from environment variables with a default value.
// Values should be comma-separated. Spaces around values are automatically trimmed.
// Returns the default value if the environment variable is not set, is empty, or contains invalid values.
//
// Example:
//
//	MAX_SIZES=1024,2048,4096
//	maxSizes := GetInt64Slice("MAX_SIZES", []int64{512, 1024})
func GetInt64Slice(key string, defaultValue []int64) []int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	slice, err := parseInt64Slice(value)
	if err != nil {
		return defaultValue
	}

	return slice
}

// parseIntSlice parses a comma-separated string into a slice of integers.
func parseIntSlice(value string) ([]int, error) {
	parts := strings.Split(value, ",")
	result := make([]int, 0, len(parts))

	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			result = append(result, 0)
			continue
		}

		parsed, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid int value at index %d: %w", i, err)
		}

		result = append(result, parsed)
	}

	return result, nil
}

// parseInt64Slice parses a comma-separated string into a slice of 64-bit integers.
func parseInt64Slice(value string) ([]int64, error) {
	parts := strings.Split(value, ",")
	result := make([]int64, 0, len(parts))

	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			result = append(result, 0)
			continue
		}

		parsed, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid int64 value at index %d: %w", i, err)
		}

		result = append(result, parsed)
	}

	return result, nil
}
