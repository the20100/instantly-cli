package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Manage Instantly webhooks",
}

// ---- webhook list ----

var (
	webhookListLimit int
	webhookListSkip  int
)

var webhookListCmd = &cobra.Command{
	Use:   "list",
	Short: "List webhooks",
	Long: `List Instantly webhooks.

Examples:
  instantly webhook list
  instantly webhook list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"limit", fmt.Sprintf("%d", webhookListLimit),
			"skip", fmt.Sprintf("%d", webhookListSkip),
		)
		items, _, err := client.ListWebhooks(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printWebhooksTable(items)
		return nil
	},
}

// ---- webhook get ----

var webhookGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details of a specific webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetWebhook(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"ID", item.ID},
			{"Name", item.Name},
			{"URL", item.URL},
			{"Active", output.FormatBool(item.Active)},
			{"Events", output.FormatLabels(item.EventTypes)},
			{"Created", output.FormatTime(item.CreatedAt)},
			{"Updated", output.FormatTime(item.UpdatedAt)},
		})
		return nil
	},
}

// ---- webhook create ----

var (
	webhookCreateURL    string
	webhookCreateEvents []string
	webhookCreateSecret string
)

var webhookCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new webhook",
	Args:  cobra.ExactArgs(1),
	Long: `Create a new Instantly webhook.

Examples:
  instantly webhook create "My Webhook" --url https://myapp.com/hook --events reply_received,email_sent
  instantly webhook event-types   # see available event types`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if webhookCreateURL == "" {
			return fmt.Errorf("--url is required")
		}
		// flatten comma-separated events
		var events []string
		for _, e := range webhookCreateEvents {
			for _, ev := range strings.Split(e, ",") {
				if ev = strings.TrimSpace(ev); ev != "" {
					events = append(events, ev)
				}
			}
		}
		payload := map[string]interface{}{
			"name": args[0],
			"url":  webhookCreateURL,
		}
		if len(events) > 0 {
			payload["event_types"] = events
		}
		if webhookCreateSecret != "" {
			payload["secret"] = webhookCreateSecret
		}
		item, err := client.CreateWebhook(payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Webhook created: %s\n", item.Name)
		fmt.Printf("ID: %s\n", item.ID)
		return nil
	},
}

// ---- webhook update ----

var (
	webhookUpdateName   string
	webhookUpdateURL    string
	webhookUpdateEvents []string
)

var webhookUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{}
		if cmd.Flags().Changed("name") {
			payload["name"] = webhookUpdateName
		}
		if cmd.Flags().Changed("url") {
			payload["url"] = webhookUpdateURL
		}
		if cmd.Flags().Changed("events") {
			var events []string
			for _, e := range webhookUpdateEvents {
				for _, ev := range strings.Split(e, ",") {
					if ev = strings.TrimSpace(ev); ev != "" {
						events = append(events, ev)
					}
				}
			}
			payload["event_types"] = events
		}
		if len(payload) == 0 {
			return fmt.Errorf("no fields to update — provide at least one flag")
		}
		item, err := client.UpdateWebhook(args[0], payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Webhook updated: %s\n", item.Name)
		return nil
	},
}

// ---- webhook delete ----

var webhookDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteWebhook(args[0]); err != nil {
			return err
		}
		fmt.Printf("Webhook %s deleted.\n", args[0])
		return nil
	},
}

// ---- webhook test ----

var webhookTestCmd = &cobra.Command{
	Use:   "test <id>",
	Short: "Send a test event to a webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.TestWebhook(args[0]); err != nil {
			return err
		}
		fmt.Printf("Test event sent to webhook %s.\n", args[0])
		return nil
	},
}

// ---- webhook resume ----

var webhookResumeCmd = &cobra.Command{
	Use:   "resume <id>",
	Short: "Resume a paused/failed webhook",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.ResumeWebhook(args[0]); err != nil {
			return err
		}
		fmt.Printf("Webhook %s resumed.\n", args[0])
		return nil
	},
}

// ---- webhook event-types ----

var webhookEventTypesCmd = &cobra.Command{
	Use:   "event-types",
	Short: "List all available webhook event types",
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := client.ListWebhookEventTypes()
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		if len(items) == 0 {
			fmt.Println("No event types found.")
			return nil
		}
		headers := []string{"NAME", "DESCRIPTION"}
		rows := make([][]string, len(items))
		for i, item := range items {
			rows[i] = []string{
				item.Name,
				output.Truncate(item.Description, 60),
			}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

func init() {
	webhookListCmd.Flags().IntVar(&webhookListLimit, "limit", 20, "Maximum number of webhooks to return")
	webhookListCmd.Flags().IntVar(&webhookListSkip, "skip", 0, "Number of webhooks to skip")

	webhookCreateCmd.Flags().StringVar(&webhookCreateURL, "url", "", "Webhook URL *(required)*")
	webhookCreateCmd.Flags().StringArrayVar(&webhookCreateEvents, "events", nil, "Event types (comma-separated or repeat flag)")
	webhookCreateCmd.Flags().StringVar(&webhookCreateSecret, "secret", "", "Webhook signing secret")

	webhookUpdateCmd.Flags().StringVar(&webhookUpdateName, "name", "", "New webhook name")
	webhookUpdateCmd.Flags().StringVar(&webhookUpdateURL, "url", "", "New webhook URL")
	webhookUpdateCmd.Flags().StringArrayVar(&webhookUpdateEvents, "events", nil, "New event types")

	webhookCmd.AddCommand(
		webhookListCmd,
		webhookGetCmd,
		webhookCreateCmd,
		webhookUpdateCmd,
		webhookDeleteCmd,
		webhookTestCmd,
		webhookResumeCmd,
		webhookEventTypesCmd,
	)
	rootCmd.AddCommand(webhookCmd)
}
