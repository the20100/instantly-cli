package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var leadCmd = &cobra.Command{
	Use:   "lead",
	Short: "Manage Instantly leads",
}

// ---- lead list ----

var (
	leadListCampaignID string
	leadListID         string
	leadListEmail      string
	leadListStatus     string
	leadListLimit      int
	leadListSkip       int
)

var leadListCmd = &cobra.Command{
	Use:   "list",
	Short: "List leads",
	Long: `List Instantly leads.

Examples:
  instantly lead list --campaign-id <id>
  instantly lead list --list-id <id>
  instantly lead list --status interested
  instantly lead list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"campaign_id", leadListCampaignID,
			"list_id", leadListID,
			"email", leadListEmail,
			"lt_interest_status", leadListStatus,
			"limit", fmt.Sprintf("%d", leadListLimit),
			"skip", fmt.Sprintf("%d", leadListSkip),
		)
		items, _, err := client.ListLeads(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printLeadsTable(items)
		return nil
	},
}

// ---- lead get ----

var leadGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details of a specific lead",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetLead(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"ID", item.ID},
			{"Email", item.Email},
			{"First Name", item.FirstName},
			{"Last Name", item.LastName},
			{"Company", item.CompanyName},
			{"Phone", item.Phone},
			{"Website", item.Website},
			{"LinkedIn", item.LinkedinURL},
			{"Status", item.Status},
			{"Campaign ID", item.CampaignID},
			{"List ID", item.ListID},
			{"Assigned To", item.AssignedTo},
			{"Created", output.FormatTime(item.CreatedAt)},
			{"Updated", output.FormatTime(item.UpdatedAt)},
		})
		return nil
	},
}

// ---- lead create ----

var (
	leadCreateFirstName   string
	leadCreateLastName    string
	leadCreateCompany     string
	leadCreatePhone       string
	leadCreateWebsite     string
	leadCreateLinkedin    string
	leadCreateCampaignID  string
	leadCreateListID      string
)

var leadCreateCmd = &cobra.Command{
	Use:   "create <email>",
	Short: "Create a new lead",
	Args:  cobra.ExactArgs(1),
	Long: `Create a new lead in Instantly.

Examples:
  instantly lead create john@acme.com --first-name John --last-name Doe --company Acme
  instantly lead create jane@corp.com --campaign-id <id> --list-id <id>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{
			"email": args[0],
		}
		if leadCreateFirstName != "" {
			payload["first_name"] = leadCreateFirstName
		}
		if leadCreateLastName != "" {
			payload["last_name"] = leadCreateLastName
		}
		if leadCreateCompany != "" {
			payload["company_name"] = leadCreateCompany
		}
		if leadCreatePhone != "" {
			payload["phone"] = leadCreatePhone
		}
		if leadCreateWebsite != "" {
			payload["website"] = leadCreateWebsite
		}
		if leadCreateLinkedin != "" {
			payload["linkedin_url"] = leadCreateLinkedin
		}
		if leadCreateCampaignID != "" {
			payload["campaign_id"] = leadCreateCampaignID
		}
		if leadCreateListID != "" {
			payload["list_id"] = leadCreateListID
		}
		item, err := client.CreateLead(payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Lead created: %s\n", item.Email)
		fmt.Printf("ID: %s\n", item.ID)
		return nil
	},
}

// ---- lead update ----

var (
	leadUpdateFirstName  string
	leadUpdateLastName   string
	leadUpdateCompany    string
	leadUpdatePhone      string
	leadUpdateWebsite    string
	leadUpdateLinkedin   string
)

var leadUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a lead",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		payload := map[string]interface{}{}
		if cmd.Flags().Changed("first-name") {
			payload["first_name"] = leadUpdateFirstName
		}
		if cmd.Flags().Changed("last-name") {
			payload["last_name"] = leadUpdateLastName
		}
		if cmd.Flags().Changed("company") {
			payload["company_name"] = leadUpdateCompany
		}
		if cmd.Flags().Changed("phone") {
			payload["phone"] = leadUpdatePhone
		}
		if cmd.Flags().Changed("website") {
			payload["website"] = leadUpdateWebsite
		}
		if cmd.Flags().Changed("linkedin") {
			payload["linkedin_url"] = leadUpdateLinkedin
		}
		if len(payload) == 0 {
			return fmt.Errorf("no fields to update — provide at least one flag")
		}
		item, err := client.UpdateLead(args[0], payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Lead updated: %s\n", item.Email)
		return nil
	},
}

