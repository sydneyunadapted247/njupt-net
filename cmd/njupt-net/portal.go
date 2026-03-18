package main

import (
	"fmt"

	portalcore "github.com/hicancan/njupt-net-cli/internal/portal"
	"github.com/spf13/cobra"
)

func newPortalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "portal",
		Short: "Portal 802 authentication commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newPortalLoginCmd())
	return cmd
}

func newPortalLoginCmd() *cobra.Command {
	var account string
	var password string
	var ip string
	var isp string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login through Portal 802",
		Long:  "Execute Portal 802 login flow with one built-in passive cleanup retry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if account == "" {
				return fmt.Errorf("portal login: --account is required")
			}
			if password == "" {
				return fmt.Errorf("portal login: --password is required")
			}
			if ip == "" {
				return fmt.Errorf("portal login: --ip is required")
			}
			if globalSession == nil {
				return fmt.Errorf("portal login: session not initialized")
			}

			client := portalcore.NewClient(globalSession)
			if err := client.Login(cmd.Context(), account, password, ip, isp); err != nil {
				return err
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "\x1b[32mPortal 802 login succeeded.\x1b[0m")
			return nil
		},
	}

	cmd.Flags().StringVarP(&account, "account", "u", "", "Portal account")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Portal password")
	cmd.Flags().StringVarP(&ip, "ip", "i", "", "Current WLAN IPv4")
	cmd.Flags().StringVar(&isp, "isp", "", "ISP suffix route: telecom|unicom|mobile")

	return cmd
}
