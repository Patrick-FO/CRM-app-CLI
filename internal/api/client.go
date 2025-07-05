package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"crm-admin/internal/config"
	"crm-admin/internal/models"
)

const (
	RequestTimeout = 30 * time.Second
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

// New creates a new API client
func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: RequestTimeout,
		},
		baseURL: config.GetBaseURL(),
	}
}

// User operations
func (c *Client) CreateUser(username, password string) (*models.User, error) {
	userReq := models.UserRequest{
		Username: username,
		Password: password,
	}

	var user models.User
	err := c.postAdmin("/user", userReq, &user)
	return &user, err
}

func (c *Client) ListUsers() ([]models.User, error) {
	var users []models.User
	err := c.getAdmin("/user", &users)
	return users, err
}

// Contact operations
func (c *Client) CreateContact(name, userID string, company, phoneNumber, contactEmail *string) (*models.Contact, error) {
	contactReq := models.AdminContactRequest{
		Name:         name,
		UserID:       userID,
		Company:      company,
		PhoneNumber:  phoneNumber,
		ContactEmail: contactEmail,
	}

	var contact models.Contact
	err := c.postAdmin("/admin/contacts", contactReq, &contact)
	return &contact, err
}

func (c *Client) ListContacts(userID string) ([]models.Contact, error) {
	var contacts []models.Contact
	url := "/admin/contacts"
	if userID != "" {
		url += fmt.Sprintf("?userId=%s", userID)
	}
	err := c.getAdmin(url, &contacts)
	return contacts, err
}

// Note operations
func (c *Client) CreateNote(title, description string, contactIDs []int, userID string) (*models.Note, error) {
	noteReq := models.AdminNoteRequest{
		ContactIDs:  contactIDs,
		Title:       title,
		Description: description,
		UserID:      userID,
	}

	var note models.Note
	err := c.postAdmin("/admin/notes", noteReq, &note)
	return &note, err
}

func (c *Client) ListNotes(contactIDs []int) ([]models.Note, error) {
	var notes []models.Note
	url := "/admin/notes"
	if len(contactIDs) > 0 {
		contactIDStrings := make([]string, len(contactIDs))
		for i, id := range contactIDs {
			contactIDStrings[i] = strconv.Itoa(id)
		}
		url += fmt.Sprintf("?contactIds=%s", strings.Join(contactIDStrings, ","))
	}
	err := c.getAdmin(url, &notes)
	return notes, err
}

// Helper methods with admin JWT authentication
func (c *Client) postAdmin(endpoint string, data interface{}, result interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewBuffer(jsonData))
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

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
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

func (c *Client) getAdmin(endpoint string, result interface{}) error {
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
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

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
