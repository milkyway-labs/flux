package parse

import (
	"github.com/spf13/cobra"
)

func NewParseCmd() *cobra.Command {
	parseCmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse some data without the need to re-syncing the whole database from scratch",
	}

	parseCmd.AddCommand(NewParseBlocksRangeCmd())

	return parseCmd
}
