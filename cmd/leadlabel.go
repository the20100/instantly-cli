package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var leadlabelCmd = &cobra.Command{
	Use:   "leadlabel",
	Short: "Manage Instantly lead labels",
}

// ---- leadlabel list ----

var (
	leadlabelListLimit int
	leadlabelListStartingAfter  string
)

var leadlabelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List lead labels",
	Long: `List Instantly lead labels.

Examples:
  instantly leadlabel list
  instantly leadlabel list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"limit", fmt.Sprintf("%d", leadlabelListLimit),
			"starting_after", leadlabelListStartingAfter,
		)
		items, _, err := client.ListLeadLabels(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printLeadLabelsTable(items)
		return nil
	},
}

// ---- leadlabel get ----

var leadlabelGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details of a specific lead label",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetLeadLabel(args[0])
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
			{"Created", output.FormatTime(item.CreatedAt)},
			{"Updated", output.FormatTime(item.UpdatedAt)},
		})
		return nil
	},
}

// ---- leadlabel create ----

var leadlabelCreateColor string

var leadlabelCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new lead label",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{"name": args[0]}
		if leadlabelCreateColor != "" {
			payload["color"] = leadlabelCreateColor
		}
		item, err := client.CreateLeadLabel(payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Lead label created: %s\n", item.Name)
		fmt.Printf("ID: %s\n", item.ID)
		return nil
	},
}

// ---- leadlabel update ----

var (
	leadlabelUpdateName  string
	leadlabelUpdateColor string
)

var leadlabelUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a lead label",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			payload["name"] = leadlabelUpdateName
		}
		if cmd.Flags().Changed("color") {
			payload["color"] = leadlabelUpdateColor
		}
		if len(payload) == 0 {
			return fmt.Errorf("no fields to update — provide --name or --color")
		}
		item, err := client.UpdateLeadLabel(args[0], payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Lead label updated: %s\n", item.Name)
		return nil
	},
}

// ---- leadlabel delete ----

var leadlabelDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a lead label",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteLeadLabel(args[0]); err != nil {
			return err
		}
		fmt.Printf("Lead label %s deleted.\n", args[0])
		return nil
	},
}

func init() {
	leadlabelListCmd.Flags().IntVar(&leadlabelListLimit, "limit", 20, "Maximum number of labels to return")
	leadlabelListCmd.Flags().StringVar(&leadlabelListStartingAfter, "starting-after", "", "Cursor for pagination: ID of the last item from the previous page")

	leadlabelCreateCmd.Flags().StringVar(&leadlabelCreateColor, "color", "", "Label color (hex or name)")

	leadlabelUpdateCmd.Flags().StringVar(&leadlabelUpdateName, "name", "", "New label name")
	leadlabelUpdateCmd.Flags().StringVar(&leadlabelUpdateColor, "color", "", "New label color")

	leadlabelCmd.AddCommand(
		leadlabelListCmd,
		leadlabelGetCmd,
		leadlabelCreateCmd,
		leadlabelUpdateCmd,
		leadlabelDeleteCmd,
	)
	rootCmd.AddCommand(leadlabelCmd)
}
