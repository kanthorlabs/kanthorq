package show

import (
	"fmt"

	"github.com/spf13/cobra"
)

func New(logo, version string) *cobra.Command {
	command := &cobra.Command{
		Use:   "show",
		Short: "Show KanthorQ system information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(logo)
			fmt.Printf("version: %s\n", version)

		},
	}
	return command
}
