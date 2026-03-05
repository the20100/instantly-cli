package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var customtagCmd = &cobra.Command{
	Use:   "customtag",
	Short: "Manage Instantly custom tags",
}

// ---- customtag list ----

var (
	customtagListLimit int
	customtagListStartingAfter  string
)

var customtagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List custom tags",
	Long: `List Instantly custom tags.

Examples:
  instantly customtag list
  instantly customtag list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"limit", fmt.Sprintf("%d", customtagListLimit),
			"starting_after", customtagListStartingAfter,
		)
		items, _, err := client.ListCustomTags(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printCustomTagsTable(items)
		return nil
	},
}

// ---- customtag get ----

var customtagGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details of a specific custom tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetCustomTag(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"ID", item.ID},
			{"Name", item.Name},
			{"Color", item.Color},
			{"Workspace ID", item.WorkspaceID},
			{"Created", output.FormatTime(item.CreatedAt)},
		})
		return nil
	},
}

// ---- customtag create ----

var customtagCreateColor string

var customtagCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new custom tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{"name": args[0]}
		if customtagCreateColor != "" {
			payload["color"] = customtagCreateColor
		}
		item, err := client.CreateCustomTag(payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Custom tag created: %s\n", item.Name)
		fmt.Printf("ID: %s\n", item.ID)
		return nil
	},
}

// ---- customtag update ----

var (
	customtagUpdateName  string
	customtagUpdateColor string
)

var customtagUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a custom tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			payload["name"] = customtagUpdateName
		}
		if cmd.Flags().Changed("color") {
			payload["color"] = customtagUpdateColor
		}
		if len(payload) == 0 {
			return fmt.Errorf("no fields to update — provide --name or --color")
		}
		item, err := client.UpdateCustomTag(args[0], payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Custom tag updated: %s\n", item.Name)
		return nil
	},
}

// ---- customtag delete ----

var customtagDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a custom tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteCustomTag(args[0]); err != nil {
			return err
		}
		fmt.Printf("Custom tag %s deleted.\n", args[0])
		return nil
	},
}

// ---- customtag toggle ----

var (
	customtagToggleTagID       string
	customtagToggleResourceID   string
	customtagToggleResourceType string
)

var customtagToggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle a custom tag on a resource",
	Long: `Add or remove a custom tag from a resource (campaign, lead, etc.).

Examples:
  instantly customtag toggle --tag-id <id> --resource-id <id> --resource-type campaign
  instantly customtag toggle --tag-id <id> --resource-id <id> --resource-type lead`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if customtagToggleTagID == "" {
			return fmt.Errorf("--tag-id is required")
		}
		if customtagToggleResourceID == "" {
			return fmt.Errorf("--resource-id is required")
		}
		if customtagToggleResourceType == "" {
			return fmt.Errorf("--resource-type is required")
		}
		if err := client.ToggleCustomTag(customtagToggleTagID, customtagToggleResourceID, customtagToggleResourceType); err != nil {
			return err
		}
		fmt.Printf("Tag %s toggled on %s %s.\n", customtagToggleTagID, customtagToggleResourceType, customtagToggleResourceID)
		return nil
	},
}

func init() {
	customtagListCmd.Flags().IntVar(&customtagListLimit, "limit", 20, "Maximum number of tags to return")
	customtagListCmd.Flags().StringVar(&customtagListStartingAfter, "starting-after", "", "Cursor for pagination: ID of the last item from the previous page")

	customtagCreateCmd.Flags().StringVar(&customtagCreateColor, "color", "", "Tag color (hex or name)")

	customtagUpdateCmd.Flags().StringVar(&customtagUpdateName, "name", "", "New tag name")
	customtagUpdateCmd.Flags().StringVar(&customtagUpdateColor, "color", "", "New tag color")

	customtagToggleCmd.Flags().StringVar(&customtagToggleTagID, "tag-id", "", "Tag ID *(required)*")
	customtagToggleCmd.Flags().StringVar(&customtagToggleResourceID, "resource-id", "", "Resource ID *(required)*")
	customtagToggleCmd.Flags().StringVar(&customtagToggleResourceType, "resource-type", "", "Resource type (campaign, lead, etc.) *(required)*")

	customtagCmd.AddCommand(
		customtagListCmd,
		customtagGetCmd,
		customtagCreateCmd,
		customtagUpdateCmd,
		customtagDeleteCmd,
		customtagToggleCmd,
	)
	rootCmd.AddCommand(customtagCmd)
}
