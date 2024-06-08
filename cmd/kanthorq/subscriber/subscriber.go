package subscriber

import "github.com/spf13/cobra"

func New() *cobra.Command {
	command := &cobra.Command{
		Use: "subscriber",
	}

	command.AddCommand(Subscribe())
	return command
}
