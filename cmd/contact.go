package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"crm-admin/internal/api"
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
		if userID == "" {
			return fmt.Errorf("user-id flag is required")
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
		return nil
	},
}

var contactListCmd = &cobra.Command{
	Use:   "list",
	Short: "List contacts",
	Long:  `Display a list of contacts. Optionally filter by user ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, _ := cmd.Flags().GetString("user-id")

		client := api.New()

		contacts, err := client.ListContacts(userID)
		if err != nil {
			return fmt.Errorf("failed to list contacts: %w", err)
		}

		if len(contacts) == 0 {
			if userID != "" {
				fmt.Printf("No contacts found for user ID %s.\n", userID)
			} else {
				fmt.Println("No contacts found.")
			}
			return nil
		}

		fmt.Println("📋 Contacts:")
		fmt.Printf("%-5s | %-20s | %-36s | %-15s | %-15s | %s\n",
			"ID", "Name", "User ID", "Company", "Phone", "Email")
		fmt.Printf("%-5s | %-20s | %-36s | %-15s | %-15s | %s\n",
			"-----", "--------------------", "------------------------------------",
			"---------------", "---------------", "---------------")

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

			fmt.Printf("%-5d | %-20s | %-36s | %-15s | %-15s | %s\n",
				contact.ID, contact.Name, contact.UserID, company, phone, email)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(contactCmd)
	contactCmd.AddCommand(contactCreateCmd)
	contactCmd.AddCommand(contactListCmd)

	// Flags for contact create
	contactCreateCmd.Flags().String("user-id", "", "ID of the user who owns this contact (required)")
	contactCreateCmd.Flags().String("company", "", "Company name (optional)")
	contactCreateCmd.Flags().String("phone", "", "Phone number (optional)")
	contactCreateCmd.Flags().String("email", "", "Contact email (optional)")
	contactCreateCmd.MarkFlagRequired("user-id")

	// Flags for contact list
	contactListCmd.Flags().String("user-id", "", "Filter contacts by user ID (optional)")
}
