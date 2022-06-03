package config

import (
	"fmt"
	"os"
)

// Get get environment variable
func Get(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists && value != "" {
		fmt.Println("Environment variable not found!", key)
		os.Exit(1)
	}
	return value
}

// Set set environment variable
func Set(key string, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		fmt.Println("Environment variable not set!")
		os.Exit(1)
	}
}

//UnSet unset environment variable
func UnSet(key string) {
	err := os.Unsetenv(key)
	if err != nil {
		fmt.Println("Environment variable not reset!")
		os.Exit(1)
	}
}
