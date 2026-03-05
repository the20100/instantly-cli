package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var campaignCmd = &cobra.Command{
	Use:   "campaign",
	Short: "Manage Instantly campaigns",
}

// ---- campaign list ----

var (
	campaignListLimit         int
	campaignListStartingAfter string
	campaignListSearch        string
	campaignListStatus        string
)

var campaignListCmd = &cobra.Command{
	Use:   "list",
	Short: "List campaigns",
	Long: `List Instantly campaigns.

Examples:
  instantly campaign list
  instantly campaign list --status active
  instantly campaign list --search "my campaign" --limit 20
  instantly campaign list --starting-after <id>
  instantly campaign list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"limit", fmt.Sprintf("%d", campaignListLimit),
			"starting_after", campaignListStartingAfter,
			"search", campaignListSearch,
			"status", campaignListStatus,
		)
		items, _, err := client.ListCampaigns(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printCampaignsTable(items)
		return nil
	},
}

// ---- campaign get ----

var campaignGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details of a specific campaign",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetCampaign(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"ID", item.ID},
			{"Name", item.Name},
			{"Status", campaignStatusLabel(item.Status)},
			{"Daily Limit", fmt.Sprintf("%d", item.DailyLimit)},
			{"Timezone", item.Timezone},
			{"Stop on Reply", output.FormatBool(item.StopOnReply)},
			{"Stop on Auto-Reply", output.FormatBool(item.StopOnAutoReply)},
			{"Link Tracking", output.FormatBool(item.LinkTracking)},
			{"Open Tracking", output.FormatBool(item.OpenTracking)},
			{"Text Only", output.FormatBool(item.TextOnly)},
			{"Email Accounts", fmt.Sprintf("%d", len(item.EmailList))},
			{"Created", output.FormatTime(item.CreatedAt)},
			{"Updated", output.FormatTime(item.UpdatedAt)},
		})
		return nil
	},
}

// ---- campaign create ----

var (
	campaignCreateTimezone        string
	campaignCreateStartHour       string
	campaignCreateEndHour         string
	campaignCreateDailyLimit      int
	campaignCreateAccounts        []string
	campaignCreateStopOnReply     bool
	campaignCreateStopOnAutoReply bool
	campaignCreateLinkTracking    bool
	campaignCreateOpenTracking    bool
	campaignCreateTextOnly        bool
)

var campaignCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new campaign",
	Args:  cobra.ExactArgs(1),
	Long: `Create a new Instantly campaign.

Examples:
  instantly campaign create "My Outreach Campaign" --accounts email1@domain.com,email2@domain.com
  instantly campaign create "Cold Outreach" --daily-limit 50 --stop-on-reply --open-tracking`,
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{
			"name": args[0],
		}
		if campaignCreateTimezone != "" {
			payload["campaign_schedule"] = map[string]interface{}{
				"timezone":   campaignCreateTimezone,
				"start_hour": campaignCreateStartHour,
				"end_hour":   campaignCreateEndHour,
			}
		}
		if campaignCreateDailyLimit > 0 {
			payload["daily_limit"] = campaignCreateDailyLimit
		}
		if len(campaignCreateAccounts) > 0 {
			// flatten comma-separated or multi-flag values
			var emails []string
			for _, a := range campaignCreateAccounts {
				for _, e := range strings.Split(a, ",") {
					if e = strings.TrimSpace(e); e != "" {
						emails = append(emails, e)
					}
				}
			}
			payload["email_list"] = emails
		}
		if cmd.Flags().Changed("stop-on-reply") {
			payload["stop_on_reply"] = campaignCreateStopOnReply
		}
		if cmd.Flags().Changed("stop-on-auto-reply") {
			payload["stop_on_auto_reply"] = campaignCreateStopOnAutoReply
		}
		if cmd.Flags().Changed("link-tracking") {
			payload["link_tracking"] = campaignCreateLinkTracking
		}
		if cmd.Flags().Changed("open-tracking") {
			payload["open_tracking"] = campaignCreateOpenTracking
		}
		if cmd.Flags().Changed("text-only") {
			payload["text_only"] = campaignCreateTextOnly
		}
		item, err := client.CreateCampaign(payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Campaign created: %s\n", item.Name)
		fmt.Printf("ID: %s\n", item.ID)
		fmt.Printf("Status: %s\n", campaignStatusLabel(item.Status))
		return nil
	},
}

// ---- campaign update ----

var (
	campaignUpdateName        string
	campaignUpdateDailyLimit  int
	campaignUpdateTimezone    string
)

var campaignUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a campaign",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			payload["name"] = campaignUpdateName
		}
		if cmd.Flags().Changed("daily-limit") {
			payload["daily_limit"] = campaignUpdateDailyLimit
		}
		if cmd.Flags().Changed("timezone") {
			payload["timezone"] = campaignUpdateTimezone
		}
		if len(payload) == 0 {
			return fmt.Errorf("no fields to update — provide at least one flag")
		}
		item, err := client.UpdateCampaign(args[0], payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Campaign updated: %s\n", item.Name)
		return nil
	},
}

// ---- campaign delete ----

var campaignDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a campaign",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteCampaign(args[0]); err != nil {
			return err
		}
		fmt.Printf("Campaign %s deleted.\n", args[0])
		return nil
	},
}

// ---- campaign activate ----

var campaignActivateCmd = &cobra.Command{
	Use:   "activate <id>",
	Short: "Activate a campaign",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.ActivateCampaign(args[0]); err != nil {
			return err
		}
		fmt.Printf("Campaign %s activated.\n", args[0])
		return nil
	},
}

// ---- campaign pause ----

var campaignPauseCmd = &cobra.Command{
	Use:   "pause <id>",
	Short: "Pause a campaign",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.PauseCampaign(args[0]); err != nil {
			return err
		}
		fmt.Printf("Campaign %s paused.\n", args[0])
		return nil
	},
}

// ---- campaign analytics ----

var (
	campaignAnalyticsStartDate string
	campaignAnalyticsEndDate   string
)

var campaignAnalyticsCmd = &cobra.Command{
	Use:   "analytics <id>",
	Short: "Get analytics for a campaign",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"start_date", campaignAnalyticsStartDate,
			"end_date", campaignAnalyticsEndDate,
		)
		item, err := client.GetCampaignAnalytics(args[0], params)
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

// ---- campaign analytics-overview ----

var campaignAnalyticsOverviewCmd = &cobra.Command{
	Use:   "analytics-overview",
	Short: "Get analytics overview for all campaigns",
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := client.GetCampaignAnalyticsOverview(nil)
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
		headers := []string{"CAMPAIGN", "SENT", "OPENED", "REPLIED", "BOUNCED"}
		rows := make([][]string, len(items))
		for i, item := range items {
			rows[i] = []string{
				output.Truncate(item.CampaignName, 36),
				fmt.Sprintf("%d", item.Sent),
				fmt.Sprintf("%d (%.1f%%)", item.Opened, item.OpenRate*100),
				fmt.Sprintf("%d (%.1f%%)", item.Replied, item.ReplyRate*100),
				fmt.Sprintf("%d", item.Bounced),
			}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

// ---- campaign duplicate ----

var campaignDuplicateCmd = &cobra.Command{
	Use:   "duplicate <id>",
	Short: "Duplicate a campaign",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.DuplicateCampaign(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Campaign duplicated: %s\n", item.Name)
		fmt.Printf("New ID: %s\n", item.ID)
		return nil
	},
}

func init() {
	// list flags
	campaignListCmd.Flags().IntVar(&campaignListLimit, "limit", 20, "Maximum number of campaigns to return")
	campaignListCmd.Flags().StringVar(&campaignListStartingAfter, "starting-after", "", "Cursor for pagination: ID of the last campaign from the previous page")
	campaignListCmd.Flags().StringVar(&campaignListSearch, "search", "", "Filter by name")
	campaignListCmd.Flags().StringVar(&campaignListStatus, "status", "", "Filter by status (active, paused, completed, draft)")

	// create flags
	campaignCreateCmd.Flags().StringVar(&campaignCreateTimezone, "timezone", "", "Campaign timezone (e.g. America/New_York)")
	campaignCreateCmd.Flags().StringVar(&campaignCreateStartHour, "start-hour", "08:00", "Daily send start hour (HH:MM)")
	campaignCreateCmd.Flags().StringVar(&campaignCreateEndHour, "end-hour", "18:00", "Daily send end hour (HH:MM)")
	campaignCreateCmd.Flags().IntVar(&campaignCreateDailyLimit, "daily-limit", 0, "Max emails per day")
	campaignCreateCmd.Flags().StringArrayVar(&campaignCreateAccounts, "accounts", nil, "Email accounts to use (comma-separated or repeat flag)")
	campaignCreateCmd.Flags().BoolVar(&campaignCreateStopOnReply, "stop-on-reply", true, "Stop sending when lead replies")
	campaignCreateCmd.Flags().BoolVar(&campaignCreateStopOnAutoReply, "stop-on-auto-reply", true, "Stop sending on auto-reply")
	campaignCreateCmd.Flags().BoolVar(&campaignCreateLinkTracking, "link-tracking", false, "Enable link tracking")
	campaignCreateCmd.Flags().BoolVar(&campaignCreateOpenTracking, "open-tracking", false, "Enable open tracking")
	campaignCreateCmd.Flags().BoolVar(&campaignCreateTextOnly, "text-only", false, "Send plain text emails only")

	// update flags
	campaignUpdateCmd.Flags().StringVar(&campaignUpdateName, "name", "", "New campaign name")
	campaignUpdateCmd.Flags().IntVar(&campaignUpdateDailyLimit, "daily-limit", 0, "New daily email limit")
	campaignUpdateCmd.Flags().StringVar(&campaignUpdateTimezone, "timezone", "", "New campaign timezone")

	// analytics flags
	campaignAnalyticsCmd.Flags().StringVar(&campaignAnalyticsStartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	campaignAnalyticsCmd.Flags().StringVar(&campaignAnalyticsEndDate, "end-date", "", "End date (YYYY-MM-DD)")

	campaignCmd.AddCommand(
		campaignListCmd,
		campaignGetCmd,
		campaignCreateCmd,
		campaignUpdateCmd,
		campaignDeleteCmd,
		campaignActivateCmd,
		campaignPauseCmd,
		campaignAnalyticsCmd,
		campaignAnalyticsOverviewCmd,
		campaignDuplicateCmd,
	)
	rootCmd.AddCommand(campaignCmd)
}
