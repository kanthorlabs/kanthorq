package benchmark

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/kanthorlabs/kanthorq/core"
	"github.com/spf13/cobra"
)

func NewCleanup() *cobra.Command {
	command := &cobra.Command{
		Use:  "cleanup",
		Args: cobra.MatchAll(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			// clean files
			storage, err := cmd.Flags().GetString("storage")
			if err != nil {
				return err
			}

			// clean files
			patterns := []string{
				fmt.Sprintf("%s/%s_*.csv", storage, namespace),
			}
			for _, pattern := range patterns {
				files, err := filepath.Glob(pattern)
				if err != nil {
					return err
				}
				for _, f := range files {
					if err := os.Remove(f); err != nil {
						return err
					}
				}
			}

			// cleanup database
			if truncate, err := cmd.Flags().GetBool("database"); err != nil && truncate {
				connection, err := cmd.Flags().GetString("connection")
				if err != nil {
					return err
				}
				conn, err := pgx.Connect(cmd.Context(), connection)
				defer conn.Close(cmd.Context())

				statement := "TRUNCATE TABLE public.%s CONTINUE IDENTITY RESTRICT;"

				_, err = conn.Exec(cmd.Context(), fmt.Sprintf(statement, core.CollectionStream))
				if err != nil {
					return err
				}
				_, err = conn.Exec(cmd.Context(), fmt.Sprintf(statement, core.CollectionConsumer))
				if err != nil {
					return err
				}
				_, err = conn.Exec(cmd.Context(), fmt.Sprintf(statement, core.CollectionConsumerJob))
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	command.Flags().BoolP("database", "", os.Getenv("TEST_BENCHMARK_CLEANUP_DATABASE") != "", "decide whether we should cleanup database or not")

	return command
}
