package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/milkyway-labs/flux/cli/parse"
	"github.com/milkyway-labs/flux/cli/root"
	"github.com/milkyway-labs/flux/cli/start"
	"github.com/milkyway-labs/flux/cli/types"
)

func NewDefaultIndexerCLI(cliCtx *types.CliContext) *cobra.Command {
	ctx := context.Background()
	rootCmd := root.NewRootCommad(ctx, cliCtx)

	// Add the sub-commands
	rootCmd.AddCommand(start.NewStartCmd())
	rootCmd.AddCommand(parse.NewParseCmd())

	return rootCmd
}
