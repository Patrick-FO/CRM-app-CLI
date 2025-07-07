package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"crm-admin/internal/api"
	"crm-admin/internal/context"
	"crm-admin/internal/models"
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

		// Check if user-id is provided or if we have context
		if userID == "" && !context.HasUserContext() {
			return fmt.Errorf("user-id flag is required (or select a user with 'crm-admin user select [user-id]')")
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
	Long:  `Display a list of notes for a user. Optionally filter by contact ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, _ := cmd.Flags().GetString("user-id")
		contactID, _ := cmd.Flags().GetInt("contact-id")

		// Check if user-id is provided or if we have context
		if userID == "" && !context.HasUserContext() {
			return fmt.Errorf("user-id flag is required (or select a user with 'crm-admin user select [user-id]')")
		}

		client := api.New()

		var notes []models.Note
		var err error

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

		if contactID > 0 {
			notes, err = client.ListNotesForContact(userID, contactID)
			if err != nil {
				return fmt.Errorf("failed to list notes for contact: %w", err)
			}
		} else {
			notes, err = client.ListNotesForUser(userID)
			if err != nil {
				return fmt.Errorf("failed to list notes for user: %w", err)
			}
		}

		if len(notes) == 0 {
			if contactID > 0 {
				fmt.Printf("No notes found for contact ID %d.\n", contactID)
			} else if username != "" {
				fmt.Printf("No notes found for %s (%s).\n", username, targetUserID)
			} else {
				fmt.Printf("No notes found for user %s.\n", targetUserID)
			}
			return nil
		}

		if contactID > 0 {
			fmt.Printf("📝 Notes for Contact %d:\n", contactID)
		} else if username != "" {
			fmt.Printf("📝 Notes for %s (%s):\n", username, targetUserID)
		} else {
			fmt.Printf("📝 Notes for User %s:\n", targetUserID)
		}

		fmt.Printf("%-5s | %-25s | %-50s | %s\n",
			"ID", "Title", "Description", "Contact IDs")
		fmt.Printf("%-5s | %-25s | %-50s | %s\n",
			"-----", "-------------------------", "--------------------------------------------------", "----------")

		for _, note := range notes {
			description := ""
			if note.Description != nil {
				desc := *note.Description
				if len(desc) > 50 {
					description = desc[:47] + "..."
				} else {
					description = desc
				}
			}

			contactIDs := ""
			for i, id := range note.ContactIDs {
				if i > 0 {
					contactIDs += ","
				}
				contactIDs += fmt.Sprintf("%d", id)
			}

			fmt.Printf("%-5d | %-25s | %-50s | %s\n",
				note.ID, note.Title, description, contactIDs)
		}

		return nil
	},
}

var noteGetCmd = &cobra.Command{
	Use:   "get [note-id]",
	Short: "Get a specific note",
	Long:  `Get detailed information about a specific note by its ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, _ := cmd.Flags().GetString("user-id")

		// Check if user-id is provided or if we have context
		if userID == "" && !context.HasUserContext() {
			return fmt.Errorf("user-id flag is required (or select a user with 'crm-admin user select [user-id]')")
		}

		noteID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid note ID '%s': %w", args[0], err)
		}

		client := api.New()

		note, err := client.GetNote(userID, noteID)
		if err != nil {
			return fmt.Errorf("failed to get note: %w", err)
		}

		fmt.Printf("📝 Note Details:\n")
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

var noteUpdateCmd = &cobra.Command{
	Use:   "update [note-id] [title] [description]",
	Short: "Update a note",
	Long:  `Update an existing note with new title and description.`,
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, _ := cmd.Flags().GetString("user-id")
		contactIDsStr, _ := cmd.Flags().GetStringSlice("contact-ids")

		// Check if user-id is provided or if we have context
		if userID == "" && !context.HasUserContext() {
			return fmt.Errorf("user-id flag is required (or select a user with 'crm-admin user select [user-id]')")
		}

		if len(contactIDsStr) == 0 {
			return fmt.Errorf("contact-ids flag is required (comma-separated list of contact IDs)")
		}

		noteID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid note ID '%s': %w", args[0], err)
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

		note, err := client.UpdateNote(userID, noteID, args[1], args[2], contactIDs)
		if err != nil {
			return fmt.Errorf("failed to update note: %w", err)
		}

		fmt.Printf("✅ Note updated successfully!\n")
		fmt.Printf("   ID: %d\n", note.ID)
		fmt.Printf("   Title: %s\n", note.Title)
		if note.Description != nil {
			fmt.Printf("   Description: %s\n", *note.Description)
		}
		fmt.Printf("   Contact IDs: %v\n", note.ContactIDs)

		return nil
	},
}

var noteDeleteCmd = &cobra.Command{
	Use:   "delete [note-id]",
	Short: "Delete a note",
	Long:  `Delete a note by its ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, _ := cmd.Flags().GetString("user-id")

		// Check if user-id is provided or if we have context
		if userID == "" && !context.HasUserContext() {
			return fmt.Errorf("user-id flag is required (or select a user with 'crm-admin user select [user-id]')")
		}

		noteID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid note ID '%s': %w", args[0], err)
		}

		client := api.New()

		err = client.DeleteNote(userID, noteID)
		if err != nil {
			return fmt.Errorf("failed to delete note: %w", err)
		}

		fmt.Printf("✅ Note %d deleted successfully!\n", noteID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(noteCmd)
	noteCmd.AddCommand(noteCreateCmd)
	noteCmd.AddCommand(noteListCmd)
	noteCmd.AddCommand(noteGetCmd)
	noteCmd.AddCommand(noteUpdateCmd)
	noteCmd.AddCommand(noteDeleteCmd)

	// Flags for note create
	noteCreateCmd.Flags().StringSlice("contact-ids", []string{}, "Comma-separated list of contact IDs this note belongs to (required)")
	noteCreateCmd.Flags().String("user-id", "", "ID of the user creating this note (optional if user is selected)")
	noteCreateCmd.MarkFlagRequired("contact-ids")

	// Flags for note list
	noteListCmd.Flags().String("user-id", "", "ID of the user whose notes to list (optional if user is selected)")
	noteListCmd.Flags().Int("contact-id", 0, "Filter notes by contact ID (optional)")

	// Flags for note get
	noteGetCmd.Flags().String("user-id", "", "ID of the user who owns the note (optional if user is selected)")

	// Flags for note update
	noteUpdateCmd.Flags().StringSlice("contact-ids", []string{}, "Comma-separated list of contact IDs this note belongs to (required)")
	noteUpdateCmd.Flags().String("user-id", "", "ID of the user who owns the note (optional if user is selected)")
	noteUpdateCmd.MarkFlagRequired("contact-ids")

	// Flags for note delete
	noteDeleteCmd.Flags().String("user-id", "", "ID of the user who owns the note (optional if user is selected)")
}
