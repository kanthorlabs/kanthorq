package migrate

import (
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/spf13/cobra"
)

func NewUp() *cobra.Command {
	command := &cobra.Command{
		Use:  "up",
		Args: cobra.MatchAll(cobra.NoArgs),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			step, err := cmd.Flags().GetInt("step")
			if err != nil {
				return err
			}

			if step < 1 {
				return errors.New("migrate up does not allow go backward")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			source, err := cmd.Flags().GetString("source")
			if err != nil {
				return err
			}
			database, err := cmd.Flags().GetString("database")
			if err != nil {
				return err
			}
			step, err := cmd.Flags().GetInt("step")
			if err != nil {
				return err
			}

			m, err := migrate.New(source, database)
			if err != nil {
				return err
			}

			if err := m.Steps(step); !errors.Is(err, migrate.ErrNoChange) && !errors.Is(err, os.ErrNotExist) {
				return err
			}
			return nil
		},
	}
	command.Flags().Int("step", 1, "step you want to go forward")

	return command
}
