package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/milkyway-labs/chain-indexer/cli/parse"
	"github.com/milkyway-labs/chain-indexer/cli/root"
	"github.com/milkyway-labs/chain-indexer/cli/start"
	"github.com/milkyway-labs/chain-indexer/cli/types"
)

func NewDefaultIndexerCLI(cliCtx *types.CliContext) *cobra.Command {
	ctx := context.Background()
	root := root.NewRootCommad(ctx, cliCtx)

	// Add the sub-commands
	root.AddCommand(start.NewStartCmd())
	root.AddCommand(parse.NewParseCmd())

	return root
}
