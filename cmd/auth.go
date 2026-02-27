package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/config"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage Instantly authentication",
}

var authSetKeyCmd = &cobra.Command{
	Use:   "set-key <api-key>",
	Short: "Save an Instantly API key to the config file",
	Long: `Save an Instantly API key to the local config file.

Get your API key from: https://app.instantly.ai/app/settings/integrations

The key is stored at:
  macOS:   ~/Library/Application Support/instantly/config.json
  Linux:   ~/.config/instantly/config.json
  Windows: %AppData%\instantly\config.json

You can also set the INSTANTLY_API_KEY env var instead of using this command.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		if len(key) < 8 {
			return fmt.Errorf("API key looks too short")
		}
		if err := config.Save(&config.Config{APIKey: key}); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
		fmt.Printf("API key saved to %s\n", config.Path())
		fmt.Printf("Key: %s\n", maskOrEmpty(key))
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		fmt.Printf("Config: %s\n\n", config.Path())
		if envKey := os.Getenv("INSTANTLY_API_KEY"); envKey != "" {
			fmt.Println("Key source: INSTANTLY_API_KEY env var (takes priority over config)")
			fmt.Printf("Key:        %s\n", maskOrEmpty(envKey))
		} else if c.APIKey != "" {
			fmt.Println("Key source: config file")
			fmt.Printf("Key:        %s\n", maskOrEmpty(c.APIKey))
		} else {
			fmt.Println("Status: not authenticated")
			fmt.Printf("\nRun: instantly auth set-key <your-api-key>\n")
			fmt.Printf("Or:  export INSTANTLY_API_KEY=<your-api-key>\n")
		}
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove the saved API key from the config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Clear(); err != nil {
			return fmt.Errorf("removing config: %w", err)
		}
		fmt.Println("API key removed from config.")
		return nil
	},
}

func init() {
	authCmd.AddCommand(authSetKeyCmd, authStatusCmd, authLogoutCmd)
	rootCmd.AddCommand(authCmd)
}
