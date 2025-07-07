package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"crm-admin/internal/config"
	"crm-admin/internal/context"
	"crm-admin/internal/models"
)

const (
	RequestTimeout = 30 * time.Second
)

type Client struct {
	httpClient    *http.Client
	baseURL       string
	contextualURL string
	userContext   *context.UserContext
}

// New creates a new API client
func New() *Client {
	baseURL := config.GetBaseURL()
	contextualURL, userContext := context.GetContextualBaseURL(baseURL)

	return &Client{
		httpClient: &http.Client{
			Timeout: RequestTimeout,
		},
		baseURL:       baseURL,
		contextualURL: contextualURL,
		userContext:   userContext,
	}
}

// User operations - use correct existing endpoints
func (c *Client) CreateUser(username, password string) (*models.User, error) {
	userReq := models.UserRequest{
		Username: username,
		Password: password,
	}

	// Using the actual endpoint from your backend
	resp, err := c.postWithAuthRaw("/api/user", userReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Your backend returns the ID in the header, not in the response body
	userID := resp.Header.Get("id")
	if userID == "" {
		return nil, fmt.Errorf("user ID not returned in response header")
	}

	return &models.User{
		ID:       userID,
		Username: username,
	}, nil
}

func (c *Client) ListUsers() ([]models.User, error) {
	var users []models.User
	err := c.getWithAuth("/api/user", &users)
	return users, err
}

func (c *Client) GetUser(userID string) (*models.User, error) {
	var user models.User
	url := fmt.Sprintf("/api/user/%s", userID)
	err := c.getWithAuth(url, &user)
	return &user, err
}

// GetBaseURL returns the base URL for display purposes
func (c *Client) GetBaseURL() string {
	return c.baseURL
}

// Contact operations - use correct existing endpoints
func (c *Client) CreateContact(name, userID string, company, phoneNumber, contactEmail *string) (*models.Contact, error) {
	// Use provided userID or fall back to context
	targetUserID := userID
	if targetUserID == "" && c.userContext != nil {
		targetUserID = c.userContext.UserID
	}
	if targetUserID == "" {
		return nil, fmt.Errorf("user ID is required (use --user-id flag or select a user first)")
	}

	contactReq := models.ContactRequest{
		Name:         name,
		Company:      company,
		PhoneNumber:  phoneNumber,
		ContactEmail: contactEmail,
	}

	var contact models.Contact

	// Use contextual URL if we have context and no explicit userID was provided
	if userID == "" && c.userContext != nil {
		err := c.postWithAuth("/contacts", contactReq, &contact)
		return &contact, err
	} else {
		url := fmt.Sprintf("/api/users/%s/contacts", targetUserID)
		err := c.postWithAuth(url, contactReq, &contact)
		return &contact, err
	}
}

func (c *Client) ListContacts(userID string) ([]models.Contact, error) {
	// Use provided userID or fall back to context
	targetUserID := userID
	if targetUserID == "" && c.userContext != nil {
		targetUserID = c.userContext.UserID
	}
	if targetUserID == "" {
		return nil, fmt.Errorf("user ID is required (use --user-id flag or select a user first)")
	}

	var contacts []models.Contact

	// Use contextual URL if we have context and no explicit userID was provided
	if userID == "" && c.userContext != nil {
		err := c.getWithAuth("/contacts", &contacts)
		return contacts, err
	} else {
		url := fmt.Sprintf("/api/users/%s/contacts", targetUserID)
		err := c.getWithAuth(url, &contacts)
		return contacts, err
	}
}

func (c *Client) GetContact(userID string, contactID int) (*models.Contact, error) {
	var contact models.Contact
	url := fmt.Sprintf("/api/users/%s/contacts/%d", userID, contactID)
	err := c.getWithAuth(url, &contact)
	return &contact, err
}

// Note operations - use correct existing endpoints
func (c *Client) CreateNote(title, description string, contactIDs []int, userID string) (*models.Note, error) {
	// Use provided userID or fall back to context
	targetUserID := userID
	if targetUserID == "" && c.userContext != nil {
		targetUserID = c.userContext.UserID
	}
	if targetUserID == "" {
		return nil, fmt.Errorf("user ID is required (use --user-id flag or select a user first)")
	}

	noteReq := models.NoteRequest{
		ContactIDs:  contactIDs,
		Title:       title,
		Description: description,
	}

	var note models.Note

	// Use contextual URL if we have context and no explicit userID was provided
	if userID == "" && c.userContext != nil {
		err := c.postWithAuth("/contacts/notes", noteReq, &note)
		return &note, err
	} else {
		url := fmt.Sprintf("/api/users/%s/contacts/notes", targetUserID)
		err := c.postWithAuth(url, noteReq, &note)
		return &note, err
	}
}

func (c *Client) ListNotesForUser(userID string) ([]models.Note, error) {
	// Use provided userID or fall back to context
	targetUserID := userID
	if targetUserID == "" && c.userContext != nil {
		targetUserID = c.userContext.UserID
	}
	if targetUserID == "" {
		return nil, fmt.Errorf("user ID is required (use --user-id flag or select a user first)")
	}

	var notes []models.Note

	// Use contextual URL if we have context and no explicit userID was provided
	if userID == "" && c.userContext != nil {
		err := c.getWithAuth("/contacts/notes", &notes)
		return notes, err
	} else {
		url := fmt.Sprintf("/api/users/%s/contacts/notes", targetUserID)
		err := c.getWithAuth(url, &notes)
		return notes, err
	}
}

func (c *Client) ListNotesForContact(userID string, contactID int) ([]models.Note, error) {
	// Use provided userID or fall back to context
	targetUserID := userID
	if targetUserID == "" && c.userContext != nil {
		targetUserID = c.userContext.UserID
	}
	if targetUserID == "" {
		return nil, fmt.Errorf("user ID is required (use --user-id flag or select a user first)")
	}

	var notes []models.Note

	// Use contextual URL if we have context and no explicit userID was provided
	if userID == "" && c.userContext != nil {
		url := fmt.Sprintf("/contacts/%d/notes", contactID)
		err := c.getWithAuth(url, &notes)
		return notes, err
	} else {
		url := fmt.Sprintf("/api/users/%s/contacts/%d/notes", targetUserID, contactID)
		err := c.getWithAuth(url, &notes)
		return notes, err
	}
}

func (c *Client) GetNote(userID string, noteID int) (*models.Note, error) {
	// Use provided userID or fall back to context
	targetUserID := userID
	if targetUserID == "" && c.userContext != nil {
		targetUserID = c.userContext.UserID
	}
	if targetUserID == "" {
		return nil, fmt.Errorf("user ID is required (use --user-id flag or select a user first)")
	}

	var note models.Note

	// Use contextual URL if we have context and no explicit userID was provided
	if userID == "" && c.userContext != nil {
		url := fmt.Sprintf("/contacts/notes/%d", noteID)
		err := c.getWithAuth(url, &note)
		return &note, err
	} else {
		url := fmt.Sprintf("/api/users/%s/contacts/notes/%d", targetUserID, noteID)
		err := c.getWithAuth(url, &note)
		return &note, err
	}
}

func (c *Client) UpdateNote(userID string, noteID int, title, description string, contactIDs []int) (*models.Note, error) {
	// Use provided userID or fall back to context
	targetUserID := userID
	if targetUserID == "" && c.userContext != nil {
		targetUserID = c.userContext.UserID
	}
	if targetUserID == "" {
		return nil, fmt.Errorf("user ID is required (use --user-id flag or select a user first)")
	}

	noteReq := models.NoteRequest{
		ContactIDs:  contactIDs,
		Title:       title,
		Description: description,
	}

	var note models.Note

	// Use contextual URL if we have context and no explicit userID was provided
	if userID == "" && c.userContext != nil {
		url := fmt.Sprintf("/contacts/notes/%d", noteID)
		err := c.putWithAuth(url, noteReq, &note)
		return &note, err
	} else {
		url := fmt.Sprintf("/api/users/%s/contacts/notes/%d", targetUserID, noteID)
		err := c.putWithAuth(url, noteReq, &note)
		return &note, err
	}
}

func (c *Client) DeleteNote(userID string, noteID int) error {
	// Use provided userID or fall back to context
	targetUserID := userID
	if targetUserID == "" && c.userContext != nil {
		targetUserID = c.userContext.UserID
	}
	if targetUserID == "" {
		return fmt.Errorf("user ID is required (use --user-id flag or select a user first)")
	}

	// Use contextual URL if we have context and no explicit userID was provided
	if userID == "" && c.userContext != nil {
		url := fmt.Sprintf("/contacts/notes/%d", noteID)
		return c.deleteWithAuth(url)
	} else {
		url := fmt.Sprintf("/api/users/%s/contacts/notes/%d", targetUserID, noteID)
		return c.deleteWithAuth(url)
	}
}

// Helper methods with admin JWT authentication
func (c *Client) postWithAuth(endpoint string, data interface{}, result interface{}) error {
	resp, err := c.postWithAuthRaw(endpoint, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) putWithAuth(endpoint string, data interface{}, result interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Choose URL based on context
	var fullURL string
	if c.userContext != nil && !isAbsoluteEndpoint(endpoint) {
		fullURL = c.contextualURL + endpoint
	} else {
		fullURL = c.baseURL + endpoint
	}

	req, err := http.NewRequest("PUT", fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.GetAdminToken())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%s): %s", resp.Status, string(body))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) deleteWithAuth(endpoint string) error {
	// Choose URL based on context
	var fullURL string
	if c.userContext != nil && !isAbsoluteEndpoint(endpoint) {
		fullURL = c.contextualURL + endpoint
	} else {
		fullURL = c.baseURL + endpoint
	}

	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.GetAdminToken())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%s): %s", resp.Status, string(body))
	}

	return nil
}

func (c *Client) postWithAuthRaw(endpoint string, data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	// Choose URL based on context
	var fullURL string
	if c.userContext != nil && !isAbsoluteEndpoint(endpoint) {
		fullURL = c.contextualURL + endpoint
	} else {
		fullURL = c.baseURL + endpoint
	}

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.GetAdminToken())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error (%s): %s", resp.Status, string(body))
	}

	return resp, nil
}

func (c *Client) getWithAuth(endpoint string, result interface{}) error {
	// Choose URL based on context
	var fullURL string
	if c.userContext != nil && !isAbsoluteEndpoint(endpoint) {
		fullURL = c.contextualURL + endpoint
	} else {
		fullURL = c.baseURL + endpoint
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.GetAdminToken())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Accept both 200 OK and 302 Found (temporary fix for backend issue)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%s): %s", resp.Status, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// isAbsoluteEndpoint checks if the endpoint starts with /api (absolute path)
func isAbsoluteEndpoint(endpoint string) bool {
	return len(endpoint) > 4 && endpoint[:4] == "/api"
}
