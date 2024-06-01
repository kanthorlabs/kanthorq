package seed

import (
	"errors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/kanthorlabs/common/commands"
	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/api"
	"github.com/kanthorlabs/kanthorq/testify"
	"github.com/spf13/cobra"
)

func NewConsumer() *cobra.Command {
	command := &cobra.Command{
		Use: "consumer",
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
			return testify.QueryTruncateConsumer()(cmd.Context(), conn)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			conn := cmd.Context().Value(connection).(*pgxpool.Pool)

			tx, err := conn.Begin(cmd.Context())
			if err != nil {
				return err
			}

			stream, err := cmd.Flags().GetString("stream")
			if err != nil {
				return err
			}
			topic, err := cmd.Flags().GetString("topic")
			if err != nil {
				return err
			}
			ensure, err := (&api.ConsumerEnsureReq{stream, idx.New("job"), topic}).Do(cmd.Context(), tx)
			if err != nil {
				return errors.Join(err, tx.Rollback(cmd.Context()))
			}

			count, err := cmd.Flags().GetInt("count")
			if err != nil {
				return err
			}
			cursor, err := (&api.ConsumerPullReq{ensure.Consumer, count}).Do(cmd.Context(), tx)
			if err != nil {
				return errors.Join(err, tx.Rollback(cmd.Context()))
			}

			if err := tx.Commit(cmd.Context()); err != nil {
				return err
			}

			if verbose, err := cmd.Flags().GetBool("verbose"); err == nil && verbose {
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"#", "Consumer", "Stream", "Topic", "Cursor", "Records"})
				t.AppendRows([]table.Row{
					{1, ensure.Consumer.Name, stream, topic, cursor, count},
				})
				t.Render()
			}
			return nil
		},
		PostRunE: commands.PostRunE(),
	}

	command.Flags().Int("count", kanthorq.ConsumerPullSize, "number of recods to seed")

	return command
}
