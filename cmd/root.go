package cmd

import (
	"fmt"
	"os"

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
  crm-admin user create "johndoe" "password123"
  crm-admin contact create "Jane Smith" --user-id "uuid-here" --company "Acme Corp" --phone "+1234567890"
  crm-admin note create "Meeting Notes" "Discussed project timeline" --contact-ids 1,2 --user-id "uuid-here"`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.crm-admin.yaml)")
}
