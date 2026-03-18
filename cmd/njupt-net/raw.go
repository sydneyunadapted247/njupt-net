package main

import (
	"fmt"
	"strings"

	"github.com/hicancan/njupt-net-cli/internal/core"
	"github.com/spf13/cobra"
)

func newRawCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raw",
		Short: "Raw probe commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newRawGetCmd())
	cmd.AddCommand(newRawPostCmd())

	return cmd
}

func newRawGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <path>",
		Short: "Send a raw GET request and print status/body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if globalSession == nil {
				return fmt.Errorf("raw get: session not initialized")
			}

			resp, err := globalSession.Get(cmd.Context(), args[0], core.RequestOptions{})
			if err != nil {
				return fmt.Errorf("raw get failed: %w", err)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Status: %d\n", resp.StatusCode)
			_, _ = cmd.OutOrStdout().Write(resp.Body)
			_, _ = fmt.Fprintln(cmd.OutOrStdout())
			return nil
		},
	}

	return cmd
}

func newRawPostCmd() *cobra.Command {
	var formPairs []string

	cmd := &cobra.Command{
		Use:   "post <path>",
		Short: "Send a raw form POST request and print status/body",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if globalSession == nil {
				return fmt.Errorf("raw post: session not initialized")
			}

			form, err := parseFormPairs(formPairs)
			if err != nil {
				return err
			}

			resp, err := globalSession.PostForm(cmd.Context(), args[0], core.RequestOptions{Form: form})
			if err != nil {
				return fmt.Errorf("raw post failed: %w", err)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Status: %d\n", resp.StatusCode)
			_, _ = cmd.OutOrStdout().Write(resp.Body)
			_, _ = fmt.Fprintln(cmd.OutOrStdout())
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&formPairs, "form", "f", nil, "Form field pair in key=value format; can be repeated")
	return cmd
}

func parseFormPairs(pairs []string) (map[string]string, error) {
	form := map[string]string{}
	for _, pair := range pairs {
		trimmed := strings.TrimSpace(pair)
		if trimmed == "" {
			continue
		}
		idx := strings.Index(trimmed, "=")
		if idx <= 0 {
			return nil, fmt.Errorf("raw post: invalid --form value %q, expected key=value", pair)
		}
		key := strings.TrimSpace(trimmed[:idx])
		value := trimmed[idx+1:]
		if key == "" {
			return nil, fmt.Errorf("raw post: empty form key in %q", pair)
		}
		form[key] = value
	}
	return form, nil
}
