package cmd

import (
	"fmt"
	"os"

	"crm-admin/internal/context"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "crm-admin",
	Short: "CRM Admin CLI Tool",
	Long: `A command line interface for managing CRM data.
	
This tool allows administrators to:
- Create and manage users
- Create and manage contacts for users  
- Create and manage notes for contacts

Examples:
  # User management
  crm-admin user create "johndoe" "password123"
  crm-admin user list
  crm-admin user select "user-uuid-here"
  
  # Contact management (with user selected)
  crm-admin contact create "Jane Smith" --company "Acme Corp" --phone "+1234567890"
  crm-admin contact list
  
  # Note management (with user selected)
  crm-admin note create "Meeting Notes" "Discussed project timeline" --contact-ids 1,2
  crm-admin note list`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Update command prompt based on user context
	updatePromptForContext()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func updatePromptForContext() {
	if context.HasUserContext() {
		userContext, err := context.LoadUserContext()
		if err == nil && userContext != nil {
			rootCmd.Use = fmt.Sprintf("crm-admin [%s]", userContext.Username)
		}
	}
}

func init() {
	// Global flags can be added here
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.crm-admin.yaml)")
}
