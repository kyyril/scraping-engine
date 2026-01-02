package utils

import (
	"distributed-scraper/internal/models"
	"fmt"
	"net/url"
	"strings"
)

// ValidateURL checks if the provided URL is valid and safe
func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("only HTTP and HTTPS URLs are allowed")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a valid host")
	}

	// Block localhost and private IPs for security
	if strings.Contains(parsedURL.Host, "localhost") ||
		strings.Contains(parsedURL.Host, "127.0.0.1") ||
		strings.Contains(parsedURL.Host, "0.0.0.0") {
		return fmt.Errorf("localhost URLs are not allowed")
	}

	return nil
}

// ValidateActionType checks if the action type is supported
func ValidateActionType(actionType string) error {
	validTypes := map[string]bool{
		string(models.ActionNavigate):   true,
		string(models.ActionClick):      true,
		string(models.ActionType_):      true,
		string(models.ActionWait):       true,
		string(models.ActionScreenshot): true,
		string(models.ActionExtract):    true,
		string(models.ActionScroll):     true,
	}

	if !validTypes[actionType] {
		return fmt.Errorf("unsupported action type: %s", actionType)
	}

	return nil
}

// ValidateSelector checks if CSS selector is not empty for actions that require it
func ValidateSelector(actionType, selector string) error {
	requiresSelector := map[string]bool{
		string(models.ActionClick):   true,
		string(models.ActionType_):   true,
		string(models.ActionExtract): true,
	}

	if requiresSelector[actionType] && selector == "" {
		return fmt.Errorf("action type %s requires a target selector", actionType)
	}

	return nil
}