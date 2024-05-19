package benchmark

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var namespace = "benchmark"

func New() *cobra.Command {
	command := &cobra.Command{
		Use: "benchmark",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			connection, err := cmd.Flags().GetString("connection")
			if err != nil {
				return err
			}

			if connection == "" {
				return errors.New(`required flag(s) "connection" not set`)
			}

			return nil
		},
	}

	command.AddCommand(NewCleanup())
	command.AddCommand(NewPrepare())
	command.AddCommand(NewSeed())

	command.PersistentFlags().StringP("storage", "", "", "path to your benchmark data storage")
	command.MarkPersistentFlagRequired("storage")
	command.PersistentFlags().StringP("connection", "", os.Getenv("TEST_DATABASE_URI"), "database connections string")
	return command
}
