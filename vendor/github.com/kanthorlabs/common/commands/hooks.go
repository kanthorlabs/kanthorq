package commands

import "github.com/spf13/cobra"

func Noop() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}
}

func PreRunE() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if cmd.Parent().PreRunE(cmd, args) != nil {
			return nil
		}
		return nil
	}
}

func PostRunE() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if cmd.Parent().PostRunE(cmd, args) != nil {
			return nil
		}
		return nil
	}
}
