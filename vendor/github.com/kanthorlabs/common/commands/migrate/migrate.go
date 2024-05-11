package migrate

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	command := &cobra.Command{
		Use:   "migrate",
		Short: "migrate our data up or down",
	}

	command.AddCommand(NewUp())
	command.AddCommand(NewDown())
	return command
}
