package logrunningoperation

import "github.com/spf13/cobra"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logrunningoperation",
		Short:   "Manage long-running operations",
		Aliases: []string{"operation", "o"},
	}
	return cmd
}
