package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/output"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage Instantly workspace and members",
}

// ---- workspace get ----

var workspaceGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get workspace details",
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetWorkspace()
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"ID", item.ID},
			{"Name", item.Name},
			{"Created", output.FormatTime(item.CreatedAt)},
			{"Updated", output.FormatTime(item.UpdatedAt)},
		})
		return nil
	},
}

// ---- workspace update ----

var workspaceUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update workspace settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name is required")
		}
		item, err := client.UpdateWorkspace(map[string]interface{}{"name": name})
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Workspace updated: %s\n", item.Name)
		return nil
	},
}

// ---- workspace member subcommands ----

var workspaceMemberCmd = &cobra.Command{
	Use:   "member",
	Short: "Manage workspace members",
}

var (
	workspaceMemberListLimit int
	workspaceMemberListSkip  int
)

var workspaceMemberListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspace members",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := buildParams(
			"limit", fmt.Sprintf("%d", workspaceMemberListLimit),
			"skip", fmt.Sprintf("%d", workspaceMemberListSkip),
		)
		items, _, err := client.ListWorkspaceMembers(params)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(items, output.IsPretty(cmd))
		}
		printWorkspaceMembersTable(items)
		return nil
	},
}

var workspaceMemberGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get details of a specific workspace member",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		item, err := client.GetWorkspaceMember(args[0])
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		output.PrintKeyValue([][]string{
			{"ID", item.ID},
			{"Email", item.Email},
			{"Role", item.Role},
			{"Status", item.Status},
			{"Created", output.FormatTime(item.CreatedAt)},
		})
		return nil
	},
}

var (
	workspaceMemberCreateEmail string
	workspaceMemberCreateRole  string
)

var workspaceMemberCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Invite a new workspace member",
	Long: `Invite a new member to the workspace.

Examples:
  instantly workspace member create --email user@domain.com --role member
  instantly workspace member create --email admin@domain.com --role admin`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if workspaceMemberCreateEmail == "" {
			return fmt.Errorf("--email is required")
		}
		payload := map[string]interface{}{
			"email": workspaceMemberCreateEmail,
		}
		if workspaceMemberCreateRole != "" {
			payload["role"] = workspaceMemberCreateRole
		}
		item, err := client.CreateWorkspaceMember(payload)
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Workspace member invited: %s\n", item.Email)
		fmt.Printf("ID: %s\n", item.ID)
		return nil
	},
}

var workspaceMemberUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a workspace member",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		role, _ := cmd.Flags().GetString("role")
		if role == "" {
			return fmt.Errorf("--role is required")
		}
		item, err := client.UpdateWorkspaceMember(args[0], map[string]interface{}{"role": role})
		if err != nil {
			return err
		}
		if output.IsJSON(cmd) {
			return output.PrintJSON(item, output.IsPretty(cmd))
		}
		fmt.Printf("Workspace member updated: %s (role: %s)\n", item.Email, item.Role)
		return nil
	},
}

var workspaceMemberDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Remove a workspace member",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteWorkspaceMember(args[0]); err != nil {
			return err
		}
		fmt.Printf("Workspace member %s removed.\n", args[0])
		return nil
	},
}

func init() {
	workspaceUpdateCmd.Flags().String("name", "", "New workspace name *(required)*")

	workspaceMemberListCmd.Flags().IntVar(&workspaceMemberListLimit, "limit", 20, "Maximum number of members to return")
	workspaceMemberListCmd.Flags().IntVar(&workspaceMemberListSkip, "skip", 0, "Number of members to skip")

	workspaceMemberCreateCmd.Flags().StringVar(&workspaceMemberCreateEmail, "email", "", "Member email *(required)*")
	workspaceMemberCreateCmd.Flags().StringVar(&workspaceMemberCreateRole, "role", "member", "Role (admin, member)")

	workspaceMemberUpdateCmd.Flags().String("role", "", "New role (admin, member) *(required)*")

	workspaceMemberCmd.AddCommand(
		workspaceMemberListCmd,
		workspaceMemberGetCmd,
		workspaceMemberCreateCmd,
		workspaceMemberUpdateCmd,
		workspaceMemberDeleteCmd,
	)
	workspaceCmd.AddCommand(
		workspaceGetCmd,
		workspaceUpdateCmd,
		workspaceMemberCmd,
	)
	rootCmd.AddCommand(workspaceCmd)
}
