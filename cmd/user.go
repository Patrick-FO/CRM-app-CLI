package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"crm-admin/internal/api"
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

		fmt.Printf("✅ User '%s' created successfully! (ID: %s)\n", user.Username, user.ID)
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

		fmt.Println("📋 Users:")
		fmt.Printf("%-36s | %s\n", "ID", "Username")
		fmt.Printf("%-36s | %s\n", "------------------------------------", "--------")

		for _, user := range users {
			fmt.Printf("%-36s | %s\n", user.ID, user.Username)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userCreateCmd)
	userCmd.AddCommand(userListCmd)
}
