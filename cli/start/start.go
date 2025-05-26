package start

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/milkyway-labs/chain-indexer/cli/types"
	"github.com/milkyway-labs/chain-indexer/indexer/builder"
	"github.com/spf13/cobra"
)

func NewStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start parsing the blockchain data",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdCtx := types.GetCmdContext(cmd)
			return startParsing(cmd.Context(), cmdCtx)
		},
	}
}

func startParsing(ctx context.Context, cmdCtx *types.CliContext) error {
	indexersBuilder := builder.NewIndexersBuilder(cmdCtx.DatabasesManager, cmdCtx.NodesManager, cmdCtx.ModulesManager)

	cfg, err := cmdCtx.LoadConfig()
	if err != nil {
		return err
	}

	indexers, err := indexersBuilder.BuildAll(ctx, cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	waitGroup := sync.WaitGroup{}
	for _, indexer := range indexers {
		err := indexer.Start(ctx, &waitGroup)
		if err != nil {
			return fmt.Errorf("staring indexer %s, %w", indexer.GetName(), err)
		}
	}

	waitGroup.Wait()
	return nil
}
