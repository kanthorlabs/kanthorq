package seed

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

type ctx string

const connection ctx = "connection"

func New() *cobra.Command {
	command := &cobra.Command{
		Use: "seed",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			database, err := cmd.Flags().GetString("database")
			if err != nil {
				return err
			}
			conn, err := pgxpool.New(cmd.Context(), database)
			if err != nil {
				return err
			}
			cmd.SetContext(context.WithValue(cmd.Context(), connection, conn))

			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			conn := cmd.Context().Value(connection).(*pgxpool.Pool)
			conn.Close()
			return nil
		},
	}
	command.AddCommand(NewStream())
	command.AddCommand(NewConsumer())

	command.PersistentFlags().Bool("clean", false, "clean the database before seeding")
	command.PersistentFlags().String("database", os.Getenv("TEST_DATABASE_URI"), "connection string of database to seed")
	command.PersistentFlags().StringP("stream", "s", os.Getenv("TEST_STREAM"), "what stream to seed")
	command.PersistentFlags().StringP("topic", "t", os.Getenv("TEST_TOPIC"), "what topic to seed")

	return command
}
