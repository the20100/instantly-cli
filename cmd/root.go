package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/the20100/instantly-cli/internal/api"
	"github.com/the20100/instantly-cli/internal/config"
)

var (
	jsonFlag   bool
	prettyFlag bool
	client     *api.Client
	cfg        *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "instantly",
	Short: "Instantly CLI — manage Instantly.ai via the API",
	Long: `instantly is a CLI tool for the Instantly.ai API.

It outputs JSON when piped (for agent use) and human-readable tables in a terminal.

Token resolution order:
  1. INSTANTLY_API_KEY env var
  2. Config file  (~/.config/instantly/config.json  via: instantly auth set-key)

Examples:
  instantly auth set-key <your-api-key>
  instantly campaign list
  instantly lead list --campaign-id <id>
  instantly account list`,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Force JSON output")
	rootCmd.PersistentFlags().BoolVar(&prettyFlag, "pretty", false, "Force pretty-printed JSON output (implies --json)")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if isAuthCommand(cmd) || cmd.Name() == "info" || cmd.Name() == "update" {
			return nil
		}
		key, err := resolveAPIKey()
		if err != nil {
			return err
		}
		client = api.NewClient(key)
		return nil
	}

	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show tool info: config path, auth status, and environment",
	Run: func(cmd *cobra.Command, args []string) {
		printInfo()
	},
}

func printInfo() {
	fmt.Printf("instantly — Instantly.ai CLI\n\n")
	exe, _ := os.Executable()
	fmt.Printf("  binary:  %s\n", exe)
	fmt.Printf("  os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()
	fmt.Println("  config paths by OS:")
	fmt.Printf("    macOS:    ~/Library/Application Support/instantly/config.json\n")
	fmt.Printf("    Linux:    ~/.config/instantly/config.json\n")
	fmt.Printf("    Windows:  %%AppData%%\\instantly\\config.json\n")
	fmt.Printf("  config:   %s\n", config.Path())
	fmt.Println()
	fmt.Printf("    INSTANTLY_API_KEY = %s\n", maskOrEmpty(os.Getenv("INSTANTLY_API_KEY")))
}

func maskOrEmpty(v string) string {
	if v == "" {
		return "(not set)"
	}
	if len(v) <= 8 {
		return "***"
	}
	return v[:4] + "..." + v[len(v)-4:]
}

func resolveAPIKey() (string, error) {
	if k := os.Getenv("INSTANTLY_API_KEY"); k != "" {
		return k, nil
	}
	var err error
	cfg, err = config.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}
	if cfg.APIKey != "" {
		return cfg.APIKey, nil
	}
	return "", fmt.Errorf("not authenticated — run: instantly auth set-key\nor set INSTANTLY_API_KEY env var")
}

func isAuthCommand(cmd *cobra.Command) bool {
	if cmd.Name() == "auth" {
		return true
	}
	p := cmd.Parent()
	for p != nil {
		if p.Name() == "auth" {
			return true
		}
		p = p.Parent()
	}
	return false
}
