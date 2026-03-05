package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var apikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Manage Instantly API keys",
}

// ---- apikey list ----

var (
	apikeyListLimit int
	apikeyListStartingAfter  string
)

var apikeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List API keys",
	Long: `List Instantly API keys for the workspace.

Examples:
  instantly apikey list
  instantly apikey list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"limit", fmt.Sprintf("%d", apikeyListLimit),
			"starting_after", apikeyListStartingAfter,
		)
		items, _, err := client.ListAPIKeys(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printAPIKeysTable(items)
		return nil
	},
}

// ---- apikey create ----

var apikeyCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.CreateAPIKey(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("API key created: %s\n", item.Name)
		fmt.Printf("ID:  %s\n", item.ID)
		fmt.Printf("Key: %s\n", item.APIKey)
		fmt.Printf("\nSave this key — it will not be shown again.\n")
		return nil
	},
}

// ---- apikey delete ----

var apikeyDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteAPIKey(args[0]); err != nil {
			return err
		}
		fmt.Printf("API key %s deleted.\n", args[0])
		return nil
	},
}

func init() {
	apikeyListCmd.Flags().IntVar(&apikeyListLimit, "limit", 20, "Maximum number of API keys to return")
	apikeyListCmd.Flags().StringVar(&apikeyListStartingAfter, "starting-after", "", "Cursor for pagination: ID of the last item from the previous page")

	apikeyCmd.AddCommand(
		apikeyListCmd,
		apikeyCreateCmd,
		apikeyDeleteCmd,
	)
	rootCmd.AddCommand(apikeyCmd)
}
