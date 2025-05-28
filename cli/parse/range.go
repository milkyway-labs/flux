package parse

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/spf13/cobra"

	clitypes "github.com/milkyway-labs/chain-indexer/cli/types"
	"github.com/milkyway-labs/chain-indexer/indexer"
	"github.com/milkyway-labs/chain-indexer/types"
)

func NewParseBlocksRangeCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "range [indexer-name] [start-height] (end-height)",
		Short: "Re-parse a range of blocks, if the end-height is not provided, we sync only the block at the start-height",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := clitypes.GetCmdContext(cmd)
			indexerName := args[0]
			startHeight, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid start-height: %w", err)
			}

			endHeight := startHeight
			if len(args) == 3 {
				parsed, err := strconv.ParseUint(args[2], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid end-height: %w", err)
				}
				endHeight = parsed
			}

			return parseBlocksRange(cmd.Context(), cliCtx, indexerName, types.Height(startHeight), types.Height(endHeight))
		},
	}

	return rootCmd
}

func parseBlocksRange(
	ctx context.Context,
	cmdCtx *clitypes.CliContext,
	indexerName string,
	startHeight types.Height,
	endHeight types.Height,
) error {
	// Load the indexer config
	cfg, err := cmdCtx.LoadConfig()
	if err != nil {
		return err
	}

	// Build the requested requestedIndexer
	requestedIndexer, err := cmdCtx.IndexersBuilder.BuildByName(ctx, cfg, indexerName)
	if err != nil {
		return err
	}

	requestedIndexer.WithCustomHeightProducer(
		indexer.NewRangeHeightProducer(startHeight, endHeight),
	)

	// Start indexing the requested range
	wg := sync.WaitGroup{}
	requestedIndexer.Start(ctx, &wg)
	wg.Wait()

	return nil
}
