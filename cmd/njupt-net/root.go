package main

import (
	"fmt"

	"github.com/hicancan/njupt-net-cli/internal/core"
	"github.com/hicancan/njupt-net-cli/internal/httpx"
	"github.com/spf13/cobra"
)

const defaultSelfBaseURL = "http://10.10.244.240:8080"

// rootOptions keeps process-wide CLI options.
// These flags are intentionally lightweight placeholders; business logic will consume them later.
type rootOptions struct {
	ConfigPath  string
	OutputMode  string
	AssumeYes   bool
	InsecureTLS bool
}

var (
	rootOpts      = &rootOptions{}
	globalSession core.SessionClient
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "njupt-net",
		Short: "NJUPT network CLI kernel",
		Long:  "Kernel-first CLI scaffold for NJUPT network workflows.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if globalSession != nil {
				return nil
			}

			session, err := httpx.NewDefaultSessionClient(defaultSelfBaseURL)
			if err != nil {
				return fmt.Errorf("initialize session client failed: %w", err)
			}
			globalSession = session
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Global placeholder flags required by the outer CLI layer.
	cmd.PersistentFlags().StringVar(&rootOpts.ConfigPath, "config", "", "Path to config file")
	cmd.PersistentFlags().StringVar(&rootOpts.OutputMode, "output", "human", "Output mode: human|json")
	cmd.PersistentFlags().BoolVar(&rootOpts.AssumeYes, "yes", false, "Skip confirmation prompts for side-effect operations")
	cmd.PersistentFlags().BoolVar(&rootOpts.InsecureTLS, "insecure-tls", true, "Allow insecure TLS certificates")

	cmd.AddCommand(newServiceCmd())
	cmd.AddCommand(newPortalCmd())
	cmd.AddCommand(newRawCmd())

	return cmd
}