// ---- lead delete ----

var leadDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a lead",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteLead(args[0]); err != nil {
			return err
		}
		fmt.Printf("Lead %s deleted.\n", args[0])
		return nil
	},
}

// ---- lead update-interest ----

var leadUpdateInterestCmd = &cobra.Command{
	Use:   "update-interest <id>",
	Short: "Update the interest status of a lead",
	Long: `Update lead interest status.

Valid statuses: interested, not_interested, meeting_booked, meeting_completed, closed

Examples:
  instantly lead update-interest <id> --status interested
  instantly lead update-interest <id> --status meeting_booked`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		if status == "" {
			return fmt.Errorf("--status is required")
		}
		item, err := client.UpdateLeadInterest(args[0], status)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Lead %s interest status updated to: %s\n", item.Email, item.Status)
		return nil
	},
}

func init() {
	// list flags
	leadListCmd.Flags().StringVar(&leadListCampaignID, "campaign-id", "", "Filter by campaign ID")
	leadListCmd.Flags().StringVar(&leadListID, "list-id", "", "Filter by lead list ID")
	leadListCmd.Flags().StringVar(&leadListEmail, "email", "", "Filter by email address")
	leadListCmd.Flags().StringVar(&leadListStatus, "status", "", "Filter by interest status")
	leadListCmd.Flags().IntVar(&leadListLimit, "limit", 20, "Maximum number of leads to return")
	leadListCmd.Flags().IntVar(&leadListSkip, "skip", 0, "Number of leads to skip")

	// create flags
	leadCreateCmd.Flags().StringVar(&leadCreateFirstName, "first-name", "", "First name")
	leadCreateCmd.Flags().StringVar(&leadCreateLastName, "last-name", "", "Last name")
	leadCreateCmd.Flags().StringVar(&leadCreateCompany, "company", "", "Company name")
	leadCreateCmd.Flags().StringVar(&leadCreatePhone, "phone", "", "Phone number")
	leadCreateCmd.Flags().StringVar(&leadCreateWebsite, "website", "", "Website URL")
	leadCreateCmd.Flags().StringVar(&leadCreateLinkedin, "linkedin", "", "LinkedIn URL")
	leadCreateCmd.Flags().StringVar(&leadCreateCampaignID, "campaign-id", "", "Campaign to add lead to")
	leadCreateCmd.Flags().StringVar(&leadCreateListID, "list-id", "", "Lead list to add lead to")

	// update flags
	leadUpdateCmd.Flags().StringVar(&leadUpdateFirstName, "first-name", "", "New first name")
	leadUpdateCmd.Flags().StringVar(&leadUpdateLastName, "last-name", "", "New last name")
	leadUpdateCmd.Flags().StringVar(&leadUpdateCompany, "company", "", "New company name")
	leadUpdateCmd.Flags().StringVar(&leadUpdatePhone, "phone", "", "New phone number")
	leadUpdateCmd.Flags().StringVar(&leadUpdateWebsite, "website", "", "New website URL")
	leadUpdateCmd.Flags().StringVar(&leadUpdateLinkedin, "linkedin", "", "New LinkedIn URL")

	// update-interest flags
	leadUpdateInterestCmd.Flags().String("status", "", "Interest status (interested, not_interested, meeting_booked, meeting_completed, closed) *(required)*")

	leadCmd.AddCommand(
		leadListCmd,
		leadGetCmd,
		leadCreateCmd,
		leadUpdateCmd,
		leadDeleteCmd,
		leadUpdateInterestCmd,
	)
	rootCmd.AddCommand(leadCmd)
}
