package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var subsequenceCmd = &cobra.Command{
	Use:   "subsequence",
	Short: "Manage Instantly campaign subsequences",
}

// ---- subsequence list ----

var (
	subsequenceListCampaignID string
	subsequenceListLimit      int
	subsequenceListSkip       int
)

var subsequenceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List campaign subsequences",
	Long: `List Instantly campaign subsequences.

Examples:
  instantly subsequence list --campaign-id <id>
  instantly subsequence list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"campaign_id", subsequenceListCampaignID,
			"limit", fmt.Sprintf("%d", subsequenceListLimit),
			"skip", fmt.Sprintf("%d", subsequenceListSkip),
		)
		items, _, err := client.ListCampaignSubsequences(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printCampaignSubsequencesTable(items)
		return nil
	},
}

// ---- subsequence get ----

var subsequenceGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details of a specific subsequence",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetCampaignSubsequence(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"ID", item.ID},
			{"Name", item.Name},
			{"Campaign ID", item.CampaignID},
			{"Status", campaignStatusLabel(item.Status)},
			{"Created", output.FormatTime(item.CreatedAt)},
			{"Updated", output.FormatTime(item.UpdatedAt)},
		})
		return nil
	},
}

// ---- subsequence create ----

var (
	subsequenceCreateCampaignID string
	subsequenceCreateName       string
)

var subsequenceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new campaign subsequence",
	Long: `Create a new campaign subsequence.

Examples:
  instantly subsequence create --campaign-id <id> --name "Follow-up Sequence"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if subsequenceCreateCampaignID == "" {
			return fmt.Errorf("--campaign-id is required")
		}
		if subsequenceCreateName == "" {
			return fmt.Errorf("--name is required")
		}
		payload := map[string]interface{}{
			"campaign_id": subsequenceCreateCampaignID,
			"name":        subsequenceCreateName,
		}
		item, err := client.CreateCampaignSubsequence(payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Subsequence created: %s\n", item.Name)
		fmt.Printf("ID: %s\n", item.ID)
		return nil
	},
}

// ---- subsequence update ----

var subsequenceUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a campaign subsequence",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name is required")
		}
		item, err := client.UpdateCampaignSubsequence(args[0], map[string]interface{}{"name": name})
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Subsequence updated: %s\n", item.Name)
		return nil
	},
}

// ---- subsequence delete ----

var subsequenceDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a campaign subsequence",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteCampaignSubsequence(args[0]); err != nil {
			return err
		}
		fmt.Printf("Subsequence %s deleted.\n", args[0])
		return nil
	},
}

// ---- subsequence pause ----

var subsequencePauseCmd = &cobra.Command{
	Use:   "pause <id>",
	Short: "Pause a campaign subsequence",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.PauseCampaignSubsequence(args[0]); err != nil {
			return err
		}
		fmt.Printf("Subsequence %s paused.\n", args[0])
		return nil
	},
}

// ---- subsequence resume ----

var subsequenceResumeCmd = &cobra.Command{
	Use:   "resume <id>",
	Short: "Resume a paused campaign subsequence",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.ResumeCampaignSubsequence(args[0]); err != nil {
			return err
		}
		fmt.Printf("Subsequence %s resumed.\n", args[0])
		return nil
	},
}

// ---- subsequence duplicate ----

var subsequenceDuplicateCmd = &cobra.Command{
	Use:   "duplicate <id>",
	Short: "Duplicate a campaign subsequence",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.DuplicateCampaignSubsequence(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Subsequence duplicated: %s\n", item.Name)
		fmt.Printf("New ID: %s\n", item.ID)
		return nil
	},
}

func init() {
	subsequenceListCmd.Flags().StringVar(&subsequenceListCampaignID, "campaign-id", "", "Filter by campaign ID")
	subsequenceListCmd.Flags().IntVar(&subsequenceListLimit, "limit", 20, "Maximum number of subsequences to return")
	subsequenceListCmd.Flags().IntVar(&subsequenceListSkip, "skip", 0, "Number of subsequences to skip")

	subsequenceCreateCmd.Flags().StringVar(&subsequenceCreateCampaignID, "campaign-id", "", "Campaign ID *(required)*")
	subsequenceCreateCmd.Flags().StringVar(&subsequenceCreateName, "name", "", "Subsequence name *(required)*")

	subsequenceUpdateCmd.Flags().String("name", "", "New subsequence name *(required)*")

	subsequenceCmd.AddCommand(
		subsequenceListCmd,
		subsequenceGetCmd,
		subsequenceCreateCmd,
		subsequenceUpdateCmd,
		subsequenceDeleteCmd,
		subsequencePauseCmd,
		subsequenceResumeCmd,
		subsequenceDuplicateCmd,
	)
	rootCmd.AddCommand(subsequenceCmd)
}
