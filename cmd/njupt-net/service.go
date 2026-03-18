package main

import (
	"fmt"
	"strings"

	"github.com/hicancan/njupt-net-cli/internal/selfservice"
	"github.com/hicancan/njupt-net-cli/internal/workflow"
	"github.com/spf13/cobra"
)

func newServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Self service commands on ZFW",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newServiceLoginCmd())
	cmd.AddCommand(newServiceBindCmd())
	cmd.AddCommand(newServiceMigrateCmd())
	return cmd
}

func newServiceLoginCmd() *cobra.Command {
	var account string
	var password string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Self service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if account == "" {
				return fmt.Errorf("service login: --account is required")
			}
			if password == "" {
				return fmt.Errorf("service login: --password is required")
			}
			if globalSession == nil {
				return fmt.Errorf("service login: session not initialized")
			}

			client := selfservice.NewClient(globalSession)
			if err := client.Login(cmd.Context(), account, password); err != nil {
				return err
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Self login succeeded.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&account, "account", "u", "", "Self service account")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Self service password")

	return cmd
}

func newServiceBindCmd() *cobra.Command {
	var fld1 string
	var fld2 string
	var fld3 string
	var fld4 string

	cmd := &cobra.Command{
		Use:   "bind",
		Short: "Bind operator credentials",
		Long:  "Submit FLDEXTRA fields and verify with readback to avoid false HTTP 200 positives.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if globalSession == nil {
				return fmt.Errorf("service bind: session not initialized")
			}

			fields := map[string]string{}
			if strings.TrimSpace(fld1) != "" {
				fields["FLDEXTRA1"] = fld1
			}
			if strings.TrimSpace(fld2) != "" {
				fields["FLDEXTRA2"] = fld2
			}
			if strings.TrimSpace(fld3) != "" {
				fields["FLDEXTRA3"] = fld3
			}
			if strings.TrimSpace(fld4) != "" {
				fields["FLDEXTRA4"] = fld4
			}

			if len(fields) == 0 {
				return fmt.Errorf("service bind: at least one of --fld1/--fld2/--fld3/--fld4 must be provided")
			}

			client := selfservice.NewClient(globalSession)
			if err := client.BindOperator(cmd.Context(), fields); err != nil {
				return err
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Service bind succeeded.")
			return nil
		},
	}

	cmd.Flags().StringVar(&fld1, "fld1", "", "Value for FLDEXTRA1 (telecom account)")
	cmd.Flags().StringVar(&fld2, "fld2", "", "Value for FLDEXTRA2 (telecom password)")
	cmd.Flags().StringVar(&fld3, "fld3", "", "Value for FLDEXTRA3 (mobile account)")
	cmd.Flags().StringVar(&fld4, "fld4", "", "Value for FLDEXTRA4 (mobile password)")

	return cmd
}

func newServiceMigrateCmd() *cobra.Command {
	var fromUser string
	var fromPwd string
	var toUser string
	var toPwd string
	var fld1 string
	var fld2 string
	var fld3 string
	var fld4 string

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate broadband binding across accounts",
		Long:  "Execute high-level migration: source login -> source unbind -> target login -> target bind.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if globalSession == nil {
				return fmt.Errorf("service migrate: session not initialized")
			}
			if strings.TrimSpace(fromUser) == "" || strings.TrimSpace(fromPwd) == "" {
				return fmt.Errorf("service migrate: --from-user and --from-password are required")
			}
			if strings.TrimSpace(toUser) == "" || strings.TrimSpace(toPwd) == "" {
				return fmt.Errorf("service migrate: --to-user and --to-password are required")
			}

			targetFields := map[string]string{}
			if strings.TrimSpace(fld1) != "" {
				targetFields["FLDEXTRA1"] = fld1
			}
			if strings.TrimSpace(fld2) != "" {
				targetFields["FLDEXTRA2"] = fld2
			}
			if strings.TrimSpace(fld3) != "" {
				targetFields["FLDEXTRA3"] = fld3
			}
			if strings.TrimSpace(fld4) != "" {
				targetFields["FLDEXTRA4"] = fld4
			}
			if len(targetFields) == 0 {
				return fmt.Errorf("service migrate: at least one of --fld1/--fld2/--fld3/--fld4 must be provided")
			}

			if err := workflow.MigrateBroadband(
				cmd.Context(),
				globalSession,
				fromUser,
				fromPwd,
				toUser,
				toPwd,
				targetFields,
			); err != nil {
				return err
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Service migrate succeeded.")
			return nil
		},
	}

	cmd.Flags().StringVar(&fromUser, "from-user", "", "Source account for unbind step")
	cmd.Flags().StringVar(&fromPwd, "from-password", "", "Source account password")
	cmd.Flags().StringVar(&toUser, "to-user", "", "Target account for bind step")
	cmd.Flags().StringVar(&toPwd, "to-password", "", "Target account password")
	cmd.Flags().StringVar(&fld1, "fld1", "", "Target FLDEXTRA1 value")
	cmd.Flags().StringVar(&fld2, "fld2", "", "Target FLDEXTRA2 value")
	cmd.Flags().StringVar(&fld3, "fld3", "", "Target FLDEXTRA3 value")
	cmd.Flags().StringVar(&fld4, "fld4", "", "Target FLDEXTRA4 value")

	return cmd
}
