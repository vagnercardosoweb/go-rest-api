package utils

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	// ErrInvalidFilePath indicates that the file path is invalid or unsafe
	ErrInvalidFilePath = errors.New("invalid or unsafe file path")

	// ErrPathTraversal indicates a path traversal attempt
	ErrPathTraversal = errors.New("path traversal attempt detected")
)

// ValidateSourceFilePath validates source file paths
// Used for stack traces and debugging - only allows valid Go files
func ValidateSourceFilePath(path string) error {
	if path == "" {
		return ErrInvalidFilePath
	}

	// Clean and normalize the path
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return ErrPathTraversal
	}

	// For stack traces, only allow Go files from the runtime and project
	if !strings.HasSuffix(cleanPath, ".go") {
		return ErrInvalidFilePath
	}

	// Check if it's a valid Go runtime or project absolute path
	if !isValidGoSourcePath(cleanPath) {
		return ErrInvalidFilePath
	}

	return nil
}

// ValidateDownloadFilePath validates paths for file downloads
// Ensures the file will only be saved in allowed directories
func ValidateDownloadFilePath(path string) error {
	if path == "" {
		return ErrInvalidFilePath
	}

	// Check for path traversal attempts before cleaning
	if strings.Contains(path, "..") {
		return ErrPathTraversal
	}

	// Clean and normalize the path
	cleanPath := filepath.Clean(path)

	// Check if the path is within allowed directories
	allowedDirs := []string{
		"/tmp/",
		"tmp/",
		"./tmp/",
		"downloads/",
		"./downloads/",
		"/var/tmp/",
	}

	isAllowed := false
	for _, allowedDir := range allowedDirs {
		if strings.HasPrefix(cleanPath, allowedDir) ||
			strings.HasPrefix(path, allowedDir) {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return ErrInvalidFilePath
	}

	return nil
}

// isValidGoSourcePath checks if the path is a valid Go file
func isValidGoSourcePath(path string) bool {
	// Allow project files
	if strings.Contains(path, "go-rest-api") {
		return true
	}

	// Allow Go runtime files
	if strings.Contains(path, "/go/src/") ||
		strings.Contains(path, "/usr/local/go/src/") ||
		strings.Contains(path, "GOROOT") {
		return true
	}

	// Allow files in GOPATH (Go modules)
	if strings.Contains(path, "/go/pkg/mod/") {
		return true
	}

	return false
}

// SanitizeFileName removes dangerous characters from file names
func SanitizeFileName(filename string) string {
	dangerous := []string{"..", "/", "\\", ":", "*", "?", "\"", "<", ">", "|"}

	sanitized := filename
	for _, char := range dangerous {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}

	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
	}

	return sanitized
}
