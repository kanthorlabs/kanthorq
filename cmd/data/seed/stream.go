package seed

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kanthorlabs/common/commands"
	"github.com/kanthorlabs/kanthorq/queries"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/spf13/cobra"
)

func NewStream() *cobra.Command {
	command := &cobra.Command{
		Use: "stream",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Parent().PreRunE(cmd, args) != nil {
				return nil
			}

			clean, err := cmd.Flags().GetBool("clean")
			if err != nil {
				return err
			}

			if !clean {
				return nil
			}

			conn := cmd.Context().Value(connection).(*pgxpool.Pool)
			return testify.QueryTruncateStream()(cmd.Context(), conn)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			database, err := cmd.Flags().GetString("database")
			if err != nil {
				return err
			}
			conn, err := pgxpool.New(cmd.Context(), database)
			if err != nil {
				return err
			}
			defer conn.Close()

			stream, err := cmd.Flags().GetString("stream")
			if err != nil {
				return err
			}
			if _, err := queries.EnsureStream(stream)(cmd.Context(), conn); err != nil {
				return err
			}

			topic, err := cmd.Flags().GetString("topic")
			if err != nil {
				return err
			}
			count, err := cmd.Flags().GetInt("count")
			if err != nil {
				return err
			}
			if err := testify.SeedStreamEvents(cmd.Context(), conn, stream, topic, count); err != nil {
				return err
			}

			if verbose, err := cmd.Flags().GetBool("verbose"); err == nil && verbose {
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"#", "Stream", "Topic", "Records"})
				t.AppendRows([]table.Row{
					{1, stream, topic, count},
				})
				t.Render()
			}
			return nil
		},
		PostRunE: commands.PostRunE(),
	}

	command.Flags().Int("count", 1000000, "number of recods to seed")

	return command
}
