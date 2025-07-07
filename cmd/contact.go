package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"crm-admin/internal/api"
	"crm-admin/internal/context"
)

var contactCmd = &cobra.Command{
	Use:   "contact",
	Short: "Manage contacts",
	Long:  `Create, list, and manage contacts in the CRM system.`,
}

var contactCreateCmd = &cobra.Command{
	Use:   "create [n]",
	Short: "Create a new contact",
	Long:  `Create a new contact with the specified name for a user.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, _ := cmd.Flags().GetString("user-id")

		// Check if user-id is provided or if we have context
		if userID == "" && !context.HasUserContext() {
			return fmt.Errorf("user-id flag is required (or select a user with 'crm-admin user select [user-id]')")
		}

		// Optional fields
		company, _ := cmd.Flags().GetString("company")
		phoneNumber, _ := cmd.Flags().GetString("phone")
		contactEmail, _ := cmd.Flags().GetString("email")

		client := api.New()

		// Convert empty strings to nil pointers
		var companyPtr, phonePtr, emailPtr *string
		if company != "" {
			companyPtr = &company
		}
		if phoneNumber != "" {
			phonePtr = &phoneNumber
		}
		if contactEmail != "" {
			emailPtr = &contactEmail
		}

		contact, err := client.CreateContact(args[0], userID, companyPtr, phonePtr, emailPtr)
		if err != nil {
			return fmt.Errorf("failed to create contact: %w", err)
		}

		fmt.Printf("✅ Contact '%s' created successfully!\n", contact.Name)
		fmt.Printf("   ID: %d\n", contact.ID)
		fmt.Printf("   User ID: %s\n", contact.UserID)
		if contact.Company != nil {
			fmt.Printf("   Company: %s\n", *contact.Company)
		}
		if contact.PhoneNumber != nil {
			fmt.Printf("   Phone: %s\n", *contact.PhoneNumber)
		}
		if contact.ContactEmail != nil {
			fmt.Printf("   Email: %s\n", *contact.ContactEmail)
		}

		fmt.Printf("\nYou can now create notes for this contact:\n")
		if context.HasUserContext() {
			fmt.Printf("   crm-admin note create \"Note Title\" \"Description\" --contact-ids %d\n", contact.ID)
		} else {
			fmt.Printf("   crm-admin note create \"Note Title\" \"Description\" --contact-ids %d --user-id %s\n", contact.ID, contact.UserID)
		}
		return nil
	},
}

var contactListCmd = &cobra.Command{
	Use:   "list",
	Short: "List contacts",
	Long:  `Display a list of contacts for a specific user.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, _ := cmd.Flags().GetString("user-id")

		// Check if user-id is provided or if we have context
		if userID == "" && !context.HasUserContext() {
			return fmt.Errorf("user-id flag is required (or select a user with 'crm-admin user select [user-id]')")
		}

		client := api.New()

		contacts, err := client.ListContacts(userID)
		if err != nil {
			return fmt.Errorf("failed to list contacts: %w", err)
		}

		// Show which user we're listing for
		targetUserID := userID
		username := ""
		if targetUserID == "" && context.HasUserContext() {
			userContext, _ := context.LoadUserContext()
			if userContext != nil {
				targetUserID = userContext.UserID
				username = userContext.Username
			}
		}

		if len(contacts) == 0 {
			if username != "" {
				fmt.Printf("No contacts found for %s (%s).\n", username, targetUserID)
			} else {
				fmt.Printf("No contacts found for user ID %s.\n", targetUserID)
			}
			return nil
		}

		if username != "" {
			fmt.Printf("📋 Contacts for %s (%s):\n", username, targetUserID)
		} else {
			fmt.Printf("📋 Contacts for User %s:\n", targetUserID)
		}
		fmt.Printf("%-5s | %-20s | %-15s | %-15s | %s\n",
			"ID", "Name", "Company", "Phone", "Email")
		fmt.Printf("%-5s | %-20s | %-15s | %-15s | %s\n",
			"-----", "--------------------", "---------------", "---------------", "---------------")

		for _, contact := range contacts {
			company := ""
			if contact.Company != nil {
				company = *contact.Company
			}
			phone := ""
			if contact.PhoneNumber != nil {
				phone = *contact.PhoneNumber
			}
			email := ""
			if contact.ContactEmail != nil {
				email = *contact.ContactEmail
			}

			fmt.Printf("%-5d | %-20s | %-15s | %-15s | %s\n",
				contact.ID, contact.Name, company, phone, email)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(contactCmd)
	contactCmd.AddCommand(contactCreateCmd)
	contactCmd.AddCommand(contactListCmd)

	// Flags for contact create
	contactCreateCmd.Flags().String("user-id", "", "ID of the user who owns this contact (optional if user is selected)")
	contactCreateCmd.Flags().String("company", "", "Company name (optional)")
	contactCreateCmd.Flags().String("phone", "", "Phone number (optional)")
	contactCreateCmd.Flags().String("email", "", "Contact email (optional)")

	// Flags for contact list
	contactListCmd.Flags().String("user-id", "", "ID of the user whose contacts to list (optional if user is selected)")
}
