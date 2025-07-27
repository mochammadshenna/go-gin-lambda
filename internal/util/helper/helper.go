package helper

import (
	"ai-service/internal/util/logger"
	"context"
	"strings"
)

// PanicOnError panics if error is not nil
func PanicOnError(err error) {
	if err != nil {
		logger.Logger.Error(err)
		panic(err)
	}
}

// PanicOnErrorContext panics if error is not nil with context
func PanicOnErrorContext(ctx context.Context, err error) {
	if err != nil {
		logger.Error(ctx, err)
		panic(err)
	}
}

// IsEmpty checks if a string is empty
func IsEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

// IsNotEmpty checks if a string is not empty
func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

// Contains checks if a slice contains a specific value
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicate values from a slice
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
