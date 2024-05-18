package migrate

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	command := &cobra.Command{
		Use:   "migrate",
		Short: "migrate our data up or down",
	}

	command.PersistentFlags().StringP("source", "s", "", "the path to a directory that contains your migration files")
	command.PersistentFlags().StringP("database", "d", "", "the database connection string")
	command.MarkPersistentFlagRequired("source")
	command.MarkPersistentFlagRequired("database")

	command.AddCommand(NewUp())
	command.AddCommand(NewDown())
	return command
}
