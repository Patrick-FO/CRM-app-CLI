package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"crm-admin/internal/api"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage notes",
	Long:  `Create, list, and manage notes for contacts in the CRM system.`,
}

var noteCreateCmd = &cobra.Command{
	Use:   "create [title] [description]",
	Short: "Create a new note",
	Long:  `Create a new note with the specified title and description for one or more contacts.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		contactIDsStr, _ := cmd.Flags().GetStringSlice("contact-ids")
		userID, _ := cmd.Flags().GetString("user-id")

		if len(contactIDsStr) == 0 {
			return fmt.Errorf("contact-ids flag is required (comma-separated list of contact IDs)")
		}
		if userID == "" {
			return fmt.Errorf("user-id flag is required")
		}

		// Parse contact IDs
		contactIDs := make([]int, len(contactIDsStr))
		for i, idStr := range contactIDsStr {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return fmt.Errorf("invalid contact ID '%s': %w", idStr, err)
			}
			contactIDs[i] = id
		}

		client := api.New()

		note, err := client.CreateNote(args[0], args[1], contactIDs, userID)
		if err != nil {
			return fmt.Errorf("failed to create note: %w", err)
		}

		fmt.Printf("✅ Note created successfully!\n")
		fmt.Printf("   ID: %d\n", note.ID)
		fmt.Printf("   Title: %s\n", note.Title)
		if note.Description != nil {
			fmt.Printf("   Description: %s\n", *note.Description)
		}
		fmt.Printf("   Contact IDs: %v\n", note.ContactIDs)
		fmt.Printf("   User ID: %s\n", note.UserID)
		return nil
	},
}

var noteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List notes",
	Long:  `Display a list of notes. Optionally filter by contact IDs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		contactIDsStr, _ := cmd.Flags().GetStringSlice("contact-ids")

		// Parse contact IDs if provided
		var contactIDs []int
		if len(contactIDsStr) > 0 {
			contactIDs = make([]int, len(contactIDsStr))
			for i, idStr := range contactIDsStr {
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return fmt.Errorf("invalid contact ID '%s': %w", idStr, err)
				}
				contactIDs[i] = id
			}
		}

		client := api.New()

		notes, err := client.ListNotes(contactIDs)
		if err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		if len(notes) == 0 {
			if len(contactIDs) > 0 {
				fmt.Printf("No notes found for contact IDs %v.\n", contactIDs)
			} else {
				fmt.Println("No notes found.")
			}
			return nil
		}

		fmt.Println("📋 Notes:")
		fmt.Printf("%-5s | %-25s | %-30s | %-15s | %s\n",
			"ID", "Title", "Description", "Contact IDs", "User ID")
		fmt.Printf("%-5s | %-25s | %-30s | %-15s | %s\n",
			"-----", "-------------------------", "------------------------------",
			"---------------", "------------------------------------")

		for _, note := range notes {
			description := ""
			if note.Description != nil {
				desc := *note.Description
				if len(desc) > 30 {
					description = desc[:27] + "..."
				} else {
					description = desc
				}
			}

			contactIDsStr := make([]string, len(note.ContactIDs))
			for i, id := range note.ContactIDs {
				contactIDsStr[i] = strconv.Itoa(id)
			}
			contactIDsDisplay := strings.Join(contactIDsStr, ",")
			if len(contactIDsDisplay) > 15 {
				contactIDsDisplay = contactIDsDisplay[:12] + "..."
			}

			fmt.Printf("%-5d | %-25s | %-30s | %-15s | %s\n",
				note.ID, note.Title, description, contactIDsDisplay, note.UserID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
	noteCmd.AddCommand(noteCreateCmd)
	noteCmd.AddCommand(noteListCmd)

	// Flags for note create
	noteCreateCmd.Flags().StringSlice("contact-ids", []string{}, "Comma-separated list of contact IDs this note belongs to (required)")
	noteCreateCmd.Flags().String("user-id", "", "ID of the user creating this note (required)")
	noteCreateCmd.MarkFlagRequired("contact-ids")
	noteCreateCmd.MarkFlagRequired("user-id")

	// Flags for note list
	noteListCmd.Flags().StringSlice("contact-ids", []string{}, "Filter notes by contact IDs (optional)")
}
