package migrate

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func NewDown() *cobra.Command {
	command := &cobra.Command{
		Use:  "down",
		Args: cobra.MatchAll(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			source, err := cmd.Flags().GetString("source")
			if err != nil {
				return err
			}
			database, err := cmd.Flags().GetString("database")
			if err != nil {
				return err
			}

			m, err := migrate.New(source, database)
			if err != nil {
				return err
			}

			if err := m.Down(); !errors.Is(err, migrate.ErrNoChange) {
				return err
			}
			return nil
		},
	}
	command.Flags().StringP("source", "s", "", "the path to a directory that contains your migration files")
	command.Flags().StringP("database", "d", "", "the database connection string")
	return command
}
