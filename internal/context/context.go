package context

import (
	"encoding/json"
	"fmt"
	"os"

	"crm-admin/internal/models"
)

const contextFileName = ".crm-context.json"

type UserContext struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

// getContextFilePath returns the path to the context file in the current directory
func getContextFilePath() string {
	return contextFileName
}

// SaveUserContext saves the selected user context to a file
func SaveUserContext(user *models.User) error {
	context := UserContext{
		UserID:   user.ID,
		Username: user.Username,
	}

	data, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal context: %w", err)
	}

	contextPath := getContextFilePath()
	err = os.WriteFile(contextPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write context file: %w", err)
	}

	return nil
}

// LoadUserContext loads the selected user context from file
func LoadUserContext() (*UserContext, error) {
	contextPath := getContextFilePath()

	if _, err := os.Stat(contextPath); os.IsNotExist(err) {
		return nil, nil // No context file exists
	}

	data, err := os.ReadFile(contextPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read context file: %w", err)
	}

	var context UserContext
	err = json.Unmarshal(data, &context)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal context: %w", err)
	}

	return &context, nil
}

// ClearUserContext removes the context file
func ClearUserContext() error {
	contextPath := getContextFilePath()

	if _, err := os.Stat(contextPath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to clear
	}

	err := os.Remove(contextPath)
	if err != nil {
		return fmt.Errorf("failed to remove context file: %w", err)
	}

	return nil
}

// HasUserContext checks if a user context exists
func HasUserContext() bool {
	contextPath := getContextFilePath()
	_, err := os.Stat(contextPath)
	return err == nil
}

// GetContextualBaseURL returns the base URL with user context if available
func GetContextualBaseURL(baseURL string) (string, *UserContext) {
	context, err := LoadUserContext()
	if err != nil || context == nil {
		return baseURL, nil
	}

	contextualURL := fmt.Sprintf("%s/api/users/%s", baseURL, context.UserID)
	return contextualURL, context
}
