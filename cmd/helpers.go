package cmd

import (
	"fmt"
	"net/url"

	"github.com/the20100/instantly-cli/internal/api"
	"github.com/the20100/instantly-cli/internal/output"
)

// buildParams creates a url.Values from alternating key/value pairs,
// skipping pairs where the value is empty.
func buildParams(pairs ...string) url.Values {
	p := url.Values{}
	for i := 0; i+1 < len(pairs); i += 2 {
		if pairs[i+1] != "" {
			p.Set(pairs[i], pairs[i+1])
		}
	}
	return p
}

// campaignStatusLabel returns a human-readable campaign status.
func campaignStatusLabel(status int) string {
	switch status {
	case 0:
		return "draft"
	case 1:
		return "active"
	case 2:
		return "paused"
	case 3:
		return "completed"
	case 4:
		return "error"
	default:
		return fmt.Sprintf("%d", status)
	}
}

// accountStatusLabel returns a human-readable account status.
func accountStatusLabel(status int) string {
	switch status {
	case 1:
		return "active"
	case 2:
		return "paused"
	case -1:
		return "error"
	default:
		return fmt.Sprintf("%d", status)
	}
}

// printCampaignsTable renders campaigns as a table.
func printCampaignsTable(items []api.Campaign) {
	if len(items) == 0 {
		fmt.Println("No campaigns found.")
		return
	}
	headers := []string{"ID", "NAME", "STATUS", "DAILY LIMIT", "CREATED"}
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = []string{
			item.ID,
			output.Truncate(item.Name, 44),
			campaignStatusLabel(item.Status),
			fmt.Sprintf("%d", item.DailyLimit),
			output.FormatTime(item.CreatedAt),
		}
	}
	output.PrintTable(headers, rows)
}

// printAccountsTable renders accounts as a table.
func printAccountsTable(items []api.Account) {
	if len(items) == 0 {
		fmt.Println("No accounts found.")
		return
	}
	headers := []string{"EMAIL", "NAME", "STATUS", "DAILY LIMIT", "WARMUP"}
	rows := make([][]string, len(items))
	for i, item := range items {
		name := item.FirstName + " " + item.LastName
		rows[i] = []string{
			item.Email,
			output.Truncate(name, 30),
			accountStatusLabel(item.Status),
			fmt.Sprintf("%d", item.DailyLimit),
			output.FormatBool(item.WarmupEnabled),
		}
	}
	output.PrintTable(headers, rows)
}

// printLeadsTable renders leads as a table.
func printLeadsTable(items []api.Lead) {
	if len(items) == 0 {
		fmt.Println("No leads found.")
		return
	}
	headers := []string{"ID", "EMAIL", "NAME", "COMPANY", "STATUS"}
	rows := make([][]string, len(items))
	for i, item := range items {
		name := item.FirstName + " " + item.LastName
		rows[i] = []string{
			item.ID,
			output.Truncate(item.Email, 36),
			output.Truncate(name, 26),
			output.Truncate(item.CompanyName, 26),
			item.Status,
		}
	}
	output.PrintTable(headers, rows)
}

// printLeadListsTable renders lead lists as a table.
func printLeadListsTable(items []api.LeadList) {
	if len(items) == 0 {
		fmt.Println("No lead lists found.")
		return
	}
	headers := []string{"ID", "NAME", "COUNT", "CREATED"}
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = []string{
			item.ID,
			output.Truncate(item.Name, 44),
			fmt.Sprintf("%d", item.Count),
			output.FormatTime(item.CreatedAt),
		}
	}
	output.PrintTable(headers, rows)
}

// printEmailsTable renders emails as a table.
func printEmailsTable(items []api.Email) {
	if len(items) == 0 {
		fmt.Println("No emails found.")
		return
	}
	headers := []string{"ID", "FROM", "TO", "SUBJECT", "TYPE", "TIME"}
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = []string{
			output.Truncate(item.ID, 24),
			output.Truncate(item.FromEmail, 28),
			output.Truncate(item.ToEmail, 28),
			output.Truncate(item.Subject, 40),
			item.Type,
			output.FormatTime(item.Timestamp),
		}
	}
	output.PrintTable(headers, rows)
}

