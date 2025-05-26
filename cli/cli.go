package cli

import (
	"context"

	"github.com/milkyway-labs/chain-indexer/cli/root"
	"github.com/milkyway-labs/chain-indexer/cli/start"
	"github.com/milkyway-labs/chain-indexer/cli/types"
	"github.com/spf13/cobra"
)

func NewSimpleCLI(cliCtx *types.CliContext) *cobra.Command {
	ctx := context.Background()
	root := root.GetRootCommad(ctx, cliCtx)
	root.AddCommand(start.NewStartCmd())

	return root
}
