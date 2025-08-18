package config

import "os"

// envOrDefault returns env value if set, otherwise fallback.
func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
