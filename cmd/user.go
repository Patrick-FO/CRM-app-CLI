package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"crm-admin/internal/api"
	"crm-admin/internal/context"
	"crm-admin/internal/models"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
	Long:  `Create, list, and manage users in the CRM system.`,
}

var userCreateCmd = &cobra.Command{
	Use:   "create [username] [password]",
	Short: "Create a new user",
	Long:  `Create a new user with the specified username and password.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.New()

		user, err := client.CreateUser(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		fmt.Printf("✅ User '%s' created successfully!\n", user.Username)
		fmt.Printf("   ID: %s\n", user.ID)
		fmt.Printf("\nYou can now select this user to work with their data:\n")
		fmt.Printf("   crm-admin user select %s\n", user.ID)
		return nil
	},
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  `Display a list of all users in the CRM system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := api.New()

		users, err := client.ListUsers()
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		if len(users) == 0 {
			fmt.Println("No users found.")
			return nil
		}

		// Show current context if any
		if context.HasUserContext() {
			userContext, _ := context.LoadUserContext()
			if userContext != nil {
				fmt.Printf("📌 Currently selected: %s (%s)\n\n", userContext.Username, userContext.UserID)
			}
		}

		fmt.Println("📋 Users:")
		fmt.Printf("%-36s | %s\n", "ID", "Username")
		fmt.Printf("%-36s | %s\n", "------------------------------------", "--------")

		for _, user := range users {
			marker := ""
			if context.HasUserContext() {
				userContext, _ := context.LoadUserContext()
				if userContext != nil && userContext.UserID == user.ID {
					marker = " ← SELECTED"
				}
			}
			fmt.Printf("%-36s | %s%s\n", user.ID, user.Username, marker)
		}

		fmt.Printf("\nTo select a user for easier management:\n")
		fmt.Printf("   crm-admin user select [user-id]\n")

		return nil
	},
}

var userSelectCmd = &cobra.Command{
	Use:   "select [user-id]",
	Short: "Select a user to work with",
	Long:  `Select a user by ID. This will make all subsequent contact and note operations default to this user.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID := args[0]

		// Create a user object with just the ID for context
		// We don't need to verify it exists - just set the context
		user := &models.User{
			ID:       userID,
			Username: userID, // Use ID as username for display until we can get the real username
		}

		// Save user context
		err := context.SaveUserContext(user)
		if err != nil {
			return fmt.Errorf("failed to save user context: %w", err)
		}

		fmt.Printf("✅ Selected user ID: %s\n", userID)
		fmt.Printf("🎯 Context set - commands will now default to this user\n")
		fmt.Println()
		fmt.Println("You can now run commands without specifying --user-id:")
		fmt.Printf("   crm-admin contact create \"John Doe\" --company \"Acme Corp\"\n")
		fmt.Printf("   crm-admin contact list\n")
		fmt.Printf("   crm-admin note create \"Meeting\" \"Important discussion\" --contact-ids 1,2\n")
		fmt.Printf("   crm-admin note list\n")
		fmt.Println()
		fmt.Printf("To switch users: crm-admin user select [other-user-id]\n")
		fmt.Printf("To exit user mode: crm-admin user exit\n")

		return nil
	},
}

var userExitCmd = &cobra.Command{
	Use:   "exit",
	Short: "Exit user selection mode",
	Long:  `Exit the current user selection and return to global mode where user-id is required for operations.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !context.HasUserContext() {
			fmt.Println("No user is currently selected.")
			return nil
		}

		userContext, err := context.LoadUserContext()
		if err != nil {
			return fmt.Errorf("failed to load user context: %w", err)
		}

		err = context.ClearUserContext()
		if err != nil {
			return fmt.Errorf("failed to clear user context: %w", err)
		}

		fmt.Printf("✅ Exited user mode for: %s\n", userContext.Username)
		fmt.Println("🌐 Back to global mode - you'll need to specify --user-id for operations")

		return nil
	},
}

var userInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show currently selected user info",
	Long:  `Display information about the currently selected user and their data.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !context.HasUserContext() {
			fmt.Println("No user currently selected.")
			fmt.Println("Use 'crm-admin user select [user-id]' to select a user.")
			return nil
		}

		userContext, err := context.LoadUserContext()
		if err != nil {
			return fmt.Errorf("failed to load user context: %w", err)
		}

		client := api.New()

		fmt.Printf("👤 Currently Selected User: %s (ID: %s)\n", userContext.Username, userContext.UserID)
		fmt.Println("=" + fmt.Sprintf("%*s", 50, "="))

		// Show contacts
		contacts, err := client.ListContacts("")
		if err != nil {
			fmt.Printf("❌ Failed to load contacts: %v\n", err)
		} else {
			fmt.Printf("\n📋 Contacts (%d):\n", len(contacts))
			if len(contacts) == 0 {
				fmt.Println("   No contacts found.")
			} else {
				for _, contact := range contacts {
					fmt.Printf("   %d. %s", contact.ID, contact.Name)
					if contact.Company != nil {
						fmt.Printf(" (%s)", *contact.Company)
					}
					fmt.Println()
				}
			}
		}

		// Show notes
		notes, err := client.ListNotesForUser("")
		if err != nil {
			fmt.Printf("❌ Failed to load notes: %v\n", err)
		} else {
			fmt.Printf("\n📝 Notes (%d):\n", len(notes))
			if len(notes) == 0 {
				fmt.Println("   No notes found.")
			} else {
				for _, note := range notes {
					fmt.Printf("   %d. %s\n", note.ID, note.Title)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userCreateCmd)
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userSelectCmd)
	userCmd.AddCommand(userExitCmd)
	userCmd.AddCommand(userInfoCmd)
}
