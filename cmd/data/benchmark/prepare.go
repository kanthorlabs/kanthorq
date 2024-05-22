package benchmark

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/kanthorq/core"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/cobra"
)

func NewPrepare() *cobra.Command {
	command := &cobra.Command{
		Use:  "prepare",
		Args: cobra.MatchAll(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			storage, err := cmd.Flags().GetString("storage")
			if err != nil {
				return err
			}
			topic, err := cmd.Flags().GetString("topic")
			if err != nil {
				return err
			}

			writer, err := cmd.Flags().GetInt("writer")
			if err != nil {
				return err
			}
			p := pool.New().WithMaxGoroutines(writer).WithErrors()

			count, err := cmd.Flags().GetInt64("count")
			if err != nil {
				return err
			}
			size, err := cmd.Flags().GetInt64("size")
			if err != nil {
				return err
			}
			var i int64
			for i < count {
				limit := utils.Min(count-i, size)
				from := i
				to := i + limit

				filename := fmt.Sprintf("%s/%s_%d_%d.csv", storage, namespace, from, to)
				p.Go(func() error {
					f, err := os.Create(filename)
					if err != nil {
						return err
					}
					defer f.Close()

					w := csv.NewWriter(f)
					w.Write((&core.Stream{}).Properties())

					rows := make([][]string, limit)
					for j := int64(0); j < limit; j++ {
						rows[j] = []string{topic, idx.New("evt")}
					}

					if err := w.WriteAll(rows); err != nil {
						return err
					}

					return nil
				})

				i += limit
			}

			return p.Wait()
		},
	}
	command.Flags().Int64P("count", "c", 1000000, "total record you want to prepare")
	command.Flags().Int64P("size", "s", 50000, "total record of each file")
	command.Flags().IntP("writer", "", 5, "set write concurrency")

	command.Flags().StringP("topic", "", os.Getenv("TEST_TOPIC"), "use custom topic")

	return command
}
