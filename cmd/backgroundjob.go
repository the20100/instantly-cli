package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var backgroundjobCmd = &cobra.Command{
	Use:   "job",
	Short: "View Instantly background jobs",
}

// ---- job list ----

var (
	backgroundjobListLimit int
	backgroundjobListStartingAfter  string
)

var backgroundjobListCmd = &cobra.Command{
	Use:   "list",
	Short: "List background jobs",
	Long: `List Instantly background jobs (bulk imports, exports, etc.).

Examples:
  instantly job list
  instantly job list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"limit", fmt.Sprintf("%d", backgroundjobListLimit),
			"starting_after", backgroundjobListStartingAfter,
		)
		items, _, err := client.ListBackgroundJobs(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printBackgroundJobsTable(items)
		return nil
	},
}

// ---- job get ----

var backgroundjobGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details of a specific background job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetBackgroundJob(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"ID", item.ID},
			{"Type", item.Type},
			{"Status", item.Status},
			{"Progress", fmt.Sprintf("%.0f%%", item.Progress*100)},
			{"Created", output.FormatTime(item.CreatedAt)},
			{"Updated", output.FormatTime(item.UpdatedAt)},
		})
		return nil
	},
}

func init() {
	backgroundjobListCmd.Flags().IntVar(&backgroundjobListLimit, "limit", 20, "Maximum number of jobs to return")
	backgroundjobListCmd.Flags().StringVar(&backgroundjobListStartingAfter, "starting-after", "", "Cursor for pagination: ID of the last item from the previous page")

	backgroundjobCmd.AddCommand(
		backgroundjobListCmd,
		backgroundjobGetCmd,
	)
	rootCmd.AddCommand(backgroundjobCmd)
}