// printWebhooksTable renders webhooks as a table.
func printWebhooksTable(items []api.Webhook) {
	if len(items) == 0 {
		fmt.Println("No webhooks found.")
		return
	}
	headers := []string{"ID", "NAME", "URL", "ACTIVE", "EVENTS"}
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = []string{
			item.ID,
			output.Truncate(item.Name, 24),
			output.Truncate(item.URL, 44),
			output.FormatBool(item.Active),
			fmt.Sprintf("%d", len(item.EventTypes)),
		}
	}
	output.PrintTable(headers, rows)
}

// printCustomTagsTable renders custom tags as a table.
func printCustomTagsTable(items []api.CustomTag) {
	if len(items) == 0 {
		fmt.Println("No custom tags found.")
		return
	}
	headers := []string{"ID", "NAME", "COLOR", "CREATED"}
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = []string{
			item.ID,
			output.Truncate(item.Name, 30),
			item.Color,
			output.FormatTime(item.CreatedAt),
		}
	}
	output.PrintTable(headers, rows)
}

// printBlocklistEntriesTable renders blocklist entries as a table.
func printBlocklistEntriesTable(items []api.BlocklistEntry) {
	if len(items) == 0 {
		fmt.Println("No blocklist entries found.")
		return
	}
	headers := []string{"ID", "VALUE", "TYPE", "CREATED"}
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = []string{
			item.ID,
			output.Truncate(item.Value, 44),
			item.Type,
			output.FormatTime(item.CreatedAt),
		}
	}
	output.PrintTable(headers, rows)
}

// printAPIKeysTable renders API keys as a table.
func printAPIKeysTable(items []api.APIKey) {
	if len(items) == 0 {
		fmt.Println("No API keys found.")
		return
	}
	headers := []string{"ID", "NAME", "KEY (MASKED)", "CREATED"}
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = []string{
			item.ID,
			output.Truncate(item.Name, 30),
			maskOrEmpty(item.APIKey),
			output.FormatTime(item.CreatedAt),
		}
	}
	output.PrintTable(headers, rows)
}

// printWorkspaceMembersTable renders workspace members as a table.
func printWorkspaceMembersTable(items []api.WorkspaceMember) {
	if len(items) == 0 {
		fmt.Println("No workspace members found.")
		return
	}
	headers := []string{"ID", "EMAIL", "ROLE", "STATUS", "CREATED"}
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = []string{
			item.ID,
			output.Truncate(item.Email, 36),
			item.Role,
			item.Status,
			output.FormatTime(item.CreatedAt),
		}
	}
	output.PrintTable(headers, rows)
}

// printBackgroundJobsTable renders background jobs as a table.
func printBackgroundJobsTable(items []api.BackgroundJob) {
	if len(items) == 0 {
		fmt.Println("No background jobs found.")
		return
	}
	headers := []string{"ID", "TYPE", "STATUS", "PROGRESS", "CREATED"}
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = []string{
			item.ID,
			output.Truncate(item.Type, 30),
			item.Status,
			fmt.Sprintf("%.0f%%", item.Progress*100),
			output.FormatTime(item.CreatedAt),
		}
	}
	output.PrintTable(headers, rows)
}

// printLeadLabelsTable renders lead labels as a table.
func printLeadLabelsTable(items []api.LeadLabel) {
	if len(items) == 0 {
		fmt.Println("No lead labels found.")
		return
	}
	headers := []string{"ID", "NAME", "COLOR", "CREATED"}
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = []string{
			item.ID,
			output.Truncate(item.Name, 30),
			item.Color,
			output.FormatTime(item.CreatedAt),
		}
	}
	output.PrintTable(headers, rows)
}

// printCampaignSubsequencesTable renders campaign subsequences as a table.
func printCampaignSubsequencesTable(items []api.CampaignSubsequence) {
	if len(items) == 0 {
		fmt.Println("No campaign subsequences found.")
		return
	}
	headers := []string{"ID", "NAME", "CAMPAIGN ID", "STATUS", "CREATED"}
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = []string{
			item.ID,
			output.Truncate(item.Name, 36),
			output.Truncate(item.CampaignID, 24),
			campaignStatusLabel(item.Status),
			output.FormatTime(item.CreatedAt),
		}
	}
	output.PrintTable(headers, rows)
}
