package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "View Instantly analytics",
}

// ---- analytics campaign ----

var (
	analyticsCampaignID        string
	analyticsCampaignStartDate string
	analyticsCampaignEndDate   string
)

var analyticsCampaignCmd = &cobra.Command{
	Use:   "campaign",
	Short: "Get campaign analytics",
	Long: `Get analytics for a specific campaign or all campaigns.

Examples:
  instantly analytics campaign --campaign-id <id>
  instantly analytics campaign --campaign-id <id> --start-date 2024-01-01 --end-date 2024-12-31
  instantly analytics campaign --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if analyticsCampaignID == "" {
			// overview mode
			items, err := client.GetCampaignAnalyticsOverview(buildParams(
				"start_date", analyticsCampaignStartDate,
				"end_date", analyticsCampaignEndDate,
			))
			if err != nil {
				return err
			}
			if output.IsJSON(cmd) {
				return output.PrintJSON(items, output.IsPretty(cmd))
			}
			if len(items) == 0 {
				fmt.Println("No analytics data found.")
				return nil
			}
			headers := []string{"CAMPAIGN", "SENT", "OPENED", "REPLIED", "BOUNCED", "OPEN%", "REPLY%"}
			rows := make([][]string, len(items))
			for i, item := range items {
				rows[i] = []string{
					output.Truncate(item.CampaignName, 30),
					fmt.Sprintf("%d", item.Sent),
					fmt.Sprintf("%d", item.Opened),
					fmt.Sprintf("%d", item.Replied),
					fmt.Sprintf("%d", item.Bounced),
					fmt.Sprintf("%.1f%%", item.OpenRate*100),
					fmt.Sprintf("%.1f%%", item.ReplyRate*100),
				}
			}
			output.PrintTable(headers, rows)
			return nil
		}
		// single campaign mode
		params := buildParams(
			"start_date", analyticsCampaignStartDate,
			"end_date", analyticsCampaignEndDate,
		)
		item, err := client.GetCampaignAnalytics(analyticsCampaignID, params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"Campaign ID", item.CampaignID},
			{"Campaign Name", item.CampaignName},
			{"Sent", fmt.Sprintf("%d", item.Sent)},
			{"Opened", fmt.Sprintf("%d (%.1f%%)", item.Opened, item.OpenRate*100)},
			{"Clicked", fmt.Sprintf("%d (%.1f%%)", item.Clicked, item.ClickRate*100)},
			{"Replied", fmt.Sprintf("%d (%.1f%%)", item.Replied, item.ReplyRate*100)},
			{"Bounced", fmt.Sprintf("%d", item.Bounced)},
			{"Unsubscribed", fmt.Sprintf("%d", item.Unsubscribed)},
			{"New Leads Contacted", fmt.Sprintf("%d", item.NewLeads)},
		})
		return nil
	},
}

// ---- analytics warmup ----

var (
	analyticsWarmupEmail     string
	analyticsWarmupStartDate string
	analyticsWarmupEndDate   string
)

var analyticsWarmupCmd = &cobra.Command{
	Use:   "warmup",
	Short: "Get warmup analytics for email accounts",
	Long: `Get warmup analytics for email accounts.

Examples:
  instantly analytics warmup
  instantly analytics warmup --email user@domain.com
  instantly analytics warmup --start-date 2024-01-01 --end-date 2024-01-31`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"email", analyticsWarmupEmail,
			"start_date", analyticsWarmupStartDate,
			"end_date", analyticsWarmupEndDate,
		)
		items, err := client.GetWarmupAnalytics(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		if len(items) == 0 {
			fmt.Println("No warmup analytics found.")
			return nil
		}
		headers := []string{"EMAIL", "DATE", "SENT", "RECEIVED"}
		rows := make([][]string, len(items))
		for i, item := range items {
			rows[i] = []string{
				output.Truncate(item.Email, 36),
				output.FormatTime(item.Date),
				fmt.Sprintf("%d", item.Sent),
				fmt.Sprintf("%d", item.Received),
			}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

func init() {
	analyticsCampaignCmd.Flags().StringVar(&analyticsCampaignID, "campaign-id", "", "Campaign ID (omit for overview of all campaigns)")
	analyticsCampaignCmd.Flags().StringVar(&analyticsCampaignStartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	analyticsCampaignCmd.Flags().StringVar(&analyticsCampaignEndDate, "end-date", "", "End date (YYYY-MM-DD)")

	analyticsWarmupCmd.Flags().StringVar(&analyticsWarmupEmail, "email", "", "Filter by account email")
	analyticsWarmupCmd.Flags().StringVar(&analyticsWarmupStartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	analyticsWarmupCmd.Flags().StringVar(&analyticsWarmupEndDate, "end-date", "", "End date (YYYY-MM-DD)")

	analyticsCmd.AddCommand(
		analyticsCampaignCmd,
		analyticsWarmupCmd,
	)
	rootCmd.AddCommand(analyticsCmd)
}
