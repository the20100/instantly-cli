package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var emailverificationCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify email addresses",
}

// ---- verify create ----

var emailverificationCreateCmd = &cobra.Command{
	Use:   "create <email>",
	Short: "Start an email verification",
	Args:  cobra.ExactArgs(1),
	Long: `Start an email address verification job.

Examples:
  instantly verify create user@domain.com
  instantly verify create user@domain.com --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.CreateEmailVerification(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Verification started for: %s\n", item.Email)
		fmt.Printf("ID:     %s\n", item.ID)
		fmt.Printf("Status: %s\n", item.Status)
		return nil
	},
}

// ---- verify check ----

var emailverificationCheckCmd = &cobra.Command{
	Use:   "check <id>",
	Short: "Check the status of an email verification",
	Args:  cobra.ExactArgs(1),
	Long: `Check the status of an email verification job.

Examples:
  instantly verify check <id>
  instantly verify check <id> --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetEmailVerification(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"ID", item.ID},
			{"Email", item.Email},
			{"Status", item.Status},
			{"Valid", output.FormatBool(item.Valid)},
		})
		return nil
	},
}

func init() {
	emailverificationCmd.AddCommand(
		emailverificationCreateCmd,
		emailverificationCheckCmd,
	)
	rootCmd.AddCommand(emailverificationCmd)
}
