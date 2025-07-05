package models

// User represents a user in the CRM system
type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

// Contact represents a contact belonging to a user
type Contact struct {
	ID           int     `json:"id,omitempty"`
	UserID       string  `json:"userId"`
	Name         string  `json:"name"`
	Company      *string `json:"company,omitempty"`
	PhoneNumber  *string `json:"phoneNumber,omitempty"`
	ContactEmail *string `json:"contactEmail,omitempty"`
}

// Note represents a note that can be linked to multiple contacts
type Note struct {
	ID          int     `json:"id,omitempty"`
	UserID      string  `json:"userId"`
	ContactIDs  []int   `json:"contactIds"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
}

// Request models for creating resources
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ContactRequest struct {
	Name         string  `json:"name"`
	Company      *string `json:"company,omitempty"`
	PhoneNumber  *string `json:"phoneNumber,omitempty"`
	ContactEmail *string `json:"contactEmail,omitempty"`
}

type NoteRequest struct {
	ContactIDs  []int  `json:"contactIds"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Admin request models for the CLI
type AdminContactRequest struct {
	Name         string  `json:"name"`
	UserID       string  `json:"userId"`
	Company      *string `json:"company,omitempty"`
	PhoneNumber  *string `json:"phoneNumber,omitempty"`
	ContactEmail *string `json:"contactEmail,omitempty"`
}

type AdminNoteRequest struct {
	ContactIDs  []int  `json:"contactIds"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      string `json:"userId"`
}
