package envconfig

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func Load() error {
	return godotenv.Load(Get("ENV_FILE", ".env"))
}

func ToList(value string, separator string) ([]string, error) {
	if strings.HasPrefix(value, "@") {
		if filedata_, err := os.ReadFile(value[1:]); err != nil {
			return nil, err
		} else {
			filedata := strings.Trim(
				strings.ReplaceAll(string(filedata_), "\r", ""),
				"\n")

			var result []string
			for _, line := range strings.Split(filedata, separator) {
				if !strings.HasPrefix(strings.Trim(line, " "), "#") {
					result = append(result, line)
				}
			}
			return result, nil
		}
	} else {
		return strings.Split(
			strings.ReplaceAll(value, "==", "\n"),
			";"), nil
	}
}

func Get(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
func GetBool(key string, defaultValue bool) bool {
	value := Get(key, strconv.FormatBool(defaultValue))
	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}

	return defaultValue
}

func GetInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func GetInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}
	return defaultValue
}
