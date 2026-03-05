package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var emailCmd = &cobra.Command{
	Use:   "email",
	Short: "Manage Instantly emails",
}

// ---- email list ----

var (
	emailListCampaignID  string
	emailListLeadEmail   string
	emailListAccountEmail string
	emailListType        string
	emailListIsUnread    bool
	emailListLimit       int
	emailListStartAfter  string
)

var emailListCmd = &cobra.Command{
	Use:   "list",
	Short: "List emails",
	Long: `List Instantly emails.

Examples:
  instantly email list --campaign-id <id>
  instantly email list --lead-email john@acme.com
  instantly email list --type reply --is-unread
  instantly email list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"campaign_id", emailListCampaignID,
			"lead_email", emailListLeadEmail,
			"account_email", emailListAccountEmail,
			"type", emailListType,
			"limit", fmt.Sprintf("%d", emailListLimit),
			"starting_after", emailListStartAfter,
		)
		if cmd.Flags().Changed("is-unread") {
			if emailListIsUnread {
				params.Set("is_unread", "true")
			}
		}
		items, _, err := client.ListEmails(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printEmailsTable(items)
		return nil
	},
}

// ---- email get ----

var emailGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details of a specific email",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetEmail(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"ID", item.ID},
			{"From", item.FromEmail},
			{"To", item.ToEmail},
			{"Subject", item.Subject},
			{"Type", item.Type},
			{"Campaign ID", item.CampaignID},
			{"Thread ID", item.ThreadID},
			{"Is Read", output.FormatBool(item.IsRead)},
			{"Timestamp", output.FormatTime(item.Timestamp)},
		})
		if item.Body != "" {
			fmt.Println()
			fmt.Println("Body:")
			fmt.Println(item.Body)
		}
		return nil
	},
}

// ---- email reply ----

var (
	emailReplySubject string
	emailReplyBody    string
	emailReplyCC      string
)

var emailReplyCmd = &cobra.Command{
	Use:   "reply <email-id>",
	Short: "Reply to an email",
	Args:  cobra.ExactArgs(1),
	Long: `Reply to an existing email thread.

Examples:
  instantly email reply <id> --body "Thanks for reaching out!"
  instantly email reply <id> --subject "Re: Your inquiry" --body "Hello..."`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if emailReplyBody == "" {
			return fmt.Errorf("--body is required")
		}
		payload := map[string]interface{}{
			"email_id": args[0],
			"body":     emailReplyBody,
		}
		if emailReplySubject != "" {
			payload["subject"] = emailReplySubject
		}
		if emailReplyCC != "" {
			payload["cc"] = emailReplyCC
		}
		if err := client.ReplyToEmail(payload); err != nil {
			return err
		}
		fmt.Printf("Reply sent to email %s.\n", args[0])
		return nil
	},
}

// ---- email forward ----

var (
	emailForwardTo   string
	emailForwardBody string
)

var emailForwardCmd = &cobra.Command{
	Use:   "forward <email-id>",
	Short: "Forward an email",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if emailForwardTo == "" {
			return fmt.Errorf("--to is required")
		}
		payload := map[string]interface{}{
			"email_id": args[0],
			"to":       emailForwardTo,
		}
		if emailForwardBody != "" {
			payload["body"] = emailForwardBody
		}
		if err := client.ForwardEmail(payload); err != nil {
			return err
		}
		fmt.Printf("Email %s forwarded to %s.\n", args[0], emailForwardTo)
		return nil
	},
}

// ---- email mark-read ----

var emailMarkReadCmd = &cobra.Command{
	Use:   "mark-read <thread-id>",
	Short: "Mark an email thread as read",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.MarkThreadAsRead(args[0]); err != nil {
			return err
		}
		fmt.Printf("Thread %s marked as read.\n", args[0])
		return nil
	},
}

func init() {
	// list flags
	emailListCmd.Flags().StringVar(&emailListCampaignID, "campaign-id", "", "Filter by campaign ID")
	emailListCmd.Flags().StringVar(&emailListLeadEmail, "lead-email", "", "Filter by lead email")
	emailListCmd.Flags().StringVar(&emailListAccountEmail, "account-email", "", "Filter by sending account email")
	emailListCmd.Flags().StringVar(&emailListType, "type", "", "Filter by type (sent, received, reply)")
	emailListCmd.Flags().BoolVar(&emailListIsUnread, "is-unread", false, "Show only unread emails")
	emailListCmd.Flags().IntVar(&emailListLimit, "limit", 20, "Maximum number of emails to return")
	emailListCmd.Flags().StringVar(&emailListStartAfter, "starting-after", "", "Cursor for pagination")

	// reply flags
	emailReplyCmd.Flags().StringVar(&emailReplySubject, "subject", "", "Reply subject")
	emailReplyCmd.Flags().StringVar(&emailReplyBody, "body", "", "Reply body *(required)*")
	emailReplyCmd.Flags().StringVar(&emailReplyCC, "cc", "", "CC email address")

	// forward flags
	emailForwardCmd.Flags().StringVar(&emailForwardTo, "to", "", "Recipient email *(required)*")
	emailForwardCmd.Flags().StringVar(&emailForwardBody, "body", "", "Additional message body")

	emailCmd.AddCommand(
		emailListCmd,
		emailGetCmd,
		emailReplyCmd,
		emailForwardCmd,
		emailMarkReadCmd,
	)
	rootCmd.AddCommand(emailCmd)
}
