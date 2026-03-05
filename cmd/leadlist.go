package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var leadlistCmd = &cobra.Command{
	Use:   "leadlist",
	Short: "Manage Instantly lead lists",
}

// ---- leadlist list ----

var (
	leadlistListLimit int
	leadlistListStartingAfter  string
)

var leadlistListCmd = &cobra.Command{
	Use:   "list",
	Short: "List lead lists",
	Long: `List Instantly lead lists.

Examples:
  instantly leadlist list
  instantly leadlist list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"limit", fmt.Sprintf("%d", leadlistListLimit),
			"starting_after", leadlistListStartingAfter,
		)
		items, _, err := client.ListLeadLists(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printLeadListsTable(items)
		return nil
	},
}

// ---- leadlist get ----

var leadlistGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details of a specific lead list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetLeadList(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"ID", item.ID},
			{"Name", item.Name},
			{"Count", fmt.Sprintf("%d", item.Count)},
			{"Created", output.FormatTime(item.CreatedAt)},
			{"Updated", output.FormatTime(item.UpdatedAt)},
		})
		return nil
	},
}

// ---- leadlist create ----

var leadlistCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new lead list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.CreateLeadList(map[string]interface{}{"name": args[0]})
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Lead list created: %s\n", item.Name)
		fmt.Printf("ID: %s\n", item.ID)
		return nil
	},
}

// ---- leadlist update ----

var leadlistUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a lead list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name is required")
		}
		item, err := client.UpdateLeadList(args[0], map[string]interface{}{"name": name})
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Lead list updated: %s\n", item.Name)
		return nil
	},
}

// ---- leadlist delete ----

var leadlistDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a lead list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteLeadList(args[0]); err != nil {
			return err
		}
		fmt.Printf("Lead list %s deleted.\n", args[0])
		return nil
	},
}

func init() {
	leadlistListCmd.Flags().IntVar(&leadlistListLimit, "limit", 20, "Maximum number of lead lists to return")
	leadlistListCmd.Flags().StringVar(&leadlistListStartingAfter, "starting-after", "", "Cursor for pagination: ID of the last item from the previous page")

	leadlistUpdateCmd.Flags().String("name", "", "New lead list name *(required)*")

	leadlistCmd.AddCommand(
		leadlistListCmd,
		leadlistGetCmd,
		leadlistCreateCmd,
		leadlistUpdateCmd,
		leadlistDeleteCmd,
	)
	rootCmd.AddCommand(leadlistCmd)
}
