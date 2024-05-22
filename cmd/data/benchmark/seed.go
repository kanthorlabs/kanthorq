package benchmark

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanthorlabs/kanthorq/core"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/cobra"
)

func NewSeed() *cobra.Command {
	command := &cobra.Command{
		Use:  "seed",
		Args: cobra.MatchAll(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			connection, err := cmd.Flags().GetString("connection")
			if err != nil {
				return err
			}
			conn, err := pgxpool.New(cmd.Context(), connection)
			if err != nil {
				return err
			}
			defer conn.Close()

			// ensure stream is exist
			streamName, err := cmd.Flags().GetString("stream")
			if err != nil {
				return err
			}
			stream, err := core.QueryStreamEnsure(streamName)(cmd.Context(), conn)
			if err != nil {
				return err
			}

			storage, err := cmd.Flags().GetString("storage")
			if err != nil {
				return err
			}
			files, err := filepath.Glob(fmt.Sprintf("%s/%s_*.csv", storage, namespace))
			if err != nil {
				return err
			}

			writer, err := cmd.Flags().GetInt("writer")
			if err != nil {
				return err
			}
			p := pool.New().WithMaxGoroutines(writer).WithErrors()

			for _, f := range files {
				statement := fmt.Sprintf(`COPY %s FROM '%s' DELIMITER ',' CSV HEADER;`, core.StreamEventCollection(stream.Name), f)
				fmt.Println(statement)
				p.Go(func() error {
					if _, err = conn.Exec(cmd.Context(), statement); err != nil {
						return err
					}

					return nil
				})
			}

			return p.Wait()
		},
	}

	command.Flags().IntP("writer", "", 5, "set write concurrency")
	command.Flags().StringP("stream", "", os.Getenv("TEST_STREAM"), "define your testing stream name")

	return command
}
