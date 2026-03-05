package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage Instantly email accounts",
}

// ---- account list ----

var (
	accountListLimit         int
	accountListStartingAfter string
	accountListSearch        string
	accountListStatus        string
)

var accountListCmd = &cobra.Command{
	Use:   "list",
	Short: "List email accounts",
	Long: `List Instantly email accounts.

Examples:
  instantly account list
  instantly account list --status active
  instantly account list --starting-after <id>
  instantly account list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"limit", fmt.Sprintf("%d", accountListLimit),
			"starting_after", accountListStartingAfter,
			"search", accountListSearch,
			"status", accountListStatus,
		)
		items, _, err := client.ListAccounts(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printAccountsTable(items)
		return nil
	},
}

// ---- account get ----

var accountGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details of a specific email account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetAccount(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"Email", item.Email},
			{"Name", item.FirstName + " " + item.LastName},
			{"Status", accountStatusLabel(item.Status)},
			{"Daily Limit", fmt.Sprintf("%d", item.DailyLimit)},
			{"Warmup Enabled", output.FormatBool(item.WarmupEnabled)},
			{"Warmup Limit", fmt.Sprintf("%d", item.WarmupLimit)},
			{"Warmup Reply Rate", fmt.Sprintf("%d%%", item.WarmupReplyRate)},
			{"SMTP Host", item.SmtpHost},
			{"SMTP Port", fmt.Sprintf("%d", item.SmtpPort)},
			{"IMAP Host", item.ImapHost},
			{"IMAP Port", fmt.Sprintf("%d", item.ImapPort)},
			{"Alias", item.Alias},
			{"Tracking Domain", item.TrackingDomain},
			{"Created", output.FormatTime(item.CreatedAt)},
			{"Updated", output.FormatTime(item.UpdatedAt)},
		})
		return nil
	},
}

// ---- account create ----

var (
	accountCreateFirstName       string
	accountCreateLastName        string
	accountCreateSmtpHost        string
	accountCreateSmtpPort        int
	accountCreateSmtpUser        string
	accountCreateSmtpPass        string
	accountCreateImapHost        string
	accountCreateImapPort        int
	accountCreateImapUser        string
	accountCreateImapPass        string
	accountCreateDailyLimit      int
	accountCreateAlias           string
	accountCreateTrackingDomain  string
)

var accountCreateCmd = &cobra.Command{
	Use:   "create <email>",
	Short: "Create/add a new email account",
	Args:  cobra.ExactArgs(1),
	Long: `Add a new email account to Instantly.

Examples:
  instantly account create user@domain.com --first-name John --last-name Doe \
    --smtp-host smtp.domain.com --smtp-port 587 --smtp-user user@domain.com --smtp-pass secret \
    --imap-host imap.domain.com --imap-port 993 --imap-user user@domain.com --imap-pass secret`,
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{
			"email": args[0],
		}
		if accountCreateFirstName != "" {
			payload["first_name"] = accountCreateFirstName
		}
		if accountCreateLastName != "" {
			payload["last_name"] = accountCreateLastName
		}
		if accountCreateSmtpHost != "" {
			payload["smtp_host"] = accountCreateSmtpHost
			payload["smtp_port"] = accountCreateSmtpPort
			payload["smtp_username"] = accountCreateSmtpUser
			payload["smtp_password"] = accountCreateSmtpPass
		}
		if accountCreateImapHost != "" {
			payload["imap_host"] = accountCreateImapHost
			payload["imap_port"] = accountCreateImapPort
			payload["imap_username"] = accountCreateImapUser
			payload["imap_password"] = accountCreateImapPass
		}
		if accountCreateDailyLimit > 0 {
			payload["daily_limit"] = accountCreateDailyLimit
		}
		if accountCreateAlias != "" {
			payload["alias"] = accountCreateAlias
		}
		if accountCreateTrackingDomain != "" {
			payload["tracking_domain_name"] = accountCreateTrackingDomain
		}
		item, err := client.CreateAccount(payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Account created: %s\n", item.Email)
		fmt.Printf("ID: %s\n", item.ID)
		return nil
	},
}

// ---- account update ----

var (
	accountUpdateFirstName      string
	accountUpdateLastName       string
	accountUpdateDailyLimit     int
	accountUpdateWarmupEnabled  bool
	accountUpdateWarmupLimit    int
	accountUpdateAlias          string
	accountUpdateTrackingDomain string
)

var accountUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an email account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{}
		if cmd.Flags().Changed("first-name") {
			payload["first_name"] = accountUpdateFirstName
		}
		if cmd.Flags().Changed("last-name") {
			payload["last_name"] = accountUpdateLastName
		}
		if cmd.Flags().Changed("daily-limit") {
			payload["daily_limit"] = accountUpdateDailyLimit
		}
		if cmd.Flags().Changed("warmup-enabled") {
			payload["warmup_enabled"] = accountUpdateWarmupEnabled
		}
		if cmd.Flags().Changed("warmup-limit") {
			payload["warmup_limit"] = accountUpdateWarmupLimit
		}
		if cmd.Flags().Changed("alias") {
			payload["alias"] = accountUpdateAlias
		}
		if cmd.Flags().Changed("tracking-domain") {
			payload["tracking_domain_name"] = accountUpdateTrackingDomain
		}
		if len(payload) == 0 {
			return fmt.Errorf("no fields to update — provide at least one flag")
		}
		item, err := client.UpdateAccount(args[0], payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Account updated: %s\n", item.Email)
		return nil
	},
}

// ---- account delete ----

var accountDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an email account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteAccount(args[0]); err != nil {
			return err
		}
		fmt.Printf("Account %s deleted.\n", args[0])
		return nil
	},
}

// ---- account warmup ----

var accountWarmupCmd = &cobra.Command{
	Use:   "warmup",
	Short: "Manage account warmup",
}

var accountWarmupEnableCmd = &cobra.Command{
	Use:   "enable <email> [email2...]",
	Short: "Enable warmup for one or more accounts",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var emails []string
		for _, a := range args {
			for _, e := range strings.Split(a, ",") {
				if e = strings.TrimSpace(e); e != "" {
					emails = append(emails, e)
				}
			}
		}
		if err := client.EnableWarmup(emails); err != nil {
			return err
		}
		fmt.Printf("Warmup enabled for %d account(s).\n", len(emails))
		return nil
	},
}

var accountWarmupDisableCmd = &cobra.Command{
	Use:   "disable <email> [email2...]",
	Short: "Disable warmup for one or more accounts",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var emails []string
		for _, a := range args {
			for _, e := range strings.Split(a, ",") {
				if e = strings.TrimSpace(e); e != "" {
					emails = append(emails, e)
				}
			}
		}
		if err := client.DisableWarmup(emails); err != nil {
			return err
		}
		fmt.Printf("Warmup disabled for %d account(s).\n", len(emails))
		return nil
	},
}

var (
	accountWarmupAnalyticsEmail     string
	accountWarmupAnalyticsStartDate string
	accountWarmupAnalyticsEndDate   string
)

var accountWarmupAnalyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "Get warmup analytics",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"email", accountWarmupAnalyticsEmail,
			"start_date", accountWarmupAnalyticsStartDate,
			"end_date", accountWarmupAnalyticsEndDate,
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
				item.Email,
				output.FormatTime(item.Date),
				fmt.Sprintf("%d", item.Sent),
				fmt.Sprintf("%d", item.Received),
			}
		}
		output.PrintTable(headers, rows)
		return nil
	},
}

// ---- account pause ----

var accountPauseCmd = &cobra.Command{
	Use:   "pause <email>",
	Short: "Pause an email account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.PauseAccount(args[0]); err != nil {
			return err
		}
		fmt.Printf("Account %s paused.\n", args[0])
		return nil
	},
}

// ---- account resume ----

var accountResumeCmd = &cobra.Command{
	Use:   "resume <email>",
	Short: "Resume a paused email account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.ResumeAccount(args[0]); err != nil {
			return err
		}
		fmt.Printf("Account %s resumed.\n", args[0])
		return nil
	},
}

func init() {
	// list flags
	accountListCmd.Flags().IntVar(&accountListLimit, "limit", 20, "Maximum number of accounts to return")
	accountListCmd.Flags().StringVar(&accountListStartingAfter, "starting-after", "", "Cursor for pagination: ID of the last account from the previous page")
	accountListCmd.Flags().StringVar(&accountListSearch, "search", "", "Filter by email or name")
	accountListCmd.Flags().StringVar(&accountListStatus, "status", "", "Filter by status")

	// create flags
	accountCreateCmd.Flags().StringVar(&accountCreateFirstName, "first-name", "", "First name")
	accountCreateCmd.Flags().StringVar(&accountCreateLastName, "last-name", "", "Last name")
	accountCreateCmd.Flags().StringVar(&accountCreateSmtpHost, "smtp-host", "", "SMTP host")
	accountCreateCmd.Flags().IntVar(&accountCreateSmtpPort, "smtp-port", 587, "SMTP port")
	accountCreateCmd.Flags().StringVar(&accountCreateSmtpUser, "smtp-user", "", "SMTP username")
	accountCreateCmd.Flags().StringVar(&accountCreateSmtpPass, "smtp-pass", "", "SMTP password")
	accountCreateCmd.Flags().StringVar(&accountCreateImapHost, "imap-host", "", "IMAP host")
	accountCreateCmd.Flags().IntVar(&accountCreateImapPort, "imap-port", 993, "IMAP port")
	accountCreateCmd.Flags().StringVar(&accountCreateImapUser, "imap-user", "", "IMAP username")
	accountCreateCmd.Flags().StringVar(&accountCreateImapPass, "imap-pass", "", "IMAP password")
	accountCreateCmd.Flags().IntVar(&accountCreateDailyLimit, "daily-limit", 0, "Max emails per day")
	accountCreateCmd.Flags().StringVar(&accountCreateAlias, "alias", "", "Display alias name")
	accountCreateCmd.Flags().StringVar(&accountCreateTrackingDomain, "tracking-domain", "", "Custom tracking domain")

	// update flags
	accountUpdateCmd.Flags().StringVar(&accountUpdateFirstName, "first-name", "", "New first name")
	accountUpdateCmd.Flags().StringVar(&accountUpdateLastName, "last-name", "", "New last name")
	accountUpdateCmd.Flags().IntVar(&accountUpdateDailyLimit, "daily-limit", 0, "New daily limit")
	accountUpdateCmd.Flags().BoolVar(&accountUpdateWarmupEnabled, "warmup-enabled", false, "Enable/disable warmup")
	accountUpdateCmd.Flags().IntVar(&accountUpdateWarmupLimit, "warmup-limit", 0, "New warmup limit")
	accountUpdateCmd.Flags().StringVar(&accountUpdateAlias, "alias", "", "New alias")
	accountUpdateCmd.Flags().StringVar(&accountUpdateTrackingDomain, "tracking-domain", "", "New tracking domain")

	// warmup analytics flags
	accountWarmupAnalyticsCmd.Flags().StringVar(&accountWarmupAnalyticsEmail, "email", "", "Filter by email")
	accountWarmupAnalyticsCmd.Flags().StringVar(&accountWarmupAnalyticsStartDate, "start-date", "", "Start date (YYYY-MM-DD)")
	accountWarmupAnalyticsCmd.Flags().StringVar(&accountWarmupAnalyticsEndDate, "end-date", "", "End date (YYYY-MM-DD)")

	accountWarmupCmd.AddCommand(accountWarmupEnableCmd, accountWarmupDisableCmd, accountWarmupAnalyticsCmd)
	accountCmd.AddCommand(
		accountListCmd,
		accountGetCmd,
		accountCreateCmd,
		accountUpdateCmd,
		accountDeleteCmd,
		accountWarmupCmd,
		accountPauseCmd,
		accountResumeCmd,
	)
	rootCmd.AddCommand(accountCmd)
}
