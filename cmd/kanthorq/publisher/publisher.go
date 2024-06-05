package publisher

import "github.com/spf13/cobra"

func New() *cobra.Command {
	command := &cobra.Command{
		Use: "publisher",
	}

	command.AddCommand(Publish())
	return command
}
