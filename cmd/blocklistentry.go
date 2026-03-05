package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var blocklistCmd = &cobra.Command{
	Use:   "blocklist",
	Short: "Manage Instantly blocklist entries",
}

// ---- blocklist list ----

var (
	blocklistListLimit int
	blocklistListStartingAfter  string
)

var blocklistListCmd = &cobra.Command{
	Use:   "list",
	Short: "List blocklist entries",
	Long: `List Instantly blocklist entries (blocked emails/domains).

Examples:
  instantly blocklist list
  instantly blocklist list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"limit", fmt.Sprintf("%d", blocklistListLimit),
			"starting_after", blocklistListStartingAfter,
		)
		items, _, err := client.ListBlocklistEntries(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printBlocklistEntriesTable(items)
		return nil
	},
}

// ---- blocklist get ----

var blocklistGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details of a specific blocklist entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetBlocklistEntry(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"ID", item.ID},
			{"Value", item.Value},
			{"Type", item.Type},
			{"Created", output.FormatTime(item.CreatedAt)},
		})
		return nil
	},
}

// ---- blocklist create ----

var (
	blocklistCreateValue string
	blocklistCreateType  string
)

var blocklistCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Add an entry to the blocklist",
	Long: `Add an email or domain to the Instantly blocklist.

Examples:
  instantly blocklist create --value spam@example.com --type email
  instantly blocklist create --value spammydomain.com --type domain`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if blocklistCreateValue == "" {
			return fmt.Errorf("--value is required")
		}
		if blocklistCreateType == "" {
			return fmt.Errorf("--type is required (email or domain)")
		}
		item, err := client.CreateBlocklistEntry(map[string]interface{}{
			"value": blocklistCreateValue,
			"type":  blocklistCreateType,
		})
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Blocklist entry created: %s (%s)\n", item.Value, item.Type)
		fmt.Printf("ID: %s\n", item.ID)
		return nil
	},
}

// ---- blocklist update ----

var (
	blocklistUpdateValue string
	blocklistUpdateType  string
)

var blocklistUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a blocklist entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{}
		if cmd.Flags().Changed("value") {
			payload["value"] = blocklistUpdateValue
		}
		if cmd.Flags().Changed("type") {
			payload["type"] = blocklistUpdateType
		}
		if len(payload) == 0 {
			return fmt.Errorf("no fields to update — provide --value or --type")
		}
		item, err := client.UpdateBlocklistEntry(args[0], payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Blocklist entry updated: %s\n", item.Value)
		return nil
	},
}

// ---- blocklist delete ----

var blocklistDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a blocklist entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteBlocklistEntry(args[0]); err != nil {
			return err
		}
		fmt.Printf("Blocklist entry %s deleted.\n", args[0])
		return nil
	},
}

func init() {
	blocklistListCmd.Flags().IntVar(&blocklistListLimit, "limit", 20, "Maximum number of entries to return")
	blocklistListCmd.Flags().StringVar(&blocklistListStartingAfter, "starting-after", "", "Cursor for pagination: ID of the last item from the previous page")

	blocklistCreateCmd.Flags().StringVar(&blocklistCreateValue, "value", "", "Email or domain to block *(required)*")
	blocklistCreateCmd.Flags().StringVar(&blocklistCreateType, "type", "", "Entry type: email or domain *(required)*")

	blocklistUpdateCmd.Flags().StringVar(&blocklistUpdateValue, "value", "", "New value")
	blocklistUpdateCmd.Flags().StringVar(&blocklistUpdateType, "type", "", "New type (email or domain)")

	blocklistCmd.AddCommand(
		blocklistListCmd,
		blocklistGetCmd,
		blocklistCreateCmd,
		blocklistUpdateCmd,
		blocklistDeleteCmd,
	)
	rootCmd.AddCommand(blocklistCmd)
}
